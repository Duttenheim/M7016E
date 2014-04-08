// assumes lib/websocket has been loaded before

function WebNode()
{
    this.id = "";
    this.superNodeId = "";
    this.currentChildren;
    this.socket;
}

WebNode.prototype.OnOpen = function()
{
    alert("WebNode.OnOpen called!");

    // create handshake and send
    var message = new Msg();
    message.Type = MsgTypeEnum.Handshake;
    var json = JSON.stringify(message);
    this.send(json);
}

WebNode.prototype.OnClose = function()
{
    alert("WebNode.OnClose called!");
}

WebNode.prototype.OnMessage = function(msg)
{
    // parse string to message
    var message = JSON.parse(msg.data);

    alert("Message type is: " + message.Type);
    if (message.Type == MsgTypeEnum.Handshake)
    {
        this.superNodeId = message.Src;
        alert(msg.data);
        alert("My Id is now: " + this.id + "!");
    }
    else if (message.Type == MsgTypeEnum.Children)
    {
        this.currentChildren = message.Payload;
    }
}

function CreateWebNode(address)
{
    var node = new WebNode();
    node.socket = OpenSocket(address, node.OnOpen, node.OnClose, node.OnMessage);
    return node;
}
