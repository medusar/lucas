package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	host := flag.String("h", "localhost", "-h")
	port := flag.Int("p", 6379, "-p")

	flag.Parse()

	address := fmt.Sprintf("%s:%d", *host, *port)

	if flag.NArg() > 0 {
		args := flag.Args()
		exeCmd(address, args)
		return
	}

	interactiveExec(address)
}

func interactiveExec(address string) {
	log.Println("connecting to", address)
	c, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")
	con := NewRedisCli(c)
	go func() {
		for {
			r, err := con.ReadReply()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(r)
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.Replace(line, "\n", "", -1)
		err = con.WriteCommand(line)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func exeCmd(address string, cmd []string) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	con := NewRedisCli(c)
	if err := con.WriteCommandArray(cmd); err != nil {
		log.Fatal(err)
	}

	r, err := con.ReadReply()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}
