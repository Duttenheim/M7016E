package protocol

import (
	"code.google.com/p/go.net/websocket"	
	"fmt"
	"encoding/json"
)


//------------------------------------------------------------------------------
/**
	This type serves as a connection to a service providing server.
	Connect this whenever you have decided to to connect to through bitverse
*/
type ServiceClient struct {
	debug 				bool
	msgChannel			chan ServiceMsg
	ws					*websocket.Conn
	Done				chan int
}

//------------------------------------------------------------------------------
/**
	Create client
*/
func MakeServiceClient(msgChannel chan ServiceMsg) *ServiceClient {
	client := new(ServiceClient)
	client.msgChannel = msgChannel
	client.Done = make(chan int)
	return client
}

//------------------------------------------------------------------------------
/**
	Connect to address, call in go-function since it will never return.
	Takes a callback which gets run whenever the client connects to the server, this is required since we might run the client in a go routine.
*/
func (client *ServiceClient) Connect(ip string, connected func(client *ServiceClient, ip string)) {
	origin := "http://localhost/"
	url := "ws://" + ip + "/service"
	
	var err error
	client.ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Printf("ServiceClient: Failed to connect to %s\n", url)
		return
	}
	
	fmt.Printf("ServiceClient: Connected to %s\n", ip)
	
	// run callback if its been set
	if (connected != nil) {
		connected(client, url)
	}
	
	// handle messages until either the message receiving fails or the application is quit
	for {
		msg := client.receive()
		if msg == nil {
			client.Done <- 1
			break
		}
		client.msgChannel <- *msg
	}
}

//------------------------------------------------------------------------------
/**
	Send message
*/
func (client *ServiceClient) Send(msg *ServiceMsg) {
	enc := json.NewEncoder(client.ws)
	err := enc.Encode(msg)
	if err != nil {
		fmt.Printf("ServiceClient: Failed to encode message\n");
		return
	}
}

//------------------------------------------------------------------------------
/**
*/
func (client *ServiceClient) receive() *ServiceMsg {
	dec := json.NewDecoder(client.ws)
	var err error
	var msg ServiceMsg

	err = dec.Decode(&msg)
	if err != nil {
		fmt.Printf("ServiceClient: Failed to decode message\n")
		return nil
	}

	return &msg
}

