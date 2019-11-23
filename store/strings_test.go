package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")

	s1 := "hello world"
	Set("s1", s1)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		exists  bool
		wantErr bool
	}{
		{"1", args{"s1"}, s1, true, false},
		{"2", args{"noexist"}, "", false, false},
		{"3", args{"hash"}, "", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.exists {
				t.Errorf("Get() exists = %v, exists %v", got1, tt.exists)
			}
		})
	}

	SetEX("s2", s1, 1)
	time.AfterFunc(time.Second, func() {
		s, exists, err := Get("s2")
		assert.Nil(t, err)
		assert.False(t, exists)
		assert.Equal(t, "", s)
	})
}

func TestSet(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")
	s1 := "hello"
	Set("s1", s1)

	s, b, e := Get("s1")
	assert.Nil(t, e)
	assert.True(t, b)
	assert.Equal(t, s1, s)

	Set("hash", s1)
	s, b, e = Get("hash")
	assert.Nil(t, e)
	assert.True(t, b)
	assert.Equal(t, s1, s)
}

func TestGetSet(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")

	s1 := "hello"
	s2 := "hello world"
	Set("s1", s1)

	Set("s5", s1)
	s, b, e := GetSet("s5", s2)
	assert.Nil(t, e)
	assert.True(t, b)
	assert.Equal(t, s1, s)
	s, b, e = Get("s5")
	assert.Nil(t, e)
	assert.True(t, b)
	assert.Equal(t, s2, s)

	type args struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		exists  bool
		wantErr bool
	}{
		{"1", args{"s1", s2}, s1, true, false},
		{"2", args{"noexists", s2}, "", false, false},
		{"3", args{"hash", s2}, "", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetSet(tt.args.key, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSet() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.exists {
				t.Errorf("GetSet() exists = %v, want %v", got1, tt.exists)
			}
		})
	}
}

func TestSetEX(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")
	Set("s1", "hello")

	Set("s2", "hello")
	SetEX("s2", "haha", 10000)
	s, b, _ := Get("s2")
	assert.True(t, b)
	assert.Equal(t, "haha", s)

	type args struct {
		key string
		val string
		ttl int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{"hash", "hello", 10000}, false},
		{"2", args{"noexists", "hello", 10000}, false},
		{"3", args{"s1", "hello", 10000}, false},
		{"3", args{"s1", "hello", 0}, true},
		{"3", args{"s1", "hello", -100}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetEX(tt.args.key, tt.args.val, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("SetEX() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetNX(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")
	Set("s1", "hello")

	type args struct {
		key string
		val string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"hash", "haha"}, false},
		{"2", args{"s1", "haha"}, false},
		{"3", args{"noexists", "haha"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetNX(tt.args.key, tt.args.val); got != tt.want {
				t.Errorf("SetNX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrLen(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")
	Set("s1", "hello")
	Set("s2", "你好")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"s1"}, 5, false},
		{"2", args{"s2"}, 6, false},
		{"3", args{"noexists"}, 0, false},
		{"4", args{"hash"}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrLen(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrLen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StrLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncr(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")
	Set("s1", "hello")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"hash"}, -1, true},
		{"2", args{"s1"}, -1, true},
		{"3", args{"s2"}, 1, false},
		{"4", args{"s2"}, 2, false},
		{"4", args{"s2"}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Incr(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Incr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Incr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncrBy(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")
	Set("s1", "hello")

	//TODO:test 01234
	//Set("s100", "01234")

	type args struct {
		key  string
		intV int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"hash", 1}, -1, true},
		{"2", args{"s1", 1}, -1, true},
		{"3", args{"s2", 1}, 1, false},
		{"4", args{"s2", 1}, 2, false},
		{"5", args{"s2", 1000}, 1002, false},
		{"6", args{"s2", -1002}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IncrBy(tt.args.key, tt.args.intV)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IncrBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppend(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")

	Set("s1", "01234")

	type args struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"noexist", "012345"}, 6, false},
		{"2", args{"s1", "01234"}, 10, false},
		{"3", args{"s1", "01234"}, 15, false},
		{"4", args{"hash", "01234"}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Append(tt.args.key, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRange(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")

	str := "0123456789"
	strLen := len(str)
	strNew := "abcdefghijklmnopqrstuvwxyz"
	newlen := len(strNew)
	Set("range", str)
	Set("range1", str)

	max := MaxStringLength
	n, e := SetRange("range1", strNew, max)
	assert.Error(t, e)
	assert.Equal(t, -1, n)

	n, e = SetRange("range1", strNew, 15)
	assert.Nil(t, e)
	assert.Equal(t, newlen+15, n)

	type args struct {
		key    string
		val    string
		offset int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"range", str, 0}, strLen, false},
		{"2", args{"noexist", strNew, 0}, newlen, false},
		{"3", args{"range", "9876543210", 5}, 5 + 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetRange(tt.args.key, tt.args.val, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetRange() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestGetRange(t *testing.T) {
	values = make(map[string]expired)
	Hset("hash", "f1", "1")
	str := "0123456789"
	Set("range", str)
	str1 := "你好g啊" //len(str1)==10
	Set("range1", str1)

	type args struct {
		key   string
		start int
		end   int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{"range", 0, 0}, str[0:1], false},
		{"2", args{"range", 0, 1}, str[0:2], false},
		{"3", args{"range", 9, 9}, str[9:], false},
		{"4", args{"range", 9, 10}, str[9:], false},
		{"5", args{"range", 9, 10}, str[9:], false},
		{"6", args{"range", -100, 100}, str, false},
		{"7", args{"range", 0, -1}, str, false},
		{"8", args{"range", 10, 9}, "", false},
		{"9", args{"range", 8, 7}, "", false},
		{"10", args{"range", -7, -100}, "", false},
		{"error", args{"hash", -7, -100}, "", true},
		{"noexists", args{"noexists", -7, -100}, "", false},

		{"range1", args{"range1", 0, 0}, str1[0:1], false},
		{"range2", args{"range1", 0, 1}, str1[0:2], false},
		{"range3", args{"range1", 9, 9}, str1[9:], false},
		{"range4", args{"range1", 9, 10}, str1[9:], false},
		{"range5", args{"range1", 9, 10}, str1[9:], false},
		{"range6", args{"range1", -100, 100}, str1, false},
		{"range7", args{"range1", 0, -1}, str1, false},
		{"range8", args{"range1", 10, 9}, "", false},
		{"range9", args{"range1", 8, 7}, "", false},
		{"range10", args{"range1", -7, -100}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRange(tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
