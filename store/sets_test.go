package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSadd(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")

	type args struct {
		key string
		els []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"set", []string{"1", "2"}}, 2, false},
		{"2", args{"set1", []string{}}, 0, false},
		{"3", args{"s1", []string{"12"}}, -1, true},
		{"4", args{"set", []string{"1", "2", "3"}}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sadd(tt.args.key, tt.args.els)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sadd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sadd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScard(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3"})

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"set"}, 3, false},
		{"2", args{"s1"}, -1, true},
		{"3", args{"noexists"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Scard(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdiff(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3"})
	Sadd("set1", []string{"2", "3"})
	Sadd("set3", []string{"4", "5", "6"})

	type args struct {
		key  string
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"1", args{"s1", []string{}}, nil, true},
		{"2", args{"set", []string{}}, []string{"1", "2", "3"}, false},
		{"3", args{"set", []string{"set1"}}, []string{"1"}, false},
		{"4", args{"set1", []string{"set3"}}, []string{"2", "3"}, false},
		{"5", args{"set1", []string{"set"}}, []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sdiff(tt.args.key, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sdiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestSdiffStore(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3"})
	Sadd("set1", []string{"2", "3"})
	Sadd("set3", []string{"4", "5", "6"})

	type args struct {
		dest string
		key  string
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"d1", "s1", []string{}}, -1, true},
		{"2", args{"d2", "set", []string{}}, 3, false},
		{"3", args{"d3", "set", []string{"set1"}}, 1, false},
		{"4", args{"d4", "set1", []string{"set3"}}, 2, false},
		{"5", args{"d5", "set1", []string{"set"}}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SdiffStore(tt.args.dest, tt.args.key, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SdiffStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SdiffStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSinter(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3"})
	Sadd("set1", []string{"2", "3"})
	Sadd("set3", []string{"4", "5", "6"})

	type args struct {
		key  string
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"1", args{"s1", []string{"set1"}}, nil, true},
		{"2", args{"set", []string{"set1"}}, []string{"2", "3"}, false},
		{"3", args{"set", []string{"set3"}}, []string{}, false},
		{"3", args{"set", []string{"set"}}, []string{"1", "2", "3"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sinter(tt.args.key, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sinter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestSinterStore(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3"})
	Sadd("set1", []string{"2", "3"})
	Sadd("set3", []string{"4", "5", "6"})

	type args struct {
		dest string
		key  string
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"d1", "set", []string{"s1", "set"}}, -1, true},
		{"2", args{"d1", "s1", []string{"set"}}, -1, true},
		{"3", args{"d1", "set", []string{"set1"}}, 2, false},
		{"4", args{"d1", "noexists", []string{"set1"}}, 0, false},
		{"5", args{"s1", "set", []string{"set1"}}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SinterStore(tt.args.dest, tt.args.key, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SinterStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SinterStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSmembers(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2"})

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"1", args{"s1"}, nil, true},
		{"2", args{"set"}, []string{"1", "2"}, false},
		{"3", args{"noexist"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Smembers(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Smembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestSismember(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2"})

	type args struct {
		key    string
		member string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"1", args{"s1", "aaa"}, false, true},
		{"2", args{"noexists", "aaa"}, false, false},
		{"3", args{"set", "1"}, true, false},
		{"4", args{"set", "2"}, true, false},
		{"5", args{"set", "5"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sismember(tt.args.key, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sismember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sismember() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpop(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3", "4", "5"})

	type args struct {
		key   string
		count int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"1", args{"s1", 1}, nil, true},
		{"2", args{"noexists", 2}, nil, false},
		{"3", args{"set", 2}, []string{"1", "2"}, false},
		{"4", args{"set", -1}, nil, false},
		{"5", args{"set", 0}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Spop(tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Spop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, len(got), len(tt.want))
		})
	}
}

func TestSrem(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3", "4", "5"})

	type args struct {
		key     string
		members []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"s1", []string{"1212"}}, -1, true},
		{"2", args{"noexists", []string{"1212"}}, 0, false},
		{"3", args{"set", []string{"1212"}}, 0, false},
		{"4", args{"set", []string{"3"}}, 1, false},
		{"5", args{"set", []string{"4", "5"}}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Srem(tt.args.key, tt.args.members)
			if (err != nil) != tt.wantErr {
				t.Errorf("Srem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Srem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSunion(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3"})
	Sadd("set1", []string{"2", "3"})
	Sadd("set3", []string{"4", "5", "6"})

	type args struct {
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"1", args{[]string{"set"}}, []string{"1", "2", "3"}, false},
		{"2", args{[]string{"s1"}}, nil, true},
		{"3", args{[]string{"noexists", "set"}}, []string{"1", "2", "3"}, false},
		{"4", args{[]string{"noexists", "set", "set1"}}, []string{"1", "2", "3"}, false},
		{"4", args{[]string{"noexists", "set", "set1", "set3"}}, []string{"1", "2", "3", "4", "5", "6"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sunion(tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sunion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestSunionStore(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3"})
	Sadd("set1", []string{"2", "3"})
	Sadd("set3", []string{"4", "5", "6"})

	type args struct {
		dest string
		key  string
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"d1", "set", []string{"set"}}, 3, false},
		{"2", args{"d2", "set", []string{"s1"}}, -1, true},
		{"3", args{"d3", "set", []string{"noexists", "set"}}, 3, false},
		{"4", args{"d4", "set", []string{"noexists", "set", "set1"}}, 3, false},
		{"4", args{"d5", "set", []string{"noexists", "set", "set1", "set3"}}, 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SunionStore(tt.args.dest, tt.args.key, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SunionStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SunionStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSmove(t *testing.T) {
	values = make(map[string]expired)
	Set("s1", "hello")
	Sadd("set", []string{"1", "2", "3"})
	Sadd("set1", []string{"2", "3"})
	Sadd("set3", []string{"4", "5", "6"})

	type args struct {
		source string
		dest   string
		member string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"s1", "set", "1"}, -1, true},
		{"2", args{"set", "set", "1"}, 1, false},
		{"3", args{"set1", "set", "2"}, 1, false},
		{"4", args{"set1", "set3", "3"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Smove(tt.args.source, tt.args.dest, tt.args.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("Smove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Smove() = %v, want %v", got, tt.want)
			}
		})
	}
}
