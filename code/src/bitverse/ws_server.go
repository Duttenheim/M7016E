package bitverse

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"net/http"
)

type wsServerType struct {
	msgChannel        chan Msg
	remoteNodeChannel chan *RemoteNode
	localNodeId       NodeId
}

func (wsServer *wsServerType) WsHandler(ws *websocket.Conn) {
	var err error
	var msg Msg
	var remoteNode *RemoteNode = nil

	for {
		dec := json.NewDecoder(ws)
		err = dec.Decode(&msg)
		if err != nil {
			debug("wsserver: Incorrect message structure.")
			break
		}
		if msg.Type == Handshake {
			remoteNode = makeRemoteNode(wsServer.remoteNodeChannel, ws, wsServer.localNodeId.String(), msg.Src)
			wsServer.remoteNodeChannel <- remoteNode

			// send our node id to the remote node so that it can also create a link
			reply := composeHandshakeMsg(wsServer.localNodeId.String())
			enc := json.NewEncoder(ws)
			enc.Encode(reply)
		} else {
			wsServer.msgChannel <- msg
		}
	}
}

func (wsServer* wsServerType) WebHandler(ws* websocket.Conn) {
	var err error
	var msg Msg

	for {
		dec := json.NewDecoder(ws)
		err = dec.Decode(&msg)
		if err != nil {
			debug("wsserver: Incorrect message structure.")
			break
		}
		wsServer.msgChannel <- msg
	}
}

/*	var in []byte
	for {
		if err := websocket.Message.Receive(ws, &in); err != nil {
			break
		}
	}
	websocket.Message.Send(ws, in)	*/
}

func makeWsServer(localNodeId NodeId, msgChannel chan Msg, remoteNodeChannel chan *RemoteNode) *wsServerType {
	wsServer := new(wsServerType)
	wsServer.msgChannel = msgChannel
	wsServer.remoteNodeChannel = remoteNodeChannel
	wsServer.localNodeId = localNodeId

	return wsServer
}

func (wsServer *wsServerType) start(port string) {
	debug("wsserver: starting a new server at port " + port)

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/node", websocket.Handler(wsServer.WsHandler))
	http.Handle("/web", websocket.Handler(wsServer.WebHandler))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("wsserver.start: " + err.Error())
	}
}
