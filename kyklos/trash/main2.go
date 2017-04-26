// package main

// import (
// 	"crypto/sha1"
// 	"fmt"
// )

// func main() {
		
// 	hasher := sha1.New()
// 	data := []byte("test")

// 	fmt.Printf("%x\n", hasher.Sum(data))
// 	fmt.Printf("%x\n", hasher.Sum(data))

// 	hasher.Write(data)
// 	fmt.Printf("%x\n", hasher.Sum(nil))

// 	hasher.Write(nil)
// 	hasher.Write(data)
// 	fmt.Printf("%x\n", hasher.Sum(nil))
// }