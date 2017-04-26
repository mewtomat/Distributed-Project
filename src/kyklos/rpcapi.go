package kyklos

type UpdateFingerTableArg struct{
	S Finger
	I int
}

type SetKeyValueArg struct{
	K string
	V string
}

func (node *Finger) GetSuccessor( dummy *int , reply *Finger) error {
	myself.finger_table_lock.Lock()
	defer myself.finger_table_lock.Unlock()
	successor := myself.finger_table.Fingers[0] 
	*reply = successor
	return nil
}

func (node *Finger) GetPredecessor( dummy *int , reply *Finger) error{
		//Debug.Println("GetPredecessor  called from outside")
	predecessor := myself.predecessor_finger
	*reply = predecessor
		//Debug.Println("Returning from getPredecessor called from outside, returning value ", *reply)
	return nil
}

func (node *Finger) Closest_preceding_finger(req *KeySpace, reply *Finger) error{
	//Debug.Println("Closest_preceding_finger called from outside")
	res, ok := myself.closest_preceding_finger(*req)
	if ok != nil {
		//Debug.Println("Closest_preceding_finger returning with error called from outside")
		return ok
	}
	*reply = res
	//Debug.Println("Closest_preceding_finger returning without error called from outside, retuning value ", *reply)
	return nil
}

func (node *Finger) FindSuccessor(req *KeySpace , reply *Finger) error{
	//Debug.Println("Request for sucessor received from outside")
	res, ok := myself.findSuccessor(*req)
	if ok !=nil{
		//Debug.Println("returned with error, from findSuccessor called from outside")
		return ok
	}
	*reply = res
	//Debug.Println("returned without error , returning ", *reply,"  and exiting findSuccessor called from outside")
	return nil
}

func (node *Finger) SetPredecessor(req *Finger,dummy *int) error{
	//acquire predecessor lock
	//Debug.Println("Set predecessor called from outside")
	myself.predecessor_finger = *req
	//Debug.Println("predecessor set to ", myself.predecessor_finger, " and returning from SetPredecessor called from outsid")
	return nil
}

func (node *Finger) Update_finger_table(req *UpdateFingerTableArg, dummy *int ) error{
	//Debug.Println("Update_finger_table called from outside")
	myself.update_finger_table(req.S, req.I)
	//Debug.Println("Returning from Update_finger_table called from outside")
	return nil
}

func (node *Finger) Notify(req *Finger, dummy *int) error{
	//Debug.Println("Notify called from outside")
	myself.notify(*req)
	//Debug.Println("Exiting notify called from outside")
	return nil
}

func (node *Finger) GetValue(req *string,reply *string) error{
	res, err := myself.getValue(*req)
	*reply = res
	return err
}

func (node *Finger) SetValue(req *SetKeyValueArg,_ *struct{}) error{
	err := myself.setValue(req.K, req.V)
	return err
}