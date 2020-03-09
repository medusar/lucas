package command

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
)

//https://redis.io/commands/publish
var publishFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'publish' command")
	}
	count := store.Publish(args[0], args[1])
	return r.WriteInteger(count)
}

//-ERR only (P)SUBSCRIBE / (P)UNSUBSCRIBE / PING / QUIT allowed in this context
var subscribeFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) < 1 {
		return r.WriteError("ERR wrong number of arguments for 'subscribe' command")
	}
	store.Subscribe(r, args)
	return nil
}
