package protocol



//------------------------------------------------------------------------------
/**
	Struct which contains an RPC invocation.
	Use for asynchronous handling of RPC call	
*/
type RpcInvocation struct {
	Id				int
	Reply 			interface{}
	Error 			error
	Done 			chan* RpcInvocation 
}

//------------------------------------------------------------------------------
/**
	Create a new RPC invocation
*/
func MakeInvocation(name string, reply interface{}, done chan* RpcInvocation) *RpcInvocation {
	invocation := new(RpcInvocation)
	invocation.Reply = reply
	invocation.Done = done
	
	return invocation
}

//------------------------------------------------------------------------------
/**
*/
func (call* RpcInvocation) done() {
	select {
	case call.Done <- call:
		// we're done
	default:
		// something went wrong
	}
}

//------------------------------------------------------------------------------
/**
	Type which holds the data used for an RPC call
*/
type RpcInvokeMessage struct {
	Rpc_function_name string
	Args interface{}
}