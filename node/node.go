package node

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Node struct {
	Acc   int
	Bak   int
	Left  chan int
	Right chan int
	Up    chan int
	Down  chan int
	last  chan int
}

func (self *Node) Call(name string, args ...string) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(strings.ToLower(args[i]))
	}
	reflect.ValueOf(self).MethodByName(strings.Title(name)).Call(inputs)
}

func (self *Node) Nop() error {
	return self.Add("nil")
}

func (self *Node) Sav() error {
	self.Bak = self.Acc
	return nil
}

func (self *Node) Swp() error {
	self.Acc, self.Bak = self.Bak, self.Acc
	return nil
}

func (self *Node) Neg() error {
	self.Acc = -self.Acc
	return nil
}

func (self *Node) recieve(src string) (int, error) {
	if src == "acc" {
		return self.Acc, nil
	}
	if src == "last" {
		if self.last != nil {
			return <-self.last, nil
		} else {
			return 0, nil
		}
	}
	if src == "any" {
		select {
		case input := <-self.Left:
			self.last = self.Left
			return input, nil
		case input := <-self.Right:
			self.last = self.Right
			return input, nil
		case input := <-self.Up:
			self.last = self.Up
			return input, nil
		case input := <-self.Down:
			self.last = self.Down
			return input, nil
		}
	}
	if src == "left" {
		return <-self.Left, nil
	}
	if src == "right" {
		return <-self.Right, nil
	}
	if src == "up" {
		return <-self.Up, nil
	}
	if src == "down" {
		return <-self.Down, nil
	}
	if src == "nil" {
		return 0, nil
	}
	return 0, errors.New("Not a src descriptor")
}

func (self *Node) transmit(dst string, num int) error {
	if dst == "acc" {
		self.Acc = num
		return nil
	}
	if dst == "any" {
		select {
		case self.Left <- num:
			return nil
		case self.Right <- num:
			return nil
		case self.Up <- num:
			return nil
		case self.Down <- num:
			return nil
		}
	}
	if dst == "last" {
		if self.last != nil {
			self.last <- num
			return nil
		} else {
			return errors.New("No last port")
		}
	}
	if dst == "left" {
		self.Left <- num
		return nil
	}
	if dst == "right" {
		self.Right <- num
		return nil
	}
	if dst == "up" {
		self.Up <- num
		return nil
	}
	if dst == "down" {
		self.Down <- num
		return nil
	}
	if dst == "nil" {
		return nil
	}
	return errors.New("Incorrect destination")
}

func (self *Node) getInput(op string) (int, error) {
	num, err := self.recieve(op)
	if err != nil {
		num, err = strconv.Atoi(op)
		if err != nil {
			return 0, err
		}
	}
	return num, nil
}

func (self *Node) Add(op string) error {
	num, err := self.getInput(op)
	if err != nil {
		return fmt.Errorf("add: getInput failed: %v", err)
	}
	self.Acc += num
	return nil
}

func (self *Node) Sub(op string) error {
	num, err := self.getInput(op)
	if err != nil {
		return fmt.Errorf("sub: getInput failed: %v", err)
	}
	self.Acc -= num
	return nil
}

func (self *Node) Mov(src string, dst string) error {
	input, err := self.getInput(src)
	if err != nil {
		return fmt.Errorf("mov: getInput failed: %v", err)
	}
	err = self.transmit(dst, input)
	if err != nil {
		return fmt.Errorf("mov: transmit failed: %v", err)
	}
	return nil
}
