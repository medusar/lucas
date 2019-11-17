package cmd

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"strconv"
)

//https://redis.io/commands/sadd
var saddFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'sadd' command")
	}
	l, err := store.Sadd(args[0], args[1:])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(l)
}

//https://redis.io/commands/scard
var scardFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'scard' command")
	}
	l, err := store.Scard(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(l)
}

//https://redis.io/commands/sdiff
var sdiffFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 1 {
		return r.WriteError("ERR wrong number of arguments for 'sdiff' command")
	}
	if len(args) == 1 {
		return smembersFunc(args, r)
	}

	set, err := store.Sdiff(args[0], args[1:]...)
	if err != nil {
		return r.WriteError(err.Error())
	}
	ret := make([]*protocol.Resp, len(set))
	for i := 0; i < len(set); i++ {
		ret[i] = protocol.NewBulk(set[i])
	}
	return r.WriteArray(ret)
}

//https://redis.io/commands/smembers
var smembersFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'smembers' command")
	}
	keys, err := store.Smembers(args[0])
	if err != nil {
		return r.WriteError(err.Error())
	}
	ret := make([]*protocol.Resp, len(keys))
	for i := 0; i < len(keys); i++ {
		ret[i] = protocol.NewBulk(keys[i])
	}
	return r.WriteArray(ret)
}

//https://redis.io/commands/sismember
var sismemberFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'sismember' command")
	}
	is, err := store.Sismember(args[0], args[1])
	if err != nil {
		return r.WriteError(err.Error())
	}
	if is {
		return r.WriteInteger(1)
	} else {
		return r.WriteInteger(0)
	}
}

//https://redis.io/commands/spop
var spopFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) > 2 || len(args) <= 0 {
		return r.WriteError("ERR wrong number of arguments for 'spop' command")
	}

	count := 1
	if len(args) == 2 {
		i, err := strconv.Atoi(args[1])
		if err != nil {
			return r.WriteError("ERR value is not an integer or out of range")
		}
		if i < 0 {
			return r.WriteError("ERR index out of range")
		}
		count = i
	}

	removed, exist, err := store.Spop(args[0], count)
	if err != nil {
		return r.WriteError(err.Error())
	}
	if exist {
		return r.WriteArray(toBulkArray(removed))
	} else {
		return r.WriteNil()
	}
}
