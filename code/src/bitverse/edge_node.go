package bitverse

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"time"
)

const HEARTBEAT_RATE time.Duration = 10
const MSG_SERVICE_GC_RATE time.Duration = 10

type EdgeNode struct {
	nodeId            NodeId
	superNode         *RemoteNode
	msgChannel        chan Msg
	remoteNodeChannel chan *RemoteNode
	transport         Transport
	msgServices       map[string]*MsgService
	repoServices      map[string]*RepoService
	bitverseObserver  BitverseObserver
	replyTable        map[string]*msgReplyType
	done			  chan int
}

func MakeEdgeNode(transport Transport, bitverseObserver BitverseObserver) (*EdgeNode, chan int) {
	edgeNode := new(EdgeNode)
	edgeNode.transport = transport
	edgeNode.bitverseObserver = bitverseObserver

	edgeNode.nodeId = generateNodeId()
	debug("edgenode: my id is " + edgeNode.Id())

	edgeNode.transport.SetLocalNodeId(edgeNode.nodeId)

	edgeNode.done = make(chan int)
	edgeNode.msgChannel = make(chan Msg)
	edgeNode.remoteNodeChannel = make(chan *RemoteNode, 10)
	edgeNode.msgServices = make(map[string]*MsgService)
	edgeNode.repoServices = make(map[string]*RepoService)
	edgeNode.replyTable = make(map[string]*msgReplyType)

	
	heartbeatTicker := time.NewTicker(time.Millisecond * HEARTBEAT_RATE * 1000)
	go func() {
		for t := range heartbeatTicker.C {
			debug("edgenode: sending heartbeat " + t.String())
			edgeNode.SendHeartbeat()
		}
	}()

	msgServiceGCTicker := time.NewTicker(time.Millisecond * MSG_SERVICE_GC_RATE * 1000)
	go func() {
		for t := range msgServiceGCTicker.C {
			for msgId, reply := range edgeNode.replyTable {
				debug("edgenode: running msg service callback listener garbage collector" + t.String())
				currentTime := int32(time.Now().Unix())
				elapsedTime := currentTime - reply.timestamp
				timeLeft := reply.timeout - elapsedTime

				if timeLeft <= 0 {
					reply.callback(errors.New("timeout"), nil) // notify the callback closure about this timeout
					delete(edgeNode.replyTable, msgId)
				}

			}
		}
	}()
	
	go func() {
		for {
			select {
			case msg := <-edgeNode.msgChannel:
				debug("edgenode: received " + msg.String())
				if msg.Dst == edgeNode.Id() && msg.Type == Data {
					msgService := edgeNode.msgServices[msg.MsgServiceName]
					if msgService == nil {
						debug("edgenode: failed to deliver message, no such service with id <" + msg.MsgServiceName + "> created")
					} else {
						msg.msgService = msgService
						observer := msgService.observer
						if observer == nil {
							debug("edgenode: failed to deliver message, no observer registered")
						} else {
							var err error
							if msg.Payload != "" && msg.PayloadType != Nil {
								if msg.Status != Error { // do not try to decrypt error response messages
									msg.Payload, err = decryptAes(msgService.aesEncryptionKey, msg.Payload)
								} else {
									err = nil
								}
							}
							if err != nil {
								info("edgenode: failed to decrypt payload, ignoring incoming msg")
                                debug(err.Error())
							} else {
								reply := edgeNode.replyTable[msg.Id]
								debug("msg <" + msg.Id + "> is a reply")
								if reply != nil {
									if msg.Status == Error {
										reply.callback(errors.New(msg.Payload), nil)
									} else {
										if msg.PayloadType == Nil {
											reply.callback(nil, nil)
										} else {
											reply.callback(nil, msg.Payload)
										}
									}
									delete(edgeNode.replyTable, msg.Id)
								} else {
									observer.OnDeliver(msgService, &msg)
								}
							}
						}
					}
				} else if msg.Type == Heartbeat {
					debug("edgenode: got heartbeat message from <" + msg.Origin + ">")
					if bitverseObserver != nil {
						if msg.Src != edgeNode.nodeId.String() {
							bitverseObserver.OnSiblingHeartbeat(edgeNode, msg.Origin) // note Src and not Payload since super node just forwards the msg
						}
					}
				} else if msg.Type == ChildJoined {
					debug("edgenode: got child joined message from <" + msg.Payload + ">")
					if bitverseObserver != nil {
						if msg.Payload != edgeNode.nodeId.String() {
							bitverseObserver.OnSiblingJoined(edgeNode, msg.Payload)
						}
					}
				} else if msg.Type == ChildLeft {
					debug("edgenode: got child left message from <" + msg.Payload + ">")
					if bitverseObserver != nil {
						if msg.Payload != edgeNode.nodeId.String() {
							bitverseObserver.OnSiblingLeft(edgeNode, msg.Payload)
						}
					}
				} else if msg.Type == Children {
					if bitverseObserver != nil {
						var children []string
						if err := json.Unmarshal([]byte(msg.Payload), &children); err != nil {
							panic(err)
						}
						bitverseObserver.OnChildrenReply(edgeNode, msg.Src, children)
					}
				} else { // ignore
				}
			case remoteNode := <-edgeNode.remoteNodeChannel:
				if remoteNode.state == Dead {
					debug("edgenode: we just lost our connection to the super node <" + remoteNode.Id() + ">")
					heartbeatTicker.Stop()
					msgServiceGCTicker.Stop()
					edgeNode.superNode = nil
					edgeNode.done <- 1
				} else {
					debug("edgenode: adding link to super node <" + remoteNode.Id() + ">")
					edgeNode.superNode = remoteNode
					if bitverseObserver != nil {
						bitverseObserver.OnConnected(edgeNode, edgeNode.superNode)
					}
				}
			}
		}
	}()

	return edgeNode, edgeNode.done
}

// DEBUG

func (edgeNode *EdgeNode) Debug() {
	debugFlag = true
}

// BITVERSE MANAGEMENT

func (edgeNode *EdgeNode) Id() string {
	return edgeNode.nodeId.String()
}

func (edgeNode *EdgeNode) Connect(remoteAddress string) {
	edgeNode.transport.ConnectToNode(remoteAddress, edgeNode.remoteNodeChannel, edgeNode.msgChannel)
}

func (edgeNode *EdgeNode) Unconnect () {
	//close(edgeNode.done)
	edgeNode.done <- 1
}
func (edgeNode *EdgeNode) Disconnect () {
	edgeNode.superNode.state = Dead
	edgeNode.remoteNodeChannel <- edgeNode.superNode	
}

func (edgeNode *EdgeNode) SendHeartbeat() {
	msg := composeHeartbeatMsg(edgeNode.Id(), edgeNode.superNode.Id())
	edgeNode.superNode.deliver(msg)
}

func (edgeNode *EdgeNode) UpdateTags(tags map[string]string) {
	
	data, err := json.Marshal(&tags)
	if (err != nil) {
		debug("edgenode: failed to marshal tags structure")
	}
	
	// compose and send message
	msg := composeUpdateTagsMsg(edgeNode.Id(), string(data))
	edgeNode.superNode.deliver(msg)	
}

// MSG SERVICE MANAGEMENT

func (edgeNode *EdgeNode) CreateMsgService(aesEncryptionKey string, serviceId string, observer MsgServiceObserver) (*MsgService, error) {
	//if serviceId == REPO_MSG_SERVICE_NAME {
	//	return nil, errors.New("service id <internal> reserved for internal usage")
	//}

	if edgeNode.msgServices[serviceId] == nil {
		msgService := composeMsgService(aesEncryptionKey, serviceId, observer, edgeNode)
		edgeNode.msgServices[serviceId] = msgService
		return msgService, nil
	} else {
		return nil, errors.New("edgenode: service id <" + serviceId + "> already exists")
	}
}

func (edgeNode *EdgeNode) GetMsgService(serviceId string) *MsgService {
	return edgeNode.msgServices[serviceId]
}

// REPO MANAGEMENT

func (edgeNode *EdgeNode) ClaimOwnership(repoId string, aesEncryptionKey string, prv *rsa.PrivateKey, pub *rsa.PublicKey, timeout int32, callback func(err error, repo interface{})) error {
	repoMsgServiceObserver := new(RepoMsgServiceObserver)

	repoMsgService, err := edgeNode.CreateMsgService(aesEncryptionKey, repoId, repoMsgServiceObserver)
	if err != nil {
		return err
	}

	pubPemKey, err := generatePublicPem(pub)
	if err != nil {
		return err
	}

	msg := composeRepoClaimMsg(edgeNode.Id(), edgeNode.superNode.Id(), repoId, pubPemKey)
	repoMsgService.sendMsgAndGetReply(msg, timeout, func(err error, reply interface{}) {
		info("edgenode: got a reply")
		if err != nil {
			info("edgenode: failed to get a claim request reply back")
			callback(err, nil)
		} else {
			info("edgenode: got a claim request reply back")
			repoService := composeRepoService(aesEncryptionKey, prv, pub, repoId, edgeNode, repoMsgService)
			callback(nil, repoService)
		}
	})

	return nil // no errors
}

/// PRIVATE

func (edgeNode *EdgeNode) registerReplyCallback(msgId string, timeout int32, callback func(err error, data interface{})) {
	reply := new(msgReplyType)
	reply.timeout = timeout
	reply.callback = callback
	reply.timestamp = int32(time.Now().Unix())
	edgeNode.replyTable[msgId] = reply
}

func (edgeNode *EdgeNode) send(msg *Msg) {
	debug("edgenode: sending from <" + msg.Src + "> to <" + msg.Dst + "> with Id <" + msg.Id + ">")	
	edgeNode.superNode.deliver(msg)
}
