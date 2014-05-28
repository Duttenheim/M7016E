package protocol

import (
	"encoding/json"
	"reflect"
	"fmt"
)

//------------------------------------------------------------------------------
/**
*/
type RpcServer struct {
	debug bool
	interfaces map[string]reflect.Value
	objects map[string]reflect.Value
	inputs map[string]reflect.Type
	outputs map[string]reflect.Type
}


//------------------------------------------------------------------------------
/**
*/
func MakeRpcServer() *RpcServer {
	server := new(RpcServer)
	server.debug = false
	server.interfaces = make(map[string]reflect.Value)
	server.objects = make(map[string]reflect.Value)
	server.inputs = make(map[string]reflect.Type)
	server.outputs = make(map[string]reflect.Type)
	return server
}

//------------------------------------------------------------------------------
/**
*/
func (server* RpcServer) SetDebugEnabled(enabled bool) {
	server.debug = enabled
}

//------------------------------------------------------------------------------
/**
*/
func (server* RpcServer) Register(v interface{}) {
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
			if (server.debug) { fmt.Printf("Method '%s' doesn't fit the criteria for RPC serialization (object, input and reply args, and error return)\n", name) }			
			continue
		}

		// get first field
		input := methodRef.Type.In(1).Elem()

		// get second field
		output := methodRef.Type.In(2).Elem()

		if (server.debug) { fmt.Printf("Method '%s' registered with input '%s' and output '%s'!\n", name, input, output) }
		
		server.interfaces[name] = methodRef.Func
		server.objects[name] = reflect.ValueOf(v)
		server.inputs[name] = input
		server.outputs[name] = output
	}
}

//------------------------------------------------------------------------------
/**
	Receives a string which must be an RpcInvokeMessage, 
	and returns a string corresponding the JSON encoded result of the execution.
	
	NOTE:
		The handling of invocing the RPC call and sending the reply must be handled somewhere else
*/
func (server* RpcServer) Process(call string) (string, error) {
	if (server.debug) { fmt.Printf("Message payload '%s'\n", call) }
	
	// create package and unmarshal
	var pkg RpcInvokeMessage
	err := json.Unmarshal([]byte(call), &pkg)

	// if this wasn't successful, just print it and return
	if (err != nil) {
		if (server.debug) { fmt.Println("Message is not an RPC message!") }
		return "", err
	}

	function, exists := server.interfaces[pkg.Rpc_function_name]
	if (!exists) {
		return "", fmt.Errorf("No interface named '%s' registered!\n", pkg.Rpc_function_name)
	}

	// get function and object
	object := server.objects[pkg.Rpc_function_name]
	inputType := server.inputs[pkg.Rpc_function_name]
	outputType := server.outputs[pkg.Rpc_function_name]

	// create input and output
	input := reflect.New(inputType)
	output := reflect.New(outputType)

	// now unmarshal arguments which should be a string	
	err = json.Unmarshal([]byte(pkg.Args.(string)), input.Interface())
	if (err != nil) {
		return "", err
	}

	// invoke that function!
	args := []reflect.Value{object, input, output}
	ret := function.Call(args)

	// if we have an error, print it and return
	if (len(ret) == 1) {
		if (ret[0].Interface() != nil) {
			return "", fmt.Errorf("RPC error: %s\n", ret[0].Interface())
		}
	}

	// now quite simply, marshal the output and send it as feedback
	var response []byte
	response, err = json.Marshal(output.Interface())
	if (err != nil) {
		return "", err
	}
	
	return string(response), nil
}

//------------------------------------------------------------------------------
/**
	Convenience function which takes function name and arguments and converts into JSON package.
	Basically, the RPC call ends up as
	{ functionName : string, args : string of { arg 1 : string, ... arg n : string} }
*/
func ComposeRPC(name string, args interface{}) (string, error) {
	var data []byte
	var err error
	
	// encode args as string
	data, err = json.Marshal(args)	
	if err != nil {
		return "", err
	}
	
	msg := &RpcInvokeMessage{name, string(data)}
	data, err = json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}