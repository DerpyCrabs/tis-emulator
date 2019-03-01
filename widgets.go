package main

import (
	"errors"
	"fmt"
	"tis/node"

	termbox "github.com/nsf/termbox-go"
)

func printString(s string, x, y int, fg, bg termbox.Attribute) error {
	width, height := termbox.Size()
	if y >= height || x+len(s) > width {
		return errors.New("Printing outside of buffer")
	}
	for i, c := range []rune(s) {
		termbox.SetCell(x+i, y, c, fg, bg)
	}
	return nil
}

func drawBorder(s string, x, y, w, h int, fg, bg termbox.Attribute) error {
	for i := x; i < (x + w); i++ {
		for j := y; j < (y + h); j++ {
			if i == x || i == (x+w-1) || j == y || j == (y+h-1) {
				termbox.SetCell(i, j, ' ', fg, bg)
			}
		}
	}
	return printString(s, x+(w/2-len(s)/2), y, bg, fg)
}

func drawNode(n node.Runner, x, y, w, h int) {
	_ = drawBorder(n.Name, x, y, w, h, termbox.ColorBlack, termbox.ColorWhite)
	_ = printString(fmt.Sprintf("Acc: %d, Bak: %d", n.Node.Acc, n.Node.Bak), x+1, y+1, termbox.ColorWhite, termbox.ColorBlue)
	for ln, line := range n.Commands {
		cur := ' '
		if ln == n.Line {
			cur = '>'
		}
		_ = printString(fmt.Sprintf("%c %s", cur, line), x+1, y+ln+2, termbox.ColorWhite, termbox.ColorBlack)
	}
}
