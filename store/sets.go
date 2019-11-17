package store

import (
	"github.com/medusar/lucas/util"
	"time"
)

var obj = &struct{}{}

type setVal struct {
	val      map[string]*struct{}
	expireAt int64
}

func (s *setVal) isAlive() bool {
	if s.expireAt == -1 {
		return true
	}
	ttl := s.expireAt - time.Now().Unix()
	return ttl >= 0
}

func (s *setVal) ttl() int {
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

func (s *setVal) setExpireAt(at int64) {
	s.expireAt = at
}

func (s *setVal) dataType() string {
	return "set"
}

func Sadd(key string, els []string) (int, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		m := make(map[string]*struct{})
		for _, el := range els {
			m[el] = obj
		}
		values[key] = &setVal{val: m, expireAt: -1}
		return len(els), nil
	}
	s, ok := v.(*setVal)
	if !ok {
		return -1, errorWrongType
	}
	t := 0
	for _, el := range els {
		_, exists := s.val[el]
		if !exists {
			s.val[el] = obj
			t++
		}
	}
	return t, nil
}

//Scard returns the cardinality (number of elements) of the set, or 0 if key does not exist.
func Scard(key string) (int, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return 0, nil
	}
	s, ok := v.(*setVal)
	if !ok {
		return -1, errorWrongType
	}
	return len(s.val), nil
}

func mapOf(key string) (map[string]*struct{}, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return nil, nil
	}
	s, ok := v.(*setVal)
	if !ok {
		return nil, errorWrongType
	}
	return s.val, nil
}

//Sdiff returns the members of the set resulting from the difference between the first set and all the successive sets.
func Sdiff(key string, keys ...string) ([]string, error) {
	s1, err := Smembers(key)
	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		s, err := Smembers(k)
		if err != nil {
			return nil, err
		}
		s1 = util.DiffArray(s1, s)
		if len(s1) == 0 {
			break
		}
	}
	return s1, nil
}

func Smembers(key string) ([]string, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return nil, nil
	}
	s, ok := v.(*setVal)
	if !ok {
		return nil, errorWrongType
	}
	ret := make([]string, len(s.val))
	i := 0
	for k, _ := range s.val {
		ret[i] = k
		i++
	}
	return ret, nil
}

func Sismember(key, member string) (bool, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return false, nil
	}
	s, ok := v.(*setVal)
	if !ok {
		return false, errorWrongType
	}
	_, is := s.val[member]
	return is, nil
}

func Spop(key string, count int) ([]string, bool, error) {
	m, err := mapOf(key)
	if err != nil {
		return nil, false, err
	}
	if m == nil {
		return nil, false, nil
	}

	r := make([]string, 0)
	i := 0
	for k := range m {
		if i >= count {
			break
		}
		delete(m, k)
		r = append(r, k)
		i++
	}
	return r, true, nil
}
