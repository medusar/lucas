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
	}
	return r.WriteInteger(0)
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

	removed, err := store.Spop(args[0], count)
	if err != nil {
		return r.WriteError(err.Error())
	}
	if removed != nil {
		return r.WriteArray(toBulkArray(*removed))
	}
	return r.WriteNil()
}

//https://redis.io/commands/sdiffstore
var sdiffStoreFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'sdiffstore' command")
	}

	var n int
	var err error

	if len(args) == 2 {
		n, err = store.SdiffStore(args[0], args[1])
	} else {
		n, err = store.SdiffStore(args[0], args[1], args[2:]...)
	}

	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/sinter
var sinterFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 1 {
		return r.WriteError("ERR wrong number of arguments for 'sinter' command")
	}
	if len(args) == 1 {
		return smembersFunc(args, r)
	}

	set, err := store.Sinter(args[0], args[1:]...)
	if err != nil {
		return r.WriteError(err.Error())
	}
	ret := make([]*protocol.Resp, len(set))
	for i := 0; i < len(set); i++ {
		ret[i] = protocol.NewBulk(set[i])
	}
	return r.WriteArray(ret)
}

//https://redis.io/commands/sinterstore
var sinterStoreFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'sinterstore' command")
	}

	var n int
	var err error

	if len(args) == 2 {
		n, err = store.SinterStore(args[0], args[1])
	} else {
		n, err = store.SinterStore(args[0], args[1], args[2:]...)
	}
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/srem
var sremFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'srem' command")
	}
	n, err := store.Srem(args[0], args[1:])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/sunion
var sunionFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 1 {
		return r.WriteError("ERR wrong number of arguments for 'sunion' command")
	}
	set, err := store.Sunion(args)
	if err != nil {
		return r.WriteError(err.Error())
	}
	ret := make([]*protocol.Resp, len(set))
	for i := 0; i < len(set); i++ {
		ret[i] = protocol.NewBulk(set[i])
	}
	return r.WriteArray(ret)
}

//https://redis.io/commands/sunionstore
var sunionStoreFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'sunionstore' command")
	}

	var n int
	var err error

	if len(args) == 2 {
		n, err = store.SunionStore(args[0], args[1])
	} else {
		n, err = store.SunionStore(args[0], args[1], args[2:]...)
	}
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/smove
var smoveFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'smove' command")
	}
	n, err := store.Smove(args[0], args[1], args[2])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(n)
}

//https://redis.io/commands/srandmember
//https://redis.io/commands/sscan
