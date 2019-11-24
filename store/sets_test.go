package store

import (
	"reflect"
	"testing"
)

func TestSadd(t *testing.T) {
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
		// TODO: Add test cases.
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
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sdiff(tt.args.key, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sdiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sdiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdiffStore(t *testing.T) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sinter(tt.args.key, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sinter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sinter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSinterStore(t *testing.T) {
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
		// TODO: Add test cases.
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
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Smembers(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Smembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Smembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSismember(t *testing.T) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Spop(tt.args.key, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Spop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Spop() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSrem(t *testing.T) {
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
		// TODO: Add test cases.
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
	type args struct {
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sunion(tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sunion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sunion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSunionStore(t *testing.T) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
