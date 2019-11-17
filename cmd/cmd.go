package cmd

import (
	"bytes"
	"fmt"
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"log"
	"sort"
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

type cmdFunc func(args []string, r *protocol.RedisConn) error

var (
	cmdFuncMap = make(map[string]cmdFunc)
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

var ttlFunc = func(args []string, r *protocol.RedisConn) error {
	var err error
	if args == nil || len(args) != 1 {
		err = r.WriteError("ERR wrong number of arguments for 'ttl' command")
		return err
	}
	ttl := store.Ttl(args[0])
	return r.WriteInteger(ttl)
}

var expireFunc = func(args []string, r *protocol.RedisConn) error {
	if args == nil || len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'expire' command")
	}
	sec, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}
	set := store.Expire(args[0], sec)
	if set {
		return r.WriteInteger(1)
	} else {
		return r.WriteInteger(0)
	}
}

var expireAtFunc = func(args []string, r *protocol.RedisConn) error {
	if args == nil || len(args) != 2 {
		return r.WriteError("ERR wrong number of arguments for 'expire' command")
	}

	timestamp, err := strconv.Atoi(args[1])
	if err != nil {
		return r.WriteError("ERR value is not an integer or out of range")
	}

	set := store.ExpireAt(args[0], int64(timestamp))
	if set {
		return r.WriteInteger(1)
	} else {
		return r.WriteInteger(0)
	}
}

var commandFunc = func(args []string, r *protocol.RedisConn) error {
	return r.WriteString("OK") //TODO
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

var keysFunc = func(args []string, r *protocol.RedisConn) error {
	if args == nil || len(args) != 1 {
		return r.WriteError("ERR wrong number of arguments for 'keys' command")
	}
	keys := store.Keys(args[0])
	if keys == nil || len(keys) == 0 {
		return r.WriteArray(nil)
	}
	resp := make([]*protocol.Resp, 0)
	for _, k := range keys {
		resp = append(resp, protocol.NewBulk(k))
	}

	return r.WriteArray(resp)
}

var existsFunc = func(args []string, r *protocol.RedisConn) error {
	if args == nil || len(args) == 0 {
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

var delFunc = func(args []string, r *protocol.RedisConn) error {
	if args == nil || len(args) == 0 {
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

func init() {
	cmdFuncMap["get"] = getFunc
	cmdFuncMap["set"] = setFunc
	cmdFuncMap["command"] = commandFunc
	cmdFuncMap["ttl"] = ttlFunc
	cmdFuncMap["getset"] = getsetFunc
	cmdFuncMap["setex"] = setexFunc
	cmdFuncMap["setnx"] = setnxFunc
	cmdFuncMap["meget"] = mgetFunc
	cmdFuncMap["mset"] = msetFunc
	cmdFuncMap["strlen"] = strlenFunc
	cmdFuncMap["incr"] = incrFunc
	cmdFuncMap["incrby"] = incrByFunc
	cmdFuncMap["decr"] = decrFun
	cmdFuncMap["decrby"] = decrByFunc
	cmdFuncMap["expire"] = expireFunc
	cmdFuncMap["keys"] = keysFunc
	cmdFuncMap["exists"] = existsFunc
	cmdFuncMap["del"] = delFunc

	keys := make([]string, 0)
	for key, _ := range cmdFuncMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	log.Printf("supported commands: %s \r\n", strings.Join(keys, ", "))
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

func ParseRequest(reqs []string) (*RedisCmd, error) {
	l := len(reqs)
	if l == 0 {
		return nil, fmt.Errorf("illegal command")
	}

	name := reqs[0]
	if name == "" {
		return nil, fmt.Errorf("illegal command, name is empty")
	}

	if l > 1 {
		args := make([]string, 0)
		for i := 1; i < l; i++ {
			args = append(args, reqs[i])
		}
		return &RedisCmd{Name: name, Args: args}, nil
	}

	return &RedisCmd{Name: name}, nil
}

func execCmd(r *protocol.RedisConn, c *RedisCmd) error {
	name := strings.ToLower(c.Name)
	f, ok := cmdFuncMap[name]
	if !ok {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("ERR unknown command `%s`, with args beginning with:", name))
		args := c.Args
		if args != nil {
			for i := 0; i < len(args); i++ {
				args[i] = fmt.Sprintf("`%s`", args[i])
			}
			buf.WriteString(strings.Join(args, ", "))
		}
		return r.WriteError(buf.String())
	}
	return f(c.Args, r)
}
