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

func InterMapKeys(m1, m2 map[string]*struct{}) map[string]*struct{} {
	r := make(map[string]*struct{})
	for k1, v := range m1 {
		if _, ok := m2[k1]; ok {
			r[k1] = v
		}
	}
	return r
}
