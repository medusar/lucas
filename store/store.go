package store

import (
	"fmt"
	"github.com/mb0/glob"
	"log"
	"time"
)

var (
	values            = make(map[string]expired)
	errorWrongType    = fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	errorInvalidInt   = fmt.Errorf("ERR value is not an integer or out of range")
	errorInvalidFloat = fmt.Errorf("ERR value is not a valid float")
)

type expired interface {
	isAlive() bool
	ttl() int
	setExpireAt(at int64)
	dataType() string
}

func Ttl(key string) int {
	v, ok := values[key]
	if !ok {
		//returns -2 if the key does not exist.
		return -2
	}
	return v.ttl()
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

func Type(key string) string {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return "none"
	}
	return v.dataType()
}
