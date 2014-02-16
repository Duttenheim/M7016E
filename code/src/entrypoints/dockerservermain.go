package main

import( 
	"docker"
	"flag"	
)

func main() {
	server := docker.DockerServer{}
	var path string
	flag.StringVar(&path, "path", "", "Path to search for docker images")
	flag.Parse()
	server.Open(7070, path)
}