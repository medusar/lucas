package store

import (
	"fmt"
	"github.com/mb0/glob"
	"log"
	"strconv"
	"time"
)

var (
	values map[string]expired = make(map[string]expired)
)

type expired interface {
	isAlive() bool
	ttl() int
	setExpireAt(at int64)
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

func (s *stringVal) setExpireAt(at int64) {
	s.expireAt = at
}

func Ttl(key string) int {
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

func Expire(key string, ttl int) bool {
	return ExpireAt(key, time.Now().Unix()+int64(ttl))
}

func ExpireAt(key string, timestamp int64) bool {
	v, ok := values[key]
	if !ok {
		return false
	}
	if !v.isAlive() {
		return false
	}

	v.setExpireAt(timestamp)
	return true
}

func Keys(pattern string) []string {
	//TODO: check pattern
	keys := make([]string, 0)
	for key, v := range values {
		if v.isAlive() && patternMatch(key, pattern) {
			keys = append(keys, key)
		}
	}
	return keys
}

func patternMatch(key, pattern string) bool {
	ok, err := glob.Match(pattern, key)
	if err != nil {
		log.Println("Failed to do glob match", err)
		return false
	}
	return ok
}

func Exists(key string) bool {
	v, ok := values[key]
	return ok && v.isAlive()
}

func Del(key string) bool {
	v, ok := values[key]

	if ok {
		alive := v.isAlive()
		delete(values, key)

		if alive {
			return true
		}
		return false
	}

	return false
}
