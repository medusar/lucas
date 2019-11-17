package store

import "time"

type setVal struct {
	val      map[string]struct{}
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


