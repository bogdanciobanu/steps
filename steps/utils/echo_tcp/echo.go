package main

import (
	"net"
	"os"
)

func main() {
	l ,err := net.Listen("tcp", "0.0.0.0:2020")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	conn,err := l.Accept()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	conn.Close()
	os.Stdout.Write(buf[:n])
}
