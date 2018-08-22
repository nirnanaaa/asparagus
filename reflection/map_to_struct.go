package reflection

import (
	"fmt"
	"reflect"
	"strconv"
)

// MapToStruct turns a map into a mapped struct
func MapToStruct(t interface{}, values map[string]interface{}) error {
	ps := reflect.ValueOf(t)
	// struct
	s := ps.Elem()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("can't assign values to struct")
	}
	for key, value := range values {
		f := s.FieldByName(key)
		if f.IsValid() {
			if f.CanSet() {
				switch v := f.Interface().(type) {
				case int:
					val, err := strconv.Atoi(value.(string))
					if err != nil {
						return err
					}
					x := int64(val)
					if !f.OverflowInt(x) {
						f.SetInt(x)
					}
				case string:
					f.SetString(value.(string))
				case bool:
					f.SetBool(value.(string) == "true")
				case []interface{}:
					input := reflect.ValueOf(value)
					f.Set(input)
				case []string:
					switch value.(type) {
					case []string:
						input := reflect.ValueOf(value.([]string))
						f.Set(input)
					case []interface{}:
						target := []string{}
						for _, v := range value.([]interface{}) {
							target = append(target, v.(string))
						}
						input := reflect.ValueOf(target)
						f.Set(input)
					}
				default:
					return fmt.Errorf("i don't know how to parse type %T", v)
				}
			}
		}
	}
	return nil
}
