package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}

	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	fmt.Println(11111)
	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	message := "server:" + string(buf[:n])
	fmt.Println(message, addr)
	conn.WriteToUDP([]byte(message), addr)
}
