
package main

import (
	"fmt"
	"bitverse"
	"protocol"
	"strconv"
	"strings"
	"os"
	"bufio"
	"time"
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
	//fmt.Printf("Sibling '%s' joined supernode!\n", id)
}

func (observer* TestObserver) OnSiblingLeft(node* bitverse.EdgeNode, id string) {
	//fmt.Printf("Sibling '%s' left supernode!\n", id)
}

func (observer* TestObserver) OnSiblingHeartbeat(node* bitverse.EdgeNode, id string) {
	//fmt.Printf("Sibling '%s' is still alive!\n", id)
}

var siblings []string
var reply chan bool

func (observer* TestObserver) OnChildrenReply(node* bitverse.EdgeNode, id string, children[] string) {
	siblings = children
	reply <- true
}

func (observer* TestObserver) OnConnected(localNode* bitverse.EdgeNode, remoteNode* bitverse.RemoteNode) {
	observer.superNode = remoteNode
}

var secret string = "3e606ad97e0a738d8da4c4c74e8cd1f1f2e016c74d85f17ac2fc3b5dab4ed6c4"

func main() {
	// create transport and channel
	transport := bitverse.MakeWSTransport()
	var done chan int
	addresses := make(chan string)
	stop := make(chan string)
	var oldID [100]string
	var numberOfId = 0
	
	// make edge node
	observer := new(TestObserver)
	node, done := bitverse.MakeEdgeNode(transport, observer)

	// make message observer
	msgObserver := new(protocol.RpcMessageObserver)
	msgObserver.Open()

	// make a docker RPC-callable object
	test := new(TestRpcInterface)
	msgObserver.Register(test)
	
	service, serviceError := node.CreateMsgService(secret, "RPCMessageService", msgObserver)
	if (serviceError != nil) {
		panic(serviceError)
	}

	// wait for input
	go func() {
		bio := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf(">>")
			line, _ := bio.ReadString('\n')
			components := strings.Split(line, " ")
			for index, comp := range components {
				components[index] = strings.Trim(comp, "\n")
			}
			if (len(components) > 0) {
				if (components[0] == "show") {
					reply = make(chan bool)
					observer.superNode.SendChildrenRequest()
					<-reply
					for index, value := range siblings {
						fmt.Printf("%d. %s\n", index, value)
					}
				} else if (components[0] == "display") {
					if (len(components) != 3) {
						fmt.Println("display address port\n")
						continue
					}else {
						s := []string{components[1], components[2]};
						var address string = strings.Join(s, ":");
						node.Unconnect()
						addresses <- address
						time.Sleep(2)
						<- addresses
						
						
					}
				}else if (components[0] == "ids") {
					for i := 0; i<numberOfId; i++ {
						fmt.Println(oldID[i])
					}
				} else if (components[0] == "run") {
					if (len(components) != 2) {
						fmt.Printf("'run' command must be supplied with a sibling node index!\n")
						continue
					}
					index, err := strconv.Atoi(components[1])
					if (err != nil) {
						fmt.Printf("Must supply a valid number to 'run'\n")
						continue
					}
					if (index < len(siblings)) {
						args := new(TestRpcInput)
						reply := new(TestRpcOutput)
						err := protocol.RpcInvokeSync("TestRpcInterface.Test", args, reply, siblings[index], service)
						if (err != nil) {
							fmt.Println(err)
							continue
						}
	
						fmt.Println(reply.Content)
					} else {
						fmt.Printf("Unknown sibling index '%d'\n", index)
					}
				} else {
					fmt.Println("Unknown command: " + components[0])
				}
			}
		}
	}()
	
	go func(){
		for {
			address := <- addresses
			//fmt.Println("address = ", address )
			oldID[numberOfId] = node.Id()
			numberOfId += 1
			if(numberOfId >= len(oldID)){
				numberOfId = 0
			}
			fmt.Println("deleting node id : ", node.Id(), "        supernode id :", observer.superNode.Id)
			node, done = bitverse.MakeEdgeNode(transport, observer)
			fmt.Println("creating node id : ", node.Id(), "        supernode id :", observer.superNode.Id)
			go node.Connect(address)
			addresses <- "OK"
			fmt.Println("connected to", address)
			
			time.Sleep(2 * time.Second)
			reply = make(chan bool)
			observer.superNode.SendChildrenRequest()
			<-reply
			for index, value := range siblings {
				show := true
				for i := 0; i<numberOfId; i++ { 
					//fmt.Println("comparing ",value," and ",oldID[i])
					if(value == oldID[i]){
						show = false
						break
					}
				}
				if(show){
					fmt.Printf("%d. %s\n", index, value)
				}
			}
			<- done
		}
	}()
	
	// connect node and wait until done (which is forever)
	go node.Connect("localhost:2020")
	//go node.Connect("130.240.233.93:2020")
	<- done
	<- stop	
}
