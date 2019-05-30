package  client

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
	msg "message_delivery_system/message"
)

func TestClientImplementsClientInterface(t *testing.T) {
	var _ ClientInterface = (*Client)(nil)
}

func TestIdReturnsId(t *testing.T) {
	id := uint64(42)
	conn := new(ClientMockedConnection)

	c := NewClient(id, conn)
	assert.Equal(t, id, c.Id())
}

func TestNextMessageReturnsIdentityMessage(t *testing.T) {
	conn := new(ClientMockedConnection)

	conn.setExpectedCommand(msg.MessageTypeIdentity)

	c := NewClient(42, conn)
	message, _ := c.NextMessage()

	assert.Equal(t, msg.MessageTypeIdentity, message.Command())
}

func TestNextMessageReturnsListMessage(t *testing.T) {
	conn := new(ClientMockedConnection)

	conn.setExpectedCommand(msg.MessageTypeList)

	c := NewClient(42, conn)
	message, _ := c.NextMessage()

	assert.Equal(t, msg.MessageTypeList, message.Command())
}

func TestNextMessageReturnsRelayMessage(t *testing.T) {
	body := "foobar 1\nfoobar 2\n\nfoobar 3"

	conn := new(ClientMockedConnection)

	conn.setExpectedCommand(fmt.Sprintf("%s\n100500,42,56\n%s", msg.MessageTypeRelay, body))

	c := NewClient(42, conn)
	message, _ := c.NextMessage()

	assert := assert.New(t)

	assert.Equal(msg.MessageTypeRelay, message.Command())
	assert.Equal(body, *message.Body())

	receivers := message.Receivers()
	assert.Len(receivers, 3)
	assert.Contains(receivers, uint64(42))
	assert.Contains(receivers, uint64(56))
	assert.Contains(receivers, uint64(100500))
}

func TestNextMessageReturnsErrorOnInvalidCommand(t *testing.T) {
	conn := new(ClientMockedConnection)

	conn.setExpectedCommand("testInvalidCommand\n100500,42,56\nfoobar")

	c := NewClient(42, conn)
	message, err := c.NextMessage()

	assert := assert.New(t)

	assert.Nil(message)
	assert.True(err.InvalidCommand())
	assert.False(err.InvalidReceivers())
	assert.False(err.ConnectionError())
}

func TestNextMessageReturnsErrorOnInvalidReceivers(t *testing.T) {
	conn := new(ClientMockedConnection)

	conn.setExpectedCommand(fmt.Sprintf("%s\n100500,4foo2,56\nfoobar", msg.MessageTypeRelay))

	c := NewClient(42, conn)
	message, err := c.NextMessage()

	assert := assert.New(t)

	assert.Nil(message)
	assert.False(err.InvalidCommand())
	assert.True(err.InvalidReceivers())
	assert.False(err.ConnectionError())
}

func TestNextMessageReturnsErrorOnReadError(t *testing.T) {
	conn := new(ClientMockedConnection)

	c := NewClient(42, conn)
	message, err := c.NextMessage()

	assert := assert.New(t)

	assert.Nil(message)
	assert.False(err.InvalidCommand())
	assert.False(err.InvalidReceivers())
	assert.True(err.ConnectionError())
}

func TestSendWritesToConnection(t *testing.T) {
	messages := []string{"testMessage1", "test\nMessage2"}

	conn := new(ClientMockedConnection)
	conn.On("Write", messages[0]).Return(nil)
	conn.On("Write", messages[1]).Return(nil)

	c := NewClient(42, conn)
	c.Send(&messages[0])
	c.Send(&messages[1])
	
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		written := conn.getWritten()
		for len(written) != 2 {
			written = conn.getWritten()
		}
	}()
	
	wg.Wait()

	conn.AssertNumberOfCalls(t, "Write", 2)
	conn.AssertExpectations(t)
	assert.Len(t, conn.getWritten(), 2)
}

/*
 * Mocks
 */

type ClientMockedConnection struct {
	mock.Mock
	command string
	written []string
	lock    sync.Mutex
}

func (c *ClientMockedConnection) Write(message string) error {
	args := c.Called(message)
	c.lock.Lock()
	defer c.lock.Unlock()

	c.written = append(c.written, message)
	return args.Error(0)
}

func (c *ClientMockedConnection) Read() (string, error) {
	if c.command != "" {
		return c.command, nil
	}

	return "", errors.New("testConnectionReadError")
}

func (c *ClientMockedConnection) Close() {
	c.Called()
}

func (c *ClientMockedConnection) setExpectedCommand(command string) {
	c.command = command
}

func (c *ClientMockedConnection) getWritten() []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.written
}
