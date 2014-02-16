package docker

import( 
	"testing"
	"net/rpc"
	"log"
	"fmt"
)

//------------------------------------------------------------------------------
/**
	Test function for server
*/
func TestDockerServer(t* testing.T) {
	t.Parallel()
	server := DockerServer{}
	server.Open(7070, "./images")
}

//------------------------------------------------------------------------------
/**
	Test function for client
*/
func TestDockerClient(t* testing.T) {
	t.Parallel()
	client, err := rpc.DialHTTP("tcp", "localhost:7070")
	if (err != nil) {
		log.Fatal("dialing:", err)
	}
	
	args := &SearchArgs{"img.dkf"}	
	reply := new(SearchReply)
		
	err = client.Call("DockerInterface.Search", args, reply)
	if (err == nil) {
		fmt.Printf("Supposedly called 'DockerInterface.Search'!\n")
	} else {
		log.Fatal(err)
	}
	fmt.Printf("\n\nContents:\n")
	fmt.Printf("%t\n", reply.Exists)
}