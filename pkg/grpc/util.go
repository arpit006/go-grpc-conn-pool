package grpc

import (
	"reflect"
	"time"
)

func IsZero[V string | int | time.Duration](v V) bool {
	value := reflect.ValueOf(v)
	kind := value.Kind()

	switch kind {
	case reflect.Int:
		return value.Int() == 0
	case reflect.String:
		return value.String() == ""
	case reflect.Struct:
		if value.Type() == reflect.TypeOf(time.Duration(0)) {
			return value.Interface().(time.Duration) == 0
		}
	}
	panic("cannot evaluate zero value of undefined type")
}

func WrapWithDefault[V string | int | time.Duration](v V, defaultVal V) V {
	if IsZero[V](v) {
		return defaultVal
	}
	return v
}
