
var tags = [];
var values = [];

var defaultTags = ['Service', 'Application', 'Location', 'Name'];

var usedTags = [];

//------------------------------------------------------------------------------
/**
*/
function UpdateTable(serverList, tableElement)
{
	var table = document.getElementById(tableElement);
	var header = table.createTHead();
	var headerRow = header.insertRow(0);
	var headerCell1 = headerRow.insertCell(0);
	var headerCell2 = headerRow.insertCell(1);
	headerCell1.innerHTML = 'Tag';
	headerCell2.innerHTML = 'Value';
	
	var row = header.insertRow(-1);	
	NewSearchRow(row, table);
}

//------------------------------------------------------------------------------
/**
*/
function NewSearchRow(row, table)
{
	var cell1 = row.insertCell(0);
	var cell2 = row.insertCell(1);
	var cell3 = row.insertCell(2);
	var div = document.createElement("div");
	var span = document.createElement("span");
	span.className = "glyphicon form-control-feedback";
	div.className = "form-group has-success has-feedback";
	
	// input field 1
	var input1 = document.createElement("input");
	input1.id = 'tag';
	input1.type = 'text';
	input1.className = "form-control";
	input1.setAttribute('list', 'services');
	input1.placeholder = 'e.g Service';
	input1.onchange = function() 
	{
		var index = usedTags.indexOf(self.value);
		if (index == -1)
		{
			self.className = "form-control";
			span.className = "glyphicon form-control-feedback glyphicon-ok";
		}
		else
		{
			self.className = "form-control has-error";
			span.className = "glyphicon form-control-feedback glyphicon-remove";
		}
	}
	
	var datalist1 = document.createElement("datalist");
	datalist1.id = 'services';
	for (var tag in defaultTags)
	{
		var listElement = document.createElement("option");
		listElement.value = defaultTags[tag];
		datalist1.appendChild(listElement);
	}
	input1.list = datalist1;
	table.appendChild(datalist1);
	div.appendChild(input1);
	div.appendChild(span);
	cell1.appendChild(div);
	
	// input field 2
	var input2 = document.createElement("input");
	input2.id = 'value';
	input2.type = 'text';	
	input2.placeholder = 'e.g Docker, MyEdgeNode...'
	cell2.appendChild(input2);
	
	var buttonAdd = document.createElement("button");
	buttonAdd.className = "btn";
	buttonAdd.innerHTML = 'Add';
	buttonAdd.onclick = function() 
	{ 		
		// only add if tag isn't occupied
		var index = usedTags.indexOf(input1.value);
		if (index == -1)
		{
			usedTags.push(input1.value);
			SetRowEnabled(row, false);
			var newRow = table.insertRow(-1);
			NewSearchRow(newRow, table); 
		}
	}
	cell3.appendChild(buttonAdd);	
	
	var buttonRemove = document.createElement("button");
	buttonRemove.className = "btn";
	buttonRemove.innerHTML = "Remove";
	buttonRemove.onclick = function()
	{
		var index = defaultTags.indexOf(input1.value);
		if (index > -1)
		{
			usedTags.splice(index, 1);
		}
		
		RemoveSearchRow(row, table);
	}	
	cell3.appendChild(buttonRemove);
}

//------------------------------------------------------------------------------
/**
*/
function RemoveSearchRow(row, table)
{
	var i = row.rowIndex;
	table.deleteRow(i);
	
	// enable the previous row
	if (i > 0)
	{
		SetRowEnabled(table.rows[i-1], true);
	}
}

//------------------------------------------------------------------------------
/**
*/
function SetRowEnabled(row, state)
{
	var elements = row.getElementsByTagName('input');
	for (var element in elements)
	{
		elements[element].disabled = !state;
	}
}

//------------------------------------------------------------------------------
/**
*/
function PerformSearch()
{

}