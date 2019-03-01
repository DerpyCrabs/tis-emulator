package main

import (
	"flag"
	"fmt"
	"time"
	"tis/node"

	termbox "github.com/nsf/termbox-go"
)

func eventLoop(quit chan error, redraw chan bool) {
	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			if event.Ch != 0 {
				//redraw <- true
			} else if event.Key == termbox.KeyEsc {
				quit <- nil
				break
			}
		}
	}
}

func redrawLoop(quit chan error, redraw chan bool, runners [][]node.Runner) {
	for {
		<-redraw
		err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		if err != nil {
			quit <- err
			break
		}
		width, height := termbox.Size()
		nodeWidth := int(width/len(runners)) - 1
		nodeHeight := int(height/len(runners[0])) - 1
		for row := range runners {
			for col := range runners[row] {
				drawNode(runners[col][row], nodeWidth*row, nodeHeight*col, nodeWidth, nodeHeight)
			}
		}
		err = termbox.Flush()
		if err != nil {
			quit <- err
			break
		}
	}
}

func redrawTicker(tick chan bool) {
	ticker := time.Tick(time.Millisecond * 100)
	for {
		<-ticker
		tick <- true
	}
}

func clockTicker(tick, start, stop chan bool) {
	var ticker <-chan bool = nil
	start_copy := start
	stop_copy := stop
	ticker_copy := make(chan bool)
	go func() {
		for {
			ticker_copy <- true
			time.Sleep(time.Millisecond * 1000)
		}
	}()
	for {
		select {
		case <-stop:
			start = start_copy
			stop = nil
			ticker = nil
		case <-start:
			start = nil
			stop = stop_copy
			ticker = ticker_copy
		case <-ticker:
			tick <- true
		}
	}
}

func initNodeRunners(rows, cols int) [][]node.Runner {
	runners := make([][]node.Runner, rows)
	for row := range runners {
		runners[row] = make([]node.Runner, cols)
		for col := range runners[row] {
			var right_ch, left_ch, down_ch, up_ch chan int
			if col != (cols - 1) {
				right_ch = make(chan int)
			}
			if col != 0 {
				left_ch = runners[row][col-1].Node.Right
			}
			if row != (rows - 1) {
				down_ch = make(chan int)
			}
			if row != 0 {
				up_ch = runners[row-1][col].Node.Down
			}
			n := node.Node{Left: left_ch, Right: right_ch, Down: down_ch, Up: up_ch}
			clock_ch := make(chan bool)
			runners[row][col] = node.Runner{Clock: clock_ch, Node: n}
			program_path := fmt.Sprintf("programs/%d-%d.tis", row+1, col+1)
			err := runners[row][col].LoadFile(program_path)
			if err != nil {
				panic(err)
			}
			go runners[row][col].RunNode()
		}
	}
	return runners
}

func main() {
	colsPtr := flag.Int("cols", 2, "a number of columns")
	rowsPtr := flag.Int("rows", 2, "a number of rows")
	flag.Parse()

	nodeRunners := initNodeRunners(*rowsPtr, *colsPtr)

	start_ch := make(chan bool)
	stop_ch := make(chan bool)
	for row := range nodeRunners {
		for _, n := range nodeRunners[row] {
			go clockTicker(n.Clock, start_ch, stop_ch)
			start_ch <- true
		}
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	quit := make(chan error)
	redraw := make(chan bool)
	go eventLoop(quit, redraw)
	go redrawTicker(redraw)
	go redrawLoop(quit, redraw, nodeRunners)

	err = <-quit
	termbox.Close()
	if err != nil {
		panic(err)
	}
}
