package store

import (
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test_scoreMemberMap_put(t *testing.T) {
	sm := new(scoreMemberMap)

	type pair struct {
		k float64
		v []string
	}

	sm.put(1, "1")
	sm.put(2, "2")
	sm.put(9, "9")
	sm.put(7, "7")
	sm.put(-100, "-100")
	sm.put(9, "01")
	sm.put(9, "02")
	sm.put(9, "03")
	sm.put(9, "03")
	log.Println(sm)

	pairs := []pair{
		{-100, []string{"-100"}},
		{1, []string{"1"}},
		{2, []string{"2"}},
		{7, []string{"7"}},
		{9, []string{"01", "02", "03", "9"}},
	}

	cur := sm.head
	for _, p := range pairs {
		assert.Equal(t, p.k, cur.score)
		assert.ElementsMatch(t, p.v, cur.members)
		cur = cur.next
	}
}

func Test_scoreMemberMap_remove(t *testing.T) {
	sm := new(scoreMemberMap)
	removed := sm.remove(1.0, "noexist")
	assert.False(t, removed)

	sm.put(1, "100")
	removed = sm.remove(1, "100")
	assert.True(t, removed)
	members := sm.get(1)
	assert.Nil(t, members)

	sm.put(1, "1")
	sm.put(2, "2")
	sm.put(9, "9")
	sm.put(7, "7")
	sm.put(-100, "-100")
	sm.put(9, "01")
	sm.put(9, "02")
	sm.put(9, "03")
	sm.put(9, "03")
	log.Println(sm)

	removed = sm.remove(1, "1")
	assert.True(t, removed)
	assert.Nil(t, sm.get(1))
	log.Println(sm)
	removed = sm.remove(1, "noexist")
	assert.False(t, removed)
	log.Println(sm)

	removed = sm.remove(-1000, "sss")
	assert.False(t, removed)
	log.Println(sm)

	removed = sm.remove(9, "03")
	assert.True(t, removed)
	assert.True(t, reflect.DeepEqual([]string{"01", "02", "9"}, sm.get(9)))
}

func Test_scoreMemberMap_get(t *testing.T) {
	sm := new(scoreMemberMap)

	assert.Nil(t, sm.get(-1))
	sm.put(1, "1")
	sm.put(2, "2")
	sm.put(9, "9")
	sm.put(7, "7")
	sm.put(-100, "-100")
	sm.put(9, "01")
	sm.put(9, "02")
	sm.put(9, "03")
	sm.put(9, "03")
	assert.True(t, reflect.DeepEqual(sm.get(1), []string{"1"}))
	assert.True(t, reflect.DeepEqual(sm.get(2), []string{"2"}))
	assert.True(t, reflect.DeepEqual(sm.get(9), []string{"01", "02", "03", "9"}))
	assert.True(t, reflect.DeepEqual(sm.get(7), []string{"7"}))
}

func Test_scoreMemberMap_count(t *testing.T) {
	sm := new(scoreMemberMap)
	assert.Equal(t, 0, sm.count(-1, 100))
	sm.put(1, "1")
	sm.put(2, "2")
	sm.put(9, "9")
	sm.put(7, "7")
	sm.put(-100, "-100")
	sm.put(9, "01")
	sm.put(9, "02")
	sm.put(9, "03")
	sm.put(9, "03")
	log.Println(sm)

	assert.Equal(t, 0, sm.count(-1, -100))
	assert.Equal(t, 1, sm.count(-100, -1))
	assert.Equal(t, 1, sm.count(-100, -100))
	assert.Equal(t, 4, sm.count(8, 100))
	assert.Equal(t, 5, sm.count(7, 100))
	assert.Equal(t, 8, sm.count(-100, 100))
}

func Test_scoreMemberMap_rangeByIndex(t *testing.T) {
	sm := new(scoreMemberMap)
	assert.Equal(t, 0, len(sm.rangeByIndex(-1, 100)))

	sm.put(1, "1")
	sm.put(2, "2")
	sm.put(9, "9")
	sm.put(7, "7")
	sm.put(-100, "-100")
	sm.put(9, "01")
	sm.put(9, "02")
	sm.put(9, "03")
	sm.put(9, "03")
	log.Println(sm)

	assert.ElementsMatch(t, sm.rangeByIndex(0, 1), []string{"-100", "1"})
	assert.ElementsMatch(t, sm.rangeByIndex(1, 3), []string{"1", "2", "7"})
	assert.ElementsMatch(t, sm.rangeByIndex(1, 1), []string{"1"})
	assert.ElementsMatch(t, sm.rangeByIndex(4, 7), []string{"01", "02", "03", "9"})
	assert.ElementsMatch(t, sm.rangeByIndex(3, 7), []string{"7", "01", "02", "03", "9"})
	assert.ElementsMatch(t, sm.rangeByIndex(-4, -1), []string{"01", "02", "03", "9"})
	assert.ElementsMatch(t, sm.rangeByIndex(-100, -8), []string{"-100"})
	assert.ElementsMatch(t, sm.rangeByIndex(-7, -6), []string{"1", "2"})
	assert.Nil(t, sm.rangeByIndex(-7, 0))
	assert.ElementsMatch(t, sm.rangeByIndex(-9, 0), []string{"-100"})
}

func Test_scoreMemberMap_rangeByIndexWithScore(t *testing.T) {
	sm := new(scoreMemberMap)
	assert.Equal(t, 0, len(sm.rangeByIndex(-1, 100)))
	sm.put(1, "1")
	sm.put(2, "2")
	sm.put(9, "9")
	sm.put(7, "7")
	sm.put(-100, "-100")
	sm.put(9, "01")
	sm.put(9, "02")
	sm.put(9, "03")
	sm.put(9, "03")
	log.Println(sm)
	assert.ElementsMatch(t, sm.rangeByIndexWithScore(0, 1), []*ZsetMember{{"-100", -100}, {"1", 1}})
	assert.Nil(t, sm.rangeByIndexWithScore(-7, 0))
}

func Test_scoreMemberMap_rangeByScore(t *testing.T) {
	sm := new(scoreMemberMap)
	assert.Nil(t, sm.rangeByScore(-1, 1000))
	sm.put(1, "1")
	sm.put(2, "2")
	sm.put(9, "9")
	sm.put(7, "7")
	sm.put(-100, "-100")
	sm.put(9, "01")
	sm.put(9, "02")
	sm.put(9, "03")
	sm.put(9, "03")

	assert.ElementsMatch(t, sm.rangeByScore(-100, 0), []string{"-100"})
	assert.Nil(t, sm.rangeByScore(1000, 2000))
	assert.ElementsMatch(t, sm.rangeByScore(2, 9), []string{"2", "7", "01", "02", "03", "9"})
}

func Test_scoreMemberMap_count1(t *testing.T) {
	sm := new(scoreMemberMap)
	assert.Equal(t, sm.count(-1, 1000), 0)

	sm.put(1, "1")
	sm.put(2, "2")
	sm.put(9, "9")
	sm.put(7, "7")
	sm.put(-100, "-100")
	sm.put(9, "01")
	sm.put(9, "02")
	sm.put(9, "03")
	sm.put(9, "03")
	assert.Equal(t, sm.count(-100, 0), 1)
	assert.Equal(t, sm.count(1000, 200), 0)
	assert.Equal(t, sm.count(1000, 2000), 0)
	assert.Equal(t, sm.count(7, 2), 0)
	assert.Equal(t, sm.count(2, 9), 6)
}
