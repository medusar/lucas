package main

import (
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

//WriteCommand write redis commands to server, commands are in one line, separated by space
func (r *RedisCli) WriteCommand(cmd string) error {
	fields := strings.Fields(cmd)
	return r.WriteCommandArray(fields)
}

//WriteCommandArray write redis commands to server, different parts of a single command are in an array
func (r *RedisCli) WriteCommandArray(cmd []string) error {
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
