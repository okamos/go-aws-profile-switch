package ui

import (
	termbox "github.com/nsf/termbox-go"
)

const (
	coldef  = termbox.ColorDefault
	help    = "[ESC] QUIT | Filter > "
	helpLen = len(help)
)

var (
	ib       InputBox
	selector Selector
)

func redrawAll() {
	w, h := termbox.Size()
	termbox.Clear(coldef, coldef)

	tbPrint(0, 0, coldef, coldef, help)
	ib.Draw(helpLen, 0, w, 0)
	selector.Draw(0, 1, w, h-1)

	termbox.Flush()
}

func tbFill(y, w int, bg termbox.Attribute) {
	for x := 0; x < w; x++ {
		termbox.SetCell(x, y, ' ', coldef, bg)
	}
}

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// Select profile
func Select() (string, error) {
	err := termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	redrawAll()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return "", nil
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				ib.moveToLeft()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				ib.moveToRight()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				ib.removeRune()
				selector.filter(string(ib.runes))
			case termbox.KeyDelete, termbox.KeyCtrlU:
				ib.clearText()
				selector.filter("")
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				selector.moveToUp()
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				selector.moveToDown()
			case termbox.KeyHome, termbox.KeyCtrlA:
				ib.cursor = 0
			case termbox.KeyEnd, termbox.KeyCtrlE:
				ib.cursor = len(ib.runes)
			case termbox.KeyEnter:
				return selector.selected(), nil
			default:
				if ev.Ch != 0 {
					ib.addRune(ev.Ch)
					selector.filter(string(ib.runes))
				}
			}
		case termbox.EventError:
			return "", ev.Err
		}
		redrawAll()
	}
}

// InitSelector set values to selector
func InitSelector(values []string) {
	selector.values = values
	selector.filter("")
}

func init() {
	ib = InputBox{
		initX: helpLen,
	}
	selector = Selector{
		index: 0,
	}
}
