package client

import (
	msg "message_delivery_system/message"
	"message_delivery_system/utils"
	"strconv"
	"strings"
)

type ClientInterface interface {
	Id() uint64
	Send(message *string)
	NextMessage() (msg.MessageInterface, *ClientError)
}

type ClientFactoryInterface interface {
	Create(connection ConnectionInterface) ClientInterface
}

type ClientFactory struct {
	sequence utils.IdSequenceInterface
}

func NewClientFactory(s utils.IdSequenceInterface) *ClientFactory {
	return &ClientFactory{sequence: s}
}

func (f *ClientFactory) Create(connection ConnectionInterface) ClientInterface {
	return NewClient(f.sequence.NextId(), connection)
}

type Client struct {
	id         uint64
	outbox     chan *string
	connection ConnectionInterface
}

func NewClient(id uint64, connection ConnectionInterface) *Client {
	return &Client{id: id, connection: connection, outbox: make(chan *string, 255)}
}

func (c *Client) Id() uint64 {
	return c.id
}

func (c *Client) Send(message *string) {
	c.outbox <- message // To send them in order
	go func() {
		message := <-c.outbox
		c.connection.Write(*message)
	}()
}

func (c *Client) NextMessage() (msg.MessageInterface, *ClientError) {
	command, err := c.readCommand()
	if err != nil {
		return nil, err
	}

	cmd := strings.SplitN(command, "\n", 3)

	switch cmd[0] {
	case msg.MessageTypeIdentity:
		return msg.NewIdentityMessage(c.id), nil
	case msg.MessageTypeList:
		return msg.NewListMessage(c.id), nil
	case msg.MessageTypeRelay:
		if len(cmd) != 3 {
			return nil, NewClientInvalidCommandError()
		}
		return c.buildRelayMessage(cmd[1], &cmd[2])
	}

	return nil, NewClientInvalidCommandError()
}

func (c *Client) buildRelayMessage(receivers string, body *string) (msg.MessageInterface, *ClientError) {
	receiverIds, err := c.parseReceivers(receivers)
	if err != nil {
		return nil, err
	}

	return msg.NewRelayMessage(c.id, receiverIds, body), nil
}

func (c *Client) readCommand() (string, *ClientError) {
	line, err := c.connection.Read()

	if err != nil {
		return "", NewClientConnectionError()
	}

	return line, nil
}

func (c *Client) parseReceivers(line string) ([]uint64, *ClientError) {
	var receivers []uint64
	for _, word := range strings.Split(line, ",") {
		word = strings.TrimSpace(word)
		id, err := strconv.ParseUint(word, 10, 64)
		if err != nil {
			return make([]uint64, 0), NewClientInvalidReceiversError()
		}

		receivers = append(receivers, id)
	}
	return receivers, nil
}
