/**
 * New node file
 */

function CreateImageList(json, node){
	var edgeNode = document.getElementById("edgeNode_id").innerHTML
	var table = document.getElementById('images_table')
	$("#images_body").empty();
	var count = 0;
	for (var i = 0; i < json.Images.length; i++)
    {
		var repo = json.Images[i].RepoTags[0].split(":");
		if(repo[0] != "<none>"){
			populateImageList(count+1, json.Images[i], edgeNode, repo, node, table);
			count++;
		}
    }
	setImagesHeaderText(count);
}

function populateImageList(nr, image, edgeNode, imageName, node, table){
	var row = document.createElement("tr");
    var cell1 = document.createElement("td");
    var cell2 = document.createElement("td");
    var cell3 = document.createElement("td");
    var cell4 = document.createElement("td");
    var cell5 = document.createElement("td");
    var cell6 = document.createElement("td");
    
    var nrtab = document.createTextNode(nr);
    var repoTab = document.createTextNode(image.RepoTags);
    var idTab = document.createTextNode(image.ID);
    var createdTab = document.createTextNode(image.Created)
    var sizeTab = document.createTextNode(image.Size)
    

    var createButton = document.createElement("a");
    createButton.href = "#test_modal_"+nr;
    createButton.setAttribute('data-toggle', 'modal');
    var buttonText = document.createTextNode("Create container");
    createButton.className = "btn btn-success"
    createButton.appendChild(buttonText);
    createButton.setAttribute('style', 'margin: 0px 2px 0px 0px');
    createButton.onclick = function()
    {
    	CreateContainerPopup(row, nr, edgeNode, image.ID, imageName);
    }
    
    var deleteButton = document.createElement("Button");
    var buttonText2 = document.createTextNode("Remove");
    deleteButton.className = "btn btn-danger"
    deleteButton.appendChild(buttonText2);
    deleteButton.onclick = function()
    {
        if (confirm("Do you really want to delete " + imageName + "?") == true) {
        	var args = new RemoveImageArgs();
	        args.name = image.ID;
	        node.CallRPCFunction("EdgeNodeHandler.RemoveImage", args, edgeNode);
	        table.deleteRow(nr);
        }
    }
    
    cell1.appendChild(nrtab);
    cell2.appendChild(idTab);
    cell3.appendChild(repoTab);
    cell4.appendChild(createdTab);
    cell5.appendChild(sizeTab);
    cell6.appendChild(createButton);
    cell6.appendChild(deleteButton);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    row.appendChild(cell4);
    row.appendChild(cell5);
    row.appendChild(cell6);
    table.tBodies.item("images_body").appendChild(row);
}

function CreateContainerList(json, node){
	var table = document.getElementById("container_table");
	var edgeNode = document.getElementById("edgeNode_id").innerHTML
	setContainerHeaderText(json.Containers.length);
	$("#containers_body").empty();
	for (var i = 0; i < json.Containers.length; i++)
    {
		populateContainerList(i+1, json.Containers[i], node, edgeNode, json, table);
    }
}

function populateContainerList(nr, container, node, edgeNode, json, table){
	var row = document.createElement("tr");
    var cell1 = document.createElement("td");
    var cell2 = document.createElement("td");
    var cell3 = document.createElement("td");
    var cell4 = document.createElement("td");
    var cell5 = document.createElement("td");
    var cell6 = document.createElement("td");
    
    var nrtab = document.createTextNode(nr);
    var idTab = document.createTextNode(container.ID);
    var imageTab = document.createTextNode(container.Image);
    var createdTab = document.createTextNode(container.Created);
    var statusTab = document.createTextNode(container.Status);
     
    //Only show start button if container is not running
    if(container.Status == "" || container.Status.indexOf("Exit") > -1 ) {
	    var startButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Start");
	    startButton.className = "btn btn-success"
	    startButton.appendChild(buttonText);
	    startButton.setAttribute('style', 'margin: 0px 2px 0px 0px');
	    startButton.onclick = function()
	    {
	        var args = new ContainerArgs();
	        args.ID = container.ID;
	        node.CallRPCFunction("EdgeNodeHandler.StartContainer", args, edgeNode);
	    }
	    cell6.appendChild(startButton);
    } else {
	    var stopButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Stop");
	    stopButton.className = "btn btn-warning"
	    stopButton.appendChild(buttonText);
	    stopButton.setAttribute('style', 'margin: 0px 2px 0px 0px');
	    stopButton.onclick = function()
	    {
	        var args = new ContainerArgs();
	        args.ID = container.ID;
	        node.CallRPCFunction("EdgeNodeHandler.StopContainer", args, edgeNode);
	    }  
	    
	    var killButton = document.createElement("Button");
	    var buttonText = document.createTextNode("Kill");
	    killButton.className = "btn btn-danger"
	    killButton.appendChild(buttonText);
	    killButton.setAttribute('style', 'margin: 0px 2px 0px 0px');
	    killButton.onclick = function()
	    {
	        var args = new ContainerArgs();
	        args.ID = container.ID;
	        node.CallRPCFunction("EdgeNodeHandler.KillContainer", args, edgeNode);            
	    }

	    cell6.appendChild(stopButton);
	    cell6.appendChild(killButton);
    }
    
    var deleteButton = document.createElement("Button");
    var buttonText = document.createTextNode("Delete");
    deleteButton.className = "btn btn-danger"
    deleteButton.appendChild(buttonText);
    deleteButton.setAttribute('style', 'margin: 0px 2px 0px 0px');
    deleteButton.onclick = function()
    {
    	if (confirm("Do you really want to delete " + container.ID + "?") == true) {
	        var args = new RemoveContainerArgs();
	        args.ID = container.ID;
	        args.RemoveVolumes = true;
	        args.Force = true;
	        node.CallRPCFunction("EdgeNodeHandler.RemoveContainer", args, edgeNode);
	        table.deleteRow(nr);
    	}
    }
    
    cell1.appendChild(nrtab);
    cell2.appendChild(idTab);
    cell3.appendChild(imageTab);
    cell4.appendChild(createdTab);
    cell5.appendChild(statusTab);
    cell6.appendChild(deleteButton);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    row.appendChild(cell4);
    row.appendChild(cell5);
    row.appendChild(cell6);
    
    table.tBodies.item("containers_body").appendChild(row);
}

function CreateContainerPopup(row, nr, edgeNode, image_ID, imageName){
	var modalDiv = document.createElement("div");
	modalDiv.className = "modal fade";
	modalDiv.id = "test_modal_"+nr;
		
	var modalHeaderDiv = document.createElement("div");
	modalHeaderDiv.className = "modal-header";
	
	var closeX = document.createElement("a");
	var closeXtext = document.createTextNode('x');
	closeX.className = "close";
	closeX.setAttribute('data-dismiss', 'modal')
	closeX.appendChild(closeXtext);
	
	var header = document.createElement("h3");
	header.innerHTML = "Create container for Image: " + imageName;
	
	modalHeaderDiv.appendChild(closeX);
	modalHeaderDiv.appendChild(header);
	
	var modalBodyDiv = document.createElement("div"); 
	modalBodyDiv.className = "modal-body";
	
	var infoTextP = document.createElement("p");
    var infoText = document.createTextNode("Type in name for container");
    infoTextP.appendChild(infoText);
	var inputDiv = document.createElement("div");
	inputDiv.className = "input-group input-group-lg";
	var inputField = document.createElement("input");
	inputField.className = "form-control";
	inputDiv.appendChild(inputField)
	inputDiv.setAttribute('placeholder', 'Name of container');

	modalBodyDiv.appendChild(infoTextP);
	modalBodyDiv.appendChild(inputDiv);
		
	var modalFooterDiv = document.createElement("div"); 
	modalFooterDiv.className = "modal-footer";
	
	var CreateButton = document.createElement("a");
	CreateButton.className = "btn btn-primary";
    var createText = document.createTextNode("Create");
    CreateButton.appendChild(createText);
    
    CreateButton.onclick = function()
    {
    	var args = new CreateContainerArgs();
        args.ContainerName = inputField.value;
        args.ImageName = image_ID;
        node.CallRPCFunction("EdgeNodeHandler.CreateContainer", args, edgeNode);
    }
    CreateButton.setAttribute('data-dismiss', 'modal');
    
	var CloseButton = document.createElement("a");
    var closeText = document.createTextNode("Cancel");
    CloseButton.appendChild(closeText);
	CloseButton.className = "btn";
	CloseButton.setAttribute('data-dismiss', 'modal');
	
	modalFooterDiv.appendChild(CreateButton);
	modalFooterDiv.appendChild(CloseButton);
	
	modalDiv.appendChild(modalHeaderDiv);
	modalDiv.appendChild(modalBodyDiv);
	modalDiv.appendChild(modalFooterDiv);
	
	row.appendChild(modalDiv);
	
}

function PullImagePopup(edgeNode, node){
	var row = document.getElementById('images_body');
	var modalDiv = document.createElement("div");
	modalDiv.className = "modal fade";
	modalDiv.id = "pull_image_modal";
		
	var modalHeaderDiv = document.createElement("div");
	modalHeaderDiv.className = "modal-header";
	
	var closeX = document.createElement("a");
	var closeXtext = document.createTextNode('x');
	closeX.className = "close";
	closeX.setAttribute('data-dismiss', 'modal')
	closeX.appendChild(closeXtext);
	
	var header = document.createElement("h3");
	header.innerHTML = "Pull new Image";
	
	modalHeaderDiv.appendChild(closeX);
	modalHeaderDiv.appendChild(header);
	
	var modalBodyDiv = document.createElement("div"); 
	modalBodyDiv.className = "modal-body";
	
	/*var addressInfoP = document.createElement("p");
    var addressInfoText = document.createTextNode("Type in address of repo <IP:PORT>");
    addressInfoP.appendChild(addressInfoText);
	var addrInputDiv = document.createElement("div");
	addrInputDiv.className = "input-group input-group-lg";
	var addrInputField = document.createElement("input");
	addrInputField.className = "form-control";
	addrInputDiv.appendChild(addrInputField)
	addrInputDiv.setAttribute('placeholder', 'IP:PORT');

	modalBodyDiv.appendChild(addressInfoP);
	modalBodyDiv.appendChild(addrInputDiv);*/
	
	var repoInfoP = document.createElement("p");
    var repoInfoText = document.createTextNode("Type in name of repo");
    repoInfoP.appendChild(repoInfoText);
	var nameInputDiv = document.createElement("div");
	nameInputDiv.className = "input-group input-group-lg";
	var nameInputField = document.createElement("input");
	nameInputField.className = "form-control";
	nameInputDiv.appendChild(nameInputField)
	nameInputDiv.setAttribute('placeholder', 'Repo name');

	modalBodyDiv.appendChild(repoInfoP);
	modalBodyDiv.appendChild(nameInputDiv);
	
		
	var modalFooterDiv = document.createElement("div"); 
	modalFooterDiv.className = "modal-footer";
	
	var pullButton = document.createElement("a");
	pullButton.className = "btn btn-primary";
    var pullButtonText = document.createTextNode("Pull Image");
    pullButton.appendChild(pullButtonText);
    pullButton.setAttribute('data-dismiss', 'modal'); 
    pullButton.onclick = function()
    {
    	var args = new ImageArgs();
        args.Repository = nameInputField.value;
        //args.Registry = addrInputField.value;
        node.CallRPCFunction("EdgeNodeHandler.PullImage", args, edgeNode);                
    }
    
	var CloseButton = document.createElement("a");
    var closeText = document.createTextNode("Cancel");
    CloseButton.appendChild(closeText);
	CloseButton.className = "btn";
	CloseButton.setAttribute('data-dismiss', 'modal');
	
	modalFooterDiv.appendChild(pullButton);
	modalFooterDiv.appendChild(CloseButton);
	
	modalDiv.appendChild(modalHeaderDiv);
	modalDiv.appendChild(modalBodyDiv);
	modalDiv.appendChild(modalFooterDiv);
	
	row.appendChild(modalDiv);
	
}

function setImagesHeaderText(count){
	document.getElementById("avail_img_head").innerHTML = "Images available: " + count;
}
function setContainerHeaderText(count){
	document.getElementById("avail_cont_head").innerHTML = "Containers available: "+ count;
}

/*
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
 */
var NodeReceiveCallback = function(reply)
{
	var node = this;
    var edgeNode = document.getElementById('edgeNode_id').innerHTML
    var json = eval ("(" + reply + ")");
    
    if(json.ReplyCode == 10){ // ListImages
    	CreateImageList(json, node)
    } else if(json.ReplyCode == 7) { // ListContainers
    	CreateContainerList(json, node)
    } else if(json.ReplyCode > 0 && json.ReplyCode < 6){ // Create, Start, Stop, Kill, Restart containers
    	listContainers(edgeNode);
    	alert(json.Content);
    } else if(json.ReplyCode == 6){ // Remove container
    	setContainerHeaderText(document.getElementById("container_table").rows.length-1) //Use -1 since the header-row is counted here
    	alert(json.Content);
    } else if(json.ReplyCode == 9){ // Remove Image
    	setImagesHeaderText(document.getElementById("images_table").rows.length-1) //Use -1 since the header-row is counted here
    	alert(json.Content);
    } else if(json.ReplyCode == 8){ // Pull image
    	listImages(edgeNode);
    	alert(json.Content);
    } else {
    	alert(json.Content);
    }
    
}