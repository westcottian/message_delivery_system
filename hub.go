package main

import (
	"fmt"
	"net"
	"os"
	"flag"
	"message_delivery_system/client"
	"message_delivery_system/utils"
	"message_delivery_system/server"
	"github.com/pkg/profile"
	)

func main() {
	// Memory profiling
	defer profile.Start(profile.MemProfile).Stop()
	listener := createListener()
	defer listener.Close()

	sequence := utils.NewIdSequence()
    	clientFactory := client.NewClientFactory(sequence)
	dispatcher := client.NewDispatcher()
	server := server.NewServer(clientFactory, dispatcher)
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
	addr := &net.TCPAddr{Port: getPort(), IP: ip}
	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		fmt.Printf("Can't start listening: %s", err.Error())
		os.Exit(1)
	}

	return listener
}

func getPort() int {
	port := flag.Int("port", 8888, "Port to listen")
	flag.Parse()

	return *port
}
