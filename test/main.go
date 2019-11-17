package main

import (
	"log"
	"reflect"
	"strconv"
)

func main() {
	buf := make([]byte, 0)
	buf = append(buf, []byte("hello")...)
	var data []byte
	data = buf[:]
	log.Println(string(data))
	buf[0] = byte('x')
	log.Println(string(data))
	// log.Println(string(buf))

	data = append(data, []byte(" boy")...)
	log.Println(string(data))
	log.Println(string(buf))

	val := (*string)(nil)
	log.Println(reflect.ValueOf(val).Kind() == reflect.Ptr)

	s, e := strconv.Unquote("*")
	if e != nil {
		log.Fatal(e)
	}
	log.Println("after unquote:", s)
}
