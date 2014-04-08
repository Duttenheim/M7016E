
var MsgTypeEnum = 
{
	Handshake : 0,
	Data : 1,
	Heartbeat : 2,
	Children : 3,
	ChildJoined : 4,
	ChildLeft : 5,
	Bye : 6	
}

function Msg ()
{
	this.Type = -1;
	this.Payload = "";
	this.PayloadType = -1;
	this.Src = "";
	this.Dst = "";
	this.Id = "";
}
