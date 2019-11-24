package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestAdd64(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"test1", args{1, 2}, 3, false},
		{"test1", args{math.MaxInt64, 2}, -1, true},
		{"test1", args{math.MaxInt64, -2}, math.MaxInt64 - 2, false},
		{"test1", args{-2, math.MaxInt64}, math.MaxInt64 - 2, false},
		{"test1", args{2, math.MaxInt64}, -1, true},
		{"test1", args{-2, math.MinInt64}, -1, true},
		{"test1", args{2, math.MinInt64}, math.MinInt64 + 2, false},
		{"test1", args{math.MinInt64, math.MinInt64}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add64(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Add64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffArray(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"d1", args{[]string{"1", "2", "3"}, []string{"1"}}, []string{"2", "3"}},
		{"d2", args{[]string{"1", "2", "3"}, []string{"1", "3", "2"}}, []string{}},
		{"d3", args{[]string{"1", "2", "3"}, nil}, []string{"1", "2", "3"}},
		{"d4", args{nil, nil}, nil},
		{"d5", args{nil, []string{"1", "2", "3"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DiffArray(tt.args.a, tt.args.b)
			assert.ElementsMatch(t, got, tt.want, fmt.Sprintf("DiffArray() = %v, want %v", got, tt.want))
		})
	}
}
