package cmd

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"strconv"
)

var getFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'get' command")
	}
	val, e := store.Get(args[0])
	if e != nil {
		return r.WriteError(e.Error())
	}
	if val == nil {
		return r.WriteNil()
	}
	return r.WriteBulk(*val)
}

var setFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'set' command")
	}
	store.Set(args[0], args[1])
	return r.WriteString("OK") //TODO:support NX, EX
}

var getsetFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'getset' command")
	}
	val, e := store.GetSet(args[0], args[1])
	if e != nil {
		return r.WriteError(e.Error())
	}
	if val == nil {
		return r.WriteNil()
	}
	return r.WriteBulk(*val)
}

var setexFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'setex' command")
	}
	key, sec, val := args[0], args[1], args[2]
	ttl, err := strconv.Atoi(sec)
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}

	err = store.SetEX(key, val, ttl)
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteString("OK")
}

var setnxFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'setnx' command")
	}
	key, val := args[0], args[1]
	if set := store.SetNX(key, val); set {
		return r.WriteInteger(1)
	} else {
		return r.WriteInteger(0)
	}
}

var setRangeFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'setrange' command")
	}
	key, offset, val := args[0], args[1], args[2]
	intOff, e := strconv.Atoi(offset)
	if e != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}
	l, e := store.SetRange(key, val, intOff)
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteInteger(l)
}

var getRangeFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 3 {
		return r.WriteError("ERR wrong number of arguments for 'getrange' command")
	}
	key, start, end := args[0], args[1], args[2]
	startInt, e := strconv.Atoi(start)
	if e != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}
	endInt, e := strconv.Atoi(end)
	if e != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}
	v, e := store.GetRange(key, startInt, endInt)
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteBulk(v)
}

var appendFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'append' command")
	}
	l, e := store.Append(args[0], args[1])
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteInteger(l)
}

var mgetFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) == 0 {
		return r.WriteError("ERR wrong number of arguments for 'mget' command")
	}
	values := store.Mget(args)
	ret := make([]*protocol.Resp, len(values))
	for i, v := range values {
		if v == nil {
			ret[i] = protocol.NewNil()
		} else {
			ret[i] = protocol.NewBulk(*v)
		}
	}
	return r.WriteArray(ret)
}

//https://redis.io/commands/mset
var msetFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) == 0 || len(args)%2 != 0 {
		return r.WriteError("ERR wrong number of arguments for 'mset' command")
	}
	store.Mset(args)
	return r.WriteString("OK")
}

var strlenFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'strlen' command")
	}
	n, e := store.StrLen(args[0])
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteInteger(n)
}

var incrFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'incr' command")
	}
	v, e := store.Incr(args[0])
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteInteger(v)
}

var incrByFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'incrby' command")
	}

	key, val := args[0], args[1]
	intV, e := strconv.Atoi(val)
	if e != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}

	v, e := store.IncrBy(key, intV)
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteInteger(v)
}

var decrFun = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'decr' command")
	}
	v, e := store.IncrBy(args[0], -1)
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteInteger(v)
}

var decrByFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'decrby' command")
	}

	key, val := args[0], args[1]
	intV, e := strconv.Atoi(val)
	if e != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}

	v, e := store.IncrBy(key, -1*intV)
	if e != nil {
		return r.WriteError(e.Error())
	}
	return r.WriteInteger(v)
}
