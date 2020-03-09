package command

import "github.com/medusar/lucas/protocol"

var pingFunc = func(args []string, r protocol.RedisRW) error {
	if len(args) != 0 && len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 1 {
		return r.WriteString(args[0])
	}
	return r.WriteString("PONG")
}

var quitFunc = func(args []string, r protocol.RedisRW) error {
	r.WriteString("OK")
	r.Close() //TODO: will cause an error in reading. should find a better way.
	return nil
}
