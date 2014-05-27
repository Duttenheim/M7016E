package main

import (
	"fmt"
	"flag"
	"protocol"
)


func main() {

	serverPort := flag.Int("serverPort", 2021, "Port of the server to listen on")
	clientAddr := flag.String("clientAddr", "localhost:2021", "Address and port to which to connect")
	onlyServer := flag.Bool("onlyServer", false, "Only starts a server")
	onlyClient := flag.Bool("onlyClient", false, "Only starts a client")
	flag.Parse()
	
	if (!*onlyClient) {
		serverChannel := make(chan protocol.ServiceMsg)
		server := protocol.MakeServiceServer(serverChannel)
		
		// listen on server-related messages
		go func() {
			for {
				select {
				case msg := <- serverChannel:
					fmt.Printf("Server: Got message with payload %s", msg.Payload)
					break
				}
			}
		}()
		
		// start server
		go server.Start(*serverPort)
	}
	
	if (!*onlyServer) {	
		clientChannel := make(chan protocol.ServiceMsg)
		client := protocol.MakeServiceClient(clientChannel)
		
		// listen on client-related messages
		go func() {
			for {
				select {
				case msg := <- clientChannel:
					fmt.Printf("Client: Got message with payload %s", msg.Payload)
					break
				}
			}
		}()
		
		go client.Connect(*clientAddr)
	}
}