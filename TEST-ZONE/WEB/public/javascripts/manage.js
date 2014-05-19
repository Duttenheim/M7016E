/**
 * New node file
 */

function CreateImageList(json, node){
	var tbody = document.getElementById("images_body");
	var header = document.getElementById("avail_img_head");
	var edgeNode = document.getElementById("edgeNode_id").innerHTML
	$("#images_body").empty();
		
	var count = 0;
	for (var i = 0; i < json.Images.length; i++)
    {
		var repo = json.Images[i].RepoTags[0].split(":");
		
		if(repo[0] != "<none>"){
			count++;
			populateImageList(tbody, i+1, json.Images[i], edgeNode, repo[0]);
		}
    }
	header.innerHTML = "Images available: " + count;
}

function populateImageList(tbody, nr, image, edgeNode, imageName){
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
    createButton.onclick = function()
    {
    	CreateContainerPopup(row, nr, edgeNode, imageName);
    }
    
    cell1.appendChild(nrtab);
    cell2.appendChild(idTab);
    cell3.appendChild(repoTab);
    cell4.appendChild(createdTab);
    cell5.appendChild(sizeTab);
    cell6.appendChild(createButton);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    row.appendChild(cell4);
    row.appendChild(cell5);
    row.appendChild(cell6);
    tbody.appendChild(row);
}

function CreateContainerList(json, node){
	console.log("Created container list");
	var tbody = document.getElementById("containers_body");
	var header = document.getElementById("avail_cont_head");
	var edgeNode = document.getElementById("edgeNode_id").innerHTML
	$("#containers_body").empty();
	header.innerHTML = "Containers available: "+ json.Containers.length;
	for (var i = 0; i < json.Containers.length; i++)
    {
		populateContainerList(tbody, i+1, json.Containers[i], node, edgeNode, json);
    }
}

function populateContainerList(tbody, nr, container, node, edgeNode, json){
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
        var args = new RemoveContainerArgs();
        args.ID = container.ID;
        args.RemoveVolumes = true;
        args.Force = true;
        node.CallRPCFunction("EdgeNodeHandler.RemoveContainer", args, edgeNode);
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
    
    tbody.appendChild(row);
}

function CreateContainerPopup(row, nr, edgeNode, imageName){
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
        args.ImageName = imageName;
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

var NodeReceiveCallback = function(reply)
{
	var node = this;
    var json = eval ("(" + reply + ")");
    
    if(json.ReplyCode == 10){
    	CreateImageList(json, node)
    } else if(json.ReplyCode == 7) {
    	CreateContainerList(json, node)
    } else {
    	alert(json.Content)
    }
    
}