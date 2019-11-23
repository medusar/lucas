package store

import (
	"github.com/medusar/lucas/util"
	"time"
)

type listVal struct {
	val      []string
	expireAt int64
}

//Only when val is not empty and ttl larger than 0
func (s *listVal) isAlive() bool {
	if s.len() == 0 { //when list is empty, it is regarded as expired
		return false
	}

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
		copy(list[1:], list[:])
		list[0] = e
	}
	s.val = list
	return len(s.val)
}

func (s *listVal) rpush(elements []string) int {
	list := s.val
	list = append(list, elements...)
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
}

func (s *listVal) rpop() string {
	l := len(s.val)
	if l > 0 {
		r := s.val[l-1]
		s.val = s.val[:l-1]
		return r
	}
	//should not happen by out caller
	return ""
}

func (s *listVal) lindex(i int) (string, bool, error) {
	l := s.len()
	if i < 0 {
		i = l + i
	}
	if i < 0 || i > l-1 {
		return "", false, nil
	}
	return s.val[i], true, nil
}

// The count argument influences the operation in the following ways:
// 1) count > 0: Remove elements equal to element moving from head to tail.
// 2) count < 0: Remove elements equal to element moving from tail to head.
// 3) count = 0: Remove all elements equal to element.
func (s *listVal) rem(count int, element string) (int, error) {
	removed := 0
	if count > 0 {
		for i, v := range s.val {
			if v == element {
				s.val = util.DeleteStringArray(i, s.val)
				removed++
			}
			if removed == count {
				break
			}
		}

	} else if count == 0 {
		i := 0
		for _, v := range s.val {
			if v != element {
				s.val[i] = v
				i++
			} else {
				removed++
			}
		}
		s.val = s.val[:i]
	} else { //count < 0
		count = -1 * count
		for i := len(s.val) - 1; i >= 0; i-- {
			if s.val[i] == element {
				s.val = util.DeleteStringArray(i, s.val)
				removed++
			}
			if removed == count {
				break
			}
		}
	}
	return removed, nil
}

func (s *listVal) set(index int, element string) error {
	l := s.len()

	if index < 0 {
		index = l + index
	}
	if index < 0 || index > l-1 {
		return errorIndexOutOfRange
	}

	s.val[index] = element
	return nil
}

func (s *listVal) lrange(start, end int) []string {
	l := s.len()

	if start < 0 {
		start = l + start
	}
	if end < 0 {
		end = l + end
	}

	if start > l-1 || end < 0 {
		return nil
	}

	//include end
	end = end + 1

	if start < 0 {
		start = 0
	}
	if end > l {
		end = l
	}

	if start >= end {
		return nil
	}

	return s.val[start:end]
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

func Lpush(key string, elements []string) (int, error) {
	lv, err := getOrCreateList(key)
	if err != nil {
		return -1, err
	}
	return lv.lpush(elements), nil
}

// LpushX inserts specified values at the head of the list stored at key,
// only if key already exists and holds a list.
// In contrary to LPUSH, no operation will be performed when key does not yet exist.
func LpushX(key string, elements []string) (int, error) {
	list, err := listOf(key)
	if err != nil {
		return -1, err
	}
	if list == nil {
		return 0, nil
	}
	return list.lpush(elements), nil
}

func Rpush(key string, elements []string) (int, error) {
	lv, err := getOrCreateList(key)
	if err != nil {
		return -1, err
	}
	return lv.rpush(elements), nil
}

// RpushX inserts specified values at the tail of the list stored at key,
// only if key already exists and holds a list.
// In contrary to RPUSH, no operation will be performed when key does not yet exist.
func RpushX(key string, elements []string) (int, error) {
	list, err := listOf(key)
	if err != nil {
		return -1, err
	}
	if list == nil {
		return 0, nil
	}
	return list.rpush(elements), nil
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

// Returns the element at index index in the list stored at key.
// The index is zero-based, so 0 means the first element, 1 the second element and so on.
// Negative indices can be used to designate elements starting at the tail of the list.
// Here, -1 means the last element, -2 means the penultimate and so forth.
//
// When the value at key is not a list, an error is returned.
func Lindex(key string, idx int) (string, bool, error) {
	lv, err := listOf(key)
	if err != nil {
		return "", false, err
	}
	if lv == nil {
		return "", false, nil
	}
	return lv.lindex(idx)
}

// Removes the first count occurrences of elements equal to element from the list stored at key.
// The count argument influences the operation in the following ways:
// 1) count > 0: Remove elements equal to element moving from head to tail.
// 2) count < 0: Remove elements equal to element moving from tail to head.
// 3) count = 0: Remove all elements equal to element.
// Note that non-existing keys are treated like empty lists,
// so when key does not exist, the command will always return 0.
func Lrem(key string, count int, element string) (int, error) {
	lv, err := listOf(key)
	if err != nil {
		return -1, err
	}
	if lv == nil {
		return 0, nil
	}
	return lv.rem(count, element)
}

// Sets the list element at index to element.
// For more information on the index argument, see Lindex.
// An error is returned for out of range indexes.
func Lset(key string, index int, element string) error {
	lv, err := listOf(key)
	if err != nil {
		return err
	}
	if lv == nil {
		return errorNoSuchKey
	}
	return lv.set(index, element)
}

// Lrange returns the specified elements of the list stored at key, and index range from start to end.
// Note: both elements at start and end are included.
// Out of range indexes will not produce an error.
// If start is larger than the end of the list, an empty list is returned.
// If stop is larger than the actual end of the list,
// it will treat it like the last element of the list.
func Lrange(key string, start, end int) ([]string, error) {
	list, err := listOf(key)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}

	return list.lrange(start, end), nil
}
