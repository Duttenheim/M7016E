var RPCReplyCode = 
{
	ErrorCode :  0,
    CreateContainer : 1,
    StartContainer : 2,
    StopContainer : 3,
    KillContainer : 4,
    RestartContainer : 5,
    RemoveContainer : 6,
    ListContainers : 7,
    PullImage : 8,
    RemoveImage : 9,
    ListImages : 10
}

function CreateContainerArgs()
{
    this.ContainerName = "";
    this.ImageName = "";
}

function ContainerArgs()
{
    this.ID = "";
}

function DockerResponse()
{
    this.Content = "";
    this.ReplyCode = 0
}