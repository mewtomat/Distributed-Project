package kyklos

import (
	"crypto/sha256"
	"strconv"
)

//Own Global State
var own_finger *Finger
var myself nodeState


// func HashFunc(Ip string , port int) KeySpace {
// 	hasher := sha256.New()
// 	combined := Ip + strconv.Itoa(port)
// 	return KeySpace{Data:hasher.Sum([]byte(combined))}
// }

func hashFunc(finger Finger) KeySpace {
	hasher := sha256.New()
	combined := finger.Ip + strconv.Itoa(finger.Port)
	return KeySpace{Data:hasher.Sum([]byte(combined))}
}
