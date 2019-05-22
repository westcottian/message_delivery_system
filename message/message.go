package message

const MessageTypeRelay string = "relay"
const MessageTypeIdentity string = "identity"
const MessageTypeList string = "list"

type MessageInterface interface {
	Command() string
	Sender() int64
	Receivers() []int64
	Body() string
}

type Message struct {
	command   string
	sender    int64
	receivers []int64
	body      string
}

func NewRelayMessage(sender int64, receivers []int64, body string) *Message {
	return &Message{command: MessageTypeRelay, body: body, sender: sender, receivers: receivers}
}

func NewIdentityMessage(sender int64) *Message {
	return &Message{command: MessageTypeIdentity, sender: sender}
}

func NewListMessage(sender int64) *Message {
	return &Message{command: MessageTypeList, sender: sender}
}

func (m *Message) Command() string {
	return m.command
}

func (m *Message) Sender() int64 {
	return m.sender
}

func (m *Message) Receivers() []int64 {
	return m.receivers
}

func (m *Message) Body() string {
	//return "42"
	return m.body
}
