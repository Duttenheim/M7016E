package main

import (
    "net/http"
	"bitverse"
	"flag"
	"fmt"
)


func main() {
	// create websocket transpot and channel
	transport := bitverse.MakeWSTransport()
	var done chan int
	//var node* bitverse.SuperNode

	// implicitly make this server work as a web server too
	http.Handle("/", http.FileServer(http.Dir(".")))

	port := flag.Int("port", 2020, "Server port, should be in the range 1023 - 49151");
	flag.Parse();

	portString := fmt.Sprintf("%d", (*port));
		
	_, done = bitverse.MakeSuperNode(transport, "127.0.0.1", portString)
	<- done
}
