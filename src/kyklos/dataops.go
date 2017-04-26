package kyklos

import(
	"errors"
)

func (node *nodeState) getValue(key string) (string, error){
	val, ok := node.store[key]
	if !ok{
		return val, errors.New("Key not found")
	}
	return val,nil
}

func (node *nodeState) setValue(key,value string)(error){
	node.store[key] = value
	return nil
}

func (node *nodeState) get(key string) (string,error) {
	handler, err := node.findSuccessor(hasher(key))
	if err !=nil{
		Error.Println("Couldn't find the node for this key")
		return "",err
	}
	value, err := handler.callRPCGetValue(key)
	if err!=nil{
		//retry
	}
	return value, err
}

func (node* nodeState) set(key,value string)(error){
	handler, err := node.findSuccessor(hasher(key))
	if err !=nil{
		Error.Println("Couldn't find the node for this key")
		return err
	}
	err = handler.callRPCSetValue(key, value)
	if err!=nil{
		Error.Println("Setting value of ", key, " failed")
		return err
	}
	return err
}