package stringz

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"regexp"
	"strconv"
	"strings"
)

func GetVersionAsInteger(version string) string {
	out := ""
	r := regexp.MustCompile(`(\d+)(?:\.(\d+))?(?:\.(\d+))?`)
	max := 999
	if ms := r.FindStringSubmatch(version); len(ms) > 0 {
		major, _ := strconv.Atoi(ms[1])
		minor, _ := strconv.Atoi(ms[2])
		patch, _ := strconv.Atoi(ms[3])
		if minor > max {
			minor = max
		}
		if patch > max {
			patch = max
		}
		if out = strings.TrimLeft(fmt.Sprintf("%d%03d%03d", major, minor, patch), "0"); out != "" {
			return out
		}
	}
	return "0"
}

// IndexOfNth
//
//	@s haystack
//	@substr needle
//	@nth nth needle from 1
func IndexOfNth(s string, substr string, nth int) int {
	if substr == "" {
		return 0
	}
	if nth < 1 {
		nth = 1
	}
	l := len(substr)
	out := -l
	for i := 0; i < nth; i++ {
		out += l
		if idx := strings.Index(s[out:], substr); idx == -1 {
			return -1
		} else {
			out += idx
		}
	}
	return out
}

// Print context is fail in go 1.21.3, we rarely use it so removed
// func PrintContextInternals(ctx interface{}, inner bool) {
// 	contextValues := reflect.ValueOf(ctx).Elem()
// 	contextKeys := reflect.TypeOf(ctx).Elem()

// 	if !inner {
// 		fmt.Printf("\n-----Fields for %s.%s-----\n", contextKeys.PkgPath(), contextKeys.Name())
// 	}
// 	if contextKeys.Kind() == reflect.Struct {
// 		for i := 0; i < contextValues.NumField(); i++ {
// 			reflectValue := contextValues.Field(i)
// 			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

// 			reflectField := contextKeys.Field(i)

// 			if reflectField.Name == "Context" {
// 				PrintContextInternals(reflectValue.Interface(), true)
// 			} else {
// 				fmt.Printf("%+v > %+v\n", reflectField.Name, reflectValue.Interface())
// 			}
// 		}
// 		fmt.Printf("-----context end-----\n")
// 	} else {
// 		fmt.Printf("-----context is empty (int)-----\n")
// 	}
// }

func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := make([]string, 0)
		for i := 0; i < len(v); i += 2 {
			if v[i] >= 0 && v[i+1] >= 0 { //Note: should always but to be safe
				groups = append(groups, str[v[i]:v[i+1]])
			}
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

func Snake2Title(s string) string {
	mUs := regexp.MustCompile(`_+`)
	mC1 := regexp.MustCompile(`(^|\s)\w`)

	return mC1.ReplaceAllStringFunc(mUs.ReplaceAllLiteralString(strings.Trim(s, "_"), " "), func(w string) string {
		return strings.ToUpper(w)
	})
}

func ToCrc32(v interface{}) string {
	return strings.ToUpper(strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v", v)))), 16))
}

func ToJson(obj interface{}, indent string) string {
	var bs []byte
	if indent != "" {
		bs, _ = json.MarshalIndent(obj, "", indent)
	} else {
		bs, _ = json.Marshal(obj)
	}
	return string(bs)
}
