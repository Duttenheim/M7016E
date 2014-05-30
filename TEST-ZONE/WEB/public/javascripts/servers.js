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
	
	img.loaded = false;
	img.status = statusElement;
	img.button = buttonElement;
	img.onload = function()
	{
		status = "<font color='green'>Online</font>";
		this.status.innerHTML = status;
		this.button.disabled = false;
		this.loaded = true;
	}
	
	var failFunction = function()
	{
		if (!this.loaded)
		{
			status = "<font color='red'>Offline</font>";
			this.status.innerHTML = status;
			this.loaded = true;
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
function RedirectToSupernode(serverip)
{
	window.location.href = "/supernode" + "?ip=" + serverip;
}