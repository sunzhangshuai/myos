package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}
	
	for {
		con, _ := conn.AcceptTCP()
		handleTCPClient(con)
	}
}

func handleTCPClient(conn *net.TCPConn) {
	var buf [512]byte
	fmt.Println(11111)
	n, err := conn.Read(buf[0:])
	if err != nil {
		return
	}
	message := "server:" + string(buf[:n])
	fmt.Println(message, conn.LocalAddr(), conn.RemoteAddr())
	conn.Write([]byte(message))
}
