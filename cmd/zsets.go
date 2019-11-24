package cmd

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"strconv"
	"strings"
)

//https://redis.io/commands/zadd
var zaddFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 3 {
		return r.WriteError("ERR wrong number of arguments for 'zadd' command")
	}
	score, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return r.WriteError("ERR value is not a valid float")
	}
	i, err := store.Zadd(args[0], score, args[2])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(i)
}

//https://redis.io/commands/zcard
var zcardFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'zcard' command")
	}
	n, e := store.Zcard(args[0])
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/zcount
var zcountFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'zcount' command")
	}

	min, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return r.WriteError("ERR value is not a valid float")
	}
	max, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return r.WriteError("ERR value is not a valid float")
	}

	n, err := store.Zcount(args[0], min, max)
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/zrange
var zrangeFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 3 && len(args) != 4 {
		return r.WriteError("ERR wrong number of arguments for 'zrange' command")
	}

	start, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}
	stop, err := strconv.Atoi(args[2])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}

	var array []string
	if len(args) == 4 {
		if strings.ToUpper(args[3]) != "WITHSCORES" {
			return r.WriteError("ERR syntax error")
		}
		array, err = store.ZrangeWithScore(args[0], start, stop)
	} else {
		array, err = store.Zrange(args[0], start, stop)
	}

	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteArray(toBulkArray(array))
}

//https://redis.io/commands/zrangebyscore
var zrangeByScoreFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 3 && len(args) != 4 {
		return r.WriteError("ERR wrong number of arguments for 'zrangebyscore' command")
	}

	min, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}
	max, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}

	var array []string
	if len(args) == 4 {
		if strings.ToUpper(args[3]) != "WITHSCORES" {
			return r.WriteError("ERR syntax error")
		}
		array, err = store.ZRangeByScoreWithScore(args[0], min, max)
	} else {
		array, err = store.ZRangeByScore(args[0], min, max)
	}
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteArray(toBulkArray(array))
}
