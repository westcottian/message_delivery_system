package main

import (
	"message_delivery_system/client"
	msg "message_delivery_system/message"
	)

var notSent = 0

func main() {
	d := client.NewDispather("testDispatcher")

	c1 := client.NewClient(100, d)
	c2 := client.NewClient(200, d)
	c3 := client.NewClient(300, d)

	message1 := msg.NewRelayMessage([]int64{200, 300, 100500}, "first message 1")
	message2 := msg.NewRelayMessage([]int64{100, 200, 300}, "second message 2")
	message3 := msg.NewRelayMessage([]int64{100, 300}, "third message 3")
	message4 := msg.NewRelayMessage([]int64{100, 200}, "fourth message 4")
	message5 := msg.NewRelayMessage([]int64{9001, 100, 300, 100, 200}, "fifth message 5")

	c1.Say(message1)

	d.Subscribe(c1)
	d.Subscribe(c2)
	d.Subscribe(c3)

	c1.Say(message2)
	c2.Say(message3)
	c3.Say(message4)
	c2.Say(message5)

	//fmt.Printf("Due for sending: %d\n", notSent)

	for notSent > 0 {
	}

	//fmt.Printf("Due for sending: %d\n", notSent)
}
