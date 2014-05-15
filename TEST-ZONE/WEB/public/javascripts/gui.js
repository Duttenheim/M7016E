/**
 * New node file
 */

function CreateContainerList(){

	var tbody = document.getElementById("containers_body");
	var header = document.getElementById("avail_cont_head");
	$("#containers_body").empty();
	header.innerHTML = "Containers available: 10";
	for (var i = 0; i < 10; i++)
    {
		populateContainerList(tbody, i+1);
    }
}

function populateContainerList(tbody, nr){
	var row = document.createElement("tr");
    var cell1 = document.createElement("td");
    var cell2 = document.createElement("td");
    var cell3 = document.createElement("td");
    var nrtab = document.createTextNode(nr);
    var nametab = document.createTextNode("Container " + nr);
    
    var startButton = document.createElement("Button");
    var buttonText = document.createTextNode("Start");
    startButton.className = "btn btn-success"
    startButton.appendChild(buttonText);
    
    var stopButton = document.createElement("Button");
    var buttonText = document.createTextNode("Stop");
    stopButton.className = "btn btn-warning"
    stopButton.appendChild(buttonText);
    
    var killButton = document.createElement("Button");
    var buttonText = document.createTextNode("Kill");
    killButton.className = "btn btn-danger"
    killButton.appendChild(buttonText);
    
    
    var deleteButton = document.createElement("Button");
    var buttonText = document.createTextNode("Delete");
    deleteButton.className = "btn btn-danger"
    deleteButton.appendChild(buttonText);
    
    cell1.appendChild(nrtab);
    cell2.appendChild(nametab);
    cell3.appendChild(startButton);
    cell3.appendChild(stopButton);
    cell3.appendChild(killButton);
    cell3.appendChild(deleteButton);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    tbody.appendChild(row);
}

function CreateImageList(){
	var tbody = document.getElementById("images_body");
	var header = document.getElementById("avail_img_head");
	$("#images_body").empty();
	header.innerHTML = "Images available: 10";
	for (var i = 0; i < 10; i++)
    {
		populateImageList(tbody, i+1);
    }
}

function populateImageList(tbody, nr){
	var row = document.createElement("tr");
    var cell1 = document.createElement("td");
    var cell2 = document.createElement("td");
    var cell3 = document.createElement("td");
    var nrtab = document.createTextNode(nr);
    var nametab = document.createTextNode("Image " + nr);

    var createButton = document.createElement("Button");
    var buttonText = document.createTextNode("Create container");
    createButton.className = "btn btn-success"
    createButton.appendChild(buttonText);
    
    cell1.appendChild(nrtab);
    cell2.appendChild(nametab);
    cell3.appendChild(createButton);
    row.appendChild(cell1);
    row.appendChild(cell2);
    row.appendChild(cell3);
    tbody.appendChild(row);
}