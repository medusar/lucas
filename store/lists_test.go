package store

import "testing"

func TestLpush(t *testing.T) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
	type args struct {
		key string
		idx int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if got1 != tt.want1 {
				t.Errorf("Lindex() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
