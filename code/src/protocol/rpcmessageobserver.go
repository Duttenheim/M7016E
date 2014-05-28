package protocol

import (
	"bitverse"
	"encoding/json"	
	"fmt"
)


//------------------------------------------------------------------------------
/**
*/
type RpcMessageObserver struct {
	debug		bool
	server		*RpcServer
}


//------------------------------------------------------------------------------
/**
	Opens the observer, run this before using it for anything
*/
func MakeRpcMessageObserver(server *RpcServer) *RpcMessageObserver{
	observer := new(RpcMessageObserver)
	observer.debug = false
	observer.server = server
	return observer
}

//------------------------------------------------------------------------------
/**
*/
func (observer* RpcMessageObserver) SetDebugEnabled(enabled bool) {
	observer.debug = enabled
}

//------------------------------------------------------------------------------
/**
	Composes an RPC message with an output and input and executes it asynchronously
*/
func RpcInvokeAsync(functionName string, args interface{}, reply interface{}, node string, service* bitverse.MsgService, done chan* RpcInvocation) *RpcInvocation {

	// create new invocation
	invocation := MakeInvocation(functionName, reply, done)
	
	// create message
	data, err := ComposeRPC(functionName, args)
	if (err != nil) {
		invocation.Error = err
		invocation.done()
	}

	// send message
	service.SendAndGetReply(node, string(data), 1, func(err error, data interface{}) {

		if (err != nil) {
			invocation.Error = err
			invocation.done()
			return
		}

		// unmarshal data
		encerr := json.Unmarshal([]byte(data.(string)), invocation.Reply)

		if (encerr != nil) {
			invocation.Error = err
			invocation.done()
			return
		}

		invocation.Error = nil
		invocation.done()
	})

	return invocation
}

//------------------------------------------------------------------------------
/**
	Invokes RPC synchronously
*/
func RpcInvokeSync(functionName string, args interface{}, reply interface{}, node string, service* bitverse.MsgService) error {
	call := <-RpcInvokeAsync(functionName, args, reply, node, service, make(chan *RpcInvocation, 1)).Done
	return call.Error
}

//------------------------------------------------------------------------------
/**
	Receives a message, decodes it as JSON and attempts to run it as an RPC call.
	Currently requires two messages and it is doubtful if it will properly work.
*/
func (observer* RpcMessageObserver) OnDeliver(service* bitverse.MsgService, msg* bitverse.Msg) {

	reply, err := observer.server.Process(msg.Payload)
	if err != nil {
		fmt.Printf("RpcMessageObserver: " + err.Error())
		return
	}
	msg.Reply(string(reply))
}
