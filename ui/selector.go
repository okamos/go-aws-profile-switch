package ui

import (
	"math"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

// Selector is visual interface to user
type Selector struct {
	index      int
	candidates []string
	values     []string
	visualY    int
}

func (s *Selector) moveToUp() {
	if s.index <= 0 {
		s.index = len(s.candidates) - 1
	} else {
		s.index--
	}
}

func (s *Selector) moveToDown() {
	if len(s.candidates) <= s.index+1 {
		s.index = 0
	} else {
		s.index++
	}
}

func (s *Selector) filter(str string) {
	s.index = 0
	if str == "" {
		s.candidates = s.values
		return
	}
	strs := make([]string, 0)
	for _, v := range s.values {
		if strings.Contains(strings.ToLower(v), strings.ToLower(str)) {
			strs = append(strs, v)
		}
	}
	s.candidates = strs
}

// Draw selector messages
func (s *Selector) Draw(x, y, w, h int) {
	var (
		first, last int
	)
	l := len(s.candidates)
	if l == 0 {
		return
	}
	last = l
	if l > h {
		page := math.Ceil(float64(s.index+1)/float64(h)) - 1
		first = int(page) * h
		last = int(page+1)*h - 1
		if last > l {
			last = l
		}
	}
	for i, m := range s.candidates[first:last] {
		tbPrint(x, y+i, coldef, coldef, m)
	}
	s.visualY = y + s.index - first
	tbFill(s.visualY, w, termbox.Attribute(5))
	tbPrint(x, s.visualY, coldef, termbox.Attribute(5), s.candidates[s.index])
}

func (s *Selector) selected() string {
	return s.candidates[s.index]
}
