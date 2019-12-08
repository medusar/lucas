package protocol

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	Delimiter = []byte("\r\n")
)

type Resp struct {
	Type byte
	Val  interface{}
	Nil  bool
}

func (resp *Resp) String() string {
	if resp.Nil {
		return "nil"
	}
	if s, ok := resp.Val.(fmt.Stringer); ok {
		return s.String()
	} else {
		return fmt.Sprintf("%v", resp.Val)
	}
}

func NewNil() *Resp {
	return &Resp{Type: '$', Nil: true}
}

func NewBulk(val string) *Resp {
	return &Resp{Type: '$', Val: val, Nil: false}
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
	if err := r.WriteBytes(data...); err != nil {
		return err
	}
	return nil
}

func (r *RedisConn) WriteBytes(data ...[]byte) error {
	for _, d := range data {
		if _, err := r.con.Write(d); err != nil {
			return err
		}
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

func (r *RedisConn) rewind(step int) error {
	newIdx := r.readIndex - step
	if newIdx < 0 {
		return fmt.Errorf("failed to rewind read index for step:%d", step)
	}
	r.readIndex = newIdx
	return nil
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
	n := 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			return n, err
		}
		if b != '\r' {
			// Can't use int(b) here, because b is ascii, if you see an ASCII table this will become clear to you.
			// Subtracting ASCII '0' (48 decimal) reduces 48 from the byte value.
			n = n*10 + int(b-'0')
			continue
		}
		_, err = r.ReadByte()
		if err != nil {
			return n, err
		}
		break
	}
	return n, nil
}

//ReadBulkString read the bulk string
func (r *RedisConn) ReadBulk() (*Resp, error) {
	n, err := r.ReadInt()
	if err != nil {
		return nil, err
	}

	if n == -1 {
		return &Resp{Type: '$', Nil: true}, nil
	}

	bs := make([]byte, n)
	for i := 0; i < n; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		bs[i] = b
	}
	//CRLF
	for i := 0; i < 2; i++ {
		if _, err := r.ReadByte(); err != nil {
			return nil, err
		}
	}

	return &Resp{Type: '$', Val: string(bs), Nil: false}, nil
}

//ReadArray read an array of data from the underlying connection
func (r *RedisConn) ReadArray() ([]interface{}, error) {
	n, err := r.ReadInt()
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rp, err := r.ReadReply()
		if err != nil {
			return nil, err
		}
		ret[i] = rp
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
func (r *RedisConn) ReadRequest() ([]string, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if b != '*' {
		err = r.rewind(1)
		if err != nil {
			return nil, err
		}
		return r.ReadInlineRequest()
	}

	n, err := r.ReadInt()
	if err != nil {
		return nil, err
	}

	ret := make([]string, n)
	for i := 0; i < n; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if b != '$' {
			return nil, fmt.Errorf("illegal request, type:%c", b)
		}

		rp, err := r.ReadBulk()
		if err != nil {
			return nil, err
		}
		//request should not contain 'nil'
		if rp.Nil == true {
			return nil, fmt.Errorf("illegal request, should not contain nil")
		}
		v, ok := rp.Val.(string)
		if !ok {
			return nil, fmt.Errorf("illegal request, bulk is not string")
		}
		ret[i] = v
	}

	return ret, nil
}

func (r *RedisConn) ReadInlineRequest() ([]string, error) {
	bytes := make([]byte, 0)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b != '\r' {
			bytes = append(bytes, b)
		} else {
			bn, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if bn == '\n' {
				break
			}
			bytes = append(bytes, b, bn)
		}
	}
	s := string(bytes)
	return strings.Split(s, " "), nil
}

func (r *RedisConn) WriteString(val string) error {
	return r.WriteBytes([]byte("+"), []byte(val), Delimiter)
}

func (r *RedisConn) WriteInteger(val int) error {
	return r.WriteBytes([]byte(":"), []byte(strconv.Itoa(val)), Delimiter)
}

func (r *RedisConn) WriteBulk(val string) error {
	data := []byte(val)
	return r.WriteBytes([]byte("$"), []byte(strconv.Itoa(len(data))), Delimiter, data, Delimiter)
}

func (r *RedisConn) WriteError(val string) error {
	return r.WriteBytes([]byte("-"), []byte(val), Delimiter)
}

func (r *RedisConn) WriteArray(val []*Resp) error {
	err := r.WriteBytes([]byte("*"), []byte(strconv.Itoa(len(val))), Delimiter)
	if err != nil {
		return err
	}

	for _, v := range val {
		t := v.Type
		switch t {
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
			if v.Nil {
				err = r.WriteNil()
			} else if s, ok := v.Val.(string); ok {
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
			return fmt.Errorf("unknown type:%c", t)
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
