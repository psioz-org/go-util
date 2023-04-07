package zz

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetItem return generic type of specified keys. Allow only map[string]interface{}/[]interface{} parent as it's used with nerdgraph structure.
func GetItem[T any](obj interface{}, keys string, sep string) (T, error) {
	out := obj
	ss := strings.Split(keys, sep)
	for _, k := range ss {
		if ki, err := strconv.Atoi(k); err == nil {
			if objArr, ok := out.([]interface{}); ok {
				out = objArr[ki]
			} else {
				var zero T
				return zero, fmt.Errorf("fail cast to array to get key %v (%v), obj %v", ki, keys, out)
			}
		} else {
			if objMap, ok := out.(map[string]interface{}); ok {
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

func Ternary[T any](b bool, vt T, vf T) T {
	if b {
		return vt
	}
	return vf
}

func To[T any](obj interface{}) (out T, err error) {
	if v, ok := obj.(T); ok {
		return v, nil
	}
	var v interface{}
	var vc128 complex128
	var vi64 int64
	var vui64 uint64
	var vf64 float64
	s := fmt.Sprintf("%v", obj)
	switch any(out).(type) {
	case bool:
		switch s {
		case "", "<nil>": //Custom to false instead of error
			v = false
		default:
			if v, err = strconv.ParseBool(s); err != nil {
				err = nil
				v = true
			}
		}
	case complex64:
		vc128, err = strconv.ParseComplex(s, 64)
		v = complex64(vc128)
	case complex128:
		v, err = strconv.ParseComplex(s, 128)
	case int:
		vi64, err = strconv.ParseInt(s, 10, strconv.IntSize)
		v = int(vi64)
	case int8:
		vi64, err = strconv.ParseInt(s, 10, 8)
		v = int8(vi64)
	case int16:
		vi64, err = strconv.ParseInt(s, 10, 16)
		v = int16(vi64)
	case int32: //=rune
		vi64, err = strconv.ParseInt(s, 10, 32)
		v = int32(vi64)
	case int64:
		v, err = strconv.ParseInt(s, 10, 64)
	case uint:
		vui64, err = strconv.ParseUint(s, 10, strconv.IntSize)
		v = uint(vui64)
	case uint8: //=byte
		vui64, err = strconv.ParseUint(s, 10, 8)
		v = uint8(vui64)
	case uint16:
		vui64, err = strconv.ParseUint(s, 10, 16)
		v = uint16(vui64)
	case uint32:
		vui64, err = strconv.ParseUint(s, 10, 32)
		v = uint32(vui64)
	case uint64:
		v, err = strconv.ParseUint(s, 10, 64)
	case float32:
		vf64, err = strconv.ParseFloat(s, 32)
		v = float32(vf64)
	case float64:
		v, err = strconv.ParseFloat(s, 64)
	case string:
		v = s
	default:
		err = fmt.Errorf("fail cast to result type %v, from %v: %v", reflect.TypeOf(out), reflect.TypeOf(obj), obj)
	}
	if err == nil {
		out = v.(T)
	}
	return
}

func ToForce[T any](obj interface{}) (out T) {
	out, _ = To[T](obj)
	return
}
