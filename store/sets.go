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

func SdiffStore(dest, key string, keys ...string) (int, error) {
	set, err := Sdiff(key, keys...)
	if err != nil {
		return 0, err
	}
	//remove
	delete(values, dest)
	if len(set) == 0 {
		return 0, nil
	}
	return Sadd(dest, set)
}

// Keys that do not exist are considered to be empty sets.
// With one of the keys being an empty set, the resulting set is also empty
// (since set intersection with an empty set always results in an empty set).
func Sinter(key string, keys ...string) ([]string, error) {
	m1, err := mapOf(key)
	if err != nil {
		return nil, err
	}

	if m1 == nil {
		return nil, nil
	}

	r := make([]string, 0)
	for mk := range m1 {
		in := true
		for _, k := range keys {
			m, err := mapOf(k)
			if err != nil {
				return nil, err
			}
			if m == nil {
				return nil, nil
			}
			if _, ok := m[mk]; !ok {
				in = false
				break
			}
		}
		if in {
			r = append(r, mk)
		}
	}
	return r, nil
}

func SinterStore(dest, key string, keys ...string) (int, error) {
	ins, err := Sinter(key, keys...)
	if err != nil {
		return -1, err
	}
	delete(values, dest)
	if len(ins) == 0 {
		return 0, nil
	}
	return Sadd(dest, ins)
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
	for k := range s.val {
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

func Srem(key string, members []string) (int, error) {
	mp, err := mapOf(key)
	if err != nil {
		return -1, err
	}
	if mp == nil {
		return 0, nil
	}

	t := 0
	for _, m := range members {
		if _, ok := mp[m]; ok {
			delete(mp, m)
			t++
		}
	}
	return t, nil
}

func Sunion(keys []string) ([]string, error) {
	mp := make(map[string]*struct{})
	for _, key := range keys {
		m, err := mapOf(key)
		if err != nil {
			return nil, err
		}
		for k := range m {
			mp[k] = obj
		}
	}

	ret := make([]string, 0)
	for k := range mp {
		ret = append(ret, k)
	}
	return ret, nil
}

func SunionStore(dest, key string, keys ...string) (int, error) {
	set, err := Sunion(append(keys, key))
	if err != nil {
		return -1, err
	}
	delete(values, dest)
	if len(set) == 0 {
		return 0, nil
	}
	return Sadd(dest, set)
}

//1 if the element is moved.
//0 if the element is not a member of source and no operation was performed.
func Smove(source, dest, member string) (int, error) {
	smp, err := mapOf(source)
	if err != nil {
		return -1, err
	}
	//If the source set does not exist or does not contain the specified element,
	// no operation is performed and 0 is returned.
	if smp == nil {
		return 0, nil
	}
	if _, ok := smp[member]; !ok {
		return 0, nil
	}
	delete(smp, member)
	Sadd(dest, []string{member})
	return 1, nil
}
