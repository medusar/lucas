package command

import "github.com/medusar/lucas/protocol"

var pingFunc = func(args []string, r *protocol.RedisConn) error {
	if len(args) != 0 && len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 1 {
		return r.WriteString(args[0])
	}
	return r.WriteString("PONG")
}
