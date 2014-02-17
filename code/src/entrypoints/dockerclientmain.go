package main

import( 
	"docker"
	"net/rpc"
	"flag"
	"fmt"
	"log"
)

func main() {

	// parse command line arguments	
	host := flag.String("host", "localhost", "Host address of docker repo server")
	port := flag.Int("port", 7070, "Host port")
	image := flag.String("image", "", "Image to search for")
	printImage := flag.Bool("print", false, "True will also fetch the image and print out its contents")
	flag.Parse()
	
	// open connection
	address := fmt.Sprintf("%s:%d", *host, *port)
	client, err := rpc.DialHTTP("tcp", address)
	if (err != nil) {
		log.Fatal(err)
	} else {
		fmt.Printf("Connected to %s\n", address)
	}
	
	args := &docker.SearchArgs{*image}	
	reply := new(docker.SearchReply)
		
	err = client.Call("DockerInterface.Search", args, reply)
	if (err != nil) {
		log.Fatal(err)
	}
	
	if (reply.Exists) {
		fmt.Printf("Image '%s' found!\n", *image)
	}
	
	if (reply.Exists && *printImage) {
		
		args := &docker.DownloadArgs{*image}
		reply := new(docker.DownloadReply)
		err = client.Call("DockerInterface.Download", args, reply)
		if (err != nil) {
			log.Fatal(err)
		}
		
		fmt.Printf("\n\n--------Contents--------\n")
		fmt.Printf("%t\n", reply.Content)
	}
	
	
}
