package server

import (
	"message_delivery_system/utils"
	"message_delivery_system/client"
	)	

type Server struct {
	idSequence utils.IdSequenceInterface
	dispatcher client.DispatcherInterface
}

func NewServer(s utils.IdSequenceInterface, d client.DispatcherInterface) *Server {
	return &Server{idSequence: s, dispatcher: d}
}

func (s *Server) Serve(connection client.ConnectionInterface) {
	client := s.createClient(connection)
	defer connection.Close()

	d := s.dispatcher
	d.Subscribe(client)

	for {
		message, err := client.NextMessage()

		switch {
		case err == nil:
			d.Dispatch(message)
		case err.ConnectionError():
			d.Unsubscribe(client)
			return
			// case err.InvalidMessage(): // Just continue
		}
	}
}

func (s *Server) getNextId() int64 {
	return s.idSequence.NextId()
}

func (s *Server) createClient(connection client.ConnectionInterface) client.ClientInterface {
	return client.NewClient(s.getNextId(), connection)
}
