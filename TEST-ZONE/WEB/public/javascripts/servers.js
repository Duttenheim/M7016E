function CheckOnline(addrid, elementid)
{
	
	var img = document.body.appendChild(document.createElement("img"));
	
	statusElement = document.getElementById(elementid);
	statusElement.innerHTML = "Unknown";
	addrElement = document.getElementById(addrid);
	
	img.onload = function()
	{
		status = "<font color='green'>Online</font>";
		statusElement.innerHTML = status;
	}
	
	failFunction = function()
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
		1000
	);
}