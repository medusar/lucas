package main

import "log"

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
}
