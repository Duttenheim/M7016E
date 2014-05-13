package main

import (
	"fmt"
	"bitverse"
	"protocol"
    "flag"
	"github.com/fsouza/go-dockerclient"
)

type EdgeNodeHandler struct {}
type ContainerArgs struct {
	ID string
}

type CreateContainerArgs struct {
	ContainerName string
	ImageName string
}

type RemoveContainerArgs struct {
    // The ID of the container.
    ID  string

    // A flag that indicates whether Docker should remove the volumes
    // associated to the container.
    RemoveVolumes bool

    // A flag that indicates whether Docker should remove the container
    // even if it is currently running.
    Force bool
}

//Replycodes
const (
    ErrorCode		= 0
    CreateContainer 	= 1
    StartContainer      = 2
    StopContainer     	= 3
    KillContainer	= 4
    RestartContainer    = 5
    RemoveContainer	= 6
    ListContainers      = 7
    PullImage		= 8
    RemoveImage		= 9
    ListImages		= 10
)

type DockerListArgs struct {
	ShowAll bool
}

type ImageArgs struct {
	ID string
	Name string
	Repository string
}

type RemoveImageArgs struct {
	Name string
}

//Reply codes are used to define what method that it returns from
// 
type RpcOutput struct {
	Content string
	ReplyCode int
}

func (obj* EdgeNodeHandler) Test(input* ContainerArgs, output* RpcOutput) error {
	output.Content = "It works!"
	return nil
}

func (obj* EdgeNodeHandler) CreateContainer(args* CreateContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	config := docker.Config{Image: args.ImageName}
	createArgs := docker.CreateContainerOptions{Name: args.ContainerName, Config: &config}	
	container, err := client.CreateContainer(createArgs)
	output.Content = ""
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	}else{
		output.Content += fmt.Sprintf("Container %s created successfully", container.ID)
		output.ReplyCode = CreateContainer
	}
	return nil

}

func (obj* EdgeNodeHandler) StartContainer(args* ContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StartContainer(args.ID, nil)
	output.Content = ""
    fmt.Printf("Container id:%s\n", args.ID)
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	}else{
		output.Content += fmt.Sprintf("Container %s started", args.ID)
		output.ReplyCode = StartContainer
	}
	return nil
	
}

func (obj* EdgeNodeHandler) StopContainer(args* ContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StopContainer(args.ID, 3)
	output.Content = ""
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	} else {
		output.Content += fmt.Sprintf("Stopped container %s", args.ID)
		output.ReplyCode = StopContainer
	}
	return nil
}

func (obj* EdgeNodeHandler) KillContainer(args* ContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	output.Content = ""
	err := client.KillContainer(docker.KillContainerOptions{ID: args.ID})
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	}else{
		output.Content += fmt.Sprintf("Container %s was killed", args.ID)
		output.ReplyCode = KillContainer
	}
	return nil
}


func (obj* EdgeNodeHandler) RestartContainer(args* ContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	output.Content = ""
	err := client.RestartContainer(args.ID, 500)
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	}else{
		output.Content += fmt.Sprintf("Container %s is restarting", args.ID)
		output.ReplyCode = RestartContainer
	}
	return nil
}

func (obj* EdgeNodeHandler) RemoveContainer(args* RemoveContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	output.Content = ""
	err := client.RemoveContainer(docker.RemoveContainerOptions{ID: args.ID, RemoveVolumes: args.RemoveVolumes, Force: args.Force})
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	}else{
		output.Content += fmt.Sprintf("Container %s was removed", args.ID)
		output.ReplyCode = RemoveContainer
	}
	return nil
}



func (obj* EdgeNodeHandler) ListContainers(args* DockerListArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, err := client.ListContainers(docker.ListContainersOptions{All: args.ShowAll})
	output.Content = ""
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	} else {
		output.Content += fmt.Sprintf("Containers found \n")
		output.ReplyCode = ListContainers
	}
	for _, img := range imgs {
		output.Content += fmt.Sprintf("ID: %s \n", img.ID)
		output.Content += fmt.Sprintf("Created: %d \n", img.Created)
	}
	return nil
}

func (obj* EdgeNodeHandler) PullImage(args* ImageArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	output.Content = ""
	err := client.PullImage(docker.PullImageOptions{Repository: args.Repository}, docker.AuthConfiguration{})
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	} else {
		output.Content += fmt.Sprintf("Pulled container")
		output.ReplyCode = PullImage
	}
	return nil
}

func (obj* EdgeNodeHandler) RemoveImage(args* RemoveImageArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	output.Content = ""
	err := client.RemoveImage(args.Name)
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	} else {
		output.Content += fmt.Sprintf("Removed %s Image", args.Name)
		output.ReplyCode = RemoveImage
	}
	return nil
}



func (obj* EdgeNodeHandler) ListImages(args* DockerListArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
        client, _ := docker.NewClient(endpoint)
        imgs, err := client.ListImages(args.ShowAll)
	if err != nil {
		output.Content += fmt.Sprintf("ERROR: %s", err)
		output.ReplyCode = ErrorCode
	} else {
		output.Content += fmt.Sprintf("Containers found \n")
		output.ReplyCode = ListImages
	}
        for _, img := range imgs {
                output.Content += fmt.Sprintf("ID: ", img.ID)
                output.Content += fmt.Sprintf("RepoTags: ", img.RepoTags)
                output.Content += fmt.Sprintf("Created: ", img.Created)
                output.Content += fmt.Sprintf("Size: ", img.Size)
                output.Content += fmt.Sprintf("VirtualSize: ", img.VirtualSize)
                output.Content += fmt.Sprintf("ParentId: ", img.ParentId)
                output.Content += fmt.Sprintf("Repository: ", img.Repository)
        }
	return nil

}

type TestObserver struct {
	superNode* bitverse.RemoteNode
}

func (observer* TestObserver) OnSiblingHeartbeat(node* bitverse.EdgeNode, id string) {
    // empty
}

func (observer* TestObserver) OnSiblingJoined(node* bitverse.EdgeNode, id string) {
    // empty
}

func (observer* TestObserver) OnSiblingLeft(node* bitverse.EdgeNode, id string) {
    // empty
}

var siblings []string
var reply chan bool

func (observer* TestObserver) OnChildrenReply(node* bitverse.EdgeNode, id string, children[] string) {
	siblings = children
	reply <- true
}

func (observer* TestObserver) OnConnected(localNode* bitverse.EdgeNode, remoteNode* bitverse.RemoteNode) {
	observer.superNode = remoteNode
}

var secret string = "3e606ad97e0a738d8da4c4c74e8cd1f1f2e016c74d85f17ac2fc3b5dab4ed6c4"

func main() {
	// create transport and channel
	transport := bitverse.MakeWSTransport()
	var done chan int
	
	// make edge node
	observer := new(TestObserver)
	node, done := bitverse.MakeEdgeNode(transport, observer)

	// make message observer
	msgObserver := new(protocol.RpcMessageObserver)
	msgObserver.Open()

	// make a docker RPC-callable object
	test := new(EdgeNodeHandler)
	msgObserver.Register(test)
	
	_, serviceError := node.CreateMsgService(secret, "RPCMessageService", msgObserver)
	if (serviceError != nil) {
		panic(serviceError)
	}

	// connect node and wait until done (which is forever)
    addr := flag.String("address", "localhost", "Bitverse server IP-address")
    flag.Parse()

	go node.Connect(*addr + ":2020")
	<- done	
}
