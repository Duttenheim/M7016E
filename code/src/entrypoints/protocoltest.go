
package main

import (
	"fmt"
	"bitverse"
	"protocol"
)

type TestRpcInterface struct {}
type TestRpcInput struct {}
type TestRpcOutput struct {
	Content string
}

func (obj* TestRpcInterface) Test(input* TestRpcInput, output* TestRpcOutput) error {
	output.Content = "It works!"
	return nil
}

type TestObserver struct {
	superNode* bitverse.RemoteNode
}

func (observer* TestObserver) OnSiblingJoined(node* bitverse.EdgeNode, id string) {
	fmt.Printf("Sibling '%s' joined supernode!\n", id)
}

func (observer* TestObserver) OnSiblingLeft(node* bitverse.EdgeNode, id string) {
	fmt.Printf("Sibling '%s' left supernode!\n", id)
}

func (observer* TestObserver) OnSiblingHeartbeat(node* bitverse.EdgeNode, id string) {
	fmt.Printf("Sibling '%s' is still alive!\n", id)
}

func (observer* TestObserver) OnChildrenReply(node* bitverse.EdgeNode, id string, children[] string) {
	fmt.Printf("Supernode has '%d' children\n", len(children))

	// run a new thread which sends locking RPC commands (needs to run in go routine since the Rpc invocation requires a reply to be sent back synchronously, meaning it will block the messaging thread)
	go func() {	
		// send them a test message
		service := node.GetMsgService("TestMessagingService")
		for _, child := range children {

			if (child == node.Id()) {
				continue
			}

			// create argument and invoke
			args := new(TestRpcInput)
			reply := new(TestRpcOutput)
			err := protocol.RpcInvokeSync("TestRpcInterface.Test", args, reply, child, service)

			if (err != nil) {
				fmt.Println(err)
				continue
			}
			fmt.Println(reply.Content)
		}
	}()
}

func (observer* TestObserver) OnConnected(localNode* bitverse.EdgeNode, remoteNode* bitverse.RemoteNode) {
	fmt.Printf("Connected to supernode '%s'\n", remoteNode.Id())
	observer.superNode = remoteNode
	//remoteNode.SendChildrenRequest()
}

var secret string = "3e606ad97e0a738d8da4c4c74e8cd1f1f2e016c74d85f17ac2fc3b5dab4ed6c4"

func main() {
	// create transport and channel
	transport := bitverse.MakeWSTransport()
	var done chan int
	
	// make edge node
	observer := new(TestObserver)
	node, done := bitverse.MakeEdgeNode(transport, observer)

	// make message observer
	msgObserver := new(protocol.RpcMessageObserver)
	msgObserver.Open()

	// make a docker RPC-callable object
	test := new(TestRpcInterface)
	msgObserver.Register(test)
	
	_, serviceError := node.CreateMsgService(secret, "TestMessagingService", msgObserver)
	if (serviceError != nil) {
		panic(serviceError)
	}

	// wait for input
	go func() {
		for {
			var c int
			fmt.Scanf("%d", &c)
			if (c == 1) {
				observer.superNode.SendChildrenRequest()
			}
		}
	}()

	// connect node and wait until done (which is forever)
	go node.Connect("localhost:2020")
	<- done	
}
