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
	"time"
	"math/rand"
	"strconv"
)

func main(){
	kyklos.Init(os.Stderr,os.Stdout, os.Stdout, os.Stderr)
	var portnum int
	for {
		fmt.Printf("Enter Portnumber: ",)
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
	pcFile,_:=os.OpenFile("2PC.log"+strconv.Itoa(portnum), os.O_RDWR | os.O_CREATE, 0666)
	scFile,_:=os.OpenFile("Consistency.log"+strconv.Itoa(portnum), os.O_RDWR | os.O_CREATE, 0666)
	kyklos.InitFileLogs(pcFile, scFile)

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
		}else if (command == "dumpdata"){
			kyklos.DumpData()
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
		} else if (command == "test"){
				kyklos.Set("1","1")
				kyklos.Set("2","2")
				kyklos.Set("3","3")
				kyklos.Set("4","4")
				kyklos.Set("5","5")
				kyklos.Set("6","6")
				kyklos.Set("7","7")
				kyklos.Set("8","8")
				kyklos.Set("9","9")
				kyklos.Set("0","0")
				avgget := time.Duration(0)
				for i:=0; i<1000; i++{
					num := rand.Int()%10
					start := time.Now()
					// cmdget(nodes, strconv.Itoa(num))
					kyklos.Get(strconv.Itoa(num))
					elapsed := time.Since(start)
					avgget += elapsed
				}
				fmt.Println("Average time taken for get operation: %v", avgget/1000.0)

				avgput := time.Duration(0)
				for i:=0; i<1000; i++{
					num1 := rand.Int()%100
					num2 := rand.Int()%100
					start := time.Now()
					err := kyklos.Set(strconv.Itoa(num1), strconv.Itoa(num2))
					if err!=nil{
						fmt.Println("test set failed for id ", num1)
					}
					elapsed := time.Since(start)
					avgput += elapsed
					time.Sleep(time.Duration(100)*time.Microsecond)
				}
				fmt.Println("Average time taken for put operation: %v", avgput/1000.0)
		}else if(command == "keytest"){
			keyVal := strconv.Itoa(42)
			for i:=0;i<5;i++{
				val := strconv.Itoa(rand.Intn(200))
				kyklos.Set(keyVal, val)
			}
		}else if(command == "get"){
			var key string
			fmt.Scanf("%s", &key)
			res, ok := kyklos.Get(key)
			if ok!=nil{
				fmt.Println("Get operation failed : ", ok)
			}else {
				fmt.Println("Got Value: ", res)
			}
		}else if (command == "help") {
			fmt.Println("createRing, join, dump, dumptable,dumpdata, get, set, help")
		}else {
			//invalid command, show list of commands to user
		}
	}
}
