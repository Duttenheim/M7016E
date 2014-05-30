
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
	
	// input field 1
	var div1 = document.createElement("div");
	var span1 = document.createElement("span");
	var label1 = document.createElement("label");
	span1.innerHTML = "";
	label1.className = "control-label";
	span1.className = "glyphicon form-control-feedback";
	div1.className = "has-feedback";
	
	var input1 = document.createElement("input");
	input1.id = 'tag';
	input1.type = 'text';
	input1.className = "form-control";
	input1.setAttribute('list', 'services');
	input1.placeholder = 'e.g Service';
	input1.onfocus = function()
	{
		div1.className = "has-feedback";
		span1.className = "glyphicon form-control-feedback";
		label1.innerHTML = "";
		label1.style.display = "none";
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
	div1.appendChild(input1);
	div1.appendChild(span1);
	div1.appendChild(label1);
	cell1.appendChild(div1);
	
	// input field 2
	var div2 = document.createElement("div");
	var span2 = document.createElement("span");
	var label2 = document.createElement("label");
	label2.className = "control-label";
	span2.className = "glyphicon form-control-feedback";
	div2.className = "has-feedback";
	
	var input2 = document.createElement("input");
	input2.id = 'value';
	input2.type = 'text';	
	input2.className = "form-control";
	input2.placeholder = 'e.g Docker'
	input2.onfocus = function()
	{
		div2.className = "has-feedback";
		span2.className = "glyphicon form-control-feedback";
		label2.innerHTML = "";
		label2.style.display = "none";
	}
	
	div2.appendChild(input2);
	div2.appendChild(span2);
	div2.appendChild(label2);
	cell2.appendChild(div2);
	
	var buttonDiv = document.createElement("span");
	var buttonAdd = document.createElement("button");
	buttonAdd.className = "btn btn-default btn-margins";
	buttonAdd.innerHTML = 'Add';
	buttonAdd.id = "add";
	
	var buttonRemove = document.createElement("button");
	buttonRemove.className = "btn btn-default btn-margins";
	buttonRemove.innerHTML = "Remove";
	buttonRemove.id = "remove";
	
	buttonAdd.onclick = function() 
	{ 		
		var index = usedTags.indexOf(input1.value);
		if (index == -1 && input1.value.length != 0)
		{
			div1.className = "has-success has-feedback";
			span1.className = "glyphicon glyphicon-ok form-control-feedback";
			label1.innerHTML = "";
			label1.style.display = "none";
		}
		else
		{
			div1.className = "has-error has-feedback";
			span1.className = "glyphicon glyphicon-remove form-control-feedback";
			label1.style.display = "inline-block";
			if (input1.value.length == 0)
			{
				label1.innerHTML = "Search tag empty";
			}
			else
			{
				label1.innerHTML = "Search tag taken";
			}
		}
		
		if (input2.value.length != 0)
		{
			div2.className = "has-success has-feedback";
			span2.className = "glyphicon form-control-feedback glyphicon-ok";
			label2.innerHTML = "";
			label2.style.display = "none";
		}
		else
		{
			div2.className = "has-error has-feedback";
			span2.className = "glyphicon form-control-feedback glyphicon-remove";
			label2.innerHTML = "Search value empty";
			label2.style.display = "inline-block";
		}
		
		// only add if tag isn't occupied
		var index = usedTags.indexOf(input1.value);
		if (index == -1 && input2.value.length != 0 && input1.value.length != 0)
		{
			usedTags.push(input1.value);
			buttonAdd.disabled = true;
			buttonRemove.disabled = false;
			SetRowEnabled(row, false);
			var newRow = table.insertRow(-1);
			NewSearchRow(newRow, table); 
		}
	}
	
	buttonRemove.onclick = function()
	{
		buttonAdd.disabled = false;
		buttonRemove.disabled = true;
		RemoveSearchRow(row, table);
	}	
	buttonDiv.appendChild(buttonAdd);
	buttonDiv.appendChild(buttonRemove);
	cell3.appendChild(buttonDiv);
	
	// finally, move focus to input1
	input1.focus();
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
		var nextRow = i-1;
		SetRowEnabled(table.rows[nextRow], true);
		ResetRow(table.rows[nextRow]);
	}
}

//------------------------------------------------------------------------------
/**
*/
function SetRowEnabled(row, state)
{
	var inputs = row.getElementsByTagName("input");
	for (var input in inputs)
	{
		inputs[input].disabled = !state;
	}
	
	var buttons = row.getElementsByTagName("button");
	for (var button in buttons)
	{
		buttons[button].disabled = !state;
		
		// always disable first remove button
		if (row.rowIndex == 0 && buttons[button].id == "remove")
		{
			buttons[button].disabled = true;
		}
	}
}

//------------------------------------------------------------------------------
/**
*/
function ResetRow(row)
{
	var inputs = row.getElementsByTagName("input");
	for (var input in inputs)
	{
		if (inputs[input].id == "tag")
		{
			var index = usedTags.indexOf(inputs[input].value);
			if (index > -1)
			{
				usedTags.splice(index, 1);
			}
		}
	}
}

//------------------------------------------------------------------------------
/**
*/
function PerformSearch()
{

}