package main

import (
	"bitverse"
)


func main() {
	// create websocket transpot and channel
	transport := bitverse.MakeWSTransport()
	var done chan int
	//var node* bitverse.SuperNode

	_, done = bitverse.MakeSuperNode(transport, "127.0.0.1", "2020")

	<- done
}
