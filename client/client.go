package main

import (
	"fmt"
	"github.com/medusar/lucas/protocol"
	"net"
	"strconv"
	"strings"
)

type RedisCli struct {
	*protocol.RedisConn
}

func NewRedisCli(con net.Conn) *RedisCli {
	return &RedisCli{RedisConn: protocol.NewRedisConn(con)}
}

//WriteCommandArray write redis commands to server, different parts of a single command are in an array
func (r *RedisCli) WriteCommandArray(cmd []string) error {
	if cmd == nil || len(cmd) == 0 {
		return nil
	}

	data := make([][]byte, 0)
	data = append(data, []byte("*"))
	data = append(data, []byte(strconv.Itoa(len(cmd))))
	data = append(data, []byte("\r\n"))

	for _, f := range cmd {
		data = append(data, []byte("$"))
		data = append(data, []byte(strconv.Itoa(len(f))))
		data = append(data, []byte("\r\n"))
		data = append(data, []byte(f))
		data = append(data, []byte("\r\n"))
	}
	return r.Write(data)
}

func ParseClientCmd(c string) ([]string, error) {
	fields := strings.Fields(c)
	for i := range fields {
		if f, err := unquote(fields[i]); err == nil {
			fields[i] = f
		} else {
			return nil, fmt.Errorf("invalid argument(s)")
		}
	}
	return fields, nil
}

//https://github.com/antirez/redis/blob/0f026af185e918a9773148f6ceaa1b084662be88/src/sds.c#L959
func unquote(s string) (string, error) {
	//TODO:
	//rs := []rune(s)
	//inQuote := false
	//inSingleQuote := false
	//for i := 0; i < len(rs); i++ {
	//
	//}
	return s, nil
}
