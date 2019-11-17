package store

import (
	"fmt"
	"strconv"
	"time"
)

const (
	//Redis Strings are limited to 512 megabytes
	MaxStringLength = 2 ^ 29 - 1
)

type stringVal struct {
	val      string
	expireAt int64
}

func (s *stringVal) isAlive() bool {
	if s.expireAt == -1 {
		return true
	}
	ttl := s.expireAt - time.Now().Unix()
	return ttl >= 0
}

func (s *stringVal) ttl() int {
	//returns -1 if the key exists but has no associated expire
	if s.expireAt == -1 {
		return -1
	}

	ttl := s.expireAt - time.Now().Unix()
	if ttl > 0 {
		return int(ttl)
	}
	//returns -2 if the key does not exist.
	return -2
}

func (s *stringVal) setExpireAt(at int64) {
	s.expireAt = at
}

func (s *stringVal) dataType() string {
	return "string"
}

func Get(key string) (string, bool, error) {
	if v, ok := values[key]; ok {
		if sv, ok := v.(*stringVal); ok {
			if !sv.isAlive() {
				delete(values, key)
				return "", false, nil
			}
			return sv.val, true, nil
		} else {
			return "", false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	} else {
		return "", false, nil
	}
}

func Set(key, val string) {
	values[key] = &stringVal{val: val, expireAt: -1}
}

func GetSet(key, val string) (string, bool, error) {
	v, ok, err := Get(key)
	if err != nil {
		return "", false, err
	}
	Set(key, val)
	return v, ok, nil
}

func SetEX(key, val string, ttl int) error {
	if ttl <= 0 {
		return fmt.Errorf("ERR invalid expire time in setex")
	}
	values[key] = &stringVal{val: val, expireAt: time.Now().Unix() + int64(ttl)}
	return nil
}

func SetNX(key, val string) bool {
	if _, ok := values[key]; ok {
		return false
	} else {
		values[key] = &stringVal{val: val, expireAt: -1}
		return true
	}
}

func StrLen(key string) (int, error) {
	v, ok, err := Get(key)
	if err != nil {
		return -1, err
	}
	if !ok {
		return 0, nil
	}
	return len(v), nil
}

func Incr(key string) (int, error) {
	return IncrBy(key, 1)
}

func IncrBy(key string, intV int) (int, error) {
	v, ok := values[key]
	if !ok {
		i := intV
		v = &stringVal{val: strconv.Itoa(i), expireAt: -1}
		values[key] = v
		return i, nil
	}

	if sv, ok := v.(*stringVal); ok {
		if !sv.isAlive() {
			i := intV
			v = &stringVal{val: strconv.Itoa(i), expireAt: -1}
			values[key] = v
			return i, nil
		}

		i, e := strconv.Atoi(sv.val)
		if e != nil {
			return -1, fmt.Errorf("ERR value is not an integer or out of range")
		}

		i = i + intV
		sv.val = strconv.Itoa(i)
		return i, nil
	} else {
		return -1, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
}

//If key already exists and is a string, this command appends the value at the end of the string.
//If key does not exist it is created and set as an empty string,
//so APPEND will be similar to SET in this special case.
func Append(key, val string) (int, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		Set(key, val)
		return len(val), nil
	}
	s, ok := v.(*stringVal)
	if !ok {
		return -1, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	s.val = s.val + val
	return len(s.val), nil
}

//offset means byte index, not rune index
func SetRange(key, val string, offset int) (int, error) {
	if offset+len(val) > MaxStringLength {
		return -1, fmt.Errorf("ERR string exceeds maximum allowed size (512MB)")
	}
	v, ok := values[key]
	if !ok || !v.isAlive() {
		rs := make([]byte, offset+len(val))
		bs := []byte(val)
		for i := 0; i < len(val); i++ {
			rs[i+offset] = bs[i]
		}
		Set(key, string(rs))
		return len(rs), nil
	}

	s, ok := v.(*stringVal)
	if !ok {
		return -1, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	var rs []byte
	bs := []byte(s.val)
	valbs := []byte(val)

	if offset+len(valbs) > len(bs) {
		rs = make([]byte, offset+len(valbs))
		copy(rs, bs)
	} else {
		rs = bs
	}

	for i := 0; i < len(val); i++ {
		rs[offset+i] = valbs[i]
	}
	s.val = string(rs)

	return len(rs), nil
}

//Returns the substring of the string value stored at key,
// determined by the offsets start and end (both are inclusive).
// Negative offsets can be used in order to provide an offset starting from the end of the string.
// So -1 means the last character, -2 the penultimate and so forth.
//
//The function handles out of range requests by limiting the resulting range to the actual length of the string.
func GetRange(key string, start, end int) (string, error) {
	s, ok, err := Get(key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	b := []byte(s)
	l := len(b)

	//check negative
	if start < 0 {
		start = l + start
	}
	if end < 0 {
		end = l + end
	}

	if start > end || start > l-1 || end < 0 {
		return "", nil
	}

	if start < 0 {
		start = 0
	}
	if end > l-1 {
		end = l - 1
	}

	rl := end - start + 1
	if rl <= 0 {
		return "", nil
	}

	rt := b[start : end+1]
	return string(rt), nil
}
