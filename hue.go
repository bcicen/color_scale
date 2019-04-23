package main

import (
	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

const (
	hueMax   = 359.0
	hueIncr  = 0.5
	hueCount = int(hueMax / hueIncr)
)

type HueBar struct {
	items  [hueCount]tcell.Color // navigation colors
	mItems []int                 // minimap sample indices
	pos    int
	width  int
	pst    tcell.Style // pointer style
	state  *State
}

func NewHueBar(s *State) *HueBar { return &HueBar{state: s} }

func (bar *HueBar) SetPos(n float64) { bar.pos = int(n / hueIncr) }
func (bar *HueBar) center() int      { return (bar.width / 2) }

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *HueBar) Draw(x, y int, s tcell.Screen) int {
	center := bar.width / 2
	boxPad := bar.width / 30
	if boxPad < 2 {
		boxPad = 2
	}
	st := tcell.StyleDefault.
		Foreground(tcell.ColorBlack)

	s.SetCell(center+x, y, bar.pst, '▾')

	// border bars
	s.SetCell(x-1, y+1, bar.pst, '│')
	s.SetCell(x-1, y+2, bar.pst, '│')
	s.SetCell(x-1, y+3, bar.pst, '│')
	s.SetCell(bar.width+x+1, y+1, bar.pst, '│')
	s.SetCell(bar.width+x+1, y+2, bar.pst, '│')
	s.SetCell(bar.width+x+1, y+3, bar.pst, '│')

	for col, color := range bar.Items() {
		s.SetCell(col+x, y+1, st.Background(color), ' ')
		s.SetCell(col+x, y+2, st.Background(color), '▁')
		//if col == center-boxPad {
		//s.SetCell(center+x-boxPad, y+1, bar.pst.Background(color), '┌')
		//}
		//if col == center+boxPad {
		//s.SetCell(center+x+boxPad, y+1, bar.pst.Background(color), '┐')
		//}
	}

	for col, color := range bar.MiniMap() {
		st = st.Background(color)
		s.SetCell(col+x, y+3, st, ' ')
	}

	s.SetCell(center+x-boxPad, y+4, bar.pst, '└')
	s.SetCell(center+x+boxPad, y+4, bar.pst, '┘')

	return 5
}

func (bar *HueBar) miniStep() int {
	n := len(bar.items) / bar.width
	if n > 13 {
		n = 13
	}
	return n
}

func (bar *HueBar) Resize(w int) {
	bar.width = w
	bar.buildMini()
}

func (bar *HueBar) Handle(change StateChange) {
	var n int
	var nc *noire.Color

	if change.Includes(SaturationChanged, ValueChanged) {
		for i := 0.0; i < hueMax; i += hueIncr {
			nc = noire.NewHSV(i, bar.state.Saturation(), bar.state.Value())
			bar.items[n] = toTColor(nc)
			n++
		}
	}

	if change.Includes(SelectedChanged, HueChanged) {
		bar.SetPos(bar.state.Hue())
	}
}

// build minimap indices
func (bar *HueBar) buildMini() {
	bar.mItems = bar.mItems[0:]
	miniStep := bar.miniStep()
	for n := 0; n < len(bar.items); n += miniStep {
		bar.mItems = append(bar.mItems, n)
	}
}

func (bar *HueBar) Items() []tcell.Color {
	l := bar.pos - bar.center()
	r := bar.pos + bar.center() + 1
	ilen := len(bar.items)
	if l < 0 {
		return append(bar.items[ilen+l:], bar.items[0:r]...)
	}
	if r > ilen {
		return append(bar.items[l:ilen-1], bar.items[0:r-ilen]...)
	}
	return bar.items[l:r]
}

func (bar *HueBar) MiniMap() []tcell.Color {
	var mpos int
	for mpos < len(bar.mItems)-1 {
		if bar.mItems[mpos+1] >= bar.pos {
			break
		}
		mpos++
	}

	l := mpos - bar.center()
	r := mpos + bar.center() + 1
	ilen := len(bar.mItems)

	var a []int
	switch {
	case l < 0:
		a = append(bar.mItems[ilen+l:], bar.mItems[0:r]...)
	case r > ilen:
		a = append(bar.mItems[l:], bar.mItems[0:r-ilen]...)
	default:
		a = bar.mItems[l:r]
	}

	var colors []tcell.Color
	for _, idx := range a {
		colors = append(colors, bar.items[idx])
	}

	return colors
}

func (bar *HueBar) Width() int { return bar.width }

func (bar *HueBar) Up(step int) {
	n := int(bar.state.Hue()) + step
	if n > hueMax {
		n -= hueMax
	}
	bar.state.SetHue(float64(n))
}

func (bar *HueBar) Down(step int) {
	n := int(bar.state.Hue()) - step
	if n < 0 {
		n += hueMax
	}
	bar.state.SetHue(float64(n))
}

func (bar *HueBar) SetPointerStyle(st tcell.Style) { bar.pst = st }
