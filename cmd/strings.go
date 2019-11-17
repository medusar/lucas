package cmd

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"strconv"
)

var getFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 1 {
		err = r.WriteError("ERR wrong number of arguments for 'get' command")
		return err
	}
	val, ok, e := store.Get(args[0])
	if e != nil {
		err = r.WriteError(e.Error())
	} else if ok {
		err = r.WriteBulk(val)
	} else {
		err = r.WriteNil()
	}
	return err
}

var setFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 2 {
		err = r.WriteError("ERR wrong number of arguments for 'set' command")
		return err
	}

	store.Set(args[0], args[1])
	err = r.WriteString("OK") //TODO:support NX, EX
	return err
}
var getsetFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 2 {
		err = r.WriteError("ERR wrong number of arguments for 'getset' command")
		return err
	}
	v, ok, e := store.GetSet(args[0], args[1])
	if e != nil {
		err = r.WriteError(e.Error())
	} else if ok {
		err = r.WriteBulk(v)
	} else {
		err = r.WriteNil()
	}
	return err
}

var setexFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 3 {
		err = r.WriteError("ERR wrong number of arguments for 'setex' command")
		return err
	}
	key, sec, val := args[0], args[1], args[2]

	ttl, err := strconv.Atoi(sec)
	if err != nil {
		err = r.WriteError("ERR value is not an integer or out of range")
		return err
	}
	err = store.SetEX(key, val, ttl)
	if err != nil {
		err = r.WriteError(err.Error())
	} else {
		err = r.WriteString("OK")
	}
	return err
}

var setnxFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 2 {
		err = r.WriteError("ERR wrong number of arguments for 'setnx' command")
		return err
	}
	key, val := args[0], args[1]
	if set := store.SetNX(key, val); set {
		err = r.WriteInteger(1)
	} else {
		err = r.WriteInteger(0)
	}
	return err
}

var setRangeFunc = func(args []string, r *protocol.RedisConn) error {
	if args == nil || len(args) != 3 {
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
	} else {
		return r.WriteInteger(l)
	}
}

var getRangeFunc = func(args []string, r *protocol.RedisConn) error {
	if args == nil || len(args) != 3 {
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
	if args == nil || len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'append' command")
	}
	l, e := store.Append(args[0], args[1])
	if e != nil {
		return r.WriteError(e.Error())
	} else {
		return r.WriteInteger(l)
	}
}

var mgetFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) == 0 {
		err = r.WriteError("ERR wrong number of arguments for 'mget' command")
		return err
	}

	ret := make([]*protocol.Resp, 0)
	for _, arg := range args {
		v, ok, e := store.Get(arg)
		if e != nil || !ok {
			ret = append(ret, protocol.NewNil())
		} else {
			ret = append(ret, protocol.NewBulk(v))
		}
	}
	return r.WriteArray(ret)
}

var msetFunc = func(args []string, r *protocol.RedisConn) error {
	if args == nil || len(args) == 0 || len(args)%2 != 0 {
		return r.WriteError("ERR wrong number of arguments for 'mset' command")
	}

	l := len(args)
	for i := 0; i < l; i = i + 2 {
		store.Set(args[i], args[i+1])
	}

	return r.WriteString("OK")
}

var strlenFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 1 {
		err = r.WriteError("ERR wrong number of arguments for 'strlen' command")
		return err
	}
	n, e := store.StrLen(args[0])
	if e != nil {
		err = r.WriteError(e.Error())
	} else {
		err = r.WriteInteger(n)
	}
	return err
}

var incrFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 1 {
		err = r.WriteError("ERR wrong number of arguments for 'incr' command")
		return err
	}
	v, e := store.Incr(args[0])
	if e != nil {
		err = r.WriteError(e.Error())
	} else {
		err = r.WriteInteger(v)
	}
	return err
}

var incrByFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 2 {
		err = r.WriteError("ERR wrong number of arguments for 'incrby' command")
		return err
	}

	key, val := args[0], args[1]
	intV, e := strconv.Atoi(val)
	if e != nil {
		err = r.WriteError("ERR value is not an integer or out of range")
		return err
	}

	v, e := store.IncrBy(key, intV)
	if e != nil {
		err = r.WriteError(e.Error())
	} else {
		err = r.WriteInteger(v)
	}
	return err
}

var decrFun = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 1 {
		err = r.WriteError("ERR wrong number of arguments for 'decr' command")
		return err
	}
	v, e := store.IncrBy(args[0], -1)
	if e != nil {
		err = r.WriteError(e.Error())
	} else {
		err = r.WriteInteger(v)
	}
	return err
}

var decrByFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 2 {
		err = r.WriteError("ERR wrong number of arguments for 'decrby' command")
		return err
	}

	key, val := args[0], args[1]
	intV, e := strconv.Atoi(val)
	if e != nil {
		err = r.WriteError("ERR value is not an integer or out of range")
		return err
	}

	v, e := store.IncrBy(key, -1*intV)
	if e != nil {
		err = r.WriteError(e.Error())
	} else {
		err = r.WriteInteger(v)
	}
	return err
}
