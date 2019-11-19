package store

import "time"

type listVal struct {
	val      []string
	expireAt int64
}

func (s *listVal) isAlive() bool {
	if s.expireAt == -1 {
		return true
	}
	ttl := s.expireAt - time.Now().Unix()
	return ttl >= 0
}

func (s *listVal) ttl() int {
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

func (s *listVal) setExpireAt(at int64) {
	s.expireAt = at
}

func (s *listVal) dataType() string {
	return "list"
}

func (s *listVal) lpush(eles []string) int {
	list := s.val
	for _, e := range eles {
		list = append(list, "")
		copy(list[1:], list[1:])
		list[0] = e
	}
	s.val = list
	return len(s.val)
}

func (s *listVal) rpush(eles []string) int {
	list := s.val
	list = append(list, eles...)
	s.val = list
	return len(s.val)
}

func (s *listVal) len() int {
	return len(s.val)
}

func (s *listVal) lpop() string {
	if len(s.val) > 0 {
		r := s.val[0]
		s.val = s.val[1:]
		return r
	}
	return ""
	//TODO: delete list when list is empty
}

func (s *listVal) rpop() string {
	l := len(s.val)
	if l > 0 {
		r := s.val[l-1]
		s.val = s.val[:l-1]
		return r
	}
	return ""
	//TODO: delete list when list is empty
}

func (s *listVal) lindex(i int) string {
	//check i ?
	return s.val[i]
}

func listOf(key string) (*listVal, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return nil, nil
	}
	lv, ok := v.(*listVal)
	if !ok {
		return nil, errorWrongType
	}
	return lv, nil
}

func getOrCreateList(key string) (*listVal, error) {
	lv, err := listOf(key)
	if err != nil {
		return nil, err
	}
	if lv == nil {
		lv = &listVal{val: make([]string, 0), expireAt: -1}
		values[key] = lv
	}
	return lv, nil
}

func Lpush(key string, eles []string) (int, error) {
	lv, err := getOrCreateList(key)
	if err != nil {
		return -1, err
	}
	return lv.lpush(eles), nil
}

func Rpush(key string, eles []string) (int, error) {
	lv, err := getOrCreateList(key)
	if err != nil {
		return -1, err
	}
	return lv.rpush(eles), nil
}

func Llen(key string) (int, error) {
	lv, err := listOf(key)
	if err != nil {
		return -1, err
	}
	if lv == nil {
		return 0, nil
	}
	return lv.len(), nil
}

func Lpop(key string) (string, bool, error) {
	lv, err := listOf(key)
	if err != nil {
		return "", false, err
	}
	if lv == nil {
		return "", false, nil
	}
	return lv.lpop(), true, nil
}

func Rpop(key string) (string, bool, error) {
	lv, err := listOf(key)
	if err != nil {
		return "", false, err
	}
	if lv == nil {
		return "", false, nil
	}
	return lv.rpop(), true, nil
}

func Lindex(key string, idx int) (string, bool, error) {
	lv, err := listOf(key)
	if err != nil {
		return "", false, err
	}
	if lv == nil {
		return "", false, nil
	}

	l := lv.len()
	if idx < 0 {
		idx = l + idx
	}

	if idx < 0 || idx > l-1 {
		return "", false, nil
	}

	return lv.lindex(idx), true, nil
}
