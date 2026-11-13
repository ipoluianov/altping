package forms

import (
	"image/color"
	"math"
	"strings"

	"github.com/ipoluianov/nui/ui"
)

// Drawing shared by the table, the tiles, the lanes and the cards, in the
// look of the current style. The canvas blends straight (not premultiplied)
// colors, so a translucent color is the color with a lower alpha.

// alpha returns the color with the alpha a of 255
func alpha(col color.RGBA, a uint8) color.RGBA {
	col.A = a
	return col
}

// classic tells the two themes the application had before the styles:
// they keep their quiet look
func (s *Style) classic() bool {
	return s.ID == styleClassicDark || s.ID == styleClassicLight
}

// stateMarkSize is the size of the mark of a host state before its name
const stateMarkSize = 8

// drawStateMark draws the mark of a host state centered at (cx, cy): a
// circle, a square or a diamond, with a soft light around it in the styles
// that glow
func drawStateMark(cnv *ui.Canvas, s *Style, cx, cy float64, size int, col color.RGBA) {
	if s.Glow {
		for i, a := range []uint8{30, 60} {
			r := size/2 + 4 - i*2
			cnv.FillRoundedRectAA(int(cx)-r, int(cy)-r, r*2, r*2, r, alpha(col, a))
		}
	}
	half := float64(size) / 2
	switch s.DotKind {
	case dotSquare:
		cnv.FillRect(int(cx-half), int(cy-half), size, size, col)
	case dotDiamond:
		p := ui.NewPath()
		p.MoveTo(cx, cy-half-1)
		p.LineTo(cx+half+1, cy)
		p.LineTo(cx, cy+half+1)
		p.LineTo(cx-half-1, cy)
		p.Close()
		cnv.FillPath(p, col)
	default:
		cnv.DrawImage(int(cx-half), int(cy-half), statusDot(col))
	}
}

// pillHeight is the height of a status plate for the font size
func pillHeight(fontSize float64) int {
	return int(fontSize*1.5 + 0.5)
}

// drawPill writes the text of a state at x, centered on cy, in the look of
// the style: a plate of a tint of the color, or of the color itself with a
// dark text. It returns the width taken.
func drawPill(cnv *ui.Canvas, s *Style, x int, cy int, text string, col color.RGBA, fontFamily string, fontSize float64) int {
	text = strings.ToUpper(text)
	size := math.Max(fontSize-3, 9)
	cnv.SetFontFamily(fontFamily)
	cnv.SetFontSize(size)
	cnv.SetFontBold(true)
	defer cnv.SetFontBold(false)
	tw, _ := cnv.MeasureTextF(text)
	h := pillHeight(size)
	dot := 0
	if s.StatusKind == statusPill {
		dot = 6 + 6
	}
	pad := int(size * 0.75)
	w := int(math.Ceil(tw)) + pad*2 + dot
	y := cy - h/2
	radius := h / 2
	if s.TileCut || s.Mono {
		radius = 2
	}
	textCol := col
	if s.StatusKind == statusSolid {
		cnv.FillRoundedRectAA(x, y, w, h, 0, col)
		textCol = color.RGBA{0, 0, 0, 255}
	} else {
		cnv.FillRoundedRectAA(x, y, w, h, radius, alpha(col, 40))
	}
	tx := x + pad
	if dot > 0 {
		cnv.FillRoundedRectAA(tx, cy-3, 6, 6, 3, col)
		tx += dot
	}
	cnv.SetColor(textCol)
	cnv.SetHAlign(ui.HAlignLeft)
	cnv.SetVAlign(ui.VAlignCenter)
	cnv.DrawText(tx, y, w, h, text)
	return w
}

// drawTrend draws the replies of the points in the rectangle, scaled to
// the slowest one, in the kind of the style; lost pings are ticks along
// the bottom
func drawTrend(cnv *ui.Canvas, s *Style, x, y, w, h int, points []trendPoint, step int, col color.RGBA) {
	top, bottom := float64(y), float64(y+h)
	if len(points) == 0 || h <= 2 {
		return
	}
	maxValue := 0.0
	for _, p := range points {
		if p.hasValue {
			maxValue = math.Max(maxValue, p.value)
		}
	}
	// Some room above the top, and no huge swings for a host that replies in no time
	maxValue = math.Max(maxValue*1.15, 1)
	yOf := func(v float64) float64 { return bottom - v/maxValue*(bottom-top) }
	xOf := func(i int) float64 { return float64(x+i*step) + float64(step)/2 }

	if s.TrendKind == trendLine {
		cnv.FillRect(x, int(bottom), w, 1, alpha(s.Muted, 90))
	}
	lost := s.Down
	for i, p := range points {
		if p.lost {
			cnv.FillRect(x+i*step, int(bottom)-trendLostHeight+1, max(step-1, 2), trendLostHeight, lost)
		}
	}

	if s.TrendKind == trendBars {
		for i, p := range points {
			if !p.hasValue {
				continue
			}
			py := math.Min(yOf(p.value), bottom-2)
			cnv.FillRect(x+i*step, int(py), max(step-1, 1), int(bottom-py), col)
		}
		return
	}

	// The runs of points with a reply, drawn as lines
	var runs [][]ui.PointF
	var run []ui.PointF
	for i, p := range points {
		if !p.hasValue {
			if len(run) > 0 {
				runs = append(runs, run)
				run = nil
			}
			continue
		}
		run = append(run, ui.PointF{X: xOf(i), Y: yOf(p.value)})
	}
	if len(run) > 0 {
		runs = append(runs, run)
	}

	if s.TrendKind == trendArea || s.TrendKind == trendGlow {
		fill := ui.LinearGradient{X1: 0, Y1: top, X2: 0, Y2: bottom,
			Stops: []ui.GradientStop{{Offset: 0, Color: alpha(col, 110)}, {Offset: 1, Color: alpha(col, 0)}}}
		for _, r := range runs {
			if len(r) < 2 {
				continue
			}
			p := ui.NewPath()
			p.MoveTo(r[0].X, bottom)
			for _, pt := range r {
				p.LineTo(pt.X, pt.Y)
			}
			p.LineTo(r[len(r)-1].X, bottom)
			p.Close()
			cnv.FillPathPaint(p, fill)
		}
	}

	width := 1.0
	if s.TrendKind != trendLine {
		width = 1.4
	}
	if s.ID == "contrast" {
		width = 2
	}
	for _, r := range runs {
		if len(r) == 1 {
			// A reply between gaps or losses has no line to it: a short dash
			cnv.FillRect(int(r[0].X)-1, int(math.Round(r[0].Y)), step, 1, col)
			continue
		}
		p := ui.NewPath()
		p.MoveTo(r[0].X, r[0].Y)
		for _, pt := range r[1:] {
			p.LineTo(pt.X, pt.Y)
		}
		if s.TrendKind == trendGlow {
			cnv.StrokePath(p, ui.StrokeStyle{Width: width + 3, Join: ui.LineJoinRound, Cap: ui.LineCapRound}, alpha(col, 50))
		}
		if s.TrendKind == trendLine && s.ID != "contrast" {
			// The thin line of the classic look, as it always was
			for i := 1; i < len(r); i++ {
				cnv.DrawLineF(r[i-1].X, r[i-1].Y, r[i].X, r[i].Y, col)
			}
			continue
		}
		cnv.StrokePath(p, ui.StrokeStyle{Width: width, Join: ui.LineJoinRound, Cap: ui.LineCapRound}, col)
	}
}

// trendStep is how many pixels a point of the small charts takes
func trendStep(s *Style) int {
	if s.TrendKind == trendBars {
		return 3
	}
	return trendBarWidth
}

// drawSignal draws four bars like the signal of a phone: more of them lit
// for a faster host, none for a failed one
func drawSignal(cnv *ui.Canvas, s *Style, x, bottom int, lit int, col color.RGBA) {
	for i := 0; i < 4; i++ {
		h := 5 + i*3
		c := alpha(s.Muted, 70)
		if i < lit {
			c = col
		}
		cnv.FillRect(x+i*5, bottom-h, 3, h, c)
	}
}

// signalBars is how many bars of the signal a host gets for its time
func signalBars(r hostRow) int {
	if r.failed || !r.processed {
		return 0
	}
	ms := float64(r.time) / 1e6
	switch {
	case ms < 20:
		return 4
	case ms < 60:
		return 3
	case ms < 150:
		return 2
	}
	return 1
}

// fillCard fills the rectangle of a card: rounded, or with two corners cut
// off in the styles that cut them; a border when border has a color
func fillCard(cnv *ui.Canvas, s *Style, x, y, w, h int, fill, border color.RGBA, borderWidth float64) {
	if s.TileCut {
		const cut = 14.0
		p := ui.NewPath()
		fx, fy, fw, fh := float64(x), float64(y), float64(w), float64(h)
		p.MoveTo(fx, fy)
		p.LineTo(fx+fw-cut, fy)
		p.LineTo(fx+fw, fy+cut)
		p.LineTo(fx+fw, fy+fh)
		p.LineTo(fx+cut, fy+fh)
		p.LineTo(fx, fy+fh-cut)
		p.Close()
		cnv.FillPath(p, fill)
		if border.A > 0 {
			cnv.StrokePath(p, ui.StrokeStyle{Width: borderWidth, Join: ui.LineJoinMiter}, border)
		}
		return
	}
	r := s.TileRadius
	if border.A > 0 {
		bw := int(math.Ceil(borderWidth))
		cnv.FillRoundedRectAA(x, y, w, h, r, border)
		cnv.FillRoundedRectAA(x+bw, y+bw, w-bw*2, h-bw*2, max(r-bw, 0), fill)
		return
	}
	cnv.FillRoundedRectAA(x, y, w, h, r, fill)
}

// stateColor is the color of the state of the host
func (s *Style) stateColor(r hostRow) color.RGBA {
	switch {
	case r.failed:
		return s.Down
	case !r.processed:
		return s.Idle
	case r.isSlow():
		return s.Slow
	}
	return s.OK
}

// stateText is the short state of the host for a plate: OK, SLOW, the error
func stateText(r hostRow) string {
	switch {
	case !r.processed:
		return "-"
	case r.failed:
		return r.details
	case r.isSlow():
		return T().StatusSlow
	}
	return T().StateOK
}

// heatColor colors a ping time by how fast it is, in the styles that do
func heatColor(s *Style, r hostRow) (color.RGBA, bool) {
	if !s.HeatTime || !r.processed || r.failed {
		return color.RGBA{}, false
	}
	ms := float64(r.time) / 1e6
	switch {
	case ms < 5:
		return s.OK, true
	case ms >= 150:
		return s.Slow, true
	}
	return color.RGBA{}, false
}
