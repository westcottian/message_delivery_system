package client

import (
	"fmt"
	msg "message_delivery_system/message"
)

type Dispatcher struct {
	name    string
	clients map[int64]*Client
}

func NewDispather(foo string) *Dispatcher {
	return &Dispatcher{name: foo, clients: make(map[int64]*Client)}
}

func (d *Dispatcher) Dispatch(message *msg.Message) {
	//d.output <- "Dispatching"

	//d.output <- fmt.Sprintf("message is '%s'", message)
	fmt.Printf("Dispatching message '%s'\n", message.Body())
	for _, id := range message.Receivers() {
		if receiver := d.clients[id]; receiver != nil {
			fmt.Printf("\tto client %d\n", receiver.id)
			receiver.Send(message.Body())
		}
	}
}

func (d *Dispatcher) Subscribe(c *Client) {
	d.clients[c.id] = c
}
