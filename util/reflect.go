package util

import (
	"reflect"
)

func ToMap(m interface{}, filter []string) map[string]interface{} {
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)

	var data = map[string]interface{}{}
	for k := 0; k < t.NumField(); k++ {

		var key = t.Field(k).Name
		var value = v.Field(k).Interface()

		var skip = false
		for _, f := range filter {
			if key == f {
				skip = true
				continue
			}
		}
		if skip {
			continue
		}

		data[key] = value
	}
	return data
}
