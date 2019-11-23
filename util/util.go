package util

import (
	"errors"
	"math"
)

var ErrOverFlow = errors.New("integer overflow")

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

func DeleteStringArray(i int, array []string) []string {
	if i < len(array)-1 {
		copy(array[i:], array[i+1:])
	}
	array[len(array)-1] = ""
	array = array[:len(array)-1]
	return array
}
