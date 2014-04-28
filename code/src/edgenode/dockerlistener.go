package main

import (
	"fmt"
	"bitverse"
	"protocol"
	"strconv"
	"strings"
	"os"
	"bufio"
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

type DockerListArgs struct {
	ShowAll bool
}

type ImageArgs struct {
	ID string
	Name string
	Repository string
}

type RpcOutput struct {
	Content string
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
		output.Content += fmt.Sprintf("Error %s", err)	
	}else{
		output.Content += fmt.Sprintf("Container %s Created", container.ID)
	}
	return nil

}

func (obj* EdgeNodeHandler) StartContainer(args* ContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StartContainer(args.ID, nil)
	output.Content = ""
	if err != nil {
		output.Content += fmt.Sprintf("Error %s", err)	
	}else{
		output.Content += fmt.Sprintf("Container %s Started", args.ID)
	}
	return nil
	
}

func (obj* EdgeNodeHandler) StopContainer(args* ContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	err := client.StopContainer(args.ID, 3)
	output.Content = ""
	if err != nil {
		output.Content += fmt.Sprintf("Error %s", err)	
	} else {
		output.Content += fmt.Sprintf("Stopped container %s", args.ID)
	}
	return nil
}

func (obj* EdgeNodeHandler) KillContainer(args* ContainerArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	output.Content = ""
	err := client.KillContainer(docker.KillContainerOptions{ID: args.ID})
	if err != nil {
		output.Content += fmt.Sprintf("Error %s", err)	
	}else{
		output.Content += fmt.Sprintf("Container %s was killed", args.ID)
	}
	return nil
}

func (obj* EdgeNodeHandler) ListContainers(args* DockerListArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, err := client.ListContainers(docker.ListContainersOptions{All: args.ShowAll})
	output.Content = ""
	if err != nil {
		output.Content += fmt.Sprintf("Error %s", err)
	} else {
		output.Content += fmt.Sprintf("Containers found \n")
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
		output.Content += fmt.Sprintf("Error %s", err)	
	} else {
		output.Content += fmt.Sprintf("Pulled container")
	}
	return nil
}


func (obj* EdgeNodeHandler) ListImages(args* DockerListArgs, output* RpcOutput) error {
	endpoint := "unix:///var/run/docker.sock"
        client, _ := docker.NewClient(endpoint)
        imgs, err := client.ListImages(args.ShowAll)
	if err != nil {
		output.Content += fmt.Sprintf("Error %s", err)
	} else {
		output.Content += fmt.Sprintf("Containers found \n")
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

func (observer* TestObserver) OnSiblingJoined(node* bitverse.EdgeNode, id string) {
	//fmt.Printf("Sibling '%s' joined supernode!\n", id)
}

func (observer* TestObserver) OnSiblingLeft(node* bitverse.EdgeNode, id string) {
	//fmt.Printf("Sibling '%s' left supernode!\n", id)
}

func (observer* TestObserver) OnSiblingHeartbeat(node* bitverse.EdgeNode, id string) {
	//fmt.Printf("Sibling '%s' is still alive!\n", id)
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
	
	service, serviceError := node.CreateMsgService(secret, "TestMessagingService", msgObserver)
	if (serviceError != nil) {
		panic(serviceError)
	}

	// wait for input
	go func() {
		bio := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf(">>")
			line, _ := bio.ReadString('\n')
			components := strings.Split(line, " ")
			for index, comp := range components {
				components[index] = strings.Trim(comp, "\n")
			}
			if (len(components) > 0) {
				if (components[0] == "show") {
					reply = make(chan bool)
					observer.superNode.SendChildrenRequest()
					<-reply
					for index, value := range siblings {
						fmt.Printf("%d. %s\n", index, value)
					}
				} else if (components[0] == "run") {
					if (len(components) != 2) {
						fmt.Printf("'run' command must be supplied with a sibling node index!\n")
						continue
					}
					index, err := strconv.Atoi(components[1])
					if (err != nil) {
						fmt.Printf("Must supply a valid number to 'run'\n")
						continue
					}
					if (index < len(siblings)) {
						args := new(ContainerArgs)
						reply := new(RpcOutput)
						err := protocol.RpcInvokeSync("EdgeNodeHandler.Test", args, reply, siblings[index], service)
						if (err != nil) {
							fmt.Println(err)
							continue
						}
	
						fmt.Println(reply.Content)
					} else {
						fmt.Printf("Unknown sibling index '%d'\n", index)
					}
				} else {
					fmt.Println("Unknown command: " + components[0])
				}
			}
		}
	}()

	// connect node and wait until done (which is forever)
	go node.Connect("localhost:2020")
	<- done	
}
