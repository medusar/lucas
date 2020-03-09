package command

import (
	"bytes"
	"fmt"
	"github.com/medusar/lucas/protocol"
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

type RedisCmd struct {
	Name string
	Args []string
}

type invoker struct {
	rc  *RedisCmd
	con protocol.RedisRW
}

type cmdFunc func(args []string, r protocol.RedisRW) error

var (
	cmdFuncMap = make(map[string]cmdFunc)
)

var commandFunc = func(args []string, r protocol.RedisRW) error {
	return r.WriteString("OK") //TODO: implement COMMAND
}

func init() {
	cmdFuncMap["command"] = WithTime(commandFunc)

	//connection
	cmdFuncMap["ping"] = pingFunc
	//quit for telnet
	cmdFuncMap["quit"] = quitFunc

	//keys
	cmdFuncMap["ttl"] = WithTime(ttlFunc)
	cmdFuncMap["expire"] = WithTime(expireFunc)
	cmdFuncMap["keys"] = WithTime(keysFunc)
	cmdFuncMap["exists"] = WithTime(existsFunc)
	cmdFuncMap["del"] = WithTime(delFunc)
	cmdFuncMap["type"] = WithTime(typeFunc)

	//string
	cmdFuncMap["get"] = WithTime(getFunc)
	cmdFuncMap["set"] = WithTime(setFunc)
	cmdFuncMap["getset"] = WithTime(getsetFunc)
	cmdFuncMap["setex"] = WithTime(setexFunc)
	cmdFuncMap["setnx"] = WithTime(setnxFunc)
	cmdFuncMap["meget"] = WithTime(mgetFunc)
	cmdFuncMap["mset"] = WithTime(msetFunc)
	cmdFuncMap["strlen"] = WithTime(strlenFunc)
	cmdFuncMap["incr"] = WithTime(incrFunc)
	cmdFuncMap["incrby"] = WithTime(incrByFunc)
	cmdFuncMap["decr"] = WithTime(decrFunc)
	cmdFuncMap["decrby"] = WithTime(decrByFunc)
	cmdFuncMap["append"] = WithTime(appendFunc)
	cmdFuncMap["setrange"] = WithTime(setRangeFunc)
	cmdFuncMap["getrange"] = WithTime(getRangeFunc)
	cmdFuncMap["setbit"] = WithTime(setbitFunc)
	cmdFuncMap["getbit"] = WithTime(getbitFunc)
	cmdFuncMap["bitcount"] = WithTime(bitcountFunc)

	//hash
	cmdFuncMap["hset"] = WithTime(hsetFunc)
	cmdFuncMap["hget"] = WithTime(hgetFunc)
	cmdFuncMap["hgetall"] = WithTime(hgetAllFunc)
	cmdFuncMap["hkeys"] = WithTime(hkeysFunc)
	cmdFuncMap["hlen"] = WithTime(hlenFunc)
	cmdFuncMap["hexists"] = WithTime(hexistsFunc)
	cmdFuncMap["hdel"] = WithTime(hdelFunc)
	cmdFuncMap["hmget"] = WithTime(hmgetFunc)
	cmdFuncMap["hmset"] = WithTime(hmsetFunc)
	cmdFuncMap["hsetnx"] = WithTime(hsetnxFunc)
	cmdFuncMap["hstrlen"] = WithTime(hstrlenFunc)
	cmdFuncMap["hvals"] = WithTime(hvalsFunc)
	cmdFuncMap["hincrby"] = WithTime(hincrByFunc)
	cmdFuncMap["hincrbyfloat"] = hincrByFloatFunc
	//HSCAN
	//cmdFuncMap["hscan"] = WithTime(hcanFunc)

	//set
	cmdFuncMap["sadd"] = WithTime(saddFunc)
	cmdFuncMap["scard"] = WithTime(scardFunc)
	cmdFuncMap["sdiff"] = sdiffFunc
	cmdFuncMap["sdiffstore"] = sdiffStoreFunc
	cmdFuncMap["sinter"] = WithTime(sinterFunc)
	cmdFuncMap["sinterstore"] = WithTime(sinterStoreFunc)
	cmdFuncMap["sismember"] = WithTime(sismemberFunc)
	cmdFuncMap["smembers"] = WithTime(smembersFunc)
	cmdFuncMap["smove"] = WithTime(smoveFunc)
	cmdFuncMap["spop"] = WithTime(spopFunc)
	cmdFuncMap["srem"] = WithTime(sremFunc)
	cmdFuncMap["sunion"] = WithTime(sunionFunc)
	cmdFuncMap["sunionstore"] = WithTime(sunionStoreFunc)

	//list
	cmdFuncMap["lpush"] = WithTime(lpushFunc)
	cmdFuncMap["rpush"] = WithTime(rpushFunc)
	cmdFuncMap["llen"] = WithTime(llenFunc)
	cmdFuncMap["lpop"] = WithTime(lpopFunc)
	cmdFuncMap["rpop"] = WithTime(rpopFunc)
	cmdFuncMap["lindex"] = WithTime(lindexFunc)
	cmdFuncMap["lrem"] = WithTime(lremFunc)
	cmdFuncMap["lset"] = WithTime(lsetFunc)
	cmdFuncMap["rpushx"] = WithTime(rpushXFunc)
	cmdFuncMap["lpushx"] = WithTime(lpushXFunc)
	cmdFuncMap["lrange"] = WithTime(lrangeFunc)

	//zset
	cmdFuncMap["zadd"] = WithTime(zaddFunc)
	cmdFuncMap["zcard"] = WithTime(zcardFunc)
	cmdFuncMap["zcount"] = WithTime(zcountFunc)
	cmdFuncMap["zrange"] = WithTime(zrangeFunc)
	cmdFuncMap["zrangebyscore"] = WithTime(zrangeByScoreFunc)
	cmdFuncMap["zrank"] = WithTime(zrankFunc)
	cmdFuncMap["zrem"] = WithTime(zremFunc)
	cmdFuncMap["zscore"] = WithTime(zscoreFunc)
	cmdFuncMap["zrevrank"] = WithTime(zrevrankFunc)
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

func Execute(r protocol.RedisRW, c *RedisCmd) error {
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

func execCmd(r protocol.RedisRW, c *RedisCmd) error {
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
