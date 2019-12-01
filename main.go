package main

import (
	"github.com/medusar/lucas/command"
	"github.com/medusar/lucas/protocol"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	l, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go http.ListenAndServe(":8080",http.DefaultServeMux)

	//start server
	go command.LoopAndInvoke()

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
		redisCmd, err := command.ParseRequest(req)
		if err != nil {
			log.Println(err)
			break
		}
		err = command.Execute(r, redisCmd)
		if err != nil {
			log.Println("Failecd to write", err)
			break
		}
	}
}
