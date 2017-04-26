// hashing based on ip + port
// find_successor			done
// find_predecessor			done
// create finger_table 		
// route the messsage
// node add
// new add let successor know
// stabalize
// node_predecessor
// transfer keys
// find successor(id)


//convert to read write locks, other packages too
//also using locks on other data variables such as predecessor_finger 

//lock on predecessor and second_succesor finger

package main

import (
	"fmt"
	"os"
	"kyklos"
)

func main(){
	kyklos.Init(os.Stderr,os.Stdout, os.Stdout, os.Stderr)
	for {
		fmt.Printf("Enter Portnumber: ",)
		var portnum int
		_, err := fmt.Scanf("%d", &portnum)
		if err !=nil{
			fmt.Println("Can't Input Portnumber")
			continue
		}
		//take input from user for portnum choice
		init_status := kyklos.InitialiseNode(portnum)
		if init_status != nil{
			fmt.Println("Node could not be initialized. Perhaps choose a new port number")
		} else {
			break
		}
	}
	//Ask for user prompt
	for{
		fmt.Printf("kyklos>", )
		var command string
		_,err := fmt.Scanf("%s", &command)
		if err !=nil{
			continue
		}
		//create a ring 
		if(command == "createRing"){
			if(!kyklos.IsPartOfRing()){
				err := kyklos.CreateRing()
				if err!=nil{
					fmt.Println("Could not create ring")
				}
			}else {
				fmt.Println("Already part of ring. Leave the existing ring first if you want to create a new ring")
			}
		// } else if(join a ring, take address and port of one of the nodes){
			//contact the node, call join function on that node, get relevant info
			// contact joining location, setup
			// mark part of ring
		} else if (command == "join"){
			var ip string
			var port int
			fmt.Scanf("%s %d", &ip, &port)
			if(!kyklos.IsPartOfRing()){
				err = kyklos.Join(ip, port)
				if err!=nil{
					fmt.Println("Could not join")
				}
			}else {
				fmt.Println("Already part of ring. Leave the existing ring first if you want to join a ring")
			}
		}else if (command == "dump"){
			kyklos.Dump()
		}else if (command == "dumptable"){
			kyklos.DumpTable()
		}else if (command == "set"){
			var key string
			var val string
			fmt.Scanf("%s %s", &key, &val)
			ok := kyklos.Set(key,val)
			if(ok !=nil){
				fmt.Println("Set operation failed: ", ok)
			}else {
				fmt.Println("Successfully set")
			}
		} else if(command == "get"){
			var key string
			fmt.Scanf("%s", &key)
			res, ok := kyklos.Get(key)
			if ok!=nil{
				fmt.Println("Get operation failed : ", ok)
			}else {
				fmt.Println("Got Value: ", res)
			}
		}else if (command == "help") {
			fmt.Println("createRing, join, dump, dumptable, get, set, help")
		}else {
			//invalid command, show list of commands to user
		}
	}
}
