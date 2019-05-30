package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"message_delivery_system/client"
	msg "message_delivery_system/message"
)

func TestServeSubscribesClientAndDispatchesMessages(t *testing.T) {
	connection := new(ServerMockedConnection)

	message1 := msg.NewIdentityMessage(42)
	message2 := msg.NewRelayMessage(456, []uint64{100500}, new(string))
	message3 := msg.NewListMessage(123)
	message4 := msg.NewIdentityMessage(789)
	message5 := msg.NewRelayMessage(7662, []uint64{234234, 4234234}, new(string))

	messages := []msg.MessageInterface{message1, message2, message3, message4, message5}
	client := newServableMockedClient(messages, 100500)

	factory := new(MockedClientFactory)
	factory.On("Create", connection).Return(client)

	dispatcher := new(MockedDispatcher)

	server := NewServer(factory, dispatcher)
	server.Serve(connection)

	assert := assert.New(t)

	assert.Len(dispatcher.subscribed, 1)
	assert.Equal(client, dispatcher.subscribed[0])

	assert.Len(dispatcher.unsubscribed, 1)
	assert.Equal(client, dispatcher.unsubscribed[0])

	assert.Len(dispatcher.dispatched, 5)
	assert.Equal(message1, dispatcher.dispatched[0])
	assert.Equal(message2, dispatcher.dispatched[1])
	assert.Equal(message3, dispatcher.dispatched[2])
	assert.Equal(message4, dispatcher.dispatched[3])
	assert.Equal(message5, dispatcher.dispatched[4])
}

func TestServeSkipsInvalidMessageErrors(t *testing.T) {
	connection := new(ServerMockedConnection)

	message1 := msg.NewIdentityMessage(42)
	message2 := msg.NewListMessage(123)

	messages := []msg.MessageInterface{message1, message2}

	client := newServableMockedClient(messages, 1) // Let's have some invalid command between messages

	factory := new(MockedClientFactory)
	factory.On("Create", connection).Return(client)

	dispatcher := new(MockedDispatcher)

	server := NewServer(factory, dispatcher)
	server.Serve(connection)

	assert := assert.New(t)

	assert.Len(dispatcher.dispatched, 2)
	assert.Equal(message1, dispatcher.dispatched[0])
	assert.Equal(message2, dispatcher.dispatched[1])
}

/*
 * Mocks
 */

/*
 * Client factory
 */

type MockedClientFactory struct {
	mock.Mock
}

func (f *MockedClientFactory) Create(connection client.ConnectionInterface) client.ClientInterface {
	args := f.Called(connection)
	return args.Get(0).(*ServableMockedClient)
}

/*
 * Client
 */

type ServableMockedClient struct {
	messages         []msg.MessageInterface
	invalidCommandAt int
	callNumber       int
}

func newServableMockedClient(expectedMessages []msg.MessageInterface, invalidCommandAt int) *ServableMockedClient {
	return &ServableMockedClient{messages: expectedMessages, invalidCommandAt: invalidCommandAt}
}


func (c *ServableMockedClient) Id() uint64 {
	return 42 // Irrelevant
}

func (c *ServableMockedClient) Send(message *string) {
	// Irrelevant
}

func (c *ServableMockedClient) NextMessage() (msg.MessageInterface, *client.ClientError) {
	if c.callNumber == c.invalidCommandAt {
		c.callNumber++
		return nil, client.NewClientInvalidCommandError()
	}

	if len(c.messages) > 0 {
		message := c.messages[0]
		c.messages = c.messages[1:]
		c.callNumber++
		return message, nil
	}

	return nil, client.NewClientConnectionError() // To disconnect
}

/*
 * Dispatcher
 */

type MockedDispatcher struct {
	subscribed, unsubscribed []client.ClientInterface
	dispatched               []msg.MessageInterface
}

func (d *MockedDispatcher) Dispatch(message msg.MessageInterface) {
	d.dispatched = append(d.dispatched, message)
}

func (d *MockedDispatcher) Subscribe(c client.ClientInterface) {
	d.subscribed = append(d.subscribed, c)
}

func (d *MockedDispatcher) Unsubscribe(c client.ClientInterface) {
	d.unsubscribed = append(d.unsubscribed, c)
}

/*
 * Connection
 */

type ServerMockedConnection struct {
}

func (c *ServerMockedConnection) Write(message string) error {
	return nil
}

func (c *ServerMockedConnection) Read() (string, error) {
	return "", nil
}

func (c *ServerMockedConnection) Close() {
}
