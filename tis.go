package main

import (
	"fmt"
	"time"
	"tis/node"
)

func main() {
	node_ch := make(chan int)
	clock_ch := make(chan bool)
	start_ch := make(chan bool)
	stop_ch := make(chan bool)
	nodes := []node.node{node.node{right: node_ch}, node.node{left: node_ch}}
	noderunners := []node.runner{node.runner{clock: clock_ch, node: nodes[0]},
		node.runner{node: nodes[1]}}
	noderunners[0].loadfile("left.tis")
	go nodeRunners[0].RunNode()
	go ClockTicker(clock_ch, start_ch, stop_ch)
	start_ch <- true
	time.Sleep(time.Second * 5)
	fmt.Println("Stopping")
	stop_ch <- true
	time.Sleep(time.Second * 3)
	fmt.Println("Starting")
	start_ch <- true
	time.Sleep(time.Second * 3)
	fmt.Println("Stopping")
	stop_ch <- true
	time.Sleep(time.Second * 3)
}
