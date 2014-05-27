package protocol

import (
	"code.google.com/p/go.net/websocket"	
	"encoding/json"
	"fmt"
)

const (
	Disconnect = iota
	Data
)

const (
	RPC = iota
	Unknown
)

//------------------------------------------------------------------------------
/**
	Type used for service messages, the biggest likelyhood is that the payload is a JSON-encoded RPC string.
	Keeps a pointer to its socket so that we might send a reply.
*/
type ServiceMsg struct {
	Type			int
	ServiceDataType int
	Payload			string
	ws				*websocket.Conn
}

//------------------------------------------------------------------------------
/**
	Sends a reply back to the sending socket, this will automatically perform a json encoding
*/
func (msg *ServiceMsg) Reply(reply *ServiceMsg) {
	enc := json.NewEncoder(msg.ws)
	err := enc.Encode(reply)
	if err != nil {
		fmt.Printf("ServiceMsg: Failed to encode message\n");
		return
	}
}