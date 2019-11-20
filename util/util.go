package util

import (
	"errors"
	"math"
	"reflect"
)

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

var ErrOverFlow = errors.New("integer overflow")

//Add64 do add operation in int64, and return error if the result overflows int64
//https://stackoverflow.com/questions/33641717/detect-signed-int-overflow-in-go
func Add64(a, b int) (int, error) {
	if a > 0 {
		if b > math.MaxInt64-a {
			return -1, ErrOverFlow
		}
	} else {
		if b < math.MinInt64-a {
			return -1, ErrOverFlow
		}
	}
	return a + b, nil
}

func Add64Float(a, b float64) (float64, error) {
	//TODO: check float overflow
	return a + b, nil
}
