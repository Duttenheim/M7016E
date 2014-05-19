//------------------------------------------------------------------------------
/**
*/
function CheckOnline(addrid, elementid, buttonid)
{
	var img = document.body.appendChild(document.createElement("img"));
	
	var addrElement = document.getElementById(addrid);
	var statusElement = document.getElementById(elementid);
	var buttonElement = document.getElementById(buttonid);
	statusElement.innerHTML = "Pending";	
	buttonElement.disabled = true;
	
	var loaded = false;
	img.onload = function()
	{
		status = "<font color='green'>Online</font>";
		statusElement.innerHTML = status;
		buttonElement.disabled = false;
		loaded = true;
	}
	
	var failFunction = function()
	{
		if (!loaded)
		{
			status = "<font color='red'>Offline</font>";
			statusElement.innerHTML = status;
			img.src = "";
		}
	}
	
	img.src = "http://" + addrElement.innerHTML + "/ping.bmp";
	img.onerror = img.onabort = failFunction;
	setTimeout
	(
		failFunction,
		3000
	);
}

//------------------------------------------------------------------------------
/**
*/
function Redirect(serverip)
{
	window.location.href = "/supernode" + "?ip=" + serverip;
}