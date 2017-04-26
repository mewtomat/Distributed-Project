package kyklos

import "sync"

type Finger struct{ 		//same as a struct Node which is visible to others
	Ip string
	Port int
	// Hash KeySpace
}

type FingerTable struct{
	Fingers []Finger
	Valid []bool
}

type KeySpace struct{ 		//define it later
	Data []byte
}

type nodeState struct{
	nodeFinger *Finger
	second_successor_finger Finger
	predecessor_finger Finger
	finger_table FingerTable
	finger_table_lock sync.Mutex
	// ring_size int 	// assume constant for all rings
	finger_table_size int
	part_of_ring bool
	hashbits int
}
