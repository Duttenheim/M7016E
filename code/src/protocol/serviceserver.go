package protocol

import (
	"code.google.com/p/go.net/websocket"	
	"fmt"
	"encoding/json"
	"net/http"
    "net"
)

//------------------------------------------------------------------------------
/**
*/
type ServiceServer struct {
	debug			bool
    bitverse        string
	msgChannel      chan ServiceMsg
	Done			chan int
}

//------------------------------------------------------------------------------
/**
	Creates new server, this doesn't start the server, but merely creates a new one, messages received are passed to the channel.
*/
func MakeServiceServer(msgChannel chan ServiceMsg, bitverseAddr string) *ServiceServer {
	server := new(ServiceServer)
    server.bitverse = bitverseAddr
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


type RequestIpInput struct {
	RequestLocal bool
}
type RequestIpOutput struct {
	IP string
}

//------------------------------------------------------------------------------
/**
	RPC-complaint call which sends back the IP of the service server
*/
func (server *ServiceServer) RequestIp(input *RequestIpInput, output *string) error {
    var reply RequestIpOutput 

	if input.RequestLocal {
		inter, err := net.InterfaceByName("eth0")
		
		var addrs []net.Addr
		addrs, err = inter.Addrs()
		if err != nil {
			return err
		}
		
		for _, addr := range addrs {
			switch ip := addr.(type) {
			case *net.IPNet:
				if ip.IP.DefaultMask() != nil {
					reply.IP = ip.IP.String()
				}
			}
		}
	} else {
		resp, err := http.Get("http://" + server.bitverse + "/globalip")
		if err != nil {
			return err
		}
		bytes := make([]byte, 255)
		num, err := resp.Body.Read(bytes)
		bytes = bytes[:num]
		defer resp.Body.Close()

		reply.IP = string(bytes)
	}

    // marshal and send
	data, err := json.Marshal(reply)
	if err != nil {
		return err
	}
	*output = string(data)

	return nil
}
