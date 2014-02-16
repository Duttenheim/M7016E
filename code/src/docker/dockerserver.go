package docker

import(
	"net"
	"net/rpc"
	"net/http"
	"fmt"
)

type DockerServer struct {}


//------------------------------------------------------------------------------
/**
	Opens a docker server on a given port and address.
	This function should be launched using a go function since it will be blocking.
	We open a TCP connection on the local machine using a provided port.
	We then register a DockerInterface object to be called using RPC.
*/
func (server* DockerServer) Open(port uint16, repoPath string) error {

	// convert port to string...
	portString := fmt.Sprintf("localhost:%d", port)
	
	// create a repo and open it
	repo := new(DockerRepo)
	repo.Load(repoPath)
	
	// create a docker interface object
	facade := new(DockerInterface)
	facade.repo = repo
	
	// register type with RPC
	rpc.Register(facade)
	
	// open socket on port
	socket, err := net.Listen("tcp", portString)
	go rpc.Accept(socket)
	
	// handle HTTP request
	rpc.HandleHTTP()
	
	// if this failed, return error
	if (err != nil) {
		return err
	}	
		
	// start serving HTTP connections
	err = http.Serve(socket, nil)
	if (err != nil) {
		return err
	}
	
	return nil
}
