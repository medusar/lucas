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
	log.Println("server stared on address:", l.Addr())

	go http.ListenAndServe(":8080", http.DefaultServeMux)

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
	//r := protocol.NewRedisConn(con)
	r := protocol.NewBufRedisConn(con)
	defer r.Close()
	for {
		req, err := r.ReadRequest()
		log.Println("req:", req)
		if err != nil {
			log.Println("Failed to read,", err)
			break
		}
		redisCmd, err := command.ParseRequest(req)
		if err != nil {
			r.WriteError(err.Error())
			continue
		}
		err = command.Execute(r, redisCmd)
		if err != nil {
			log.Println("Failecd to write", err)
			break
		}
	}
}
