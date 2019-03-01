package node

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Runner struct {
	Name     string
	Line     int
	Commands []string
	Node     Node
	Clock    chan bool
	labels   map[string]int
}

func (nr *Runner) RunNode() {
	ta := regexp.MustCompile("^(?:\\s?)(\\w+) (\\w+|-+\\d+) (\\w+|-+\\d+)((\\s*)#(.*)$|(\\s*)$)")
	oa := regexp.MustCompile("^(?:\\s?)(\\w+) (\\w+|-+\\d+)((\\s*)#(.*)$|(\\s*)$)")
	na := regexp.MustCompile("^(?:\\s?)(\\w+)((\\s*)#(.*)$|(\\s*)$)")
	for {
		<-nr.Clock
		//fmt.Printf("%s: Acc - %d, Bak - %d\n", nr.Name, nr.Node.Acc, nr.Node.Bak)
		if nr.Line >= len(nr.Commands) {
			nr.Line = 0
		}
		match := ta.FindStringSubmatch(nr.Commands[nr.Line])
		if match != nil {
			nr.Node.Call(match[1], match[2], match[3])
		} else {
			match = oa.FindStringSubmatch(nr.Commands[nr.Line])
			if match != nil {
				if strings.ToLower(match[1]) == "jmp" {
					nr.Line = nr.labels[strings.ToLower(match[2])]
					continue
				} else if strings.ToLower(match[1]) == "jez" {
					if nr.Node.Acc == 0 {
						nr.Line = nr.labels[strings.ToLower(match[2])]
					}
					continue
				} else if strings.ToLower(match[1]) == "jnz" {
					if nr.Node.Acc != 0 {
						nr.Line = nr.labels[strings.ToLower(match[2])]
					}
					continue
				} else if strings.ToLower(match[1]) == "jgz" {
					if nr.Node.Acc > 0 {
						nr.Line = nr.labels[strings.ToLower(match[2])]
					}
					continue
				} else if strings.ToLower(match[1]) == "jlz" {
					if nr.Node.Acc < 0 {
						nr.Line = nr.labels[strings.ToLower(match[2])]
					}
					continue
				} else if strings.ToLower(match[1]) == "jro" {
					num, err := strconv.Atoi(match[2])
					if err != nil {
						nr.Line += 0
					} else {
						nr.Line += num
					}
					continue
				}
				nr.Node.Call(match[1], match[2])
			} else {
				match = na.FindStringSubmatch(nr.Commands[nr.Line])
				if match != nil {
					nr.Node.Call(match[1])
				}
			}
		}
		nr.Line += 1
	}
}

func (nr *Runner) readLabels() {
	r := regexp.MustCompile("^(\\w+):(.*)$")
	nr.labels = make(map[string]int)
	for i, line := range nr.Commands {
		match := r.FindStringSubmatch(line)
		if match != nil {
			nr.Commands[i] = match[2]
			nr.labels[strings.ToLower(match[1])] = i
		}
	}
}

func (nr *Runner) stripComments() {
	name := regexp.MustCompile("^##(?:\\s*)(.*)(\\s*)$")
	match := name.FindStringSubmatch(nr.Commands[0])
	if match != nil {
		nr.Name = match[1]
	}
	comment, _ := regexp.Compile("^#(.*)$")
	temp := []string{}
	for _, line := range nr.Commands {
		match := comment.FindStringSubmatch(line)
		if match == nil && line != "" {
			temp = append(temp, line)
		}
	}
	nr.Commands = temp
}

func (nr *Runner) LoadFile(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	nr.Commands = strings.Split(string(dat), "\n")
	nr.stripComments()
	nr.readLabels()
	nr.Line = 0
	return nil
}
