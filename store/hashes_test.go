package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestHset(t *testing.T) {
	values = make(map[string]expired)
	Set("str1", "str1")

	type args struct {
		key   string
		field string
		val   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"hset1", args{"key1", "f1", "v1"}, true, false},
		{"hset2", args{"key1", "f1", "v1"}, false, false},
		{"hset3", args{"key2", "f1", "v2"}, true, false},
		{"hset4", args{"key2", "f2", "v2"}, true, false},
		{"str1", args{"str1", "f2", "v2"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hset(tt.args.key, tt.args.field, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hset() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestHget(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "v2")
	Set("str1", "string1")

	type args struct {
		key   string
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{"hget1", args{"hash1", "f1"}, "v1", true, false},
		{"hget2", args{"hash1", "f2"}, "v2", true, false},
		{"hget3", args{"hash9", "f2"}, "", false, false},
		{"hget4", args{"str1", "f2"}, "", false, true},
		{"hget5", args{"hash1", "f3"}, "", false, false},
		{"hget6", args{"noexist", "f1"}, "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Hget(tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hget() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Hget() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHgetall(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "v2")
	Hset("hash2", "f2", "v2")
	Set("str1", "string1")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"hgetall1", args{"hash1"}, map[string]string{"f1": "v1", "f2": "v2"}, false},
		{"hgetall2", args{"hash2"}, map[string]string{"f2": "v2"}, false},
		{"hgetall3", args{"str1"}, nil, true},
		{"hgetall4", args{"noexist"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hgetall(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hgetall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hgetall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHkeys(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "v2")
	Hset("hash2", "f2", "v2")
	Set("str1", "string1")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"hkeys1", args{"hash1"}, []string{"f1", "f2"}, false},
		{"hkeys2", args{"hash2"}, []string{"f2"}, false},
		{"hkeys3", args{"str1"}, nil, true},
		{"hkeys4", args{"noexist"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hkeys(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hkeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestHlen(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "v2")
	Hset("hash2", "f2", "v2")
	Set("str1", "string1")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"hlen1", args{"hash1"}, 2, false},
		{"hlen2", args{"hash2"}, 1, false},
		{"hlen3", args{"str1"}, -1, true},
		{"hlen4", args{"noexist"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hlen(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hlen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hlen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexists(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "v2")
	Hset("hash2", "f2", "v2")
	Set("str1", "string1")

	type args struct {
		key   string
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"hexist1", args{"hash1", "f1"}, 1, false},
		{"hexist2", args{"hash1", "f2"}, 1, false},
		{"hexist3", args{"hash2", "f2"}, 1, false},
		{"hexist6", args{"hash2", "f1"}, 0, false},
		{"hexist4", args{"str1", "f1"}, -1, true},
		{"hexist5", args{"noexists", "f1"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hexists(tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hexists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hexists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHdel(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "v2")

	Hset("hash2", "f2", "v2")
	Hset("hash2", "f1", "v2")

	Hset("hash3", "f1", "v1")
	Set("str1", "string1")

	type args struct {
		key    string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"hdel1", args{"hash1", []string{"f1"}}, 1, false},
		{"hdel2", args{"hash1", []string{"f2"}}, 1, false},
		{"hdel3", args{"hash2", []string{"f2", "f1"}}, 2, false},
		{"hdel4", args{"hash3", []string{"f2", "f1", "f5"}}, 1, false},
		{"hdel5", args{"str1", []string{"f2", "f1", "f5"}}, -1, true},
		{"hdel6", args{"noexists", []string{"f2", "f1", "f5"}}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hdel(tt.args.key, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hdel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hdel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHsetNX(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "v2")
	Set("str1", "string1")

	type args struct {
		key   string
		field string
		val   string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"hn1", args{"hash1", "f1", "v2"}, 0, false},
		{"hn2", args{"hash1", "f2", "v2"}, 0, false},
		{"hn3", args{"hash1", "f3", "v3"}, 1, false},
		{"hn4", args{"str1", "f3", "v3"}, -1, true},
		{"hn5", args{"noexists", "f3", "v3"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HsetNX(tt.args.key, tt.args.field, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("HsetNX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HsetNX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHstrLen(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "123456789")
	Hset("hash1", "f3", "你好")
	Set("str1", "string1")

	type args struct {
		key   string
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"hl1", args{"hash1", "f1"}, 2, false},
		{"hl2", args{"hash1", "f2"}, 9, false},
		{"hl3", args{"hash1", "f3"}, 6, false},
		{"hl4", args{"hash1", "f4"}, 0, false},
		{"hl5", args{"str1", "f4"}, -1, true},
		{"hl6", args{"noexists", "f4"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HstrLen(tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("HstrLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HstrLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHvals(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash1", "f1", "v1")
	Hset("hash1", "f2", "v2")

	Hset("hash2", "f2", "v2")

	Hset("hash3", "f1", "你好")
	Set("str1", "string1")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"hv1", args{"hash1"}, []string{"v1", "v2"}, false},
		{"hv2", args{"hash2"}, []string{"v2"}, false},
		{"hv3", args{"hash3"}, []string{"你好"}, false},
		{"hv4", args{"noexists"}, nil, false},
		{"hv5", args{"str1"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hvals(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hvals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestHincrBy(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	type args struct {
		key   string
		field string
		delta string
	}
	maxIntStr := strconv.Itoa(math.MaxInt64)
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"hinc1", args{"hash1", "f1", "1"}, 1, false},
		{"hinc2", args{"hash1", "f1", "1"}, 2, false},
		{"hinc3", args{"hash1", "f2", "-1000"}, -1000, false},
		{"hinc4", args{"hash1", "f3", maxIntStr}, math.MaxInt64, false},
		{"hinc5", args{"hash1", "f3", maxIntStr}, -1, true},
		{"hinc6", args{"s1", "f3", maxIntStr}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HincrBy(tt.args.key, tt.args.field, tt.args.delta)
			if (err != nil) != tt.wantErr {
				t.Errorf("HincrBy() error = %v, wantErr %v, name:%s", err, tt.wantErr, tt.name)
				return
			}
			if got != tt.want {
				t.Errorf("HincrBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHincrByFloat(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	maxFloat := fmt.Sprintf("%f", math.MaxFloat64)
	type args struct {
		key   string
		field string
		delta string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"hinc1", args{"hash1", "f1", "1"}, "1", false},
		{"hinc2", args{"hash1", "f1", "1.054"}, "2.054", false},
		{"hinc3", args{"hash1", "f2", "-1000.9801"}, "-1000.9801", false},
		{"hinc4", args{"s1", "f3", "1212"}, "", true},
		{"hinc5", args{"hash1", "f3", maxFloat}, maxFloat, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HincrByFloat(tt.args.key, tt.args.field, tt.args.delta)
			if tt.wantErr {
				if err == nil {
					t.Errorf("HincrByFloat() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				real, _ := strconv.ParseFloat(got, 64)
				want, _ := strconv.ParseFloat(tt.want, 64)
				if real != want {
					t.Errorf("HincrByFloat() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
