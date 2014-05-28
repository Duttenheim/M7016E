package main

import (
	"fmt"
	"flag"
	"protocol"
)

type Test struct {
	// empty
}

type TestInput struct {}
type TestOutput struct {
	Message string
}

func (test *Test) Hest(input *TestInput, output *TestOutput) error {
	output.Message = "It works!"
	return nil
}

func main() {

	serverPort := flag.Int("serverPort", 2021, "Port of the server to listen on")
	clientAddr := flag.String("clientAddr", "localhost:2021", "Address and port to which to connect")
	openServer := flag.Bool("server", false, "Starts a server")
	openClient := flag.Bool("client", false, "Starts a client")
	var server *protocol.ServiceServer
	var client *protocol.ServiceClient
	flag.Parse()
	
	// setup rpc server
	rpc := protocol.MakeRpcServer()
	
	test := new(Test)
	rpc.Register(test)
	
	// open server if requested
	if (*openServer) {
		serverChannel := make(chan protocol.ServiceMsg)
		server = protocol.MakeServiceServer(serverChannel)
		rpc.Register(server)
		
		// listen on server-related messages
		go func() {
			for {
				select {
				case msg := <- serverChannel:
					if msg.ServiceDataType == protocol.RPC {
						reply, err := rpc.Process(msg.Payload)
						if err != nil {
							fmt.Printf("Server: RPC failed with %s\n", err.Error())
							continue
						}
						
						// create response and reply
						response := new(protocol.ServiceMsg)
						response.Type = protocol.Data
						response.ServiceDataType = protocol.Reply
						response.Payload = reply
						msg.Reply(response)
					} else {
						fmt.Printf("Server: Received %s\n", msg.ToString())
					}
				}
			}
		}()
		
		// start server
		go server.Start(*serverPort)
	}
	
	// open client if requested
	if (*openClient) {	
		clientChannel := make(chan protocol.ServiceMsg)
		client = protocol.MakeServiceClient(clientChannel)
		
		// listen on client-related messages
		go func() {
			for {
				select {
				case msg := <- clientChannel:
					if msg.ServiceDataType == protocol.RPC {
						reply, err := rpc.Process(msg.Payload)
						if err != nil {
							fmt.Printf("Client: RPC failed with %s\n", err.Error())
							continue
						}
						
						// create response and reply
						response := new(protocol.ServiceMsg)
						response.Type = protocol.Data
						response.ServiceDataType = protocol.Reply
						response.Payload = reply
						msg.Reply(response)
					} else {
						fmt.Printf("Client: Received %s\n", msg.ToString())
					}					
				}
			}
		}()
		
		go client.Connect(*clientAddr, func(client *protocol.ServiceClient, ip string) {
			msg := new(protocol.ServiceMsg)
			msg.Type = protocol.Data
			msg.ServiceDataType = protocol.RPC
			
			// create message
			args := new(TestInput)
			data, err := protocol.ComposeRPC("Test.Hest", args)
			if err != nil {
				fmt.Printf("RPC failed with: %s\n", err.Error())
				return
			}
			msg.Payload = data
			client.Send(msg)
		})
	}
	
	if (*openClient) {
		<- client.Done
	} 
	if (*openServer) {
		<- server.Done
	}
}