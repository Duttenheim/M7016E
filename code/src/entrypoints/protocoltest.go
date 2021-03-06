
package main

import (
	"fmt"
	"bitverse"
	"protocol"
	"strconv"
	"strings"
	"os"
	"bufio"
	"flag"
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
	
	// make edge node
	observer := new(TestObserver)
	node, done := bitverse.MakeEdgeNode(transport, observer)
	
	// create our special RPC server
	server := protocol.MakeRpcServer()

	// make a docker RPC-callable object
	test := new(TestRpcInterface)
	server.Register(test)

	// make message observer
	msgObserver := protocol.MakeRpcMessageObserver(server)	
	
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

		// get address and port
    addr := flag.String("address", "localhost", "Bitverse server IP-address")
	port := flag.Int("port", 2020, "Bitverse server port")
    flag.Parse()
	
	portString := fmt.Sprintf("%d", (*port));

	// connect and wait for connection to die
	go node.Connect(*addr + ":" + portString)
	<- done	
}
