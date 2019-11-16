package cmd

import (
	"fmt"
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"strconv"
	"strings"
	"time"
)

var (
	GetInfo = &RedisCmdInfo{Name: "get", Arity: 2, Flags: []string{"readonly", "fast"}, FirstKey: 1, LastKey: 1, Step: 1}
	SetInfo = &RedisCmdInfo{"set", -3, []string{"write", "denyoom"}, 1, 1, 1}
	//TODO
	invokerChan = make(chan *invoker, 1024*1024)
)

//https://redis.io/commands/command
type RedisCmdInfo struct {
	Name     string
	Arity    int
	Flags    []string
	FirstKey int
	LastKey  int
	Step     int
}

//
//func (r RedisCmdInfo) Encode() *protocol.RespArray {
//	resp := &protocol.RespArray{
//		Data: []protocol.RespData{
//			&protocol.RespSimpleString{Data: r.Name},
//			&protocol.RespInteger{Data: r.Arity},
//			&protocol.RespArray{
//				Data: encodeFlags(r.Flags),
//			},
//			&protocol.RespInteger{Data: r.FirstKey},
//			&protocol.RespInteger{Data: r.LastKey},
//			&protocol.RespInteger{Data: r.Step},
//		},
//	}
//	return resp
//}
//
//func encodeFlags(flags []string) []protocol.RespData {
//	var resp []protocol.RespData
//	for _, f := range flags {
//		resp = append(resp, &protocol.RespSimpleString{Data: f})
//	}
//	return resp
//}

type RedisCmd struct {
	Name string
	Args []string
}

type invoker struct {
	rc  *RedisCmd
	con *protocol.RedisConn
}

func LoopAndInvoke() {
	for in := range invokerChan {
		if in.con.IsClosed() {
			continue
		}
		if err := execCmd(in.con, in.rc); err != nil {
			in.con.Close()
		}
	}
}

func Execute(r *protocol.RedisConn, c *RedisCmd) error {
	select {
	case invokerChan <- &invoker{rc: c, con: r}:
		return nil
	case <-time.After(time.Millisecond * 100):
		return fmt.Errorf("server too busy")
	}
}

func ParseRequest(reqs []interface{}) (*RedisCmd, error) {
	l := len(reqs)
	if l == 0 {
		return nil, fmt.Errorf("illegal command")
	}

	name := ""
	if n, ok := reqs[0].(string); ok {
		name = n
	} else {
		return nil, fmt.Errorf("illegal command, not string")
	}

	if name == "" {
		return nil, fmt.Errorf("illegal command, name is empty")
	}

	if l > 1 {
		args := make([]string, 0)
		for i := 1; i < l; i++ {
			if arg, ok := reqs[i].(string); ok {
				args = append(args, arg)
			} else {
				return nil, fmt.Errorf("illegal command, param:%d is not string", i)
			}
		}

		return &RedisCmd{Name: name, Args: args}, nil
	}

	return &RedisCmd{Name: name}, nil
}

func execCmd(r *protocol.RedisConn, c *RedisCmd) error {
	name := strings.ToLower(c.Name)
	var err error
	switch name {
	case "ttl":
		args := c.Args
		if args == nil || len(args) != 1 {
			err = r.WriteError("ERR wrong number of arguments for 'ttl' command")
			break
		}
		ttl := store.Ttl(args[0])
		err = r.WriteInteger(ttl)
	case "get":
		args := c.Args
		if args == nil || len(args) != 1 {
			err = r.WriteError("ERR wrong number of arguments for 'get' command")
			break
		}
		val, ok, e := store.Get(args[0])
		if e != nil {
			err = r.WriteError(e.Error())
		} else if ok {
			err = r.WriteBulk(val)
		} else {
			err = r.WriteNil()
		}
	case "set":
		args := c.Args
		if args == nil || len(args) != 2 {
			err = r.WriteError("ERR wrong number of arguments for 'set' command")
			break
		}

		store.Set(args[0], args[1])
		err = r.WriteString("OK") //TODO:support NX, EX
	case "setex":
		args := c.Args
		if args == nil || len(args) != 3 {
			err = r.WriteError("ERR wrong number of arguments for 'setex' command")
			break
		}
		key, sec, val := args[0], args[1], args[2]

		ttl, err := strconv.Atoi(sec)
		if err != nil {
			err = r.WriteError("ERR value is not an integer or out of range")
			break
		}
		err = store.SetEX(key, val, ttl)
		if err != nil {
			err = r.WriteError(err.Error())
		} else {
			err = r.WriteString("OK")
		}
	case "setnx":
		args := c.Args
		if args == nil || len(args) != 2 {
			err = r.WriteError("ERR wrong number of arguments for 'setnx' command")
			break
		}
		key, val := args[0], args[1]
		if set := store.SetNX(key, val); set {
			err = r.WriteInteger(1)
		} else {
			err = r.WriteInteger(0)
		}
	case "setrange":
		//TODO:
	case "strlen":
		args := c.Args
		if args == nil || len(args) != 1 {
			err = r.WriteError("ERR wrong number of arguments for 'strlen' command")
			break
		}
		n, e := store.StrLen(args[0])
		if e != nil {
			err = r.WriteError(e.Error())
		} else {
			err = r.WriteInteger(n)
		}
	case "incr":
		args := c.Args
		if args == nil || len(args) != 1 {
			err = r.WriteError("ERR wrong number of arguments for 'incr' command")
			break
		}
		v, e := store.Incr(args[0])
		if e != nil {
			err = r.WriteError(e.Error())
		} else {
			err = r.WriteInteger(v)
		}
		//TODO:
	case "command":
		//TODO
		err = r.WriteString("OK")
	}
	return err
}
