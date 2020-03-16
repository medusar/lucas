package main

import (
	"github.com/medusar/lucas/command"
	"github.com/medusar/lucas/protocol"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	initEnv()

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
	r := protocol.NewBufRedisConn(con)
	defer r.Close()
	for {
		req, err := r.ReadRequest()
		if err != nil {
			break
		}
		redisCmd, err := command.ParseRequest(req)
		if err != nil {
			r.WriteError(err.Error())
			continue
		}
		err = command.Execute(r, redisCmd)
		if err != nil {
			break
		}
	}
}
