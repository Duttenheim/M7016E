package docker

import(
	"errors"
	"fmt"
)

//------------------------------------------------------------------------------
/**
	@type DockerInterface
	
	The docker interface exposes certain functions using RPC.
	All functions which uses a DockerInterface should follow the RPC calling model.
	The types below are the types which should be used to send/receive data from the functions called through RPC.
*/
type DockerInterface struct {
	repo* DockerRepo
}

//------------------------------------------------------------------------------
/**
	@type DownloadArgs
	@type DownloadReply
	
	Argument and reply types which are used with the Download() function.
*/
type DownloadArgs struct {
	Name string
}

type DownloadReply struct {
	Content string
}

//------------------------------------------------------------------------------
/**
	RPC function which accepts a string as download argument, and replies with a string content.
*/
func (facade* DockerInterface) Download(args* DownloadArgs, reply* DownloadReply) error { 
	
	// just check the repo is setup
	if (facade.repo == nil) {
		return errors.New("Remote repository not yet setup!")
	}
	
	// find image
	image := facade.repo.Find(args.Name)
	
	// set data
	if (image != nil) {
		reply.Content = image.content
		return nil
	}
	
	// return error saying the image doesn't exist
	msg := fmt.Sprintf("Image with name '%s' does not exist!", args.Name)
	return errors.New(msg)
}

//------------------------------------------------------------------------------
/**
	@type FindImageArgs
	@type FindImageReply
*/
type SearchArgs struct {
	Name string
}

type SearchReply struct {
	Exists bool
}

//------------------------------------------------------------------------------
/**
	RPC function which accepts a string as search argument, and replies with a boolean whether or not this image exists or not.
*/
func (facade* DockerInterface) Search(args* SearchArgs, reply* SearchReply) error {

	// just check the repo is setup
	if (facade.repo == nil) {
		return errors.New("Remote repository not yet setup!")
	}
	
	// see if image exists
	exists := facade.repo.Exists(args.Name)
	reply.Exists = exists
	
	// this function cannot cause an error
	return nil
}