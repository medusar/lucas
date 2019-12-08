package command

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
)

// https://redis.io/commands/hset
var hsetFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) < 3 || len(args)%2 != 1 {
		return r.WriteError("ERR wrong number of arguments for 'hset' command")
	}
	added := 0
	key := args[0]
	for i := 1; i < len(args); i = i + 2 {
		a, err := store.Hset(key, args[i], args[i+1])
		if err != nil {
			return r.WriteError(err.Error())
		}
		if a {
			added = added + 1
		}
	}
	return r.WriteInteger(added)
}

//https://redis.io/commands/hget
var hgetFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'hget' command")
	}
	key, field := args[0], args[1]
	v, exist, err := store.Hget(key, field)
	if err != nil {
		return r.WriteError(err.Error())
	}
	if exist {
		return r.WriteBulk(v)
	}
	return r.WriteNil()
}

//https://redis.io/commands/hgetall
var hgetAllFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'hgetall' command")
	}
	m, err := store.Hgetall(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	ret := make([]*protocol.Resp, 0)
	for k, v := range m {
		ret = append(ret, protocol.NewBulk(k), protocol.NewBulk(v))
	}
	return r.WriteArray(ret)
}

//https://redis.io/commands/hkeys
var hkeysFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'hkeys' command")
	}
	keys, err := store.Hkeys(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	ret := make([]*protocol.Resp, 0)
	for _, key := range keys {
		ret = append(ret, protocol.NewBulk(key))
	}
	return r.WriteArray(ret)
}

//https://redis.io/commands/hlen
var hlenFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'hlen' command")
	}
	l, err := store.Hlen(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(l)
}

//https://redis.io/commands/hexists
var hexistsFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'hexists' command")
	}
	l, err := store.Hexists(args[0], args[1])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(l)
}

//https://redis.io/commands/hdel
var hdelFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'hdel' command")
	}
	total, err := store.Hdel(args[0], args[1:])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(total)
}

//https://redis.io/commands/hmget
var hmgetFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'hmget' command")
	}
	ret := make([]*protocol.Resp, 0)
	for i := 1; i < len(args); i++ {
		v, exists, err := store.Hget(args[0], args[i])
		if err != nil {
			return r.WriteError(err.Error())
		}
		if exists {
			ret = append(ret, protocol.NewBulk(v))
		} else {
			ret = append(ret, protocol.NewNil())
		}
	}
	return r.WriteArray(ret)
}

//https://redis.io/commands/hmset
var hmsetFunc = func(args []string,r protocol.RedisRW) error {
	return r.WriteError("ERR 'hmset' is considered deprecated, please use 'hset' instead")
}

//https://redis.io/commands/hsetnx
var hsetnxFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'hsetnx' command")
	}
	l, err := store.HsetNX(args[0], args[1], args[2])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(l)
}

//https://redis.io/commands/hstrlen
var hstrlenFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'hstrlen' command")
	}
	l, err := store.HstrLen(args[0], args[1])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(l)
}

//https://redis.io/commands/hvals
var hvalsFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'hvals' command")
	}
	vals, err := store.Hvals(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteArray(toBulkArray(vals))
}

//https://redis.io/commands/hincrby
var hincrByFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'hincrby' command")
	}
	v, err := store.HincrBy(args[0], args[1], args[2])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(v)
}

//https://redis.io/commands/hincrbyfloat
var hincrByFloatFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'hincrbyfloat' command")
	}
	v, err := store.HincrByFloat(args[0], args[1], args[2])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteBulk(v)
}
