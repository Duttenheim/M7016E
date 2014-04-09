
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

function Msg()
{
	this.Type;
	this.Payload;
	this.PayloadType;
	this.Src;
	this.Dst;
	this.Id;
    this.MsgServiceName;
}
