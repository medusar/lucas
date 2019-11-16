package main

import (
	"github.com/medusar/lucas/cmd"
	"github.com/medusar/lucas/protocol"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	//start server
	go cmd.LoopAndInvoke()

	for {
		con, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(con)
	}
}

func serve(con net.Conn) {
	defer con.Close()
	r := protocol.NewRedisConn(con)
	for {
		req, err := r.ReadRequest()
		if err != nil {
			log.Println("Failed to read,", err)
			break
		}
		redisCmd, err := cmd.ParseRequest(req)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(redisCmd)

		err = cmd.Execute(r, redisCmd)
		if err != nil {
			log.Println("Failecd to write", err)
			break
		}
	}
}
