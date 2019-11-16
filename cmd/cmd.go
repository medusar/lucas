package cmd

import "fmt"

var (
	GetInfo = &RedisCmdInfo{Name: "get", Arity: 2, Flags: []string{"readonly", "fast"}, FirstKey: 1, LastKey: 1, Step: 1}
	SetInfo = &RedisCmdInfo{"set", -3, []string{"write", "denyoom"}, 1, 1, 1}
	//TODO
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
