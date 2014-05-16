function CheckOnline(addrid, elementid, buttonid)
{
	var img = document.body.appendChild(document.createElement("img"));
	
	var addrElement = document.getElementById(addrid);
	var statusElement = document.getElementById(elementid);
	var buttonElement = document.getElementById(buttonid);
	statusElement.innerHTML = "Pending";	
	buttonElement.disabled = true;
	
	img.onload = function()
	{
		status = "<font color='green'>Online</font>";
		statusElement.innerHTML = status;
		buttonElement.enabled = true;
	}
	
	var failFunction = function()
	{
		status = "<font color='red'>Offline</font>";
		statusElement.innerHTML = status;
		img.src = "";
	}
	
	img.src = "http://" + addrElement.innerHTML + "/ping.bmp";
	img.onerror = img.onabort = failFunction;
	setTimeout
	(
		failFunction,
		3000
	);
}
