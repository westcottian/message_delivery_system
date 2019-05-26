package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	msg "message_delivery_system/message"
)

func TestDispatcherImplementsDispatcherInterface(t *testing.T) {
	var _ DispatcherInterface = (*Dispatcher)(nil)
}

func TestDispatchRespondsWithIdentity(t *testing.T) {
	clientId := uint64(42)
	client := newDispatcheableMockedClient(clientId)

	d := NewDispatcher()
	d.Subscribe(client)
	d.Dispatch(msg.NewIdentityMessage(clientId))

	assert := assert.New(t)
	assert.Len(client.received, 1)
	assertSentMessageContainsClient(assert, client.received[0], clientId)
}

func TestDispatchNotSendingIdentityToDisconnectedClient(t *testing.T) {
	clientId := uint64(42)
	client := newDispatcheableMockedClient(clientId)

	d := NewDispatcher()
	d.Subscribe(client)
	d.Unsubscribe(client)

	d.Dispatch(msg.NewIdentityMessage(clientId))

	assert.Len(t, client.received, 0)
}

func TestDispatchRespondsWithList(t *testing.T) {
	clientId1 := uint64(42)
	clientId2 := uint64(100500)
	clientId3 := uint64(9001)

	client1 := newDispatcheableMockedClient(clientId1)
	client2 := newDispatcheableMockedClient(clientId2)
	client3 := newDispatcheableMockedClient(clientId3)

	assert := assert.New(t)

	d := NewDispatcher()
	d.Subscribe(client1)
	d.Dispatch(msg.NewListMessage(clientId1))

	var m string

	assert.Len(client1.received, 1)
	assert.Len(client2.received, 0) // Nothing sent there
	assert.Len(client3.received, 0)
	m = client1.received[0]
	assertSentMessageNotContainsClient(assert, m, clientId1)
	assertSentMessageNotContainsClient(assert, m, clientId2)
	assertSentMessageNotContainsClient(assert, m, clientId3)

	// Now lets connect more clients

	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Dispatch(msg.NewListMessage(clientId1))

	assert.Len(client1.received, 2)
	assert.Len(client2.received, 0) // Still nothing sent there
	assert.Len(client3.received, 0)
	m = client1.received[1]
	assertSentMessageNotContainsClient(assert, m, clientId1)
	assertSentMessageContainsClient(assert, m, clientId2)
	assertSentMessageContainsClient(assert, m, clientId3)

	d.Dispatch(msg.NewListMessage(clientId2))

	assert.Len(client1.received, 2)
	assert.Len(client2.received, 1)
	assert.Len(client3.received, 0)
	m = client2.received[0]
	assertSentMessageContainsClient(assert, m, clientId1)
	assertSentMessageNotContainsClient(assert, m, clientId2)
	assertSentMessageContainsClient(assert, m, clientId3)
}

func TestDispatchRespondsWithShorterListAfterUnsubscribing(t *testing.T) {
	clientId1 := uint64(42)
	clientId2 := uint64(100500)
	clientId3 := uint64(9001)

	client1 := newDispatcheableMockedClient(clientId1)
	client2 := newDispatcheableMockedClient(clientId2)
	client3 := newDispatcheableMockedClient(clientId3)

	assert := assert.New(t)

	d := NewDispatcher()
	d.Subscribe(client1)
	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Dispatch(msg.NewListMessage(clientId1))

	var m string

	assert.Len(client1.received, 1)
	m = client1.received[0]
	assertSentMessageContainsClient(assert, m, clientId2)
	assertSentMessageContainsClient(assert, m, clientId3)

	// Now lets disconnect client2

	d.Unsubscribe(client2)

	d.Dispatch(msg.NewListMessage(clientId1))

	assert.Len(client1.received, 2)
	m = client1.received[1]
	assertSentMessageNotContainsClient(assert, m, clientId2) // Must go away
	assertSentMessageContainsClient(assert, m, clientId3)
}

func TestDispatchSendsRelayMessageToReceivers(t *testing.T) {
	body1 := "test message\nbody 1"
	body2 := "test message\nbody 2"
	body3 := "test message\nbody 3"

	clientId1 := uint64(42)
	clientId2 := uint64(100500)
	clientId3 := uint64(9001)

	client1 := newDispatcheableMockedClient(clientId1)
	client2 := newDispatcheableMockedClient(clientId2)
	client3 := newDispatcheableMockedClient(clientId3)

	assert := assert.New(t)

	d := NewDispatcher()
	d.Subscribe(client1)
	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Dispatch(msg.NewRelayMessage(clientId1, []uint64{clientId2, clientId3, clientId1}, &body1))
	d.Dispatch(msg.NewRelayMessage(clientId1, []uint64{clientId2}, &body2))
	d.Dispatch(msg.NewRelayMessage(clientId3, []uint64{clientId2, clientId1}, &body3))

	assert.Len(client1.received, 2)
	assert.Equal(body1, client1.received[0])
	assert.Equal(body3, client1.received[1])

	assert.Len(client2.received, 3)
	assert.Equal(body1, client2.received[0])
	assert.Equal(body2, client2.received[1])
	assert.Equal(body3, client2.received[2])

	assert.Len(client3.received, 1)
	assert.Equal(body1, client3.received[0])
}

func TestDispatchIgnoresNonExistingReceivers(t *testing.T) {
	body1 := "test message\nbody 1"

	clientId1 := uint64(42)
	clientId2 := uint64(100500)
	clientId3 := uint64(9001)

	client1 := newDispatcheableMockedClient(clientId1)
	client2 := newDispatcheableMockedClient(clientId2)
	client3 := newDispatcheableMockedClient(clientId3)

	assert := assert.New(t)

	d := NewDispatcher()
	d.Subscribe(client1)
	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Dispatch(msg.NewRelayMessage(clientId1, []uint64{clientId2, 987654321, clientId3}, &body1))

	assert.Len(client1.received, 0)

	assert.Len(client2.received, 1)
	assert.Equal(body1, client2.received[0])

	assert.Len(client3.received, 1)
	assert.Equal(body1, client3.received[0])
}

func TestDispatchNotSendingToUnsubscribedClients(t *testing.T) {
	body1 := "test message\nbody 1"

	clientId1 := uint64(42)
	clientId2 := uint64(100500)
	clientId3 := uint64(9001)

	client1 := newDispatcheableMockedClient(clientId1)
	client2 := newDispatcheableMockedClient(clientId2)
	client3 := newDispatcheableMockedClient(clientId3)

	assert := assert.New(t)

	d := NewDispatcher()
	d.Subscribe(client1)
	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Unsubscribe(client2)

	d.Dispatch(msg.NewRelayMessage(clientId1, []uint64{clientId2, clientId3}, &body1))

	assert.Len(client1.received, 0)
	assert.Len(client2.received, 0)

	assert.Len(client3.received, 1)
	assert.Equal(body1, client3.received[0])
}

/*
 * Helpers
 */

func assertSentMessageContainsClient(assert *assert.Assertions, sent string, clientId uint64) {
	assert.Contains(sent, fmt.Sprintf("%d", clientId))
}

func assertSentMessageNotContainsClient(assert *assert.Assertions, sent string, clientId uint64) {
	assert.NotContains(sent, fmt.Sprintf("%d", clientId))
}

/*
 * Mocks
 */

type DispatcherMockedClient struct {
	id uint64
	mock.Mock
	received []string
}

func newDispatcheableMockedClient(id uint64) *DispatcherMockedClient {
	return &DispatcherMockedClient{id: id}
}

func (c *DispatcherMockedClient) Id() uint64 {
	return c.id
}

func (c *DispatcherMockedClient) Send(message *string) {
	c.received = append(c.received, *message)
}

func (c *DispatcherMockedClient) NextMessage() (msg.MessageInterface, *ClientError) {
	return nil, nil // Irrelevant for this test
}
