var GuiChildrenReceivedCallback = function(children)
{
    var body = document.getElementById("children");
    body.innerHTML = "<h1>Children</h1>";

    var node = this;

    // create table and body
    var table = document.createElement("table");
    var tableBody = document.createElement("tbody");
    for (var i = 0; i < children.length; i++)
    {
        var nodeString = children[i];

        var row = document.createElement("tr");
        var cell1 = document.createElement("td");
        var cell1Contents = document.createElement("button");
        var buttonContents = document.createTextNode("Node: " + nodeString);
        cell1Contents.appendChild(buttonContents);
        cell1.appendChild(cell1Contents);

        cell1Contents.onclick = function() 
        {   
            SetupNodeUI(node, nodeString);
            //var args = new ContainerArgs();            
            //node.CallRPCFunction("EdgeNodeHandler.StartContainer", args, nodeString);
        }

        row.appendChild(cell1);
        tableBody.appendChild(row);
    }   
    table.appendChild(tableBody);   
    body.appendChild(table);
}

var GuiMessageReceivedCallback = function(reply)
{
    alert(reply.Content);
}


function SetupNodeUI(node, nodeName)
{
    $("#node").load("node.html", function() 
    {
        $("#node_name").html(nodeName);

        var startButton = $("#start").get(0);
        var stopButton = $("#stop").get(0);
        var createButton = $("#create_container").get(0);

        var startStopContainer = $("#start_stop_containername").get(0);
        var createContainer = $("#create_containername").get(0);
        var imageName = $("#image_name").get(0);

        startButton.onclick = function()
        {
            var args = new ContainerArgs();
            args.ID = startStopContainer.value;
            node.CallRPCFunction("EdgeNodeHandler.StartContainer", args, nodeName);            
        }

        stopButton.onclick = function()
        {
            var args = new ContainerArgs();
            args.ID = startStopContainer.value;
            node.CallRPCFunction("EdgeNodeHandler.StopContainer", args, nodeName);            
        }

        createButton.onclick = function()
        {
            var args = new CreateContainerArgs();
            args.ContainerName = createContainer.value;
            args.ImageName = imageName.value;
            node.CallRPCFunction("EdgeNodeHandler.CreateContainer", args, nodeName);                        
        }

    });
}
