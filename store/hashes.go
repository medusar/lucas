package store

import (
	"fmt"
	"github.com/medusar/lucas/util"
	"math"
	"strconv"
	"time"
)

type hashVal struct {
	val      map[string]string
	expireAt int64
}

func (s *hashVal) isAlive() bool {
	if s.expireAt == -1 {
		return true
	}
	ttl := s.expireAt - time.Now().Unix()
	return ttl >= 0
}

func (s *hashVal) ttl() int {
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

func (s *hashVal) setExpireAt(at int64) {
	s.expireAt = at
}

func (s *hashVal) dataType() string {
	return "hash"
}

//Hset set a field to a hash, return true if the field doesn't exist before
func Hset(key, field, val string) (bool, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		m := &hashVal{val: make(map[string]string), expireAt: -1}
		m.val[field] = val
		values[key] = m
		return true, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return false, errorWrongType
	}
	_, exists := h.val[field]
	h.val[field] = val
	return !exists, nil
}

func Hget(key, field string) (string, bool, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return "", false, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return "", false, errorWrongType
	}
	val, exist := h.val[field]
	return val, exist, nil
}

func Hgetall(key string) (map[string]string, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return nil, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return nil, errorWrongType
	}
	return h.val, nil
}

func Hkeys(key string) ([]string, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return nil, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return nil, errorWrongType
	}
	keys := make([]string, len(h.val))
	i := 0
	for k, _ := range h.val {
		keys[i] = k
		i++
	}
	return keys, nil
}

//Hlen return number of fields in the hash, or 0 when key does not exist.
func Hlen(key string) (int, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return 0, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return -1, errorWrongType
	}
	return len(h.val), nil
}

//Returns if field is an existing field in the hash stored at key.
//1 if the hash contains field.
//0 if the hash does not contain field, or key does not exist.
func Hexists(key, field string) (int, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return 0, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return -1, errorWrongType
	}
	_, exists := h.val[field]
	if exists {
		return 1, nil
	}
	return 0, nil
}

//Hdel return the number of fields that were removed from the hash, not including specified but non existing fields.
// If key does not exist, it is treated as an empty hash and this command returns 0.
func Hdel(key string, fields []string) (int, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return 0, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return -1, errorWrongType
	}

	t := 0
	for _, f := range fields {
		if _, exists := h.val[f]; exists {
			delete(h.val, f)
			t++
		}
	}
	return t, nil
}

//HsetNX sets field in the hash stored at key to value, only if field does not yet exist.
// If key does not exist, a new key holding a hash is created.
// If field already exists, this operation has no effect.
func HsetNX(key, field, val string) (int, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		m := &hashVal{val: make(map[string]string), expireAt: -1}
		m.val[field] = val
		values[key] = m
		return 1, nil
	}

	h, ok := v.(*hashVal)
	if !ok {
		return -1, errorWrongType
	}

	_, exists := h.val[field]
	if exists {
		return 0, nil
	}
	h.val[field] = val
	return 1, nil
}

//HstrLen return the string length of the value associated with field,
// or zero when field is not present in the hash or key does not exist at all.
func HstrLen(key, field string) (int, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return 0, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return -1, errorWrongType
	}
	val, exists := h.val[field]
	if !exists {
		return 0, nil
	}
	return len(val), nil
}

func Hvals(key string) ([]string, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return nil, nil
	}
	h, ok := v.(*hashVal)
	if !ok {
		return nil, errorWrongType
	}
	ret := make([]string, 0)
	for _, val := range h.val {
		ret = append(ret, val)
	}
	return ret, nil
}

//The range of values supported by HINCRBY is limited to 64 bit signed integers.
// Return the value at field after the increment operation.
func HincrBy(key, field, delta string) (int, error) {
	incr, err := strconv.Atoi(delta)
	if err != nil {
		return -1, errorInvalidInt
	}

	v, ok := values[key]
	if !ok || !v.isAlive() {
		m := &hashVal{val: make(map[string]string), expireAt: -1}
		m.val[field] = delta
		values[key] = m
		return incr, nil
	}

	h, ok := v.(*hashVal)
	if !ok {
		return -1, errorWrongType
	}
	val, exists := h.val[field]
	if !exists {
		h.val[field] = delta
		return incr, nil
	}
	old, err := strconv.Atoi(val)
	if err != nil {
		return -1, fmt.Errorf("ERR hash value is not an integer")
	}

	newVal, err := util.Add64(old, incr)
	if err != nil {
		return -1, fmt.Errorf("ERR increment or decrement would overflow")
	}

	h.val[field] = strconv.Itoa(newVal)
	return newVal, nil
}

func HincrByFloat(key, field, delta string) (string, error) {
	incr, err := strconv.ParseFloat(delta, 64)
	if err != nil {
		return "", errorInvalidFloat
	}

	v, ok := values[key]
	if !ok || !v.isAlive() {
		m := &hashVal{val: make(map[string]string), expireAt: -1}
		m.val[field] = delta
		values[key] = m
		return delta, nil
	}

	h, ok := v.(*hashVal)
	if !ok {
		return "", errorWrongType
	}

	val, exists := h.val[field]
	if !exists {
		h.val[field] = delta
		return delta, nil
	}

	old, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "", fmt.Errorf("ERR hash value is not a float")
	}

	newVal := old + incr
	//TODO: maybe invalid check?
	if newVal > math.MaxFloat64 || newVal < math.SmallestNonzeroFloat64 {
		return "", fmt.Errorf("ERR increment or decrement would overflow")
	}

	fieldVal := fmt.Sprintf("%f", newVal)
	h.val[field] = fieldVal
	return fieldVal, nil
}
