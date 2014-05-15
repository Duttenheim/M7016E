package main

import (
	"bitverse"
	"fmt"
    "os"
    "strconv"
)


func main() {
	// create websocket transpot and channel
	transport := bitverse.MakeWSTransport()
	var done chan int
	var port int = 0
	//var node* bitverse.SuperNode

	if (len(os.Args) > 1){
		port,_ = strconv.Atoi(os.Args[1])
	}
	if ((port > 1023) && (port < 49151)){
		_, done = bitverse.MakeSuperNode(transport, "127.0.0.1", os.Args[1])

		<- done
	}
	fmt.Println("usage : ./bitverseserver port   with 1023<port<49151")
	os.Exit(2)
}
