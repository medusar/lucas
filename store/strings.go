package store

import (
	"fmt"
	"strconv"
	"time"
)

const (
	//MaxStringLength is the limit of bytes that a redis string can hold, it is limited to 512 megabytes
	MaxStringLength = 536870911 //2^29-1
	//MaxBitOffset is the max offset of bit operation
	MaxBitOffset = 4294967295 //2^32-1
	//TODO:check string size
)

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

func (s *stringVal) dataType() string {
	return "string"
}

func (s *stringVal) getRange(start, end int) string {
	l := len(s.val)
	//check negative
	if start < 0 {
		start = l + start
	}
	if end < 0 {
		end = l + end
	}

	if start > end || start > l-1 || end < 0 {
		return ""
	}

	if start < 0 {
		start = 0
	}
	if end > l-1 {
		end = l - 1
	}

	rl := end - start + 1
	if rl <= 0 {
		return ""
	}

	return s.val[start : end+1]
}

func (s *stringVal) setBit(offset, bit int) int {
	byteIndex := offset / 8

	bytes := []byte(s.val)
	if byteIndex > len(bytes)-1 {
		totalBytes := byteIndex + 1
		bytes = make([]byte, totalBytes)
		copy(bytes, s.val)
	}

	i := int(bytes[byteIndex])
	bitIndex := uint(7 - offset%8)

	has := hasBit(i, bitIndex)
	if bit == 1 {
		i = setBit(i, bitIndex)
	} else {
		i = clearBit(i, bitIndex)
	}
	bytes[byteIndex] = byte(i)
	s.val = string(bytes)

	if has {
		return 1
	}
	return 0
}

func (s *stringVal) getBit(offset int) int {
	byteIndex := offset / 8
	bytes := []byte(s.val)
	if byteIndex > len(bytes)-1 {
		return 0
	}
	i := int(bytes[byteIndex])
	bitIndex := uint(7 - offset%8)
	has := hasBit(i, bitIndex)
	if has {
		return 1
	}
	return 0
}

func (s *stringVal) countBit(start, end int) int {
	l := len(s.val) * 8
	//check negative
	if start < 0 {
		start = l + start
	}
	if end < 0 {
		end = l + end
	}
	if start > end || start > l-1 || end < 0 {
		return 0
	}
	if start < 0 {
		start = 0
	}
	if end > l-1 {
		end = l - 1
	}

	rl := end - start + 1
	if rl <= 0 {
		return 0
	}

	startB := start / 8
	endB := end / 8

	total := 0
	for i := startB; i <= endB; i++ {
		b := int(s.val[i])
		for j := 0; j < 8; j++ {
			offset := i*8 + j
			if offset >= start && offset <= end {
				if hasBit(b, uint(7-j)) {
					total++
				}
			} else if offset > end {
				break
			}
		}
	}
	return total
}

func setBit(n int, pos uint) int {
	n |= 1 << pos
	return n
}
func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}
func hasBit(n int, pos uint) bool {
	return n&(1<<pos) > 0
}

func stringOf(key string) (*stringVal, error) {
	v, ok := values[key]
	if !ok || !v.isAlive() {
		return nil, nil
	}
	str, ok := v.(*stringVal)
	if !ok {
		return nil, errorWrongType
	}
	return str, nil
}

func Get(key string) (*string, error) {
	str, err := stringOf(key)
	if err != nil {
		return nil, err
	}

	if str == nil {
		return nil, nil
	}
	return &str.val, nil
}

func Set(key, val string) {
	values[key] = &stringVal{val: val, expireAt: -1}
}

func GetSet(key, val string) (*string, error) {
	str, err := stringOf(key)
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, nil
	}
	old := str.val
	str.val = val
	return &old, nil
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
	str, err := stringOf(key)
	if err != nil {
		return -1, err
	}
	if str == nil {
		return 0, nil
	}
	return len(str.val), nil
}

func Incr(key string) (int, error) {
	return IncrBy(key, 1)
}

//FIXME: "01234" can't be converted to in redis server.
func IncrBy(key string, intV int) (int, error) {
	str, err := stringOf(key)
	if err != nil {
		return -1, err
	}
	if str == nil {
		i := intV
		Set(key, strconv.Itoa(i))
		return i, nil
	}

	i, e := strconv.Atoi(str.val)
	if e != nil {
		return -1, errorInvalidInt
	}

	i = i + intV
	str.val = strconv.Itoa(i)
	return i, nil
}

//If key already exists and is a string, this command appends the value at the end of the string.
//If key does not exist it is created and set as an empty string,
//so APPEND will be similar to SET in this special case.
func Append(key, val string) (int, error) {
	str, err := stringOf(key)
	if err != nil {
		return -1, err
	}
	if str == nil {
		Set(key, val)
		return len(val), nil
	}
	str.val = str.val + val
	return len(str.val), nil
}

//offset means byte index, not rune index
func SetRange(key, val string, offset int) (int, error) {
	str, err := stringOf(key)
	if err != nil {
		return -1, err
	}

	if offset+len(val) > MaxStringLength {
		return -1, fmt.Errorf("ERR string exceeds maximum allowed size (512MB)")
	}

	if str == nil {
		rs := make([]byte, offset+len(val))
		for i := 0; i < len(val); i++ {
			rs[i+offset] = val[i]
		}
		Set(key, string(rs))
		return len(rs), nil
	}

	var rs []byte

	bs := []byte(str.val)
	newLen := offset + len(val)

	if newLen > len(bs) {
		rs = make([]byte, newLen)
		copy(rs, bs)
	} else {
		rs = bs
	}

	for i := 0; i < len(val); i++ {
		rs[offset+i] = val[i]
	}
	str.val = string(rs)
	return len(rs), nil
}

// GetRange returns the substring of the string value stored at key,
// determined by the offsets start and end (both are inclusive).
// Negative offsets can be used in order to provide an offset starting from the end of the string.
// So -1 means the last character, -2 the penultimate and so forth.
//
//The function handles out of range requests by limiting the resulting range to the actual length of the string.
func GetRange(key string, start, end int) (string, error) {
	str, err := stringOf(key)
	if err != nil {
		return "", err
	}
	if str == nil {
		return "", nil
	}
	return str.getRange(start, end), nil
}

// Mget returns the values of all specified keys.
// For every key that does not hold a string value or does not exist,
// the special value nil is returned.
// Because of this, the operation never fails.
func Mget(keys []string) []*string {
	if len(keys) == 0 {
		return nil
	}
	ret := make([]*string, len(keys))
	for i, k := range keys {
		str, err := stringOf(k)
		if err != nil || str == nil {
			ret[i] = nil
		} else {
			ret[i] = &str.val
		}
	}
	return ret
}

func Mset(kvs []string) {
	for i := 0; i < len(kvs); i = i + 2 {
		Set(kvs[i], kvs[i+1])
	}
}

// SetBit implements redis setbit commands.
// https://redis.io/commands/setbit
// Sets or clears the bit at offset in the string value stored at key.
// The bit is either set or cleared depending on value, which can be either 0 or 1.
// When key does not exist, a new string value is created.
// The string is grown to make sure it can hold a bit at offset.
// The offset argument is required to be greater than or equal to 0,
// and smaller than 232 (this limits bitmaps to 512MB).
// When the string at key is grown, added bits are set to 0.
// Return the original bit value stored at offset.
func SetBit(key string, offset, bit int) (int, error) {
	str, err := stringOf(key)
	if err != nil {
		return -1, err
	}
	if offset > MaxBitOffset {
		return -1, fmt.Errorf("ERR offset is not an integer or out of range")
	}
	if bit != 1 && bit != 0 {
		return -1, fmt.Errorf("ERR bit is not an integer or out of range")
	}
	if str == nil {
		str = &stringVal{val: "", expireAt: -1}
		values[key] = str
	}
	return str.setBit(offset, bit), nil
}

// GetBit returns the bit value at offset in the string value stored at key.
// When offset is beyond the string length, the string is assumed to be a contiguous space with 0 bits.
// When key does not exist it is assumed to be an empty string,
// so offset is always out of range and the value is also assumed to be a contiguous space with 0 bits.
func GetBit(key string, offset int) (int, error) {
	str, err := stringOf(key)
	if err != nil {
		return -1, err
	}
	if offset > MaxBitOffset {
		return -1, fmt.Errorf("ERR offset is not an integer or out of range")
	}
	if str == nil {
		return 0, nil
	}
	return str.getBit(offset), nil
}

func BitCount(key string, start, end int) (int, error) {
	str, err := stringOf(key)
	if err != nil {
		return -1, err
	}
	if str == nil {
		return 0, nil
	}
	return str.countBit(start, end), nil
}
