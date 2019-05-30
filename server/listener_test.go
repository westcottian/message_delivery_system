package server

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"sync"
	msg "message_delivery_system/message"
)

func TestListsAndIdentities(t *testing.T) {
        port := 1234

        l := NewListener(port)
        go l.Listen()

        var wg sync.WaitGroup

        input1, output1 := getHelperClient(port)
        input2, output2 := getHelperClient(port)
        input3, output3 := getHelperClient(port)

        input1 <- msg.MessageTypeIdentity

        input2 <- msg.MessageTypeList
        input2 <- msg.MessageTypeList

        input3 <- msg.MessageTypeList
        input3 <- msg.MessageTypeIdentity
        input3 <- msg.MessageTypeIdentity

        assert := assert.New(t)

        wg.Add(3)

        go func() {
                defer wg.Done()
                r := <-output1
                assert.Contains(r, "Your ID is 2")
        }()
 	go func() {
                defer wg.Done()
                r := <-output2
                assert.Contains(r, "Client IDs are")
                assert.Contains(r, "1")
                assert.NotContains(r, "3")
                assert.Contains(r, "2")

                r = <-output2
                assert.Contains(r, "Client IDs are")
                assert.Contains(r, "1")
                assert.NotContains(r, "3")
                assert.Contains(r, "2")

        }()

        go func() {
                defer wg.Done()
                r := <-output3
                assert.Contains(r, "Client IDs are")
                assert.NotContains(r, "1")

                r = <-output3
                assert.Contains(r, "Your ID is 1")

                r = <-output3
                assert.Contains(r, "Your ID is 1")

        }()

        wg.Wait()

}	

func TestRelays(t *testing.T) {
        port := 1235

        l := NewListener(port)
        go l.Listen()

        var wg sync.WaitGroup

        in1, out1 := getHelperClient(port)
        in2, out2 := getHelperClient(port)
        in3, out3 := getHelperClient(port)

        body1 := "test message 1"
        body2 := "test END message 2\numad?"
        body3 := "test message 3\n\nEND"

        message1 := fmt.Sprintf("%s\n3\n%s", msg.MessageTypeRelay, body1)
        message2 := fmt.Sprintf("%s\n2,1\n%s", msg.MessageTypeRelay, body2)
        message3 := fmt.Sprintf("%s\ninvalid receivers\nignored body", msg.MessageTypeRelay)
        message4 := fmt.Sprintf("invalid type\n1\nignored body")
        message5 := fmt.Sprintf("%s\n1,100500\n%s", msg.MessageTypeRelay, body3)

	in1 <- message1
        in1 <- message3

        in2 <- message4

        in3 <- message2
        in3 <- message5
	
        assert := assert.New(t)

        wg.Add(1)

	go func() {
                defer wg.Done()
                r := <-out1
                assert.Contains(r, body2)

        }()
 

        go func() {
                defer wg.Done()
                r := <-out2
                assert.Contains(r, body1)
        }()

        go func() {
                defer wg.Done()
                r := <-out3
                assert.Contains(r, body2)

		r = <-out3
		assert.Contains(r, body3)
        }()

        wg.Wait()		
}


/*
 * Helpers
 */

func getHelperClient(port int) (chan string, chan string) {
	input := make(chan string, 255)
	output := make(chan string, 255)
	
	bytes := make([]byte, 1024)
	var len int
	go func() {
		c, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
                 fmt.Println(err)
                 return 
                }
		defer c.Close()
		//reader := bufio.NewReader(c)
		for {
			reader := bufio.NewReader(c)
			// Sending
			m := <-input
			fmt.Fprintf(c, m+"\n")
			// Reading response
			len, err = reader.Read(bytes)
			if err != nil {
                		 fmt.Println(err)
                 		return
                	}	
			output <- string(bytes[:len])
		}
	}()
	return input, output
}

