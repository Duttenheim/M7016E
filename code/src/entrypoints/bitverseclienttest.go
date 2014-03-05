package main

import (
	"fmt"
	"bitverse"
)

type TestObserver struct {}

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
	
	// send them a test message
	service := node.GetMsgService("TestMessagingService")
	for _, child := range children {
		service.Send(child, "Hello!")
	}
}

func (observer* TestObserver) OnConnected(localNode* bitverse.EdgeNode, remoteNode* bitverse.RemoteNode) {
	fmt.Printf("Connected to supernode '%s'\n", remoteNode.Id())
	remoteNode.SendChildrenRequest()
}

type TestMsgObserver struct {}

func (observer* TestMsgObserver) OnDeliver(service* bitverse.MsgService, msg* bitverse.Msg) {
	fmt.Println(msg.Payload)
}

var secret string = "3e606ad97e0a738d8da4c4c74e8cd1f1f2e016c74d85f17ac2fc3b5dab4ed6c4"

func main() {
	// create transport and channel
	transport := bitverse.MakeWSTransport()
	var done chan int
	
	// make edge node
	node, done := bitverse.MakeEdgeNode(transport, new(TestObserver))

	_, serviceError := node.CreateMsgService(secret, "TestMessagingService", new(TestMsgObserver))
	if (serviceError != nil) {
		panic(serviceError)
	}

	// connect node and wait until done (which is forever)
	go node.Connect("localhost:2020")
	<- done	
}