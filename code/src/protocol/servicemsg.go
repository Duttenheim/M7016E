package protocol

import (
	"code.google.com/p/go.net/websocket"	
	"encoding/json"
	"fmt"
	"sync"
)

const (
	Disconnect = iota
	Data
)

const (
	RPC = iota
	Reply
	Unknown
)

var MsgSequenceNumber int = 0
var mutex sync.Mutex

//------------------------------------------------------------------------------
/**
	Type used for service messages, the biggest likelyhood is that the payload is a JSON-encoded RPC string.
	Keeps a pointer to its socket so that we might send a reply.
*/
type ServiceMsg struct {
	Id				int
	Type			int
	ServiceDataType int
	Payload			string
	ws				*websocket.Conn
}

//------------------------------------------------------------------------------
/**
	Create new message. Messages must be created with this function to propery work!
*/
func MakeServiceMsg(msgType int, service int, payload string) *ServiceMsg {
	msg := new(ServiceMsg)
	msg.Type = msgType
	msg.ServiceDataType = service
	msg.Payload = payload
	
	mutex.Lock()
	msg.Id = MsgSequenceNumber
	MsgSequenceNumber++
	mutex.Unlock()
	
	return msg
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

//------------------------------------------------------------------------------
/**
	Converts to readable string
*/
func (msg *ServiceMsg) ToString() string {
	var retval string
	retval = "Message [type:"
	
	if msg.Type == Disconnect {
		retval += "Disconnect"
	} else if msg.Type == Data {
		retval += "Data"	
	} else {
		retval += "Unknown"
	}
	
	retval += " service:"
	if msg.ServiceDataType == RPC {
		retval += "RPC"
	} else if msg.ServiceDataType == Reply {
		retval += "Reply"
	} else {
		retval += "Unknown"
	}
	
	retval += " payload:" + msg.Payload
	retval += "]"
	return retval
}