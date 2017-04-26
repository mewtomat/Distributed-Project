package kyklos

func (node *Finger) getSuccessor( _ , reply *Finger) error {
	myself.finger_table_lock.Lock()
	defer myself.finger_table_lock.Unlock()
	successor := myself.finger_table.Fingers[0] 
	*reply = successor
	return nil
}

func (node *Finger) getPredecessor( _ , reply *Finger) error{
	predecessor := myself.predecessor_finger
	*reply = predecessor
	return nil
}

func (node *Finger) closest_preceding_finger(req *KeySpace, reply *Finger) error{
	res, ok := myself.closest_preceding_finger(*req)
	if ok != nil {
		return ok
	}
	*reply = res
	return nil
}

func (node *Finger) findSuccessor(req *KeySpace , reply *Finger) error{
	res, ok := myself.findSuccessor(*req)
	if ok !=nil{
		return ok
	}
	*reply = res
	return nil
}

func (node *Finger) setPredecessor(req *Finger, _ *struct{}) error{
	//acquire predecessor lock
	myself.predecessor_finger = *req
	return nil
}

func (node *Finger) update_finger_table(req *UpdateFingerTableArg, _ *struct{} ) error{
	myself.update_finger_table(req.S, req.I)
	return nil
}

func (node *Finger) notify(req *Finger, _ *struct{}) error{
	myself.notify(*req)
	return nil
}



func (node *Finger) isnil() bool {
	return node.Ip == "" && node.Port ==0
}
