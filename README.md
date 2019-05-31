# message_delivery_system
A message delivery system with a server/hub which relays messages to the clients/receivers. 

This is a TCP server accepting (using a simple client for testing purposes)
1 - Echo message
2 - Identify yourself
3 - List clients connected
4 - Relay message to other clients (up to 1024KB binary-encoded payload and 255 clients)

There are below major components here:

### HUB

The Hub acts like a server, listening to one port and accepting all the new connected clients. When a client connects, starts a goroutine to process the connection:

After the initialization, the goroutine enters in a for loop, constantly reading incoming Requests from the client. When a new request is received, the Hub starts a new goroutine to handle the operation.
  
Hub relays incoming message bodies to receivers based on user ID(s) defined in the message. Hub assigns unique user id to the client once its connected.

    user_id - unsigned 64 bit integer
    Connection to hub is done using pure TCP.
    

### CLIENT

Clients are users who are connected to the hub. 
Like Hub, the Client connects to a given address and starts a goroutine.
    
Client may send three types of messages which are described below.

###### Identity message
Client can send a identity message which the hub will answer with the user_id of the connected user.

![Identity](https://raw.githubusercontent.com/Everyplay/developer-assignment-backend/master/identity.seq.png)

###### List message
Client can send a list message which the hub will answer with the list of all connected client user_id:s (excluding the requesting client).

![List](https://raw.githubusercontent.com/Everyplay/developer-assignment-backend/master/list.seq.png)

###### Relay message
Client can send a relay messages which body is relayed to receivers marked in the message.

![Relay](https://raw.githubusercontent.com/Everyplay/developer-assignment-backend/master/relay.seq.png)

*Relay example: receivers: 2 and 3, body: foobar*


## Standard Points:
This repository has been made public with all commit history included;

To run this program, you will need first to install:

		Git
		Go
		
## Installation
Clone the project under the directory **$GOPATH/src/.
Go to the dir message_delivery_system
Compile using below command:
		$ godep go build 
		
It will create executable binary named message_delivery_system.

## Usage

Start the server:

    ./message_delivery_system -port=1234

Now connect to it:

    nc 127.0.0.1 1234
    
Commands:
    
To get your ID:

    identity
    
To get the list of connected clients:

    list

Relay message format:

    relay // Type
    42,100500,9001 // Receivers
    foo bar
    umad?

## Test
Run the test script to run all the tests and benchmarks

```sh runtestandcodecoverage.sh```
