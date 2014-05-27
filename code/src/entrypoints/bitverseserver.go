package main

import (
    "net/http"
	"bitverse"
	"flag"
	"fmt"
	"bufio"
	"os"
	"strings"
)


func main() {
	// create websocket transpot and channel
	transport := bitverse.MakeWSTransport()
	// var done chan int

	//var node* bitverse.SuperNode

	// implicitly make this server work as a web server too
	http.Handle("/", http.FileServer(http.Dir(".")))

	port := flag.Int("port", 2020, "Server port, should be in the range 1023 - 49151");
	
	ringAddress := flag.String("ringAddress", "127.0.0.1", "Supernodes to connect to, separate with comma");
	ringPort := flag.Int("ringPort", 2020, "Supernode port for all supernodes")
	flag.Parse();

	portString := fmt.Sprintf("%d", (*port));
	ringPortString := fmt.Sprintf("%d", (*ringPort));
		
	node, done := bitverse.MakeSuperNode(transport, "127.0.0.1", portString)
	
	// wait for input, whenever enter gets pressed, we connect this supernode to another
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	
	addrList := strings.Split(*ringAddress, ",")
	node.ConnectSuccessor(addrList, ringPortString)
	<- done
}
