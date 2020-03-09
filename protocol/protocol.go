package protocol

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

var (
	Delimiter = []byte("\r\n")
	Nil       = []byte("$-1\r\n")
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

type RedisRW interface {
	ReadByte() (byte, error)
	ReadLine() (string, error)
	ReadInt() (int, error)
	ReadBulk() (*Resp, error)
	ReadArray() ([]interface{}, error)
	ReadReply() (interface{}, error)
	ReadRequest() ([]string, error)
	ReadInlineRequest() ([]string, error)

	WriteString(val string) error
	WriteInteger(val int) error
	WriteBulk(val string) error
	WriteError(val string) error
	WriteArray(val []*Resp) error
	WriteNil() error

	Close()
	IsClosed() bool
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
	for {
		bytes := make([]byte, 0)
		for {
			b, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if b != '\n' && b != '\r' {
				bytes = append(bytes, b)
			} else if b == '\n' {
				break
			} else {
				//\r\n
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
		if len(bytes) == 0 {
			continue
		}
		s := strings.TrimSpace(string(bytes))
		if s == "" {
			continue
		}
		return strings.Fields(s), nil
	}
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

type BufRedisConn struct {
	con    net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	closed bool
}

func (c *BufRedisConn) Write(data [][]byte) error {
	for i := range data {
		if _, err := c.writer.Write(data[i]); err != nil {
			return err
		}
	}
	return c.writer.Flush()
}

func (c *BufRedisConn) writeBytes(data ...[]byte) error {
	for i := range data {
		if _, err := c.writer.Write(data[i]); err != nil {
			return err
		}
	}
	return c.writer.Flush()
}

func (c *BufRedisConn) ReadByte() (byte, error) {
	return c.reader.ReadByte()
}

func (c *BufRedisConn) ReadLine() (string, error) {
	bytes, err := c.reader.ReadBytes('\r')
	if err != nil {
		return "", err
	}
	//_, err = c.reader.ReadByte() //read next '\n'
	//if err != nil {
	//	return "", err
	//}
	c.reader.Discard(1)
	l := len(bytes)
	return string(bytes[:l-1]), nil
}

func (c *BufRedisConn) ReadInt() (int, error) {
	bytes, err := c.reader.ReadBytes('\r')
	if err != nil {
		return -1, err
	}
	c.reader.Discard(1)
	n := 0
	for i := 0; i < len(bytes)-1; i++ {
		n = n*10 + int(bytes[i]-'0')
	}
	return n, nil
}

func (c *BufRedisConn) ReadBulk() (*Resp, error) {
	n, err := c.ReadInt()
	if err != nil {
		return nil, err
	}
	if n == -1 {
		return &Resp{Type: '$', Nil: true}, nil
	}

	buf := make([]byte, n)
	_, err = io.ReadFull(c.reader, buf)
	if err != nil {
		return nil, err
	}
	c.reader.Discard(2)
	return &Resp{Type: '$', Val: string(buf), Nil: false}, nil
}

func (c *BufRedisConn) ReadArray() ([]interface{}, error) {
	n, err := c.ReadInt()
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, n)
	for i := 0; i < n; i++ {
		rp, err := c.ReadReply()
		if err != nil {
			return nil, err
		}
		ret[i] = rp
	}
	return ret, nil
}

func (c *BufRedisConn) ReadReply() (interface{}, error) {
	b, err := c.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case '+':
		return c.ReadLine()
	case '-':
		return c.ReadLine()
	case ':':
		return c.ReadInt()
	case '$':
		return c.ReadBulk()
	case '*':
		return c.ReadArray()
	default:
		return nil, fmt.Errorf("unknown type:%c", b)
	}
}

func (c *BufRedisConn) ReadRequest() ([]string, error) {
	bs, err := c.reader.Peek(1)
	if err != nil {
		return nil, err
	}
	if bs[0] != '*' {
		return c.ReadInlineRequest()
	}

	c.reader.Discard(1) //discard bs[0]
	n, err := c.ReadInt()
	if err != nil {
		return nil, err
	}

	ret := make([]string, n)
	for i := 0; i < n; i++ {
		b, err := c.ReadByte()
		if err != nil {
			return nil, err
		}

		if b != '$' {
			return nil, fmt.Errorf("illegal request, type:%c", b)
		}

		rp, err := c.ReadBulk()
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

func (c *BufRedisConn) ReadInlineRequest() ([]string, error) {
	for {
		//support \r\n and \n in `netcat`
		bytes, err := c.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		//remove last `\n`
		bytes = bytes[:len(bytes)-1]
		//remove last `\r`
		if l := len(bytes); l > 0 && bytes[l-1] == '\r' {
			bytes = bytes[:l-1]
		}

		if len(bytes) == 0 {
			continue
		}
		str := strings.TrimSpace(string(bytes))
		//ignore empty blank
		if str == "" {
			continue
		}
		return strings.Fields(str), nil
	}
}

func (c *BufRedisConn) WriteString(val string) error {
	c.writeBytes([]byte("+"), []byte(val), Delimiter)
	return c.writer.Flush()
}

func (c *BufRedisConn) WriteInteger(val int) error {
	c.writeBytes([]byte(":"), []byte(strconv.Itoa(val)), Delimiter)
	return c.writer.Flush()
}

func (c *BufRedisConn) WriteBulk(val string) error {
	data := []byte(val)
	c.writeBytes([]byte("$"), []byte(strconv.Itoa(len(data))), Delimiter, data, Delimiter)
	return c.writer.Flush()
}

func (c *BufRedisConn) WriteError(val string) error {
	c.writeBytes([]byte("-"), []byte(val), Delimiter)
	return c.writer.Flush()
}

func (c *BufRedisConn) WriteNil() error {
	c.writeBytes(Nil)
	return c.writer.Flush()
}

func (c *BufRedisConn) WriteArray(val []*Resp) error {
	c.writeArray(val)
	return c.writer.Flush()
}

func (c *BufRedisConn) writeArray(val []*Resp) error {
	err := c.writeBytes([]byte("*"), []byte(strconv.Itoa(len(val))), Delimiter)
	if err != nil {
		return err
	}
	for _, v := range val {
		t := v.Type
		switch t {
		case '+':
			if s, ok := v.Val.(string); ok {
				err = c.writeBytes([]byte("+"), []byte(s), Delimiter)
			} else {
				return fmt.Errorf("illegal simple string type")
			}
		case '-':
			if s, ok := v.Val.(string); ok {
				err = c.writeBytes([]byte("+"), []byte(s), Delimiter)
			} else {
				return fmt.Errorf("illegal simple error type")
			}
		case ':':
			if s, ok := v.Val.(int); ok {
				err = c.writeBytes([]byte(":"), []byte(strconv.Itoa(s)), Delimiter)
			} else {
				return fmt.Errorf("illegal simple integer type")
			}
		case '$':
			if v.Nil {
				err = c.writeBytes(Nil)
			} else if s, ok := v.Val.(string); ok {
				data := []byte(s)
				err = c.writeBytes([]byte("$"), []byte(strconv.Itoa(len(data))), Delimiter, data, Delimiter)
			} else {
				return fmt.Errorf("illegal simple bulk string type")
			}
		case '*':
			if s, ok := v.Val.([]*Resp); ok {
				err = c.writeArray(s)
				if err != nil {
					return err
				}
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

func (c *BufRedisConn) Close() {
	c.closed = true
	c.con.Close()
}

func (c *BufRedisConn) IsClosed() bool {
	return c.closed
}

func NewBufRedisConn(con net.Conn) *BufRedisConn {
	return &BufRedisConn{con: con, reader: bufio.NewReader(con), writer: bufio.NewWriter(con)}
}
