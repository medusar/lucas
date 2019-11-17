package main

import "reflect"

//IsPointer checks if val is a pointer type
func IsPointer(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Ptr
}
