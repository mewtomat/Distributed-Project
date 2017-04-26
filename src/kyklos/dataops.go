package kyklos

import(
	"math/rand"
)

func (node *nodeState) getValue(key string) (string, error){
	val, ok := node.store[key]
	return val,ok
}

func (node *nodeState) setValue(key,value string)(error){
	node.store[key] = val
	return nil
}

func (node *nodeState) get(key string) (string,error) {
	handler, err := node.findSuccessor(key)
	if err !=nil{
		Error.Println("Couldn't find the node for this key")
		return err
	}
	value, err := node.callRPCGetValue(key)
	if err!=nil{
		//retry
	}
	return value, err
}

func (node* nodeState) set(key,value string)(error){
	handler, err := node.findSuccessor(key)
	if err !=nil{
		Error.Println("Couldn't find the node for this key")
		return err
	}
	err = node.callRPCSetValue(key)
	if err!=nil{
		Error.Println("Setting value of ", key, " failed")
		return err
	}
	return err
}