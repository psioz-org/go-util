package mapz

import (
	"encoding/json"

	"github.com/zev-zakaryan/go-util/conv"
)

func Join[TKey comparable, TVal comparable](out map[TKey]TVal, excludeEmpty bool, in ...map[TKey]TVal) {
	var empty TVal
	for _, in1 := range in {
		for k, v := range in1 {
			if _, ok := out[k]; !ok && (!excludeEmpty || v != empty) {
				out[k] = v
			}
		}
	}
}

func JoinAny[TKey comparable, TVal any](out map[TKey]TVal, in ...map[TKey]TVal) {
	for _, in1 := range in {
		for k, v := range in1 {
			if _, ok := out[k]; !ok {
				out[k] = v
			}
		}
	}
}

func ToMap(obj interface{}) map[string]interface{} {
	var out map[string]interface{}
	var objJ []byte
	switch v := obj.(type) {
	case string:
		objJ = []byte(v)
	case []byte:
		objJ = v
	default:
		objJ, _ = json.Marshal(obj)
	}
	json.Unmarshal(objJ, &out)
	return out
}

func ToStringMap(obj interface{}) map[string]string {
	var outI map[string]interface{}
	var objJ []byte
	switch v := obj.(type) {
	case string:
		objJ = []byte(v)
	case []byte:
		objJ = v
	default:
		objJ, _ = json.Marshal(obj)
	}
	json.Unmarshal(objJ, &outI)
	out := make(map[string]string)
	for k, v := range outI {
		out[k] = conv.ToForce[string](v)
	}
	return out
}
