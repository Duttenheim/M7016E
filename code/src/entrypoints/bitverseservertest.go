package main

import (
    "net/http"
	"bitverse"
)


func main() {
	// create websocket transpot and channel
	transport := bitverse.MakeWSTransport()
	var done chan int

    // implicitly make this server work as a web server too
    http.Handle("/", http.FileServer(http.Dir("web/")))
	_, done = bitverse.MakeSuperNode(transport, "127.0.0.1", "2020")

	<- done
}
