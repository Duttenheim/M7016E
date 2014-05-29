package bitverse

import (
	"encoding/json"
	"fmt"
)

type repokey_t struct {
	repoId string
	key    string
}

type SuperNode struct {
	nodeId                 NodeId
	children               map[string]*RemoteNode
	msgChannel             chan Msg
	remoteNodeChannel      chan *RemoteNode
	seqNumberCounter       int
	localAddr              string
	localPort              string
	transport              Transport
			
	repoAutenticationTable map[string]*string    // repoid:public key
	repositories           map[repokey_t]*string // global key-value store
	tags				   map[string]map[string]string
}

func MakeSuperNode(transport Transport, localAddress string, localPort string) (*SuperNode, chan int) {
	superNode := new(SuperNode)

	superNode.localPort = localPort
	superNode.children = make(map[string]*RemoteNode)
	superNode.transport = transport
	superNode.tags = make(map[string]map[string]string)

	superNode.repoAutenticationTable = make(map[string]*string)
	superNode.repositories = make(map[repokey_t]*string)

	superNode.nodeId = generateNodeId()
	debug("supernode: my id is " + superNode.Id())

	superNode.transport.SetLocalNodeId(superNode.nodeId)

	done := make(chan int)
	superNode.msgChannel = make(chan Msg)
	superNode.remoteNodeChannel = make(chan *RemoteNode, 10)

	go superNode.transport.Listen(localAddress, localPort, superNode.remoteNodeChannel, superNode.msgChannel)

	go func() {
		for {
			select {
			case msg := <-superNode.msgChannel:
				//debug("supernode: received " + msg.String())
				if msg.Dst == superNode.Id() && msg.Type == Data {
					// ignore, not supported

				} else if msg.Type == Data && msg.ServiceType == Repo && msg.RepoCmd == Claim { // repo claim request
					// REPO CLAIM REQUEST
					repoId := msg.RepoId
					pubKeyPem := msg.Signature
					debug("supernode: got a repo claim request for repo " + repoId + " with public key <" + pubKeyPem + ">")

					if superNode.repoAutenticationTable[repoId] == nil {
						// it is free, claim it!
						superNode.repoAutenticationTable[repoId] = &pubKeyPem // XXX is this safe?
						msg.Status = Ok

					} else {
						// already claimed
						if pubKeyPem == *superNode.repoAutenticationTable[repoId] {
							// but, it the same owner
							msg.Status = Ok
						} else {
							msg.Status = Error
							msg.Payload = "repo already claimed"
						}
					}

					childId := msg.Src
					msg.Dst = childId
					msg.Src = superNode.Id()
					superNode.sendToChild(msg)

				} else if msg.Type == Data && msg.ServiceType == Repo && msg.RepoCmd == Store {
					// REPO STORE REQUEST
					debug("supernode: got a repo store request repo <" + msg.RepoId + "> with key <" + msg.RepoKey + "> value <" + msg.RepoValue + "> with signature <" + msg.Signature + ">")

					repoId := msg.RepoId

					if superNode.repoAutenticationTable[repoId] == nil {
						msg.Status = Error
						msg.Payload = "no such repo " + repoId
					} else {
						key := msg.RepoKey
						value := msg.RepoValue
						signature := msg.Signature

						pubPemKey := superNode.repoAutenticationTable[repoId]
						if pubPemKey == nil {
							errMsg := "failed to receive public key for repo <" + repoId + ">"
							info("supernode: ERROR " + errMsg)
							msg.Status = Error
							msg.Payload = errMsg
						} else {
							//_, pub, importErr := ImportPem("cert2")
							_, pub, importErr := importKeyFromString(*pubPemKey)
							if importErr != nil {
								errMsg := "failed to convert pem public key for repo <" + repoId + ">"
								info("supernode: ERROR " + errMsg)
								msg.Status = Error
								msg.Payload = errMsg
							} else {
								verfErr := verify(pub, value, signature) // the key and value are aes encrypted
								if verfErr != nil {
									errMsg := "failed to verify signature for repo <" + repoId + ">"
									info("supernode: ERROR " + errMsg)
									msg.Status = Error
									msg.Payload = errMsg
								} else {
									oldValue := superNode.repositories[repokey_t{repoId, key}]
									superNode.repositories[repokey_t{repoId, key}] = &value
									if oldValue == nil {
										info("supernode: setting key <" + key + "> to value <" + value + ">")
										msg.Status = Ok
										msg.PayloadType = Nil
									} else {
										info("supernode: replacing key <" + key + "> with value <" + value + ">, old value was <" + *oldValue + ">")
										msg.Status = Ok
										msg.Payload = *oldValue
									}
								}
							}
						}
					}

					// now it is time to send a reply back depending of the outcome
					childId := msg.Src
					msg.Dst = childId
					msg.Src = superNode.Id()
					superNode.sendToChild(msg)
				} else if msg.Type == Data && msg.ServiceType == Repo && msg.RepoCmd == Lookup {
					// REPO LOOKUP REQUEST
					debug("supernode: got a repo look request repo <" + msg.RepoId + "> with key <" + msg.RepoKey + "> with signature <" + msg.Signature + ">")

					repoId := msg.RepoId
					if superNode.repoAutenticationTable[repoId] == nil {
						msg.Status = Error
						msg.Payload = "no such repo " + repoId
					} else {
						key := msg.RepoKey
						signature := msg.Signature

						pubPemKey := superNode.repoAutenticationTable[repoId]
						if pubPemKey == nil {
							errMsg := "failed to receive public key for repo <" + repoId + ">"
							info("supernode: ERROR " + errMsg)
							msg.Status = Error
							msg.Payload = errMsg
						} else {
							_, pub, importErr := importKeyFromString(*pubPemKey)
							if importErr != nil {
								errMsg := "failed to convert pem public key for repo <" + repoId + ">"
								info("supernode: ERROR " + errMsg)
								msg.Status = Error
								msg.Payload = errMsg
							} else {
								verfErr := verify(pub, key, signature) // the key is aes encrypted
								if verfErr != nil {
									errMsg := "failed to verify signature for repo <" + repoId + ">"
									info("supernode: ERROR " + errMsg)
									msg.Status = Error
									msg.Payload = errMsg
								} else {
									value := superNode.repositories[repokey_t{repoId, key}]
									if value == nil {
										msg.Status = Ok
										msg.PayloadType = Nil
									} else {
										msg.Status = Ok
										msg.Payload = *value
									}
								}
							}
						}

						// now it is time to send a reply back depending of the outcome
						childId := msg.Src
						msg.Dst = childId
						msg.Src = superNode.Id()
						superNode.sendToChild(msg)
					}

				} else if msg.Type == Heartbeat || msg.Type == ChildJoined || msg.Type == ChildLeft {
					superNode.forwardToChildren(msg)
					
				} else if msg.Type == UpdateTags {
					superNode.updateTags(msg)
				
				} else if msg.Type == SearchTags {
					superNode.searchTags(msg)				
					
				} else if msg.Type == GetTags {
					superNode.getTags(msg)
					
				} else if msg.Type == MakeImposter {
					superNode.makeSupernode(msg)

				} else if msg.Type == Children {
					superNode.sendChildrenReply(msg.Src)

				} else {
					superNode.sendToChild(msg)
				}
			case remoteNode := <-superNode.remoteNodeChannel:
				if remoteNode.state == Dead {
					delete(superNode.children, remoteNode.Id())
					delete(superNode.tags, remoteNode.Id())

					str := fmt.Sprintf("supernode: removing remote node %s, number of remote nodes are now %d", remoteNode.Id(), len(superNode.children))
					fmt.Println(str)

					msg := composeChildLeft(superNode.nodeId.String(), remoteNode.Id())
					superNode.forwardToChildren(*msg)
				} else {
					superNode.children[remoteNode.Id()] = remoteNode
					superNode.tags[remoteNode.Id()] = make(map[string]string)

					str := fmt.Sprintf("supernode: adding remote node %s, number of remote nodes are now %d", remoteNode.Id(), len(superNode.children))
					info(str)

					msg := composeChildJoin(superNode.nodeId.String(), remoteNode.Id())
					superNode.forwardToChildren(*msg)
				}
			}
		}
	}()

	return superNode, done
}

// Connects this supernode to another in order to form a P2P ring, call in separate go-function
func(superNode* SuperNode) ConnectSuccessor(addrs []string, port string) {

	for _, addr := range addrs {
		transport := MakeWSTransport()
		transport.SetLocalNodeId(superNode.nodeId)
		msgChannel := make(chan Msg)
		nodeChannel := make(chan *RemoteNode, 10)	
		var SN *RemoteNode
	
		go func() {
			for {
				select {
				case msg := <-msgChannel:
					debug("supernode: relaying " + msg.String())
					if msg.Type == ChildJoined {
						superNode.children[msg.Origin] = SN
						debug("supernode: remote child joined on another supernode, mapped " + msg.Origin + " to: " + SN.Id())
					} else if msg.Type == ChildLeft {
						delete(superNode.children, msg.Origin)
						debug("supernode: remote child " + msg.Origin)
					} else if msg.Type == Heartbeat {
						// keep this super node local
					} else {
						msg.Src = msg.Origin
						superNode.msgChannel <- msg
					}
				case remoteNode := <-nodeChannel:
					if remoteNode.state == Dead {
						debug("supernode: we just lost our connection to the successor <" + remoteNode.Id() + ">")
						break
					} else {
						debug("supernode: adding link to successor node <" + remoteNode.Id() + ">")
						SN = remoteNode
						
						// make this supernode an imposter
						msg := composeMakeImposterMsg(superNode.Id())
						remoteNode.deliver(msg)
					}
				}
			}
		}()
	
		go transport.ConnectToNode(addr + ":" + port, nodeChannel, msgChannel)	
	}
}

// BITVERSE MANAGEMENT

func (superNode *SuperNode) Id() string {
	return superNode.nodeId.String()
}

// DEBUG

func (superNode *SuperNode) Debug() {
	debugFlag = true
}

/// PRIVATE

func (superNode *SuperNode) sendChildrenReply(nodeId string) {
	debug("supernode: sending children reply to " + nodeId)
	var childIds []string
	i := 0
	for childNodeId, node := range superNode.children {
		if node.imposter == false {
			childIds = append(childIds, childNodeId)
			i++
		}
	}

	json, _ := json.Marshal(childIds)
	reply := composeChildrenReplyMsg(superNode.Id(), nodeId, string(json))

	if remoteNode, ok := superNode.children[nodeId]; ok {
		remoteNode.deliver(reply)
	}
}

func (superNode *SuperNode) sendToChild(msg Msg) {
	if val, ok := superNode.children[msg.Dst]; ok {
		debug("supernode: forwarding " + msg.String() + " to " + val.Id())
		msg.Src = superNode.Id()
		val.deliver(&msg)
	}
}

func (superNode *SuperNode) forwardToChildren(msg Msg) {
	for _, remoteNode := range superNode.children {
		if msg.Src != remoteNode.Id() && msg.Origin != remoteNode.Id() { // do not forward messages to a remote node where it came from
			debug("supernode: forwarding " + msg.String() + " to " + remoteNode.Id())
			msg.Src = superNode.Id()
			remoteNode.deliver(&msg)
		}
	}
}

//------------------------------------------------------------------------------
/**
	Tags
*/
func (superNode *SuperNode) updateTags(msg Msg) {

	// find node
	if val, ok := superNode.children[msg.Src]; ok {
	
		// decode message contents into tag dictionary
		tags := make(map[string]string)
		err := json.Unmarshal([]byte(msg.Payload), &tags)
		
		if (err != nil) {
			debug(err.Error())
			return
		}
		
		// notify
		debug("supernode: updated tags for " + msg.Src + " with " + msg.Payload)
		
		// set tags for node
		superNode.tags[val.Id()] = tags
	} else {
		debug("supernode: failed to update tags for " + msg.Src)
	}
}

//------------------------------------------------------------------------------
/**
*/
func (superNode *SuperNode) searchTags(msg Msg) {

	// decode match tags
	search := new(SearchTagsType)
	err := json.Unmarshal([]byte(msg.Payload), search)
	if err != nil {
		debug("supernode: failed to decode search tags criteria")
		return
	}

	// get tags by sending a message to all imposters
	spam := new(Msg)
	spam.Type = GetTags
	
	// create map of tags and of matching nodes
	nodeTags := 		make(map[string]map[string]string)
	var matchingNodes 	[]string
	
	for id, node := range superNode.children {
		if node.imposter == true {
		
			// deliver message
			node.deliver(spam)
			
			// get response
			response := node.receive()
			
			var nodeTags map[string]map[string]string
			err := json.Unmarshal([]byte(response.Payload), nodeTags)
			
			if err != nil {
				debug("supernode: failed to decode tags from " + id)
				continue
			}
			
			// merge with current
			for val, key := range nodeTags {
				nodeTags[val] = key
			}
		}
	}
	
	// merge other node tags with ours
	for val, key := range superNode.tags {
		nodeTags[val] = key
	}
	
	// go through nodes
	for node, tags := range nodeTags {
		// go through tags in node
		for key, val := range tags {
			
			// see if tags exist in the search criteria
			if tag, ok := search.tags[key]; ok {
				
				// if the tag is found and the value matches
				if tag == val {
					matchingNodes = append(matchingNodes, node)
				}
			}
		}
	}
	
	// encode to json again and send to the rest of the children
	data, err := json.Marshal(matchingNodes)
	if (err != nil) {
		debug("supernode: failed to encode to string, this should never happen")
		return
	}
	
	// reply
	msg.Reply(string(data))
}

//------------------------------------------------------------------------------
/**
*/
func (superNode *SuperNode) getTags(msg Msg) {
	// decode tags from message
	tags := make(map[string]map[string]string)
	
	// if we have a specific node we want to get the tags for
	if len(msg.Dst) == 0 {
		tags = superNode.tags
	} else {
		// find node
		if val, ok := superNode.tags[msg.Dst]; ok {
			tags[msg.Dst] = val
		} else {
			debug("supernode: no node named " + msg.Dst + " found")
			return
		}
	}
	
	// encode to json again and send to the rest of the children
	data, err := json.Marshal(tags)
	if (err != nil) {
		debug("supernode: failed to reencode to string, this should never happen")
		return
	}
	
	// notify
	debug("supernode: replying with tags " + string(data) + " for node " + msg.Dst)
	
	// get [ay;pad
	msg.Payload = string(data)
	
	// send back
	superNode.children[msg.Src].deliver(&msg)
}

//------------------------------------------------------------------------------
/**
*/
func (superNode *SuperNode) makeSupernode(msg Msg) {
	// find node
	if val, ok := superNode.children[msg.Src]; ok {
		debug("supernode: making " + msg.Src + " into an imposter remote node")
		val.imposter = true
	} else {
		debug("supernode: no supernode named " + msg.Src + " found")
		return
	}
}
