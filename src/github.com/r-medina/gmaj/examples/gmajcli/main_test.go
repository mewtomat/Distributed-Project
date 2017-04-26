package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	// "strings"
	"time"
	"math/rand"
	"strconv"

	// "github.com/Bowery/prompt"
	"github.com/r-medina/gmaj"
	"github.com/r-medina/gmaj/gmajpb"
)

const promptStr = "gmaj> "

func main() {
	count := flag.Int(
		"count", 1, "Total number of Chord nodes to start up in this process",
	)
	parentAddr := flag.String(
		"parent-addr", "", "Address of a node in the Chord ring you wish to join",
	)
	parentID := flag.String(
		"parent-id", "", "ID of a node in the Chord ring you wish to join",
	)
	addr := flag.String(
		"addr", "", "Address to listen on",
	)

	flag.Parse()

	var parent *gmajpb.Node
	if *parentAddr == "" {
		parent = nil
	} else {
		val := big.NewInt(0)
		val.SetString(*parentID, 10)
		parent = &gmajpb.Node{
			Id:   val.Bytes(),
			Addr: *parentAddr,
		}
		fmt.Printf(
			"Attach this node to id:%v, addr:%v\n",
			gmaj.IDToString(parent.Id), parent.Addr,
		)
	}

	nodes := make([]*gmaj.Node, *count)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c

		shutdown(nodes)

		os.Exit(1)
	}()

	var err error
	for i := range nodes {
		nodes[i], err = gmaj.NewNode(parent, gmaj.WithAddress(*addr))
		if err != nil {
			fmt.Println("Unable to create new node!")
			log.Fatal(err)
		}
		parent = nodes[i].Node

		fmt.Printf(
			"Created -id %v -addr %v\n",
			gmaj.IDToString(nodes[i].Id), nodes[i].Addr,
		)
	}

	cmds["help"](nil)

	avgget := time.Duration(0)
	for i:=0; i<1000; i++{
		cmdput := cmds["put"]
		cmdput(nodes, "1", "1")
		cmdput(nodes, "2", "2")
		cmdput(nodes, "3", "3")
		cmdput(nodes, "4", "4")
		cmdput(nodes, "5", "5")
		cmdput(nodes, "6", "6")
		cmdput(nodes, "7", "7")
		cmdput(nodes, "8", "8")
		cmdput(nodes, "9", "9")
		cmdput(nodes, "0", "0")

		num := rand.Int()%10
		cmdget := cmds["get"]
		start := time.Now()
		cmdget(nodes, strconv.Itoa(num))
		elapsed := time.Since(start)
		avgget += elapsed
	}
	fmt.Println("Average time taken for get operation: %v", avgget/1000.0)

	avgput := time.Duration(0)
	for i:=0; i<1000; i++{
		num1 := rand.Int()%100
		num2 := rand.Int()%100
		cmdput := cmds["put"]
		start := time.Now()
		cmdput(nodes, strconv.Itoa(num1), strconv.Itoa(num2))
		elapsed := time.Since(start)
		avgput += elapsed
	}
	fmt.Println("Average time taken for put operation: %v", avgput/1000.0)

	shutdown(nodes)
}

func shutdown(nodes []*gmaj.Node) {
	fmt.Println("shutting down...")

	for _, node := range nodes {
		node.Shutdown()
	}
}
