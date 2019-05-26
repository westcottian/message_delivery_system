package client

import (
	"fmt"
	msg "message_delivery_system/message"
	"strings"
	"sync"
)

type DispatcherInterface interface {
	Dispatch(message msg.MessageInterface)
	Subscribe(client ClientInterface)
	Unsubscribe(client ClientInterface)
}

type Dispatcher struct {
	clients      map[uint64]ClientInterface
	clientsMutex sync.Mutex
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{clients: make(map[uint64]ClientInterface)}
}

func (d *Dispatcher) Dispatch(message msg.MessageInterface) {
	switch message.Command() {
	case msg.MessageTypeRelay:
		d.relay(message)
	case msg.MessageTypeIdentity:
		d.identify(message.Sender())
	case msg.MessageTypeList:
		d.list(message.Sender())
	}

}

func (d *Dispatcher) identify(sender uint64) {
	body := fmt.Sprintf("[Server] Your ID is %d\n", sender)
	d.sendBody(sender, &body)
}

func (d *Dispatcher) list(sender uint64) {
	var clientList []string

	d.lockClients()
	for id, _ := range d.clients {
		if id != sender {
			clientList = append(clientList, fmt.Sprintf("%d", id))
		}
	}
	d.unlockClients()

	body := fmt.Sprintf("[Server] Client IDs are %s\n", strings.Join(clientList, ", "))
	d.sendBody(sender, &body)
}

func (d *Dispatcher) relay(message msg.MessageInterface) {
	for _, id := range message.Receivers() {
		d.sendBody(id, message.Body())
	}
}

func (d *Dispatcher) Subscribe(c ClientInterface) {
	d.lockClients()
	defer d.unlockClients()

	d.clients[c.Id()] = c
}

func (d *Dispatcher) Unsubscribe(c ClientInterface) {
	d.lockClients()
	defer d.unlockClients()

	delete(d.clients, c.Id())
}

func (d *Dispatcher) sendBody(receiver uint64, body *string) {
	client := d.client(receiver)

	if client != nil {
		client.Send(body)
	}
}

func (d *Dispatcher) client(id uint64) ClientInterface {
	d.lockClients()
	defer d.unlockClients()

	return d.clients[id]
}

func (d *Dispatcher) lockClients() {
	d.clientsMutex.Lock()
}

func (d *Dispatcher) unlockClients() {
	d.clientsMutex.Unlock()
}
