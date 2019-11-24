package store

import (
	"fmt"
	"sort"
	"time"
)

type zsetMember struct {
	Member string
	Score  float64
}

type scoreMember struct {
	next    *scoreMember
	pre     *scoreMember
	score   float64
	members []string //members are ordered lexicographically
}

// sorted by score from low to high
// members with same score are order lexicographically
type scoreMemberMap struct {
	head *scoreMember
	tail *scoreMember
	size int
}

func (sm *scoreMemberMap) String() string {
	if sm.head == nil {
		return "nil"
	}

	logElement := func(m *scoreMember) string {
		if m == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", m.members)
	}

	sb := logElement(sm.head)
	for cur := sm.head.next; cur != nil; cur = cur.next {
		sb = sb + "<->" + logElement(cur)
	}
	return sb
}

// remove the member associated with score
func (sm *scoreMemberMap) remove(score float64, member string) bool {
	cur := sm.tail
	for cur != nil {
		if cur.score == score {
			break
		}
		cur = cur.pre
	}
	if cur == nil {
		return false
	}

	if len(cur.members) == 1 {
		if cur.next != nil {
			cur.next.pre = cur.pre
		} else {
			sm.tail = cur.pre
		}
		if cur.pre != nil {
			cur.pre.next = cur.next
		} else {
			sm.head = cur.next
		}
		cur = nil
		sm.size -= 1
		return true
	}

	i := -1
	for j, v := range cur.members {
		if v == member {
			i = j
			break
		}
	}
	cur.members = append(cur.members[0:i], cur.members[i+1:]...)
	sm.size -= 1
	return true
}

func (sm *scoreMemberMap) put(score float64, member string) {
	if sm.head == nil {
		sm.head = &scoreMember{score: score, members: []string{member}}
		sm.tail = sm.head
		sm.size += 1
		return
	}
	if sm.tail.score < score {
		element := &scoreMember{score: score, members: []string{member}}
		element.pre = sm.tail
		sm.tail.next = element
		sm.tail = element
		sm.size += 1
		return
	}
	cur := sm.tail
	for cur != nil {
		if cur.score <= score {
			break
		}
		cur = cur.pre
	}
	if cur == nil {
		element := &scoreMember{score: score, members: []string{member}}
		element.next = sm.head
		sm.head.pre = element
		sm.head = element
		sm.size += 1
	} else if cur.score < score {
		element := &scoreMember{score: score, members: []string{member}}
		cur.next.pre = element
		element.next = cur.next
		element.pre = cur
		cur.next = element
		sm.size += 1
	} else { //cur.score == score
		//check if member exists
		for _, v := range cur.members {
			if v == member {
				return
			}
		}
		sm.size += 1
		cur.members = append(cur.members, member)
		sort.Strings(cur.members)
	}
}

func (sm *scoreMemberMap) get(score float64) []string {
	for cur := sm.head; cur != nil; cur = cur.next {
		if cur.score == score {
			return cur.members
		}
	}
	return nil
}

// count returns the number of elements in the sorted set at key with a score between min and max.
func (sm *scoreMemberMap) count(min, max float64) int {
	n := 0
	for cur := sm.head; cur != nil && cur.score <= max; cur = cur.next {
		if cur.score >= min {
			n += len(cur.members)
		}
	}
	return n
}

func (sm *scoreMemberMap) rangeByScore(min, max float64) []string {
	var ret []string
	for cur := sm.head; cur != nil && cur.score <= max; cur = cur.next {
		if cur.score >= min {
			ret = append(ret, cur.members...)
		}
	}
	return ret
}

func (sm *scoreMemberMap) doRange(start, stop int, apply func(score float64, member string)) {
	size := sm.size
	if start < 0 {
		start = size + start
	}
	if stop < 0 {
		stop = size + stop
	}
	if start < 0 {
		start = 0
	}
	if start > stop || start >= size {
		return
	}
	if stop >= size {
		stop = size - 1
	}
	i := 0
	for cur := sm.head; cur != nil; cur = cur.next {
		for _, m := range cur.members {
			if i < start {
				i++
				continue
			}
			if i > stop {
				break
			}
			apply(cur.score, m)
			i++
		}
		if i > stop {
			break
		}
	}
}

func (sm *scoreMemberMap) rangeByIndex(start, stop int) []string {
	var ret []string
	sm.doRange(start, stop, func(_ float64, member string) {
		ret = append(ret, member)
	})
	return ret
}

func (sm *scoreMemberMap) rangeByIndexWithScore(start, stop int) []*zsetMember {
	var ret []*zsetMember
	sm.doRange(start, stop, func(score float64, member string) {
		ret = append(ret, &zsetMember{Member: member, Score: score})
	})
	return ret
}

func (sm *scoreMemberMap) rangeByScoreWithScore(min, max float64) []*zsetMember {
	var ret []*zsetMember
	for cur := sm.head; cur != nil && cur.score <= max; cur = cur.next {
		if cur.score >= min {
			for _, m := range cur.members {
				ret = append(ret, &zsetMember{Member: m, Score: cur.score})
			}
		}
	}
	return ret
}

func (sm *scoreMemberMap) rank(member string) int {
	i := 0
	for cur := sm.head; cur != nil; cur = cur.next {
		for _, m := range cur.members {
			if m == member {
				return i
			}
			i++
		}
	}
	return -1
}

func (sm *scoreMemberMap) revrank(member string) int {
	n := 0
	for cur := sm.tail; cur != nil; cur = cur.pre {
		for i := len(cur.members) - 1; i >= 0; i-- {
			if cur.members[i] == member {
				return n
			}
			n++
		}
	}
	return -1
}

type zsetVal struct {
	msMap    map[string]float64 //key:member,value:score
	smMap    *scoreMemberMap
	expireAt int64
}

func newZset() *zsetVal {
	return &zsetVal{msMap: make(map[string]float64), smMap: new(scoreMemberMap), expireAt: -1}
}

func (s *zsetVal) isAlive() bool {
	if s.expireAt == -1 {
		return true
	}
	ttl := s.expireAt - time.Now().Unix()
	return ttl >= 0
}

func (s *zsetVal) ttl() int {
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

func (s *zsetVal) setExpireAt(at int64) {
	s.expireAt = at
}

func (s *zsetVal) dataType() string {
	return "zset"
}

// return the number of elements added to the sorted set,
// not including elements already existing for which the score was updated.
func (s *zsetVal) add(score float64, member string) int {
	oldScore, exist := s.msMap[member]
	if exist {
		if oldScore != score { //update only when score changed
			s.msMap[member] = score
			s.smMap.remove(oldScore, member)
			s.smMap.put(score, member)
		}
		return 0
	}

	//add new member
	s.msMap[member] = score
	s.smMap.put(score, member)
	return 1
}

func (s *zsetVal) card() int {
	return len(s.msMap)
}

func (s *zsetVal) count(min, max float64) int {
	return s.smMap.count(min, max)
}

func (s *zsetVal) rangeByIndex(start, stop int) []string {
	return s.smMap.rangeByIndex(start, stop)
}

func (s *zsetVal) rangeByIndexWithScore(start, stop int) []string {
	array := s.smMap.rangeByIndexWithScore(start, stop)
	var ret []string
	for _, m := range array {
		ret = append(ret, m.Member, fmt.Sprintf("%f", m.Score))
	}
	return ret
}

func (s *zsetVal) rangeByScore(min, max float64) []string {
	return s.smMap.rangeByScore(min, max)
}

func (s *zsetVal) rangeByScoreWithScore(min, max float64) []string {
	array := s.smMap.rangeByScoreWithScore(min, max)
	var ret []string
	for _, m := range array {
		ret = append(ret, m.Member, fmt.Sprintf("%f", m.Score))
	}
	return ret
}

func (s *zsetVal) rank(member string) *int {
	_, exist := s.msMap[member]
	if !exist {
		return nil
	}
	i := s.smMap.rank(member)
	if i < 0 {
		panic("illegal state, rank is little than 0")
	}
	return &i
}

func (s *zsetVal) revrank(member string) *int {
	_, exist := s.msMap[member]
	if !exist {
		return nil
	}
	i := s.smMap.revrank(member)
	if i < 0 {
		panic("illegal state, rank is little than 0")
	}
	return &i
}

func (s *zsetVal) remove(members []string) int {
	n := 0
	for _, m := range members {
		score, exist := s.msMap[m]
		if !exist {
			continue
		}
		s.smMap.remove(score, m)
		n++
	}
	return n
}

func (s *zsetVal) score(member string) *string {
	score, exist := s.msMap[member]
	if !exist {
		return nil
	}
	scoreStr := fmt.Sprintf("%f", score)
	return &scoreStr
}

func zsetOf(key string) (*zsetVal, error) {
	z, ok := values[key]
	if !ok || !z.isAlive() {
		return nil, nil
	}
	zset, ok := z.(*zsetVal)
	if !ok {
		return nil, errorWrongType
	}
	return zset, nil
}

func Zadd(key string, score float64, member string) (int, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return -1, err
	}
	if zset == nil {
		zset = newZset()
		values[key] = zset
	}
	return zset.add(score, member), nil
}

// Zcard returns the cardinality (number of elements) of the sorted set, or 0 if key does not exist.
func Zcard(key string) (int, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return -1, err
	}
	if zset == nil {
		return 0, nil
	}
	return zset.card(), nil
}

// Zcount returns the number of elements in the sorted set at key with a score between min and max.
func Zcount(key string, min, max float64) (int, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return -1, err
	}
	if zset == nil {
		return 0, nil
	}
	return zset.count(min, max), nil
}

func ZrangeWithScore(key string, start, stop int) ([]string, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return nil, err
	}
	if zset == nil {
		return nil, nil
	}
	return zset.rangeByIndexWithScore(start, stop), nil
}

func Zrange(key string, start, stop int) ([]string, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return nil, err
	}
	if zset == nil {
		return nil, nil
	}
	return zset.rangeByIndex(start, stop), nil
}

func ZRangeByScore(key string, min, max float64) ([]string, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return nil, err
	}
	if zset == nil {
		return nil, nil
	}
	return zset.rangeByScore(min, max), nil
}

func ZRangeByScoreWithScore(key string, min, max float64) ([]string, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return nil, err
	}
	if zset == nil {
		return nil, nil
	}
	return zset.rangeByScoreWithScore(min, max), nil
}

func Zrank(key, member string) (*int, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return nil, err
	}
	if zset == nil {
		return nil, nil
	}
	return zset.rank(member), nil
}

func Zrevrank(key, member string) (*int, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return nil, err
	}
	if zset == nil {
		return nil, nil
	}
	return zset.revrank(member), nil
}

func Zrem(key string, members []string) (int, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return -1, err
	}
	if zset == nil {
		return 0, nil
	}
	return zset.remove(members), nil
}

func Zscore(key, member string) (*string, error) {
	zset, err := zsetOf(key)
	if err != nil {
		return nil, err
	}
	if zset == nil {
		return nil, nil
	}
	return zset.score(member), nil
}
