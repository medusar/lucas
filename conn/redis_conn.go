package conn

import (
	"errors"
	"github.com/medusar/lucas/protocol"
	"log"
	"net"
	"sync"
)

type RedisConn struct {
	conn      net.Conn
	rw        protocol.RedisRW
	rcvChan   chan []string
	sndChan   chan []byte
	closeChan chan struct{}
	closed    bool
	sync.Once
}

func InitRedisConn(conn net.Conn) (*RedisConn, error) {
	redisConn := &RedisConn{
		conn:      conn,
		rw:        protocol.NewBufRedisConn(conn),
		rcvChan:   make(chan []string, 4096),
		sndChan:   make(chan []byte, 4096),
		closeChan: make(chan struct{}, 0),
		closed:    false,
	}
	go redisConn.writeLoop()
	go redisConn.readLoop()
	return redisConn, nil
}

func (r *RedisConn) ReadMsg() ([]string, error) {
	select {
	case msg := <-r.rcvChan:
		return msg, nil
	case <-r.closeChan:
		return nil, errors.New("connection is closed")
	}
}

func (r *RedisConn) WriteMsg(msg []byte) error {
	select {
	case r.sndChan <- msg:
		return nil
	case <-r.closeChan:
		return errors.New("connection is closed")
	}
}

func (r *RedisConn) Close() {
	r.Once.Do(func() {
		r.closed = true
		close(r.closeChan)
		r.conn.Close()
	})
}

func (r *RedisConn) readRequest() ([]string, error) {
	return r.rw.ReadRequest()
}

func (r *RedisConn) readLoop() {
	for {
		request, err := r.readRequest()
		if err != nil {
			log.Println("failed to read data", err)
			r.Close()
			return
		}
		select {
		case r.rcvChan <- request:
		case <-r.closeChan:
			return
		}
	}
}

func (r *RedisConn) writeLoop() {
	for {
		var msg []byte
		select {
		case msg = <-r.sndChan:
		case <-r.closeChan:
			return
		}

		if _, err := r.conn.Write(msg); err != nil {
			log.Println("failed to write data", err)
			r.Close()
			return
		}
	}
}
