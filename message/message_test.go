package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageImplementsMessageInterface(t *testing.T) {
	var _ MessageInterface = (*Message)(nil)
}

func TestRelayCommandReturnsRelay(t *testing.T) {
	body := "testBody"
	m := NewRelayMessage(9001, []uint64{100500, 42}, &body)

	assert.Equal(t, MessageTypeRelay, m.Command())
}

func TestRelaySenderReturnsSender(t *testing.T) {
	sender := uint64(42)
	body := "testBody"
	m := NewRelayMessage(sender, []uint64{100500, 42}, &body)

	assert.Equal(t, sender, m.Sender())
}

func TestRelayBodyReturnsBody(t *testing.T) {
	body := "testBody\nline2 foobar"
	m := NewRelayMessage(9001, []uint64{100500, 42}, &body)

	assert.Equal(t, body, *m.Body())
}

func TestRelayReceiversReturnsReceivers(t *testing.T) {
	body := "testBody"
	receivers := []uint64{100500, 42}
	m := NewRelayMessage(9001, receivers, &body)
	r := m.Receivers()

	assert.Len(t, r, 2)
	assert.EqualValues(t, 100500, r[0])
	assert.EqualValues(t, 42, r[1])
}

func TestIdentityCommandReturnsIdentity(t *testing.T) {
	m := NewIdentityMessage(42)

	assert.Equal(t, MessageTypeIdentity, m.Command())
}

func TestIdentitySenderReturnsSender(t *testing.T) {
	sender := uint64(42)
	m := NewIdentityMessage(sender)

	assert.Equal(t, sender, m.Sender())
}

func TestIdentityBodyIsEmpty(t *testing.T) {
	m := NewIdentityMessage(42)

	assert.Empty(t, m.Body())
}

func TestIdentityReceiversAreEmpty(t *testing.T) {
	m := NewIdentityMessage(42)
	r := m.Receivers()

	assert.Empty(t, r)
}

func TestListCommandReturnsList(t *testing.T) {
	m := NewListMessage(42)

	assert.Equal(t, MessageTypeList, m.Command())
}

func TestListSenderReturnsSender(t *testing.T) {
	sender := uint64(42)
	m := NewListMessage(sender)

	assert.Equal(t, sender, m.Sender())
}

func TestListBodyIsEmpty(t *testing.T) {
	m := NewListMessage(42)

	assert.Empty(t, m.Body())
}

func TestListReceiversAreEmpty(t *testing.T) {
	m := NewListMessage(42)
	r := m.Receivers()

	assert.Empty(t, r)
}
