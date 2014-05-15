package protocol

import (
	"bitverse"
	"encoding/json"
	"reflect"
	"fmt"
)


//------------------------------------------------------------------------------
/**
*/
type RpcMessageObserver struct {
	debug bool
	interfaces map[string]reflect.Value
	objects map[string]reflect.Value
	inputs map[string]reflect.Type
	outputs map[string]reflect.Type
}

//------------------------------------------------------------------------------
/**
	Struct which contains an RPC invocation.
	Use for asynchronous handling of RPC call	
*/
type RpcInvocation struct {
	Service string
	Args interface{}
	Reply interface{}
	Error error
	Done chan* RpcInvocation 
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
	Opens the observer, run this before using it for anything
*/
func (observer* RpcMessageObserver) Open() {
	observer.debug = false
	observer.interfaces = make(map[string]reflect.Value)
	observer.objects = make(map[string]reflect.Value)
	observer.inputs = make(map[string]reflect.Type)
	observer.outputs = make(map[string]reflect.Type)
}

//------------------------------------------------------------------------------
/**
*/
func (observer* RpcMessageObserver) SetDebugEnabled(enabled bool) {
	observer.debug = enabled
}

//------------------------------------------------------------------------------
/**
*/
type RpcInvokeMessage struct {
	Rpc_function_name string
	Args interface{}
}

//------------------------------------------------------------------------------
/**
*/
func (observer* RpcMessageObserver) Register(v interface{}) {
	// reflect type to value
	typeRef := reflect.TypeOf(v)

	// loop through methods and register them
	numMethods := typeRef.NumMethod()
	for i := 0; i < numMethods; i++ {

		// get method
		methodRef := typeRef.Method(i)
		name := fmt.Sprintf("%s.%s", typeRef.Elem().Name(), methodRef.Name)

		// assert that method has two exported fields (input and output)
		if (methodRef.Type.NumIn() != 3 || methodRef.Type.NumOut() != 1) {
			fmt.Printf("Method '%s' doesn't fit the criteria for RPC serialization (object, input and reply args, and error return)\n", name)			
			continue
		}

		// get first field
		input := methodRef.Type.In(1).Elem()

		// get second field
		output := methodRef.Type.In(2).Elem()

		if (observer.debug) { fmt.Printf("Method '%s' registered with input '%s' and output '%s'!\n", name, input, output) }
		
		observer.interfaces[name] = methodRef.Func
		observer.objects[name] = reflect.ValueOf(v)
		observer.inputs[name] = input
		observer.outputs[name] = output
	}
}

//------------------------------------------------------------------------------
/**
	Composes an RPC message with an output and input and executes it asynchronously
*/
func RpcInvokeAsync(functionName string, args interface{}, reply interface{}, node string, service* bitverse.MsgService, done chan* RpcInvocation) *RpcInvocation {

	// create message
	msg := &RpcInvokeMessage{functionName, args}

	// create new invocation
	invocation := new(RpcInvocation)
	invocation.Service = functionName
	invocation.Args = args
	invocation.Reply = reply
	invocation.Done = done

	// encode as JSON
	data, err := json.Marshal(msg)
	if (err != nil) {
		invocation.Error = err
		invocation.done()
	}

	// send message
	service.SendAndGetReply(node, string(data), 1, func(err error, data interface{}) {

		if (err != nil) {
			invocation.Error = err
			invocation.done()
			return
		}

		// unmarshal data
		encerr := json.Unmarshal([]byte(data.(string)), invocation.Reply)

		if (encerr != nil) {
			invocation.Error = err
			invocation.done()
			return
		}

		invocation.Error = nil
		invocation.done()
	})

	return invocation
}

//------------------------------------------------------------------------------
/**
	Invokes RPC synchronously
*/
func RpcInvokeSync(functionName string, args interface{}, reply interface{}, node string, service* bitverse.MsgService) error {
	call := <-RpcInvokeAsync(functionName, args, reply, node, service, make(chan *RpcInvocation, 1)).Done
	return call.Error
}

//------------------------------------------------------------------------------
/**
	Receives a message, decodes it as JSON and attempts to run it as an RPC call.
	Currently requires two messages and it is doubtful if it will properly work.
*/
func (observer* RpcMessageObserver) OnDeliver(service* bitverse.MsgService, msg* bitverse.Msg) {
		
	if (observer.debug) { fmt.Printf("Message payload '%s'\n", msg.Payload) }
	
	// create package and unmarshal
	var pkg RpcInvokeMessage
	err := json.Unmarshal([]byte(msg.Payload), &pkg)

	// if this wasn't successful, just print it and return
	if (err != nil) {
		if (observer.debug) { fmt.Println("Message is not an RPC message!") }
		return	
	}

	_, exists := observer.interfaces[pkg.Rpc_function_name]
	if (!exists) {
		if (observer.debug) { fmt.Printf("No interface named '%s' registered!\n", pkg.Rpc_function_name) }
		return
	}

	// get function and object
	function := observer.interfaces[pkg.Rpc_function_name]
	object := observer.objects[pkg.Rpc_function_name]
	inputType := observer.inputs[pkg.Rpc_function_name]
	outputType := observer.outputs[pkg.Rpc_function_name]

	// create input and output
	input := reflect.New(inputType)
	output := reflect.New(outputType)

	// now unmarshal arguments which should be a string
	err = json.Unmarshal([]byte(pkg.Args.(string)), input.Interface())
	if (err != nil) {
		if (observer.debug) { fmt.Println("Message did not receive a proper input structure") }
		return
	}

	// invoke that function!
	args := []reflect.Value{object, input, output}
	ret := function.Call(args)

	// if we have an error, print it and return
	if (len(ret) == 1) {
		if (ret[0].Interface() != nil) {
			fmt.Printf("RPC error: %s\n", ret[0].Interface())
			return
		}
	}

	// now quite simply, marshal the output and send it as feedback
	var response []byte
	response, err = json.Marshal(output.Interface())
	if (err != nil) {
		if (observer.debug) { fmt.Println("Output could not be encoded!") }
		return
	}

	// send the message back
	msg.Reply(string(response))
}
