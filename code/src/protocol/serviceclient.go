package protocol

import (
	"code.google.com/p/go.net/websocket"	
	"fmt"
	"encoding/json"
)


//------------------------------------------------------------------------------
/**
*/
type ServiceClient struct {
	debug 				bool
	msgChannel			chan ServiceMsg
	ws					*websocket.Conn
}

//------------------------------------------------------------------------------
/**
	Create client
*/
func MakeServiceClient(msgChannel chan ServiceMsg) *ServiceClient {
	client := new(ServiceClient)
	client.msgChannel = msgChannel
	return client
}

//------------------------------------------------------------------------------
/**
	Connect to address, call in go-function since it will never return
*/
func (client *ServiceClient) Connect(ip string) {
	origin := "http://localhost/"
	url := "ws://" + ip + "/service"
	
	var err error
	client.ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Printf("ServiceClient: Failed to connect to %s\n", url)
		return
	}
	
	// handle messages until either the message receiving fails or the application is quit
	for {
		msg := client.receive()
		if msg == nil {
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

