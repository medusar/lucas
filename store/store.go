package store

import (
	"fmt"
	"time"
)

var (
	values map[string]expired = make(map[string]expired)
)

type expired interface {
	isAlive() bool
	ttl() int
}

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

func Ttl(key string) int {
	//TODO: support other types
	v, ok := values[key]
	if !ok {
		//returns -2 if the key does not exist.
		return -2
	}
	return v.ttl()
}

func Get(key string) (string, bool, error) {
	if v, ok := values[key]; ok {
		if sv, ok := v.(*stringVal); ok {
			if !sv.isAlive() {
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
	//TODO:
	return 0, nil
}
