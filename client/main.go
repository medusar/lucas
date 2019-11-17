package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/medusar/lucas/protocol"
	"log"
	"net"
	"os"
	"strconv"
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
			outputReply(r)
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.Replace(line, "\n", "", -1)
		cs, err := ParseClientCmd(line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		err = con.WriteCommandArray(cs)
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
	outputReply(r)
}

func outputReply(rep interface{}) {
	fmt.Println(toString(rep))
}

func toString(rep interface{}) string {
	switch rep.(type) {
	case string:
		return strconv.Quote(rep.(string))
	case int:
		return strconv.Itoa(rep.(int))
	case *protocol.Resp:
		resp := rep.(*protocol.Resp)
		if resp.Nil {
			return "(nil)"
		} else {
			return toString(resp.Val)
		}
	case []interface{}:
		array := rep.([]interface{})
		if array == nil || len(array) == 0 {
			return "(empty list or set)"
		}
		var buffer bytes.Buffer
		for i, v := range array {
			buffer.WriteString(fmt.Sprintf("%d) %s\r\n", i, toString(v)))
		}
		return buffer.String()
	default:
		return fmt.Sprintf("%v", rep)
	}
}
