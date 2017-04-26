package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/r-medina/gmaj"
)

type command func(nodes []*gmaj.Node, args ...string) bool

var allCmds []string

const IND = 10

var cmds = map[string]command{
	"quit": func(nodes []*gmaj.Node, args ...string) (stop bool) {
		fmt.Println("goodbye")

		stop = true
		return
	},

	"node": func(nodes []*gmaj.Node, args ...string) (stop bool) {
		for _, node := range nodes {
			fmt.Println(node)
		}
		return
	},

	"table": func(nodes []*gmaj.Node, args ...string) (stop bool) {
		for _, node := range nodes {
			fmt.Printf("%s: %s\n", gmaj.IDToString(node.Id), node.FingerTableString())
		}
		return
	},

	"addr": func(nodes []*gmaj.Node, args ...string) (stop bool) {
		for _, node := range nodes {
			fmt.Println(node.Addr)
		}
		return
	},

	"data": func(nodes []*gmaj.Node, args ...string) (stop bool) {
		for _, node := range nodes {
			fmt.Println(node.DatastoreString())
		}
		return
	},

	"get": func(nodes []*gmaj.Node, args ...string) (stop bool) {
		if len(args) > 0 {
			for {
				num := strconv.Itoa(rand.Int()%IND)
				val, err := gmaj.Get(nodes[0], args[0]+"_"+num)
				if err != nil {
					fmt.Println(err)
					break
				} else if val == nil {
					n := rand.Int()%100
					time.Sleep(time.Duration(n)*time.Millisecond)
				} else {
					fmt.Printf("%s\n", val)
					break
				}
			}
		}
		return
	},

	"put": func(nodes []*gmaj.Node, args ...string) (stop bool) {
		if len(args) > 1 {
			for {
				// Phase 1
				final := true
				var ver int32
				ver = 0
				for i:=0; i<IND; i++ {
					num := strconv.Itoa(i)
					vote, _ := gmaj.Put(nodes[0], args[0]+"_"+num, []byte(args[1]), ver)
					if vote == false {
						final = false
					}
				}

				// Phase 2
				if final == false {
					ver = 2
				} else {
					ver = 1
				}
				for i:=0; i<IND; i++ {
					num := strconv.Itoa(i)
					gmaj.Put(nodes[0], args[0]+"_"+num, []byte(args[1]), ver)
				}

				// Repeat
				if ver == 1 {
					break
				} else {
					n := rand.Int()%100
					time.Sleep(time.Duration(n)*time.Millisecond)
				}
			}
		}
		return
	},

	"help": func(nodes []*gmaj.Node, args ...string) (stop bool) {
		fmt.Printf("available commands: %v\n", allCmds)
		return
	},
}

func init() {
	allCmds = commands()
}

func commands() []string {
	out := make([]string, 0, len(cmds))
	for cmd := range cmds {
		out = append(out, cmd)
	}

	return out
}
