package command

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"strconv"
)

var ttlFunc = func(args []string, r protocol.RedisRW) error {
	var err error
	if len(args) != 1 {
		err = r.WriteError("ERR wrong number of arguments for 'ttl' command")
		return err
	}
	ttl := store.Ttl(args[0])
	return r.WriteInteger(ttl)
}

var expireFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'expire' command")
	}
	sec, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}
	set := store.Expire(args[0], sec)
	if set {
		return r.WriteInteger(1)
	}
	return r.WriteInteger(0)
}

var expireAtFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'expire' command")
	}

	timestamp, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}

	set := store.ExpireAt(args[0], int64(timestamp))
	if set {
		return r.WriteInteger(1)
	}
	return r.WriteInteger(0)
}

var keysFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'keys' command")
	}
	keys := store.Keys(args[0])
	if keys == nil || len(keys) == 0 {
		return r.WriteArray(nil)
	}
	return r.WriteArray(toBulkArray(keys))
}

var existsFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) == 0 {
		return r.WriteError("ERR wrong number of arguments for 'exists' command")
	}
	total := 0

	for _, key := range args {
		if store.Exists(key) {
			total = total + 1
		}
	}

	return r.WriteInteger(total)
}

var delFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) == 0 {
		return r.WriteError("ERR wrong number of arguments for 'del' command")
	}
	total := 0
	for _, key := range args {
		if store.Del(key) {
			total = total + 1
		}
	}
	return r.WriteInteger(total)
}

var typeFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'type' command")
	}
	t := store.Type(args[0])
	return r.WriteString(t)
}
