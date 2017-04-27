package kyklos

import(
	"errors"
	"math/rand"
	"strconv"
	"time"
)

func (node *nodeState) getValue(key string) (string, error){
	node.sm.RLock()
	val, ok := node.store[key]
	node.sm.RUnlock()
	if !ok{
		return "", errors.New("Key not found")
	}

	node.tm.RLock()
	_, ok = node.tempstore[key]
	node.tm.RUnlock()
	if ok{
		return "",errors.New("Key unavailable")
	}
	Consis.Println("Read: Key: ", key, " Value: ", val)
	return val,nil
}

func (node *nodeState) setValue(key,value string, phase int)(bool, error){
	if phase == 0 {
		node.tm.RLock()
		_, ok := node.tempstore[key]
		node.tm.RUnlock()
		if ok {
			return false, errors.New("Key undergoing another write")
		}

		node.tm.Lock()
		node.tempstore[key] = value
		node.tm.Unlock()

		return true, nil
	} else if phase == 1 {
		node.tm.Lock()
		delete(node.tempstore, key)
		node.tm.Unlock()

		node.sm.Lock()
		node.store[key] = value
		node.sm.Unlock()

		Consis.Println("Write: Key: ", key, " Value: ", value)
		return true, nil
	} else if phase == 2 {
		node.tm.Lock()
		delete(node.tempstore, key)
		node.tm.Unlock()

		return true, nil	
	} else {
		return false, errors.New("Unidentified phase in 2pc")
	}

	return false, nil
}

func (node *nodeState) get(key string) (string,error) {
	rf:= node.rf
	idxs := rand.Perm(rf)
	for idx:=range idxs{
		i := idx
		// Debug.Println("Checking on node with idx ", i)
		combined_key := key + "_" + strconv.Itoa(i)
		handler, err := node.findSuccessor(hasherFunc(combined_key))
		// Debug.Println(handler)
		if err !=nil{
			// Error.Println("Couldn't find the node for this key")
			continue
		}
		value, err := handler.callRPCGetValue(combined_key)
		if err!=nil{
			continue
		}else{
			//update the other replicas which failed till now
			for failedidx :=0;failedidx<idx ; failedidx ++{
				combined_key2 := key + "_" + strconv.Itoa(failedidx)
				handler, _ := node.findSuccessor(hasherFunc(combined_key2))
				handler.callRPCSetValue(combined_key2, value, 1)
			}
			return value,err
		}
	}
	return "", errors.New("Key not found in database")
}

func (node* nodeState) sendAbort(key,value string, idx int){
	for i:=0;i<idx;i++{
		combined_key := key + "_" + strconv.Itoa(i)
		handler, _ := node.findSuccessor(hasherFunc(combined_key))
		go handler.callRPCSetValue(combined_key, value,2) 			//is it okay to do this?
	}
}

func (node* nodeState) set(key,value string)(error){

	TwoPC.Println("Initiating Phase 0 (ask) for Key: ", key, " Value: ",value )
	for {
		voteChannel :=make(chan bool, node.rf+5)
		for i:=0;i<node.rf;i++{
			combined_key := key + "_" + strconv.Itoa(i)

			go func(voteChannel chan bool, combined_key string){
				handler, err := node.findSuccessor(hasherFunc(combined_key))
				if err !=nil{
					// Error.Println("Couldn't find the node for this key")
					voteChannel<-false
					return 
				}
				vote,err := handler.callRPCSetValue(combined_key, value, 0)
				voteChannel <-vote
				if err!=nil{
					// Error.Println("Setting value of ", key , ", idx", i, " in phase 0 failed")
					return 
				}
			}(voteChannel, combined_key)

		}

		//wait on channel to collect the votes
		//If somebody says no send abort to everyone

		proceed:=true
		for i:=0;i<node.rf;i++{
			vote :=<-voteChannel
			if!vote{
				proceed = false
				break
			}
		}
		if proceed{
			TwoPC.Println("Proceeding to Phase 1 (Commit) for Key: ", key, " Value: ",value )
			//reached here -> all replied yes, send commit to everyone
			for i:=0;i<node.rf;i++{
				combined_key := key + "_" + strconv.Itoa(i)

				go func(combined_key string){
					// Debug.Println(combined_key)
					handler, err := node.findSuccessor(hasherFunc(combined_key))
					// Debug.Println(handler)
					if err !=nil{
						// Error.Println("Couldn't find the node for this key")
					}
					_,err = handler.callRPCSetValue(combined_key, value, 1)
					if err!=nil{
						// Error.Println("Setting value of ", key , ", idx", i, " in phase 1 failed")
					}
				}(combined_key)
			}
			break
		} else {
			TwoPC.Println("Negative vote for Key: ", key, " Value: ",value )
			TwoPC.Println("Aborting for now" )
			node.sendAbort(key, value,node.rf)
		}
		backoff := rand.Intn(31)
		backoff = backoff+10
		TwoPC.Println("Going into backoff for ", backoff,"ms for Key: ", key, " Value: ",value )
		time.Sleep(time.Duration(backoff)*time.Millisecond)
	}
	return nil
}

//periodic replica maintainance?