package forms

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/nui/ui"
)

// stylePicker is the popup of the style button in the toolbar: a card per
// style with a small picture of how it looks and its name. A click applies
// the style; the arrows and Enter do the same from the keyboard.
type stylePicker struct {
	ui.Widget
	hover  int
	active int // chosen with the keyboard, -1 - none
}

const (
	pickerColumns = 4
	pickerThumbW  = 150
	pickerThumbH  = 86
	pickerGap     = 10
	pickerPad     = 12
)

// showStylePicker opens the picker under the button, its right edge at the
// button's
func showStylePicker(button ui.Widgeter) {
	form := lastCreatedMainWidget.Form()
	if form == nil {
		return
	}
	p := &stylePicker{hover: -1, active: -1}
	p.InitWidget()
	p.SetAbsolutePositioning(true)
	p.SetRole("popup")
	p.SetAutoFillBackground(true)
	p.SetCanBeFocused(true)
	p.SetOnPaint(p.paint)
	p.SetOnMouseMove(func(x, y int, mods ui.KeyModifiers) bool {
		if i := p.cardAt(x, y); i != p.hover {
			p.hover = i
			p.Update()
		}
		return true
	})
	p.SetOnMouseLeave(func() {
		p.hover = -1
		p.Update()
	})
	p.SetOnMouseDown(func(b ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
		if i := p.cardAt(x, y); b == ui.MouseButtonLeft && i >= 0 {
			p.pick(i)
		}
		return true
	})
	p.OnKeyDown.Connect(p.keyDown)

	w, h := p.size()
	p.SetSize(w, h)
	bx, by := button.RectClientAreaOnWindow()
	x := max(bx+button.Width()-w, 0)
	p.SetPosition(x, by+button.Height())
	form.OpenPopup(p)
	p.Focus()
}

func (p *stylePicker) titleHeight() int {
	return ui.ThemeLineHeight() + pickerPad
}

func (p *stylePicker) nameHeight() int {
	return ui.ThemeLineHeight() + 6
}

func (p *stylePicker) size() (int, int) {
	rows := (len(styles) + pickerColumns - 1) / pickerColumns
	w := pickerPad*2 + pickerColumns*pickerThumbW + (pickerColumns-1)*pickerGap
	h := p.titleHeight() + rows*(pickerThumbH+p.nameHeight()) + (rows-1)*pickerGap + pickerPad
	return w, h
}

// cardRect is the rectangle of the card of the style number i: the picture and the name
func (p *stylePicker) cardRect(i int) image.Rectangle {
	col, row := i%pickerColumns, i/pickerColumns
	x := pickerPad + col*(pickerThumbW+pickerGap)
	y := p.titleHeight() + row*(pickerThumbH+p.nameHeight()+pickerGap)
	return image.Rect(x, y, x+pickerThumbW, y+pickerThumbH+p.nameHeight())
}

func (p *stylePicker) cardAt(x, y int) int {
	for i := range styles {
		if image.Pt(x, y).In(p.cardRect(i)) {
			return i
		}
	}
	return -1
}

func (p *stylePicker) pick(i int) {
	form := p.Form()
	s := config.GetSettings()
	s.Theme = styles[i].ID
	if form != nil {
		form.Panel().CloseAllPopup()
	}
	lastCreatedMainWidget.ApplySettings(s)
	lastCreatedLeftWidget.FocusTable()
}

func (p *stylePicker) keyDown(e *ui.KeyEvent) {
	i := p.active
	if i < 0 {
		i = currentStyleIndex()
	}
	switch e.Key {
	case ui.KeyArrowLeft:
		i--
	case ui.KeyArrowRight:
		i++
	case ui.KeyArrowUp:
		i -= pickerColumns
	case ui.KeyArrowDown:
		i += pickerColumns
	case ui.KeyEnter, ui.KeySpace:
		p.pick(i)
		e.Handled = true
		return
	default:
		return
	}
	p.active = min(max(i, 0), len(styles)-1)
	p.Update()
	e.Handled = true
}

func currentStyleIndex() int {
	for i, s := range styles {
		if s == currentStyle {
			return i
		}
	}
	return 0
}

func (p *stylePicker) paint(cnv *ui.Canvas) {
	pal := ui.CurrentPalette()
	cnv.FillRect(0, 0, p.Width(), p.Height(), pal.PopupBase)
	cnv.SetColor(pal.Border)
	cnv.DrawRect(0, 0, p.Width(), p.Height())

	family := ui.ThemeFontFamily()
	fs := ui.ThemeFontSize()
	cnv.SetFontFamily(family)
	cnv.SetFontSize(fs)
	cnv.SetFontBold(true)
	cnv.SetColor(pal.Text)
	cnv.SetHAlign(ui.HAlignLeft)
	cnv.SetVAlign(ui.VAlignCenter)
	cnv.DrawText(pickerPad, 2, p.Width()-pickerPad*2, p.titleHeight()-2, withoutHotkey(T().ToolStyle))
	cnv.SetFontBold(false)

	current := currentStyleIndex()
	for i, s := range styles {
		r := p.cardRect(i)
		thumb := image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+pickerThumbH)
		switch {
		case i == current:
			cnv.FillRoundedRectAA(thumb.Min.X-3, thumb.Min.Y-3, thumb.Dx()+6, thumb.Dy()+6, 9, pal.Highlight)
		case i == p.hover || i == p.active:
			cnv.FillRoundedRectAA(thumb.Min.X-3, thumb.Min.Y-3, thumb.Dx()+6, thumb.Dy()+6, 9, ui.MixColors(pal.PopupBase, pal.Highlight, 0.5))
		}
		drawStylePreview(cnv, s, thumb)
		cnv.SetFontFamily(family)
		cnv.SetFontSize(fs - 1)
		cnv.SetHAlign(ui.HAlignCenter)
		cnv.SetVAlign(ui.VAlignCenter)
		cnv.SetColor(pal.Text)
		if i == current {
			cnv.SetFontBold(true)
		}
		cnv.DrawText(r.Min.X, thumb.Max.Y+2, r.Dx(), p.nameHeight()-2, fitText(family, fs-1, s.DisplayName(), r.Dx()))
		cnv.SetFontBold(false)
	}
}

// drawStylePreview draws a small picture of the main window in the style:
// its colors, the view of the hosts, the small charts and the plates
func drawStylePreview(cnv *ui.Canvas, s *Style, r image.Rectangle) {
	p := s.Palette
	x, y, w, h := r.Min.X, r.Min.Y, r.Dx(), r.Dy()
	cnv.FillRoundedRectAA(x, y, w, h, 6, p.Border)
	cnv.FillRoundedRectAA(x+1, y+1, w-2, h-2, 5, p.Window)
	cnv.Save()
	cnv.TranslateAndClip(x+1, y+1, w-2, h-2)
	defer cnv.Restore()
	w, h = w-2, h-2

	// The toolbar: a few icons and the accent of the chart button
	for i := 0; i < 3; i++ {
		cnv.FillRect(6+i*9, 6, 5, 5, alpha(p.Text, 120))
	}
	cnv.FillRoundedRectAA(w-50, 4, 9, 9, 2, p.Highlight)

	top := 17
	if s.KPIs {
		cw := (w - 8 - 3*4) / 4
		for i, col := range []color.RGBA{p.Text, s.OK, s.Slow, s.Down} {
			cx := 4 + i*(cw+4)
			cnv.FillRoundedRectAA(cx, top, cw, 12, 2, p.Base)
			cnv.FillRect(cx, top+1, 2, 10, col)
			cnv.FillRect(cx+5, top+4, cw/3, 4, col)
		}
		top += 15
	}
	listW := w * 64 / 100
	rnd := rand.New(rand.NewSource(int64(len(s.Name) + int(s.Trend.R))))
	states := []color.RGBA{s.OK, s.OK, s.Down, s.OK, s.Slow, s.OK}

	switch s.View {
	case viewTiles:
		cols, rows := 3, 2
		tw := (listW - 4 - (cols-1)*3) / cols
		th := (h - top - 4 - (rows-1)*3) / rows
		for i := 0; i < cols*rows; i++ {
			tx := 4 + (i%cols)*(tw+3)
			ty := top + (i/cols)*(th+3)
			state := states[i%len(states)]
			fill := p.Base
			if s.TileSolid && state == s.Down {
				fill = s.Down
			}
			mini := *s
			mini.TileRadius = min(s.TileRadius/3, 4)
			fillCard(cnv, &mini, tx, ty, tw, th, fill, alpha(p.Border, 255), 1)
			cnv.FillRect(tx+4, ty+4, tw/2, 2, alpha(p.Text, 150))
			cnv.FillRect(tx+4, ty+9, tw/3, 5, state)
			col := s.Trend
			if s.TrendByState || state != s.OK {
				col = state
			}
			miniTrend(cnv, s, tx+3, ty+th-10, tw-6, 7, col, rnd)
		}
	case viewLanes:
		n := 6
		lh := (h - top - 2) / n
		for i := 0; i < n; i++ {
			ly := top + i*lh
			state := states[i%len(states)]
			cnv.FillRect(4, ly+lh-1, listW-4, 1, p.Divider)
			cnv.FillRect(6, ly+lh/2-2, 14, 3, alpha(p.Text, 150))
			base := float64(ly) + float64(lh)*0.65
			path := ui.NewPath()
			path.MoveTo(24, base)
			for bx := 28.0; bx < float64(listW-18); bx += 9 {
				a := float64(lh) * (0.35 + 0.4*rnd.Float64())
				if state == s.Down {
					a = 0
				}
				path.LineTo(bx, base)
				path.LineTo(bx+1.5, base-a)
				path.LineTo(bx+3, base+2)
				path.LineTo(bx+4, base)
			}
			path.LineTo(float64(listW-16), base)
			cnv.StrokePath(path, ui.StrokeStyle{Width: 1, Join: ui.LineJoinRound}, state)
			cnv.FillRect(listW-12, ly+lh/2-3, 8, 5, state)
		}
	default:
		cnv.FillRect(4, top, listW-4, h-top-4, p.Base)
		cnv.FillRect(4, top, listW-4, 6, p.Button)
		n := 6
		rh := (h - top - 10) / n
		for i := 0; i < n; i++ {
			ry := top + 7 + i*rh
			state := states[i%len(states)]
			if i == 1 {
				cnv.FillRect(4, ry, listW-4, rh, p.Selection)
			}
			mark := *s
			mark.Glow = false
			drawStateMark(cnv, &mark, 8, float64(ry)+float64(rh)/2, 3, state)
			cnv.FillRect(12, ry+rh/2-1, 16, 2, alpha(p.Text, 150))
			col := s.Trend
			if !s.classic() && state == s.Slow {
				col = state
			}
			miniTrend(cnv, s, 32, ry+1, 26, rh-2, col, rnd)
			switch s.StatusKind {
			case statusText:
				cnv.FillRect(listW-18, ry+rh/2-1, 10, 2, state)
			case statusSolid:
				cnv.FillRect(listW-20, ry+1, 14, rh-2, state)
			default:
				cnv.FillRoundedRectAA(listW-20, ry+1, 14, rh-2, (rh-2)/2, alpha(state, 70))
			}
		}
	}

	// The history chart on the right
	cx := listW + 4
	cw := w - cx - 4
	cnv.FillRect(cx, top, cw, h-top-4, p.Base)
	points := make([]trendPoint, cw/2)
	v := 0.5
	for i := range points {
		v = math.Max(0.15, math.Min(1, v+rnd.NormFloat64()*0.12))
		points[i] = trendPoint{hasValue: true, value: v}
		if rnd.Float64() < 0.04 {
			points[i].value = v * 1.9
		}
	}
	chart := *s
	if s.ChartFill > 0 && chart.TrendKind == trendLine {
		chart.TrendKind = trendArea
	}
	drawTrend(cnv, &chart, cx+2, top+6, cw-4, h-top-12, points, 2, s.Trend)
}

// miniTrend draws a small made-up chart in the trend kind of the style
func miniTrend(cnv *ui.Canvas, s *Style, x, y, w, h int, col color.RGBA, rnd *rand.Rand) {
	step := 2
	if s.TrendKind == trendBars {
		step = 3
	}
	points := make([]trendPoint, max(w/step, 1))
	for i := range points {
		points[i] = trendPoint{hasValue: true, value: 0.3 + rnd.Float64()*0.4}
		if rnd.Float64() < 0.1 {
			points[i].value = 1
		}
	}
	mini := *s
	mini.Glow = false
	drawTrend(cnv, &mini, x, y, w, h, points, step, col)
}
