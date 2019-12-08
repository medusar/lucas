package command

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"strconv"
)

// https://redis.io/commands/lpush
var lpushFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'lpush' command")
	}
	n, err := store.Lpush(args[0], args[1:])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/rpush
var rpushFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'rpush' command")
	}
	n, err := store.Rpush(args[0], args[1:])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/llen
var llenFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'llen' command")
	}
	n, err := store.Llen(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/lpop
var lpopFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'lpop' command")
	}
	v, exists, err := store.Lpop(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	if exists {
		return r.WriteBulk(v)
	}
	return r.WriteNil()
}

//https://redis.io/commands/rpop
var rpopFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'rpop' command")
	}
	v, exists, err := store.Rpop(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	if exists {
		return r.WriteBulk(v)
	}
	return r.WriteNil()
}

//https://redis.io/commands/lindex
var lindexFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'lindex' command")
	}

	//TODO: return error only when the key exists
	idx, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of rang")
	}
	v, exists, err := store.Lindex(args[0], idx)
	if err != nil {
		return r.WriteError(err.Error())
	}
	if exists {
		return r.WriteBulk(v)
	}
	return r.WriteNil()
}

//https://redis.io/commands/lrem
var lremFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'lrem' command")
	}

	//TODO: return error only when the key exists
	count, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of rang")
	}

	n, err := store.Lrem(args[0], count, args[2])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/lset
var lsetFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'lrem' command")
	}

	//TODO: return error only when the key exists
	index, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of rang")
	}

	err = store.Lset(args[0], index, args[1])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteString("OK")
}

//https://redis.io/commands/rpushx
var rpushXFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'rpushx' command")
	}
	n, err := store.RpushX(args[0], args[1:])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/lpushx
var lpushXFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'lpushx' command")
	}
	n, err := store.LpushX(args[0], args[1:])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/lrange
var lrangeFunc = func(args []string,r protocol.RedisRW) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'lrange' command")
	}

	//TODO: return error only when the key exists
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of rang")
	}

	stop, err := strconv.Atoi(args[2])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of rang")
	}

	list, err := store.Lrange(args[0], start, stop)
	if err != nil {
		return r.WriteError(err.Error())
	}
	array := toBulkArray(list)
	return r.WriteArray(array)
}

//https://redis.io/commands/rpoplpush
