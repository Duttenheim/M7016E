var GuiChildrenReceivedCallback = function(children)
{
    var body = document.getElementById("children");
    body.innerHTML = "";

    var node = this;

    // create table and body
    var table = document.createElement("table");
    var tableBody = document.createElement("tbody");
    for (var i = 0; i < children.length; i++)
    {
        var row = document.createElement("tr");
        var cell1 = document.createElement("td");
        var cell1Contents = document.createTextNode("Node id: " + children[i]);
        cell1.appendChild(cell1Contents);

        var cell2 = document.createElement("td");
        var cell2Contents = document.createElement("button");

        cell2Contents.onclick = function() 
        {   
            var args = new ContainerArgs();            
            node.CallRPCFunction("EdgeNodeHandler.StartContainer", args, children[i]);
        }

        var buttonContents = document.createTextNode("Start container");
        cell2Contents.appendChild(buttonContents);
        cell2.appendChild(cell2Contents);

        row.appendChild(cell1);
        row.appendChild(cell2);

        tableBody.appendChild(row);
    }   
    table.setAttribute("border", "1");
    table.appendChild(tableBody);   
    body.appendChild(table);
}

