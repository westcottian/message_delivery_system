package main

import (
	"fmt"
	"net"
	"os"
	"message_delivery_system/client"
//	msg "message_delivery_system/message"
	"message_delivery_system/utils"
	"message_delivery_system/server"
	)

const serverPort = 8888

func main() {
	listener := createListener()
	defer listener.Close()

	sequence := utils.NewIdSequence()
	dispatcher := client.NewDispather()
	server := server.NewServer(sequence, dispatcher)
	acceptConnections(server, listener)
}

func acceptConnections(server *server.Server, listener *net.TCPListener) {
	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Can't accept connection: %s", err.Error())
		} else {
			conn := client.NewConnection(connection)
			go server.Serve(conn)
		}
	}
}

func createListener() *net.TCPListener {
	ip := net.IPv4(127, 0, 0, 1)
	addr := &net.TCPAddr{Port: serverPort, IP: ip}
	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		fmt.Printf("Can't start listening: %s", err.Error())
		os.Exit(1)
	}

	return listener
}
