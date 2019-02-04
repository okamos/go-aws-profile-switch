package ui

import termbox "github.com/nsf/termbox-go"

// InputBox captures user input
type InputBox struct {
	runes  []rune
	cursor int
	initX  int
	x, y   int // InputBox position in the screen
}

func (ib *InputBox) addRune(r rune) {
	if len(ib.runes) == ib.cursor {
		ib.runes = append(ib.runes, r)
	} else {
		ib.runes = append(ib.runes[:ib.cursor], append([]rune{r}, ib.runes[ib.cursor:]...)...)
	}
	ib.cursor++
}

func (ib *InputBox) removeRune() {
	l := len(ib.runes)
	if l < 1 {
		return
	}
	if l == ib.cursor {
		ib.runes = ib.runes[:l-1]
	} else {
		ib.runes = append(ib.runes[:ib.cursor-1], ib.runes[ib.cursor:]...)
	}
	ib.cursor--
}

func (ib *InputBox) clearText() {
	ib.runes = []rune{}
	ib.cursor = 0
}

func (ib *InputBox) moveToLeft() {
	if ib.cursor > 0 {
		ib.cursor--
	}
}

func (ib *InputBox) moveToRight() {
	if ib.cursor < len(ib.runes) {
		ib.cursor++
	}
}

// Draw user input
func (ib *InputBox) Draw(x, y, w, h int) {
	for _, r := range ib.runes {
		termbox.SetCell(x, y, r, coldef, coldef)
		x++
	}
	if len(ib.runes) <= ib.cursor {
		termbox.SetCell(ib.initX+ib.cursor, y, ' ', coldef, termbox.Attribute(3))
	} else {
		termbox.SetCell(ib.initX+ib.cursor, y, ib.runes[ib.cursor], coldef, termbox.Attribute(3))
	}
}
