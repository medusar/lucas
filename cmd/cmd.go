package cmd

import (
	"bytes"
	"fmt"
	"github.com/medusar/lucas/protocol"
	"log"
	"sort"
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

var commandFunc = func(args []string, r *protocol.RedisConn) error {
	return r.WriteString("OK") //TODO
}

func init() {
	cmdFuncMap["command"] = commandFunc

	//keys
	cmdFuncMap["ttl"] = ttlFunc
	cmdFuncMap["expire"] = expireFunc
	cmdFuncMap["keys"] = keysFunc
	cmdFuncMap["exists"] = existsFunc
	cmdFuncMap["del"] = delFunc
	cmdFuncMap["type"] = typeFunc

	//string
	cmdFuncMap["get"] = getFunc
	cmdFuncMap["set"] = setFunc
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
	cmdFuncMap["append"] = appendFunc
	cmdFuncMap["setrange"] = setRangeFunc
	cmdFuncMap["getrange"] = getRangeFunc

	//hash
	cmdFuncMap["hset"] = hsetFunc
	cmdFuncMap["hget"] = hgetFunc
	cmdFuncMap["hgetall"] = hgetAllFunc
	cmdFuncMap["hkeys"] = hkeysFunc
	cmdFuncMap["hlen"] = hlenFunc
	cmdFuncMap["hexists"] = hexistsFunc
	cmdFuncMap["hdel"] = hdelFunc
	cmdFuncMap["hmget"] = hmgetFunc
	cmdFuncMap["hmset"] = hmsetFunc
	cmdFuncMap["hsetnx"] = hsetnxFunc
	cmdFuncMap["hstrlen"] = hstrlenFunc
	cmdFuncMap["hvals"] = hvalsFunc
	cmdFuncMap["hincrby"] = hincrByFunc
	cmdFuncMap["hincrbyfloat"] = hincrByFloatFunc
	//HSCAN
	//cmdFuncMap["hscan"] = hcanFunc

	//set
	cmdFuncMap["sadd"] = saddFunc
	cmdFuncMap["scard"] = scardFunc
	cmdFuncMap["sdiff"] = sdiffFunc
	cmdFuncMap["sdiffstore"] = sdiffStoreFunc
	cmdFuncMap["sinter"] = sinterFunc
	cmdFuncMap["sinterstore"] = sinterStoreFunc
	cmdFuncMap["sismember"] = sismemberFunc
	cmdFuncMap["smembers"] = smembersFunc
	cmdFuncMap["smove"] = smoveFunc
	cmdFuncMap["spop"] = spopFunc
	cmdFuncMap["srem"] = sremFunc
	cmdFuncMap["sunion"] = sunionFunc
	cmdFuncMap["sunionstore"] = sunionStoreFunc

	//list
	cmdFuncMap["lpush"] = lpushFunc
	cmdFuncMap["rpush"] = rpushFunc
	cmdFuncMap["llen"] = llenFunc
	cmdFuncMap["lpop"] = lpopFunc
	cmdFuncMap["rpop"] = rpopFunc
	cmdFuncMap["lindex"] = lindexFunc
	cmdFuncMap["lrem"] = lremFunc
	cmdFuncMap["lset"] = lsetFunc
	cmdFuncMap["rpushx"] = rpushXFunc
	cmdFuncMap["lpushx"] = lpushXFunc
	cmdFuncMap["lrange"] = lrangeFunc

	keys := make([]string, 0)
	for key := range cmdFuncMap {
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

func toBulkArray(val []string) []*protocol.Resp {
	r := make([]*protocol.Resp, len(val))
	for i := 0; i < len(val); i++ {
		r[i] = protocol.NewBulk(val[i])
	}
	return r
}
