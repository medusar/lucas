package util

import "reflect"

//IsPointer checks if val is a pointer type
func IsPointer(val interface{}) bool {
	return reflect.ValueOf(val).Kind() == reflect.Ptr
}

//DiffArray returns the elements in `a` that aren't in `b`.
func DiffArray(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
