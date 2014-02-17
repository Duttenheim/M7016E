package main

import( 
	"docker"
	"flag"	
)

func main() {
	server := docker.DockerServer{}
	path := flag.String("path", "", "Path to search for docker images")
	port := flag.Int("port", 7070, "Host port")
	flag.Parse()
	server.Open(uint16(*port), *path)
}