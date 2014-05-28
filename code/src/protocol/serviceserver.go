package protocol

import (
	"code.google.com/p/go.net/websocket"	
	"fmt"
	"encoding/json"
	"net/http"
)

//------------------------------------------------------------------------------
/**
*/
type ServiceServer struct {
	debug			bool
	msgChannel      chan ServiceMsg
	ws				*websocket.Conn
	Done			chan int
}

//------------------------------------------------------------------------------
/**
	Creates new server, this doesn't start the server, but merely creates a new one, messages received are passed to the channel.
*/
func MakeServiceServer(msgChannel chan ServiceMsg) *ServiceServer {
	server := new(ServiceServer)
	server.msgChannel = msgChannel
	server.Done = make(chan int)
	return server
}

//------------------------------------------------------------------------------
/**
	Internal handler for per-client connections
*/
func (server *ServiceServer) Handler(ws *websocket.Conn) {
	var err error
	var msg ServiceMsg
	
	fmt.Printf("ServiceServer: New connection established from %s\n", ws.Config().Origin.Host)
	
	for {
		dec := json.NewDecoder(ws)
		err = dec.Decode(&msg)
		msg.ws = ws
			
		if err != nil {
			fmt.Printf("ServiceServer: Lost connection to %s\n", ws.Config().Origin.Host)
			break
		}
		if msg.Type == Disconnect {
			// graceful disconnect, say goodbye
			fmt.Printf("ServiceServer: Client %s disconnected gracefully\n", ws.Config().Origin.Host)
			break
		} else {
			server.msgChannel <- msg
		}
	}
}

//------------------------------------------------------------------------------
/**
	Starts server listening on port
*/
func (server *ServiceServer) Start(port int) {
	fmt.Printf("ServiceServer: Opening listener on port %d\n", port)
	
	http.Handle("/service", websocket.Handler(server.Handler))
	
	portString := fmt.Sprintf("%d", port)
	err := http.ListenAndServe(":" + portString, nil)
	if err != nil {
		fmt.Printf("ServiceServer: Failed to open on port %d\n", port)
		return
	}
	
	server.Done <- 1
}