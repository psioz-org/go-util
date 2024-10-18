package conv

//
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/zev-zakaryan/go-util/stringz"
)

// GetItem return generic type of specified keys. Allow only map[string]any/[]any parent as it's used with nerdgraph structure.
func GetItem[T any](obj any, keys string, sep string) (T, error) {
	out := obj
	ss := strings.Split(keys, sep)
	for _, k := range ss {
		if ki, err := strconv.Atoi(k); err == nil {
			if objArr, ok := out.([]any); ok {
				out = objArr[ki]
			} else {
				var zero T
				return zero, fmt.Errorf("fail cast to array to get key %v (%v), obj %v", ki, keys, out)
			}
		} else {
			if objMap, ok := out.(map[string]any); ok {
				out = objMap[k]
			} else {
				var zero T
				return zero, fmt.Errorf("fail cast to map to get key %v (%v), obj %v", k, keys, out)
			}
		}

	}

	if outT, ok := out.(T); ok {
		return outT, nil
	} else {
		return To[T](out)
	}
}

func GetItems(obj any, keys string, opts ...Option) []any {
	optsMap := make(map[Option]struct{}, 0)
	for _, opt := range opts {
		optsMap[opt] = struct{}{}
	}
	out := make([]any, 0)
	getItems(obj, strings.Split(keys, "."), &out, optsMap)
	return out
}

type Option string

const (
	DotAlternative        = "․"
	OptOmitNoValue Option = "omit no value" //We ignore undefined if exists such as javascript, in go completely ignore null
)

func getKeys(mapKeys []reflect.Value) (mKeys []string) {
	for _, e := range mapKeys {
		mKeys = append(mKeys, e.Interface().(string))
	}
	sort.Strings(mKeys) //Need sort in golang for stable iteration order
	return
}

func getItemsSlice(obj any, keys []string, out *[]any, opts map[Option]struct{}, key string) {
	sl := reflect.ValueOf(obj)
	if key == "#" || key == "#v" {
		for i := 0; i < sl.Len(); i++ {
			getItems(sl.Index(i).Interface(), keys, out, opts)
		}
		return
	}
	if key == "#k" {
		for i := 0; i < sl.Len(); i++ {
			getItems(i, keys, out, opts)
		}
		return
	}
	if i, err := strconv.Atoi(key); err == nil && i < sl.Len() {
		getItems(sl.Index(i).Interface(), keys, out, opts)
	} else if len(keys) == 0 { //get null for final missing by default
		getItems(nil, keys, out, opts)
	}
}
func getItemsMap(obj any, keys []string, out *[]any, opts map[Option]struct{}, key string) {
	m := reflect.ValueOf(obj)
	if key == "#" || key == "#v" {
		mKeys := getKeys(m.MapKeys()) //Need sort in golang for stable iteration order
		for _, k := range mKeys {
			getItems(m.MapIndex(reflect.ValueOf(k)).Interface(), keys, out, opts)
		}
		return
	}
	if key == "#k" {
		mKeys := getKeys(m.MapKeys()) //Need sort in golang for stable iteration order
		for _, k := range mKeys {
			getItems(k, keys, out, opts)
		}
		return
	}
	if strings.HasPrefix(key, "^") {
		key = strings.ReplaceAll(key, DotAlternative, ".")
		reg, err := regexp.Compile(key)
		if err != nil {
			return
		}
		mKeys := getKeys(m.MapKeys()) //Need sort in golang for stable iteration order
		for _, k := range mKeys {
			if reg.Match([]byte(k)) {
				getItems(m.MapIndex(reflect.ValueOf(k)).Interface(), keys, out, opts)
			}
		}
		return
	}
	v := m.MapIndex(reflect.ValueOf(key))
	if v.IsValid() {
		getItems(v.Interface(), keys, out, opts)
	} else if len(keys) == 0 { //get null for final missing by default
		getItems(nil, keys, out, opts)
	}
}

func getItems(obj any, keys []string, out *[]any, opts map[Option]struct{}) {
	if len(keys) == 0 {
		if _, exists := opts[OptOmitNoValue]; exists && obj == nil {
			return
		}
		*out = append(*out, obj)
		return
	}
	key := keys[0]
	keys = keys[1:]
	//Note: If you have a pointer to a slice instead of a slice, you'll need to use Elem() to get the underlying value. e.g. reflect.TypeOf(reflect.ValueOf(t).Elem().Interface()).Kind()
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice:
		getItemsSlice(obj, keys, out, opts, key)
		return
	case reflect.Map:
		getItemsMap(obj, keys, out, opts, key)
		return
	}
}

func Ternary[T any](b bool, vt T, vf T) T {
	if b {
		return vt
	}
	return vf
}

func toBool(obj any) any {
	s := str(obj)
	switch s {
	case "", "<nil>": //Custom to false instead of error
		return false
	default:
		if v, err := strconv.ParseBool(s); err == nil {
			return v
		}
		return true
	}
}

func toString(obj any) (v any) {
	//case func(),struct{} match nothing even using before any. any match map,struct{},error,any (e.g. int also match) so it's not useful
	switch objV := obj.(type) {
	case nil, bool, complex64, complex128, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, string:
		v = fmt.Sprintf("%v", obj)
	case error:
		v = objV.Error()
	case []byte:
		// Can't convert to array because len must be constant so we change type
		objBs := make([]uint16, len(objV))
		for i1, v1 := range objV {
			objBs[i1] = uint16(v1)
		}
		v = stringz.ToJson(objBs, "")
	default: //any=any. []any will match only []any (also slice only, not array). The same as map[any]any.
		rv := reflect.ValueOf(obj)
		switch rv.Kind() {
		case reflect.Array, reflect.Map, reflect.Slice: //Beware marshal of byte slice will be base64, we handle above
			v = stringz.ToJson(obj, "")
		case reflect.Func:
			if v = getFuncBodyString(getFuncAST(getFuncInfo(obj))); v == "" {
				v = fmt.Sprintf("%T", obj)
			}
		case reflect.Struct: //match struct{} (instance)
			v = stringz.ToJson(obj, "")
		default: //reflect.Interface cant be send as param
			if v = stringz.ToJson(obj, ""); v == "" {
				v = fmt.Sprintf("%+v", obj)
			}
		}
	}
	return
}

func toObject[T any](obj any, out T) (v any, err error) {
	s := str(obj)
	switch any(&out).(type) {
	case *error: //Unlike other instance, "out error" starts with nil so any(out).(type) = nil. We still can check pointer.
		v = errors.New(s)
	default:
		if s == "" {
			v = out
			return
		}
		if s == "<nil>" {
			s = "null"
		}
		if err = json.Unmarshal([]byte(s), &out); err != nil {
			err = fmt.Errorf("fail cast to result type %v, from %v: %v", reflect.TypeOf(&out), reflect.TypeOf(obj), obj)
		} else {
			v = out
		}
	}
	return
}

func To[T any](obj any) (out T, err error) {
	if v, ok := obj.(T); ok {
		return v, nil
	}
	var v any
	var vc128 complex128
	var vi64 int64
	var vui64 uint64
	var vf64 float64

	switch any(out).(type) {
	case bool:
		v = toBool(obj)
	case complex64:
		vc128, err = strconv.ParseComplex(str(obj), 64)
		v = complex64(vc128)
	case complex128:
		v, err = strconv.ParseComplex(str(obj), 128)
	case float32:
		vf64, err = strconv.ParseFloat(str(obj), 32)
		v = float32(vf64)
	case float64:
		v, err = strconv.ParseFloat(str(obj), 64)
	case int:
		vi64, err = strconv.ParseInt(str(obj), 10, strconv.IntSize)
		v = int(vi64)
	case int8:
		vi64, err = strconv.ParseInt(str(obj), 10, 8)
		v = int8(vi64)
	case int16:
		vi64, err = strconv.ParseInt(str(obj), 10, 16)
		v = int16(vi64)
	case int32: //=rune
		vi64, err = strconv.ParseInt(str(obj), 10, 32)
		v = int32(vi64)
	case int64:
		v, err = strconv.ParseInt(str(obj), 10, 64)
	case uint:
		vui64, err = strconv.ParseUint(str(obj), 10, strconv.IntSize)
		v = uint(vui64)
	case uint8: //=byte
		vui64, err = strconv.ParseUint(str(obj), 10, 8)
		v = uint8(vui64)
	case uint16:
		vui64, err = strconv.ParseUint(str(obj), 10, 16)
		v = uint16(vui64)
	case uint32:
		vui64, err = strconv.ParseUint(str(obj), 10, 32)
		v = uint32(vui64)
	case uint64:
		v, err = strconv.ParseUint(str(obj), 10, 64)
	case string:
		v = toString(obj)
	case uintptr: //Can't be default, will error with Unmarshal "&out"
		err = fmt.Errorf("fail cast to result type %v, from %v: %v", reflect.TypeOf(&out), reflect.TypeOf(obj), obj)
	default: //case nil (error)&case <no match> (any instance). We ignore uintptr
		v, err = toObject(obj, out)
	}
	if err == nil && v != nil { //cant cast v=nil to (any), just dont set
		out = v.(T) //We don't check ok, to make sure not assign invalid type to v and always set err if error
	}
	return
}

func ToForce[T any](obj any) (out T) {
	out, _ = To[T](obj)
	return
}

func str(obj any) string {
	switch objV := obj.(type) {
	case string:
		return objV
	default:
		return fmt.Sprintf("%v", obj)
	}
}

func getFuncInfo(f any) (name, file string) {
	pc := reflect.ValueOf(f).Pointer()
	if fn := runtime.FuncForPC(pc); fn != nil {
		name = fn.Name()
		if i := strings.LastIndexByte(name, '/'); i >= 0 { //Path always sep by '/' irrespective of the OS
			name = name[i+1:]
		}
		//Format is package(.some parent in file).funcname e.g. area function of Shape interface will be Shape.area,
		//area function of Rect struct will be Rect.area-fm. For star, it'll be (*Rect).area2-fm
		if i := strings.IndexByte(name, '.'); i >= 0 {
			name = name[i+1:]
		}
		name = strings.TrimSuffix(name, "-fm")
		file, _ = fn.FileLine(pc)
	}
	return
}

func getFuncAST(funcname, filename string) (_ *ast.FuncDecl, _ *token.FileSet) {
	//Note function from object instance did not return file and line in new golang version
	//We still handle in case it's back as old behavior
	if filename == "<autogenerated>" {
		return
	}
	fs := token.NewFileSet()
	funcname = "." + funcname
	var file *ast.File
	var err error
	if file, err = parser.ParseFile(fs, filename, nil, 0); err != nil {
		return
	}
	for _, d := range file.Decls {
		if f, ok := d.(*ast.FuncDecl); ok && strings.HasSuffix(funcname, "."+f.Name.Name) && (f.Recv == nil || matchStructFunc(f, funcname)) {
			return f, fs
		}
	}
	return
}

func matchStructFunc(f *ast.FuncDecl, funcname string) bool {
	for _, l := range f.Recv.List { // fn is *ast.FuncDecl
		var fullname string
		switch vToken := l.Type.(type) {
		case *ast.StarExpr:
			fullname = fmt.Sprintf(".(*%v).%v", vToken.X, f.Name.Name)
		case *ast.Ident:
			fullname = fmt.Sprintf(".%v.%v", vToken.Name, f.Name.Name)
		}
		if fullname == funcname {
			return true
		}
	}
	return false
}

func getFuncBodyString(f any, fs *token.FileSet) string {
	if fs == nil {
		return ""
	}
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fs, f); err != nil {
		return ""
	}
	return buf.String()
}
