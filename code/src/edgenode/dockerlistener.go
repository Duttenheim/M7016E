package main

import (
	"fmt"
	"bitverse"
	"protocol"
	"encoding/json"
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
	Registry string
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

type Container struct {
	ID string
	Image string
	Created string
	Status string
}

type ContainerCollection struct {
	ReplyCode int
	Containers []Container
}

type Image struct {
	ID string
	Created string
	Size string
	VirtualSize string
	RepoTags []string
}

type ImageCollection struct {
	ReplyCode int
	Images []Image
}

func (obj* EdgeNodeHandler) Test(input* ContainerArgs, output* RpcOutput) error {
	output.Content = "It works!"
	return nil
}

func (obj* EdgeNodeHandler) CreateContainer(args* CreateContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	config := docker.Config{Image: args.ImageName}
	createArgs := docker.CreateContainerOptions{Name: args.ContainerName, Config: &config}	
	container, err := client.CreateContainer(createArgs)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s created successfully", container.ID)
		rpcOutput.ReplyCode = CreateContainer
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil

}

func (obj* EdgeNodeHandler) StartContainer(args* ContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StartContainer(args.ID, nil)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s started", args.ID)
		rpcOutput.ReplyCode = StartContainer
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
	
}

func (obj* EdgeNodeHandler) StopContainer(args* ContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StopContainer(args.ID, 3)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	} else {
		rpcOutput.Content += fmt.Sprintf("Stopped container %s", args.ID)
		rpcOutput.ReplyCode = StopContainer
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}

func (obj* EdgeNodeHandler) KillContainer(args* ContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.KillContainer(docker.KillContainerOptions{ID: args.ID})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s was killed", args.ID)
		rpcOutput.ReplyCode = KillContainer
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}


func (obj* EdgeNodeHandler) RestartContainer(args* ContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.RestartContainer(args.ID, 500)
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s is restarting", args.ID)
		rpcOutput.ReplyCode = RestartContainer
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}

func (obj* EdgeNodeHandler) RemoveContainer(args* RemoveContainerArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.RemoveContainer(docker.RemoveContainerOptions{ID: args.ID, RemoveVolumes: args.RemoveVolumes, Force: args.Force})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	}else{
		rpcOutput.Content += fmt.Sprintf("Container %s was removed", args.ID)
		rpcOutput.ReplyCode = RemoveContainer
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}



func (obj* EdgeNodeHandler) ListContainers(args* DockerListArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, err := client.ListContainers(docker.ListContainersOptions{All: args.ShowAll})
	rpcOutput := RpcOutput{}
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
		b, _ := json.Marshal(rpcOutput)
		*output += fmt.Sprintf(string(b))
	} else {
		
		list := []Container{}
		for _, img := range imgs {
			cont := Container{}
			cont.ID += fmt.Sprintf(img.ID)
			cont.Image += fmt.Sprintf(img.Image)
			cont.Created += fmt.Sprintf("%v",img.Created)
			cont.Status += fmt.Sprintf(img.Status)
			list = append(list,cont)
		}
		containerColl := ContainerCollection{ReplyCode: ListContainers, Containers: list}
		b, _ := json.Marshal(containerColl)
		*output += fmt.Sprintf(string(b))
	}
	
	return nil
}

func (obj* EdgeNodeHandler) PullImage(args* ImageArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.PullImage(docker.PullImageOptions{Repository: args.Repository, Registry: args.Registry}, docker.AuthConfiguration{})
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	} else {
		rpcOutput.Content += fmt.Sprintf("Pulled Image: " + args.Registry+"/"+args.Repository)
		rpcOutput.ReplyCode = PullImage
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}

func (obj* EdgeNodeHandler) RemoveImage(args* RemoveImageArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	rpcOutput := RpcOutput{}
	rpcOutput.Content = ""
	err := client.RemoveImage(args.Name)
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
	} else {
		rpcOutput.Content += fmt.Sprintf("Removed %s Image", args.Name)
		rpcOutput.ReplyCode = RemoveImage
	}
	b, _ := json.Marshal(rpcOutput)
	*output += fmt.Sprintf(string(b))
	return nil
}



func (obj* EdgeNodeHandler) ListImages(args* DockerListArgs, output* string) error {
	endpoint := "unix:///var/run/docker.sock"
        client, _ := docker.NewClient(endpoint)
        imgs, err := client.ListImages(args.ShowAll)
	rpcOutput := RpcOutput{}
	if err != nil {
		rpcOutput.Content += fmt.Sprintf("ERROR: %s", err)
		rpcOutput.ReplyCode = ErrorCode
		b, _ := json.Marshal(rpcOutput)
		*output += fmt.Sprintf(string(b))
	} else {
		list := []Image{}
		for _, img := range imgs {
			image := Image{}
		        image.ID += fmt.Sprintf(img.ID)
		        image.RepoTags = img.RepoTags
		        image.Created += fmt.Sprintf("%v", img.Created)
		        image.Size += fmt.Sprintf("%d",img.Size)
		        image.VirtualSize += fmt.Sprintf("%d", img.VirtualSize)
			list = append(list, image)
        	}
		imageColl := ImageCollection{ReplyCode: ListImages, Images: list}
		b, _ := json.Marshal(imageColl)
		*output += fmt.Sprintf(string(b))
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