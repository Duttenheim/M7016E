// assumes lib/websocket has been loaded before

var secret = "3e606ad97e0a738d8da4c4c74e8cd1f1f2e016c74d85f17ac2fc3b5dab4ed6c4";

//------------------------------------------------------------------------------
/**
*/
function WebNode()
{
    // generate UUID    
    var uid = generateUUID();
    uid = CryptoJS.SHA1(uid);
    this.id = uid.toString(CryptoJS.enc.Hex);

    // generate new uuid and encode as sha1
    this.superNodeId = "";
    this.socket;

    this.childrenRetreivedCallback = function(children)
    {
        alert(children);
    }

    this.messageRetreivedCallback = function(message)
    {
        alert(message);
    }
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.OnOpen = function()
{
    // create handshake and send
    var message = new Msg();
    message.Type = MsgTypeEnum.Handshake;
    message.Src = this.id;
    var json = JSON.stringify(message);
    this.socket.send(json);
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.OnClose = function()
{
    //alert("WebNode.OnClose called!");
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.OnMessage = function(msg)
{
    // parse string to message
    var message = JSON.parse(msg.data);

    if (message.Type == MsgTypeEnum.Handshake)
    {
        this.superNodeId = message.Src;
    }
    else if (message.Type == MsgTypeEnum.Children)
    {
        var children = JSON.parse(message.Payload);
        var index = children.indexOf(this.id);
        if (index > -1) { children.splice(index, 1); }
        this.childrenRetreivedCallback(children);
    }
    else if (message.Type == MsgTypeEnum.Data)
    {
        var decrypted = CryptoJS.AES.decrypt(message.Payload, secret);
        this.messageRetreivedCallback(CryptoJS.enc.Base64.parse(decrypted));
    }
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.GetSiblings = function()
{
    // create message
    var message = new Msg();
    message.Type = MsgTypeEnum.Children;

    // encode and send
    var json = JSON.stringify(message);
    this.socket.send(json);
}

//------------------------------------------------------------------------------
/**
*/
function EncryptAES(params)
{
    var base64 = CryptoJS.enc.Base64.parse(params);
    var encrypted = CryptoJS.AES.encrypt(base64.toString(), secret, { mode: CryptoJS.mode.CFB });
    base64 = CryptoJS.enc.Base64.parse(encrypted.iv.toString() + encrypted.ciphertext.toString());
    return base64.toString();
}

//------------------------------------------------------------------------------
/**
*/
WebNode.prototype.CallRPCFunction = function(name, args, node)
{
    // create invocation, the first field is the function name, the second is a JSON encoded list of arguments
    var rpcInvoke =
    {
        Rpc_function_name : name,
        Args : JSON.stringify(args)
    };

    // create message
    var message = new Msg();
    var jsonPayload = JSON.stringify(rpcInvoke);
    var encrypted = EncryptAES(jsonPayload);
    message.Payload = encrypted;
    message.Type = MsgTypeEnum.Data;
    message.MsgServiceName = "RPCMessageService";
    message.Dst = node;
    message.Src = this.id;

    // encode message and send
    var json = JSON.stringify(message);
    this.socket.send(json);
}

//------------------------------------------------------------------------------
/**
*/
function CreateWebNode(address)
{
    var node = new WebNode();
    node.socket = OpenSocket(address, 
        function() { node.OnOpen(); }, 
        function() { node.OnClose(); }, 
        function(msg) { node.OnMessage(msg) }
    );
    return node;
}
