package server

import (
	_ "message_delivery_system/utils"
	"message_delivery_system/client"
	)	

type Server struct {
	clientFactory client.ClientFactoryInterface
	dispatcher client.DispatcherInterface
}

func NewServer(f client.ClientFactoryInterface, d client.DispatcherInterface) *Server {
	return &Server{clientFactory: f, dispatcher: d}
}

func (s *Server) Serve(connection client.ConnectionInterface) {
	client := s.createClient(connection)
	defer connection.Close()

	d := s.dispatcher
	d.Subscribe(client)
	defer d.Unsubscribe(client)

	for {
		message, err := client.NextMessage()

		switch {
		case err == nil:
			d.Dispatch(message)
		case err.ConnectionError():
			return
			// case err.InvalidMessage(): // Just continue
		}
	}
}

func (s *Server) createClient(connection client.ConnectionInterface) client.ClientInterface {
	return s.clientFactory.Create(connection)
}
