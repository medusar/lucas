package main

import (
	"fmt"
	"log"
	"math"
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
		log.Println(e)
	}
	log.Println("after unquote:", s)

	fmt.Println(len("哈哈大家好"))

	fmt.Println(math.MaxFloat64)

	v()
	v(1)

	m := make(map[string]int)
	m["1"] = 1
	m["2"] = 2
	for k := range m {
		delete(m, k)
	}
	log.Println(m)
}

func v(i ...int) {
	log.Println(i == nil)
	log.Println(len(i))
}
