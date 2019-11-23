package store

import (
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"strconv"
	"testing"
)

func TestLpush(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	//TODO: rpush empty list

	for i := 6; i < 10; i++ {
		Lpush("list0", []string{strconv.Itoa(i)})
	}

	list, err := listOf("list0")
	assert.Nil(t, err)
	for i := 0; i < 4; i++ {
		assert.Equal(t, strconv.Itoa(9-i), list.val[i])
	}

	type args struct {
		key  string
		eles []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"list1", []string{"1"}}, 1, false},
		{"2", args{"list1", []string{"2", "3", "4"}}, 4, false},
		{"3", args{"list1", []string{}}, 4, false},
		{"4", args{"list1", nil}, 4, false},
		{"5", args{"s1", []string{"1", "2"}}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Lpush(tt.args.key, tt.args.eles)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lpush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Lpush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRpush(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	for i := 0; i < 4; i++ {
		Rpush("list0", []string{strconv.Itoa(i)})
	}

	list, err := listOf("list0")
	assert.Nil(t, err)
	for i := 0; i < 4; i++ {
		assert.Equal(t, strconv.Itoa(i), list.val[i])
	}

	type args struct {
		key  string
		eles []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"list1", []string{"1"}}, 1, false},
		{"2", args{"list1", []string{"2", "3", "4"}}, 4, false},
		{"3", args{"list1", []string{}}, 4, false},
		{"4", args{"list1", nil}, 4, false},
		{"5", args{"s1", []string{"1", "2"}}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Rpush(tt.args.key, tt.args.eles)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rpush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rpush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLlen(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	for i := 0; i < 10; i++ {
		Rpush("list1", []string{strconv.Itoa(i)})
	}
	val, err := listOf("list1")
	assert.Nil(t, err)
	assert.Equal(t, 10, val.len())

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"list1"}, 10, false},
		{"2", args{"s1"}, -1, true},
		{"2", args{"noexists"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Llen(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Llen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Llen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLpop(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	for i := 6; i < 10; i++ {
		Lpush("list1", []string{strconv.Itoa(i)})
	}

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{"1", args{"list1"}, "9", true, false},
		{"2", args{"list1"}, "8", true, false},
		{"3", args{"list1"}, "7", true, false},
		{"4", args{"list1"}, "6", true, false},
		{"5", args{"list1"}, "", false, false},
		{"6", args{"noexists"}, "", false, false},
		{"7", args{"s1"}, "", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Lpop(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lpop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Lpop() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Lpop() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRpop(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	for i := 0; i < 4; i++ {
		s := strconv.Itoa(i)
		Lpush("list1", []string{s})
		Rpush("list2", []string{s})
	}

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{"1", args{"list1"}, "0", true, false},
		{"2", args{"list1"}, "1", true, false},
		{"3", args{"list1"}, "2", true, false},
		{"4", args{"list1"}, "3", true, false},
		{"5", args{"list1"}, "", false, false},
		{"6", args{"noexists"}, "", false, false},

		{"7", args{"list2"}, "3", true, false},
		{"8", args{"list2"}, "2", true, false},
		{"9", args{"list2"}, "1", true, false},
		{"10", args{"list2"}, "0", true, false},
		{"11", args{"list2"}, "", false, false},
		{"12", args{"s1"}, "", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Rpop(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rpop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rpop() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Rpop() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLindex(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	for i := 0; i < 4; i++ {
		s := strconv.Itoa(i)
		Lpush("list1", []string{s})
	}

	type args struct {
		key string
		idx int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		exist   bool
		wantErr bool
	}{
		{"1", args{"list1", 0}, "3", true, false},
		{"2", args{"list1", 1}, "2", true, false},
		{"3", args{"list1", 2}, "1", true, false},
		{"4", args{"list1", 3}, "0", true, false},
		{"5", args{"list1", 4}, "", false, false},
		{"6", args{"list1", -1}, "0", true, false},
		{"7", args{"list1", -2}, "1", true, false},
		{"8", args{"list1", -3}, "2", true, false},
		{"9", args{"list1", -4}, "3", true, false},
		{"9", args{"list1", -10}, "", false, false},
		{"10", args{"s1", -10}, "", false, true},
		{"11", args{"noexists", 0}, "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Lindex(tt.args.key, tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lindex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Lindex() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.exist {
				t.Errorf("Lindex() got1 = %v, exist %v", got1, tt.exist)
			}
		})
	}
}

func TestLrem(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")

	Lpush("list1", []string{"0", "1", "2", "3"})

	// count > 0
	i, err := Lrem("list1", 1, "2")
	assert.Nil(t, err)
	assert.Equal(t, 1, i)

	i, err = Lrem("list1", 2, "0")
	assert.Nil(t, err)
	assert.Equal(t, 1, i)

	i, err = Lrem("list1", 1, "3")
	assert.Nil(t, err)
	assert.Equal(t, 1, i)

	list, err := listOf("list1")
	assert.Nil(t, err)
	assert.Equal(t, 1, list.len())

	//count <0
	Rpush("list2", []string{"0", "1", "2", "3", "1", "0", "2", "3", "3", "4", "5"})
	i, err = Lrem("list2", -1, "1")
	assert.Nil(t, err)
	assert.Equal(t, 1, i)

	s, exist, err := Lindex("list2", 4)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.Equal(t, "0", s)

	//remain "0", "1", "2", "3", "0", "2", "3", "3", "4", "5"
	i, err = Lrem("list2", -2, "3")
	assert.Nil(t, err)
	assert.Equal(t, 2, i)

	//remain "0", "1", "2", "3", "0", "2", "4", "5"
	s, exist, err = Lindex("list2", -3)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.Equal(t, "2", s)

	s, exist, err = Lindex("list2", -4)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.Equal(t, "0", s)

	i, err = Lrem("list2", -1, "noexist")
	assert.Nil(t, err)
	assert.Equal(t, 0, i)

	//remain "0", "1", "2", "3", "0", "2", "4", "5"
	i, err = Lrem("list2", -100, "2")
	assert.Nil(t, err)
	assert.Equal(t, 2, i)

	//remain "0", "1", "3", "0", "4", "5"
	i, err = Lrem("list2", -100, "0")
	assert.Nil(t, err)
	assert.Equal(t, 2, i)

	//remain "1", "3", "4", "5"
	s, exist, err = Lindex("list2", -4)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.Equal(t, "1", s)

	Rpush("list4", []string{"0", "1", "2", "3", "1", "0", "2", "3", "3", "4", "5"})
	i, err = Lrem("list4", -5, "3")
	assert.Nil(t, err)
	assert.Equal(t, 3, i)

	// count == 0
	Rpush("list3", []string{"0", "1", "2", "3", "1", "0", "2", "3", "3", "4", "5"})
	i, err = Lrem("list3", 0, "3")
	assert.Nil(t, err)
	assert.Equal(t, 3, i)

	list, _ = listOf("list3")
	log.Println(list.val)
	//remain  "0", "1", "2", "1", "0", "2", "4", "5"
	s, exist, err = Lindex("list3", -4)
	assert.Nil(t, err)
	assert.True(t, exist)
	assert.Equal(t, "0", s)

	i, err = Lrem("noexists", 0, "11")
	assert.Nil(t, err)
	assert.Equal(t, 0, i)

	i, err = Lrem("s1", 0, "11")
	assert.Error(t, err)
	assert.Equal(t, -1, i)

}

func TestLset(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")
	Rpush("list", []string{"0", "1", "2", "3"})

	type args struct {
		key     string
		index   int
		element string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{"list", -1, "5"}, false},
		{"2", args{"list", -10, "5"}, true},
		{"3", args{"list", -4, "0/-4"}, false},
		{"4", args{"list", 0, "0/-4"}, false},
		{"4", args{"list", 4, "0/-4"}, true},
		{"4", args{"list", 3, "0/-4"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Lset(tt.args.key, tt.args.index, tt.args.element); (err != nil) != tt.wantErr {
				t.Errorf("Lset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	Rpush("list1", []string{"0", "1", "2", "3"})
	assert.Nil(t, Lset("list1", 0, "-4/0"))
	assert.Nil(t, Lset("list1", 3, "-1/3"))
	assert.Nil(t, Lset("list1", -2, "-2/2"))
	assert.Nil(t, Lset("list1", -3, "-3/1"))
	list, err := listOf("list1")
	assert.Nil(t, err)
	assert.Equal(t, "-4/0", list.val[0])
	assert.Equal(t, "-1/3", list.val[3])
	assert.Equal(t, "-2/2", list.val[2])
	assert.Equal(t, "-3/1", list.val[1])
}

func TestRpushX(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")
	Rpush("list", []string{"0"})

	type args struct {
		key      string
		elements []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"list", []string{"1", "2", "3"}}, 4, false},
		{"2", args{"noexists", []string{"1", "2", "3"}}, 0, false},
		{"3", args{"s1", []string{"1", "2", "3"}}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RpushX(tt.args.key, tt.args.elements)
			if (err != nil) != tt.wantErr {
				t.Errorf("RpushX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RpushX() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLpushX(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")
	Rpush("list", []string{"0"})

	type args struct {
		key      string
		elements []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"list", []string{"1", "2", "3"}}, 4, false},
		{"2", args{"noexists", []string{"1", "2", "3"}}, 0, false},
		{"3", args{"s1", []string{"1", "2", "3"}}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LpushX(tt.args.key, tt.args.elements)
			if (err != nil) != tt.wantErr {
				t.Errorf("LpushX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LpushX() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLrange(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "s1")
	array := make([]string, 10)
	for i := 0; i < 10; i++ {
		array[i] = strconv.Itoa(i)
	}
	Rpush("list", array)

	type args struct {
		key   string
		start int
		end   int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"1", args{"list", 0, 0}, array[0:1], false},
		{"2", args{"list", 0, 1}, array[0:2], false},
		{"3", args{"list", 9, 9}, array[9:], false},
		{"4", args{"list", 9, 10}, array[9:], false},
		{"5", args{"list", 9, 10}, array[9:], false},
		{"6", args{"list", -100, 100}, array, false},
		{"7", args{"list", 0, -1}, array, false},
		{"8", args{"list", 10, 9}, nil, false},
		{"9", args{"list", 8, 7}, nil, false},
		{"10", args{"list", -7, -100}, nil, false},
		{"error", args{"s1", -7, -100}, nil, true},
		{"noexists", args{"noexists", -7, -100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Lrange(tt.args.key, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lrange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lrange() got = %v, want %v", got, tt.want)
			}
		})
	}
}
