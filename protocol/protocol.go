package protocol

import (
	"fmt"
	"net"
	"strconv"
)

var (
	Delimiter = []byte("\r\n")
)

type Resp struct {
	Type byte
	Val  interface{}
}

//RedisConn represents a connection establish between client and server
type RedisConn struct {
	con net.Conn

	buf       []byte
	limit     int
	readIndex int
	closed    bool
}

func NewRedisConn(con net.Conn) *RedisConn {
	return &RedisConn{con: con, buf: make([]byte, 1024), limit: 0, readIndex: 0}
}

func (r *RedisConn) Write(data [][]byte) error {
	for i := range data {
		if err := r.WriteBytes(data[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisConn) WriteBytes(data []byte) error {
	if _, err := r.con.Write(data); err != nil {
		return err
	}
	return nil
}

func (r *RedisConn) read() error {
	if r.buf == nil {
		r.buf = make([]byte, 1024)
	}
	if r.readIndex >= r.limit {
		n, err := r.con.Read(r.buf)
		if err != nil {
			return err
		}
		r.readIndex = 0
		r.limit = n
	}
	return nil
}

func (r *RedisConn) readToCRLF() ([]byte, error) {
	ret := make([]byte, 0)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return ret, err
		}
		if b != '\r' {
			ret = append(ret, b)
		} else {
			bn, err := r.ReadByte()
			if err != nil {
				return ret, err
			}
			if bn == '\n' {
				break
			}
			ret = append(ret, b, bn)
		}
	}
	return ret, nil
}

//ReadByte read a single byte from the connection
func (r *RedisConn) ReadByte() (byte, error) {
	if err := r.read(); err != nil {
		return 0, err
	}
	b := r.buf[r.readIndex]
	r.readIndex = r.readIndex + 1
	return b, nil
}

//ReadLine read data until a '\r\n` is found, return data before `\r\n`
func (r *RedisConn) ReadLine() (string, error) {
	b, err := r.readToCRLF()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//ReadInt read a int from the underlying connection
func (r *RedisConn) ReadInt() (int, error) {
	s, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return n, nil
}

//ReadBulkString read the bulk string
func (r *RedisConn) ReadBulk() (string, error) {
	n, err := r.ReadInt()
	if err != nil {
		return "", err
	}

	if n == -1 {
		return "nil", nil
	}

	bs := make([]byte, 0)
	for i := 0; i < n; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		bs = append(bs, b)
	}
	//CRLF
	for i := 0; i < 2; i++ {
		if _, err := r.ReadByte(); err != nil {
			return "", err
		}
	}
	return string(bs), nil
}

//ReadArray read an array of data from the underlying connection
func (r *RedisConn) ReadArray() ([]interface{}, error) {
	n, err := r.ReadInt()
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, 0)
	for i := 0; i < n; i++ {
		rp, err := r.ReadReply()
		if err != nil {
			return nil, err
		}
		ret = append(ret, rp)
	}
	return ret, nil
}

//ReadReply is used to read response from server, usually after client send a request to the server
func (r *RedisConn) ReadReply() (interface{}, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case '+':
		return r.ReadLine()
	case '-':
		return r.ReadLine()
	case ':':
		return r.ReadInt()
	case '$':
		return r.ReadBulk()
	case '*':
		return r.ReadArray()
	default:
		return nil, fmt.Errorf("unknown type:%c", b)
	}
}

//ReadRequest read commands sent from a client
func (r *RedisConn) ReadRequest() ([]interface{}, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != '*' {
		return nil, fmt.Errorf("illegal request, type:%c", b)
	}

	return r.ReadArray()
}

func (r *RedisConn) WriteString(val string) error {
	b := []byte("+")
	b = append(b, []byte(val)...)
	b = append(b, Delimiter...)
	return r.WriteBytes(b)
}

func (r *RedisConn) WriteInteger(val int) error {
	b := []byte(":")
	b = append(b, []byte(strconv.Itoa(val))...)
	b = append(b, Delimiter...)
	return r.WriteBytes(b)
}

func (r *RedisConn) WriteBulk(val string) error {
	b := []byte("$")
	data := []byte(val)
	b = append(b, []byte(strconv.Itoa(len(data)))...)
	b = append(b, Delimiter...)
	b = append(b, data...)
	b = append(b, Delimiter...)
	return r.WriteBytes(b)
}

func (r *RedisConn) WriteError(val string) error {
	b := []byte("-")
	b = append(b, []byte(val)...)
	b = append(b, Delimiter...)
	return r.WriteBytes(b)
}

func (r *RedisConn) WriteArray(val []*Resp) error {
	b := []byte("*")
	b = append(b, []byte(strconv.Itoa(len(val)))...)
	b = append(b, Delimiter...)

	err := r.WriteBytes(b)
	if err != nil {
		return err
	}

	for _, v := range val {
		switch v.Type {
		case '+':
			if s, ok := v.Val.(string); ok {
				err = r.WriteString(s)
			} else {
				return fmt.Errorf("illegal simple string type")
			}
		case '-':
			if s, ok := v.Val.(string); ok {
				err = r.WriteString(s)
			} else {
				return fmt.Errorf("illegal simple error type")
			}
		case ':':
			if s, ok := v.Val.(int); ok {
				err = r.WriteInteger(s)
			} else {
				return fmt.Errorf("illegal simple integer type")
			}
		case '$':
			if s, ok := v.Val.(string); ok {
				err = r.WriteBulk(s)
			} else {
				return fmt.Errorf("illegal simple bulk string type")
			}
		case '*':
			if s, ok := v.Val.([]*Resp); ok {
				err = r.WriteArray(s)
			}
			return fmt.Errorf("illegal simple array type")
		default:
			return fmt.Errorf("unknown type:%c", b)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisConn) Close() {
	r.closed = true
	r.con.Close()
}

func (r *RedisConn) IsClosed() bool {
	return r.closed
}

func (r *RedisConn) WriteNil() error {
	b := []byte("$-1\r\n")
	return r.WriteBytes(b)
}
