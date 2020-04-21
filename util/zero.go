package util

import "reflect"

func IsZero(i interface{}) bool {
	return i == nil || reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface())
}

