package kyklos

import (
	"math/rand"
	"net"
	"errors"
	"net/rpc"
	"strconv"
	"time"
)

func InitialiseNode(portnum int) error{
	return myself.initialiseNode(portnum)
}

func IsPartOfRing()bool{
	return myself.part_of_ring
}

func CreateRing() error{
	return myself.createRing()
}

func Join(ip string, port int) error{
	return myself.join(Finger{Ip:ip, Port:port})
}

func Dump(){
	Debug.Println("OwnFinger: ", *own_finger)
	Debug.Println("Successor: ", myself.finger_table.Fingers[0])
	Debug.Println("Second_Successor: ", myself.second_successor_finger)
	Debug.Println("Predecessor: ", myself.predecessor_finger)
}

func DumpTable(){
	Debug.Println(myself.finger_table)
}

func DumpData(){
	Debug.Println(myself.store)
}

func Set(key, val string) error{
	return myself.set(key,val)
}

func Get(key string)(string, error){
	return myself.get(key)
}

func startServer(addr *net.TCPAddr) error{
	listener, err := net.ListenTCP("tcp",addr)
	CheckError(err)

	for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        rpc.ServeConn(conn)
    }
}

func startProtocols(){
	go doStabilize()
	go doFixFingers()
	go doFailureHandling()
}

func doStabilize(){
 	for{
	    go myself.stabilize()
	    time.Sleep(time.Duration(500)*time.Millisecond)
	}
}

func doFixFingers(){
	for{
	    go myself.fix_fingers()
	    time.Sleep(time.Duration(50)*time.Millisecond)
	}
}

func doFailureHandling(){
	for{
	    go myself.failureHandler()
	    time.Sleep(time.Duration(1000)*time.Millisecond)
	}
}


func (node *nodeState) initialiseNode(portnum int) error{
	//get own's ip address

	addrs, err := net.InterfaceAddrs()
    if err != nil {
    	Error.Println("Couldnt get Interface Addresses")
        return err
    }
	myIP := ""
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                myIP = ipnet.IP.String()
      			break
            }
        }
    }

	if myIP == "" {
		Error.Println("Own IP could not be determined")
		return errors.New("Own IP could not be determined")
	}
    tcpAddr, err2 := net.ResolveTCPAddr("tcp", myIP + ":" + strconv.Itoa(portnum))
    if err2 != nil{
    	Error.Println("tcp resolution failed in initialisation")
    	return err2
    }

    go startServer(tcpAddr)
	own_finger = new(Finger)
	own_finger.Ip=myIP 
	own_finger.Port = portnum


	node.nodeFinger = own_finger
	
	node.second_successor_finger = Finger{}
	node.predecessor_finger = Finger{}
	node.finger_table_size = 256+1
	node.hashbits = 256
	node.rf = 5
	node.store = make(map[string]string)
	node.tempstore = make(map[string]string)
	
	node.finger_table = FingerTable{Fingers: make([]Finger,node.finger_table_size), Valid: make([]bool, node.finger_table_size)} 
	for i:=0;i<node.finger_table_size;i++{
		node.finger_table.Valid[i] = false
	}

	node.part_of_ring = false

	Info.Println("Node Initialized : OwnFinger:",myself.nodeFinger,"\n part of ring: ",myself.part_of_ring,"\n" )
    err3 := rpc.Register(own_finger)
    if err3 !=nil{
    	Error.Println("Could not register rpc in initialize")
    	CheckError(err3)
    	return err3
    }
	// go startProtocols()
	return nil
}

func (node *nodeState) createRing() error{
	node.finger_table_lock.Lock()
	defer node.finger_table_lock.Unlock()
	for i:=0;i<node.finger_table_size;i++ {
		node.finger_table.Fingers[i] = *node.nodeFinger
		node.finger_table.Valid[i] = true
	}
	node.second_successor_finger = *node.nodeFinger
	node.predecessor_finger = *node.nodeFinger
	node.part_of_ring = true
	go startProtocols()
	Info.Println("Create Ring Completed : part of ring: ", node.part_of_ring, "\n predecessor : " ,node.predecessor_finger, "\n successor: ",node.finger_table.Fingers[0])
	return nil
}

func (node *nodeState) getSuccessor() (Finger, error){
	node.finger_table_lock.Lock()
	defer node.finger_table_lock.Unlock()
	return node.finger_table.Fingers[0], nil
}

func (node *nodeState) findSuccessor(key KeySpace) (Finger,error){ 		//used when called by other node, this method is exported
	//Debug.Println("node level find successor called")
	predecessor_of_key, err := node.findPredecessor(key) 	
	if err!=nil{
		Error.Println("Error in finding successor of ", key, ": ", err)
		return Finger{},err
	}
	// res , err2 := predecessor_of_key.getSuccessor() // getSuccessor is an RPC call here
	res , err2 := predecessor_of_key.callRPCGetSuccessor() // getSuccessor is an RPC call here
	if err2!=nil{
		Error.Println("In findSuccessor: Predecessor of key failed to return its successor")
		return Finger{},err2
	}
	return res,nil
}

func (node *nodeState) findPredecessor(key KeySpace) (Finger, error){
	//Debug.Println("node level find Predecessor called")
	currNode := *node.nodeFinger
	// prevNode := *node.nodeFinger
	// for ;not(key>currNode && key<=currNode.successor) ; { 
	for {
		//Debug.Println("In loop of findPredecessor node level")
		successor, err := currNode.callRPCGetSuccessor() 			//rpc call
		if err !=nil {
			Error.Println("findPredecessor on node failed")
			return Finger{}, err
		}
		if (betweenRightIncl(hashFunc(currNode), hashFunc(successor), key)){
			break
		} 
		// 1. comparision among keys in above line should be accordance to sha256 values comparision
		// 2. currNode. successor is an RPC call
		// 3. closes_preceding_finger is also an RPC call
		// 4. How are RPC calls handles on own server 
		
		//if the above calls relating to successor fail, that is the currNode doesn;t actually exist, 
		//send msg to prevNode to declare currNode invalid  and revertback currNode to prevNode
		//If the currnode, after reverting back, fails again, it means that this node failed when we were
		// on next Node. In this case even the prevNode will be set to currNode.
		//In this case nothing can be done except to go into this loop from starting

		//else 
		// prevNode = currNode
		// currNode = currNode.closest_preceding_finger(key)  	//rpc call which calls the below closes_preceding_finger
		currNode,err = currNode.callRPCClosestPrecedingFinger(key)  	//rpc call which calls the below closes_preceding_finger
		// If for some reason the call on 
	} 
	return currNode,nil;
}

func (node *nodeState) closest_preceding_finger(key KeySpace) (Finger, error){
	own_id := hashFunc(*node.nodeFinger)

	node.finger_table_lock.Lock()
	defer node.finger_table_lock.Unlock()
	
	for i := node.finger_table_size-1; i>=0;i--{
		node_id := hashFunc(node.finger_table.Fingers[i])
		if (node.finger_table.Valid[i]) && between(own_id, key, node_id){ //(node_id lies between own_id and key){ 		//handle comparision
			return node.finger_table.Fingers[i],nil
		}
	} 
	return *node.nodeFinger,nil
}

func (node *nodeState) join(participant Finger ) error{
	err := node.init_finger_table(participant)
	//Debug.Println("In join: Finger Table initialized")
	//Debug.Println(node.finger_table)
	if err !=nil{
		Error.Println("join in node failed: Call to init_finger_table failed")
		return err
	}

	err = node.getKeys()
	if err!=nil{
		Error.Println("getKeys in node failed")
		return err
	}

	err2 := node.update_others()
	//Debug.Println("Completed asking others to update their finger_tables")
	if err2 !=nil{
		Error.Println("join in node failed: Call to update_others failed")
		return err2
	}
	node.part_of_ring = true
	go startProtocols()
	// move keys in (predecessor, n] from succesor
	// let the successor know that now node is its predecessor
	// successor sets its predecessor node and let's its older predecessor know 
	// to set node as its successor
	// Info.Println("Join Completed : Finger Table: \n", node.finger_table)
	Info.Println("Join Completed ")
	return nil
}

func (node *nodeState) getKeys() error {
	node.finger_table_lock.Lock()
	successor := node.finger_table.Fingers[0]
	node.finger_table_lock.Unlock()
	err := successor.callRPCGetKeys(node.predecessor_finger, *node.nodeFinger)
	if err!=nil{
		Error.Println("get keys failed ", err)
		return err
	}
	return nil
}

func (node *nodeState) init_finger_table(participant Finger) error {
	//Debug.Println("Entering init_finger_table")
	successor,err:=participant.callRPCFindSuccessor(hashFunc(*node.nodeFinger)) 			//rpc call
	//Debug.Println("host returned the finger to contact")
	if err!=nil{
		Error.Println("init_finger_table in node failed : Failed call to callRPCFindSuccessor")
		return err
	}
	node.finger_table_lock.Lock()

	node.finger_table.Valid[0] = true
	node.finger_table.Fingers[0] = successor

	node.finger_table_lock.Unlock()

	second_successor, err2 := successor.callRPCGetSuccessor()
	if(err2 != nil){
		Error.Println("init_finger_table in node failed : Failed call to callRPCGetSuccessor for setting second_successor : ", err2)
		return err2
	}

	node.second_successor_finger = second_successor

	//Debug.Println("Successor set, asking successor for his predecessor")
	node.predecessor_finger, err = successor.callRPCGetPredecessor()		//rpc call
	//Debug.Println("call for predecessor completed")
	if err !=nil {
		Error.Println("init_finger_table in node failed: Failed call to callRPCGetPredecessor")
		return err
	}

	//Debug.Println("Letting successor know to update his predecessor ")
	successor.callRPCSetPredecessor(*node.nodeFinger)		//rpc call
	//Debug.Println("Completed call successor know to update his predecessor ")

	//Debug.Println("Filling up the finger table")
	my_id := hashFunc(*node.nodeFinger)
	for i:=1;i<node.finger_table_size;i++{
		//Debug.Println("Filling entry ", i)
		// target_id := my_id + power(2,i)
		target_id := powerOffset(my_id, i, node.hashbits)
		// if  (target_id lies between my_id and node.finger_table.Fingers[i-1].Hash) {
		if  (betweenLeftIncl(my_id,hashFunc(node.finger_table.Fingers[i-1]), target_id)) {
			node.finger_table.Fingers[i] = node.finger_table.Fingers[i-1]
		} else{
			//Debug.Println("rpc calling for successor")
			node.finger_table.Fingers[i],err = participant.callRPCFindSuccessor(target_id)
			//Debug.Println("rpc calling for successor completed")
		}
		node.finger_table.Valid[i] = true
	}
	//Debug.Println("Exiting finger_table without error")
	return nil
}

func (node *nodeState) update_others() error{
	for i:=0;i<node.finger_table_size;i++{
		// p:= node.findPredecessor(node.nodeFinger.Hash - power(2,i)) 	//p is a finger
		p,err:= node.findPredecessor(negativePowerOffset(hashFunc(*node.nodeFinger),i,node.hashbits)) 	//p is a finger
		if err !=nil{
			Error.Println("update_others in node failed: Call to findPredecessor failed")
		}
		// p.update_finger_table(*node.nodeFinger,i)			//rpc call which calls below defined function
		err = p.callRPCUpdateFingerTable(*node.nodeFinger,i)			//rpc call which calls below defined function
		if err !=nil{
			Error.Println("update_others in node failed: Call to callRPCUpdateFingerTable failed")
		}
	}
	return nil
}

func (node *nodeState) update_finger_table(s Finger,i int) error{
	my_id := hashFunc(*node.nodeFinger)
	node.finger_table_lock.Lock()
	if  (node.finger_table.Valid[i]){
		node.finger_table_lock.Unlock()
		return nil
	}
	//Debug.Println("Entered CS in update_finger_table")
	target_id := hashFunc(node.finger_table.Fingers[i])
	// if (s.Hash >= my_id && s.Hash < target_id) && (node.finger_table.Valid[i]){ 			//Check if this validity needs to be checked or not
	if (betweenLeftIncl(my_id, target_id, hashFunc(s))){ 			//Check if this validity needs to be checked or not
		node.finger_table.Valid[i] = true
		node.finger_table.Fingers[i] = s
	}
	//Debug.Println("Exiting CS in update_finger_table")
	node.finger_table_lock.Unlock()
	// node.predecessor_finger.update_finger_table(s,i) 			//RPC Call
	//Debug.Println("Exiting CS in update_finger_table")
	err := node.predecessor_finger.callRPCUpdateFingerTable(s,i) 			//RPC Call
	if err !=nil{
			Error.Println("update_finger_table in node failed: Call to callRPCUpdateFingerTable failed")
	}
	Info.Println("Finger Table Updated: s: ", s ," \n i: ", i)
	return nil
}

// func (node *nodeState) handleFailure() error{
// 	//ping the successor by asking him of his successor
// 	//if ping fails, proceed with handling failure
// 	var successor Finger
// 	node.finger_table_lock.Lock()
// 	successor = node.finger_table.Fingers[0]
// 	node.finger_table_lock.Unlock()

// 	x,err:= successor.callRPCGetPredecessor()
// 	if (x == nil) || (err != nil){
// 		second_successor := node.second_successor_finger
// 		new_second_successor, err := second_successor.callRPCGetSuccessor()
// 		if err!=nil{

// 		}
// 		err2 := second_successor.
// 	}
// 	return nil
// }

func (node *nodeState) failureHandler() error{
	var successor Finger
	node.finger_table_lock.Lock()
	successor = node.finger_table.Fingers[0]
	node.finger_table_lock.Unlock()

	nx, err2 := successor.callRPCGetSuccessor()
	if (fingerEquals(nx, Finger{}) )|| (err2 != nil ){
			Error.Println("Call to Successor Failed. Handle Failure")
			node.finger_table_lock.Lock()
			node.finger_table.Fingers[0] = node.second_successor_finger
			node.finger_table_lock.Unlock()
			err3 := node.second_successor_finger.callRPCSetPredecessor(*own_finger)
			if err3 !=nil{
				Error.Println("Error in setting predecessor of new successor: ", err3)
				return err3
			}
			return err2
	} else {
		node.second_successor_finger = nx
	}
	return nil
}

//stabilize runs periodically
func (node *nodeState) stabilize() error {
	var successor Finger
	node.finger_table_lock.Lock()
	successor = node.finger_table.Fingers[0]
	node.finger_table_lock.Unlock()
	// x:= succesor.getPredecessor()		//rpc call
	x,err:= successor.callRPCGetPredecessor()		//rpc call
	if err!=nil{
		Error.Println("Predecessor call on successor failed: ", err)
		return err
	}  
	// Debug.Println("stabilize : pred of successor = ", x)
	// if (x.Hash lies between node.nodeFinger.Hash and successor.Hash){
	if (between(hashFunc(*node.nodeFinger),hashFunc(successor),hashFunc(x))){
		node.finger_table_lock.Lock()
		node.finger_table.Fingers[0] = x
		successor = x
		node.finger_table_lock.Unlock()
	}
	// successor.notify(*node.nodeFinger)		//rpc call which calls below defined function
	err  = successor.callRPCNotify(*node.nodeFinger)		//rpc call which calls below defined function
	if err !=nil{
			Error.Println("stabilize in node failed: Call to callRPCNotify failed")
			return err
	}

	return nil
}

func (node *nodeState) notify( prospective Finger) error{
	if (node.predecessor_finger.isnil()) ||(between(hashFunc(node.predecessor_finger),hashFunc(*node.nodeFinger),hashFunc(prospective))){
			node.predecessor_finger = prospective
		}
	return nil
}

//fix_fingers runs periodically to refresh finger table entries
func (node *nodeState) fix_fingers() error{
	// i:= (rand number between 1 and node.finger_table_size-1)
	i:= rand.Intn(node.finger_table_size-1)+1
	// target_id = node.nodeFinger.Hash + power(2,i)
	target_id := powerOffset(hashFunc(*node.nodeFinger), i, node.hashbits)
	fresh_successor,err := node.findSuccessor(target_id)
	if err !=nil{
			Error.Println("fix_fingers in node failed: Call to findSuccessor failed: ", err)
			return err
	}
	node.finger_table_lock.Lock()
	node.finger_table.Fingers[i] = fresh_successor
	node.finger_table_lock.Unlock()
	return nil
}

func (node *nodeState) sendKeys(trgt_pred, trgt Finger ) error{
	toDelete := []string{}
	for key,val := range node.store{
		hashedKey:= hasherFunc(key)

		if betweenRightIncl(hashedKey, hashFunc(trgt_pred), hashFunc(trgt) ){
			if _, err := trgt.callRPCSetValue(key,val,1);err!=nil{
				return err
			}
			toDelete = append(toDelete, key)
		}
	}
	for _,key := range toDelete{
		delete(node.store, key)
	}
	return nil
}