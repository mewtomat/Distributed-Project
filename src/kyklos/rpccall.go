package kyklos

import (
	"strconv"
	"net/rpc"
)

func fingerEquals(a, b Finger) bool{
	return (a.Ip == b.Ip) && (a.Port == b.Port)
}

func (node *Finger) callRPCGetSuccessor() (Finger,error){
	if fingerEquals(*myself.nodeFinger,*node){
		//calling on own Server
		var reply Finger
		err := myself.nodeFinger.getSuccessor(nil,&reply)
		return reply, err
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return Finger{}, err
	}
	defer client.Close()
	var reply Finger
	dummy:=0
	err = client.Call("Finger.GetSuccessor",&dummy,&reply)
	if err!=nil{
		Error.Println("callRPCGetSuccessor Failed ", err)
		return Finger{}, err
	}
	return reply,nil
}

func (node *Finger) callRPCGetPredecessor() (Finger,error){
	if fingerEquals(*myself.nodeFinger,*node){
		//calling on own Server
		var reply Finger
		res := myself.nodeFinger.getPredecessor(nil, &reply)
		return reply, res
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return Finger{}, err
	}
	defer client.Close()
	var reply Finger
	//Debug.Println(reply)
	dummy := 0
	err = client.Call("Finger.GetPredecessor",&dummy,&reply)
	if err!=nil{
		Error.Println("callRPCGetPredecessor Failed ", err)
		return Finger{}, err
	}
	return reply,nil
}

func (node *Finger) callRPCClosestPrecedingFinger(key KeySpace) (Finger,error){
	if fingerEquals(*myself.nodeFinger,*node){
		//calling on own Server
		var reply Finger
		err := myself.nodeFinger.closest_preceding_finger(&key,&reply)
		return reply, err
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return Finger{}, err
	}
	defer client.Close()
	var reply Finger
	err = client.Call("Finger.Closest_preceding_finger",&key,&reply)
	if err!=nil{
		Error.Println("callRPCClosestPrecedingFinger Failed: ", err , " Dialed on: ", addr)
		return Finger{}, err
	}
	return reply,nil
}

func (node *Finger) callRPCFindSuccessor(key KeySpace) (Finger,error){
	if fingerEquals(*myself.nodeFinger,*node){
		//calling on own Server
		var reply Finger
		err := myself.nodeFinger.findSuccessor(&key, &reply)
		return reply, err
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return Finger{}, err
	}
	defer client.Close()
	var reply Finger
	err = client.Call("Finger.FindSuccessor",&key,&reply)
	if err!=nil{
		//Debug.Println("received successor", reply)
		Error.Println("callRPCFindSuccessor Failed ", err)
		return Finger{}, err
	}
	//Debug.Println("received successor", reply)
	return reply,nil
}

func (node *Finger) callRPCSetPredecessor(pred Finger) (error){
	//Debug.Println("Making callRPCSetPredecessor call")
	if fingerEquals(*myself.nodeFinger,*node){
		//calling on own Server
		res := myself.nodeFinger.setPredecessor(&pred,nil)
		return res
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return err
	}
	defer client.Close()
	dummy :=0
	err = client.Call("Finger.SetPredecessor",&pred,&dummy)
	if err!=nil{
		Error.Println("callRPCSetPredecessor Failed ", err)
		return err
	}
	//Debug.Println("Successfully returned from callRPCSetPredecessor ")
	return nil
}

func (node *Finger) callRPCUpdateFingerTable(s Finger, i int) (error){
	arg := UpdateFingerTableArg{S:s, I:i}
	if fingerEquals(*myself.nodeFinger,*node){
		//calling on own Server
		res := myself.nodeFinger.update_finger_table(&arg, nil)
		return res
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return err
	}
	defer client.Close()
	dummy :=0
	err = client.Call("Finger.Update_finger_table",&arg,&dummy)
	if err!=nil{
		Error.Println("callRPCUpdateFingerTable Failed ", err)
		return err
	}
	return nil
}

func (node *Finger) callRPCNotify(req Finger) (error){
	if fingerEquals(*myself.nodeFinger,*node){
		//calling on own Server
		res := myself.nodeFinger.notify(&req, nil)
		return res
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return err
	}
	defer client.Close()
	dummy :=0 
	err = client.Call("Finger.Notify",&req,&dummy)
	if err!=nil{
		Error.Println("RPCCallNotify Failed ", err)
		return err
	}
	return nil
}

func (node *Finger) callRPCGetValue(key string) (string,error){
	if fingerEquals(*myself.nodeFinger,*node){
		res,err := myself.getValue(key)
		return res,err
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return "",err
	}
	defer client.Close()
	var reply string
	err = client.Call("Finger.GetValue",&key,&reply)
	if err!=nil{
		Error.Println("RPCCallGetValue Failed : ", err)
		return "",err
	}
	return reply,err
}

func (node *Finger) callRPCSetValue(key,value string, phase int) (bool, error){
	if fingerEquals(*myself.nodeFinger,*node){
		vote, err := myself.setValue(key, value, phase)
		return vote, err
	}
	//else rpc call
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return false, err
	}
	defer client.Close()
	arg := SetKeyValueArg{K:key, V:value, P:phase}
	var reply bool
	err = client.Call("Finger.SetValue",&arg,&reply)
	if err!=nil{
		Error.Println("RPCCallSetValue Failed ", err)
		return false, err
	}
	return reply, err
}

func (node *Finger) callRPCGetKeys(pred, me Finger) (error){
	addr := node.Ip + ":" + strconv.Itoa(node.Port)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		Error.Println("dialing", err)
		return  err
	}
	defer client.Close()
	arg := GetKeysArg{P:pred, M: me}
	err = client.Call("Finger.GetKeys",&arg,nil)
	if err!=nil{
		Error.Println("callRPCGetKeys Failed ", err)
		return  err
	}
	return err
}

