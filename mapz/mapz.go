package mapz

import (
	"encoding/json"

	"github.com/zev-zakaryan/go-util/zz"
)

func Join(out map[string]string, excludeEmpty bool, in ...map[string]string) {
	for _, in1 := range in {
		for k, v := range in1 {
			if _, ok := out[k]; !ok && (!excludeEmpty || v != "") {
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
		out[k] = zz.ToForce[string](v)
	}
	return out
}
