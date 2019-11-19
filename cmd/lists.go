package cmd

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
)

// https://redis.io/commands/lpush
var lpushFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) < 2 {
		return r.WriteError("ERR wrong number of arguments for 'lpush' command")
	}
	l, err := store.Lpush(args[0], args[1:])
	if err != nil {
		return r.WriteError(err.Error())
	}
	return r.WriteInteger(l)
}