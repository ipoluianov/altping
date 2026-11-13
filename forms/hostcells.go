package forms

import (
	"image"
	"image/color"
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

// The host table draws two cells itself: the name, after a dot in the color of
// the host state, and the trend, a small chart of the last minutes. The other
// cells are plain text; only the values that mean something are colored.

const (
	cellPaddingX  = 8
	statusDotSize = 8
	statusDotGap  = 8
	notifyIconGap = 4
	trendPaddingY = 5
	// A lost ping is a tick of this height at the bottom of the trend chart
	trendLostHeight = 5
	trendBarWidth   = 2 // pixels per point of the trend chart
	trendCacheLife  = time.Second
)

// textColor is the color of the ordinary text of the table
func textColor() color.RGBA {
	return ui.CurrentPalette().Text
}

// fieldColor colors the value of a field: the state colors only where they tell
// something, e.g. a slow time or lost pings; the rest is the ordinary text
func (r hostRow) fieldColor(field int, text string) color.Color {
	if !r.processed || text == "-" {
		return colorMuted.get()
	}
	switch field {
	case colDetails:
		if r.failed {
			return colorFailed.get()
		}
		return colorMuted.get()
	case colTime:
		if r.slowMs > 0 && r.time > time.Duration(r.slowMs)*time.Millisecond {
			return colorSlow.get()
		}
		if col, ok := heatColor(currentStyle, r); ok {
			return col
		}
	case colAvg:
		if r.isSlow() {
			return colorSlow.get()
		}
	case colLoss:
		switch {
		case r.stats.Lost == 0:
		case r.stats.Lost == r.stats.Sent:
			return colorFailed.get()
		default:
			return colorSlow.get()
		}
	}
	return textColor()
}

// statusDots are the dots before the names, by color
var statusDots = map[color.RGBA]image.Image{}

// statusDot returns a smooth round dot of the color, statusDotSize across
func statusDot(col color.RGBA) image.Image {
	if img, ok := statusDots[col]; ok {
		return img
	}
	const samples = 4 // per pixel side, for the smooth edge
	img := image.NewNRGBA(image.Rect(0, 0, statusDotSize, statusDotSize))
	center := float64(statusDotSize) / 2
	for y := 0; y < statusDotSize; y++ {
		for x := 0; x < statusDotSize; x++ {
			inside := 0
			for sy := 0; sy < samples; sy++ {
				for sx := 0; sx < samples; sx++ {
					dx := float64(x) + (float64(sx)+0.5)/samples - center
					dy := float64(y) + (float64(sy)+0.5)/samples - center
					if dx*dx+dy*dy <= center*center {
						inside++
					}
				}
			}
			a := uint8(int(col.A) * inside / (samples * samples))
			img.SetNRGBA(x, y, color.NRGBA{R: col.R, G: col.G, B: col.B, A: a})
		}
	}
	statusDots[col] = img
	return img
}

// drawNameCell draws the mark of the host state, the bell of a host that
// beeps, and the name; a selected row gets a bar of the accent on its left
func (c *LeftWidget) drawNameCell(cnv *ui.Canvas, width, height int, r hostRow, notify, selected bool) {
	st := currentStyle
	if selected && !st.classic() {
		cnv.FillRect(0, 0, 3, height, st.Palette.Highlight)
	}
	x := cellPaddingX
	drawStateMark(cnv, st, float64(x)+statusDotSize/2, float64(height)/2, statusDotSize, colorToRGBA(r.color()))
	x += statusDotSize + statusDotGap
	if notify && c.iconNotify != nil {
		cnv.DrawImage(x, (height-notifyIconSize)/2, c.iconNotify)
	}
	x += notifyIconSize + notifyIconGap
	if c.shareSlot {
		if icon := c.shareIcon(r.share); icon != nil {
			cnv.DrawImage(x, (height-notifyIconSize)/2, icon)
		}
		x += notifyIconSize + notifyIconGap
	}

	cnv.SetFontFamily(c.lvItems.FontFamily())
	cnv.SetFontSize(c.lvItems.FontSize())
	cnv.SetHAlign(ui.HAlignLeft)
	cnv.SetVAlign(ui.VAlignCenter)
	if r.processed {
		cnv.SetColor(textColor())
	} else {
		cnv.SetColor(colorMuted.get())
	}
	cnv.DrawText(x, 0, width-x-cellPaddingX, height, fitText(c.lvItems.FontFamily(), c.lvItems.FontSize(), r.name, width-x-cellPaddingX))
}

// drawStatusCell writes the state of the host on a plate, in the styles
// that have them
func (c *LeftWidget) drawStatusCell(cnv *ui.Canvas, height int, r hostRow) {
	st := currentStyle
	col := st.stateColor(r)
	if !r.processed {
		col = st.Muted
	}
	drawPill(cnv, st, cellPaddingX, height/2, stateText(r), col, c.lvItems.FontFamily(), c.lvItems.FontSize())
}

// fitText shortens the text with an ellipsis to fit the width
func fitText(fontFamily string, fontSize float64, text string, width int) string {
	if w, _, err := ui.MeasureText(fontFamily, fontSize, text); err != nil || w <= width {
		return text
	}
	runes := []rune(text)
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		shortened := strings.TrimRight(string(runes), " ") + "…"
		if w, _, _ := ui.MeasureText(fontFamily, fontSize, shortened); w <= width {
			return shortened
		}
	}
	return ""
}

// shareState is how the values of a host go to u00.io, shown by the mark before its name
type shareState int

const (
	shareOff     shareState = iota // not shared: no mark
	shareIdle                      // shared, but nothing goes now: pinging is stopped or nothing sent yet
	shareActive                    // the values get there
	shareFailing                   // the last value did not get there
)

// shareActiveFor is how long after the last value got there the host still counts as sending
const shareActiveFor = 15 * time.Second

func hostShareState(h *config.ConfigHost) shareState {
	if !h.Share || h.ShareURL() == "" {
		return shareOff
	}
	if !system.Get().IsRunning() {
		return shareIdle
	}
	status := system.GetShareStatus(h.ShareKey)
	switch {
	case status.At.IsZero():
		return shareIdle
	case status.Err != nil:
		return shareFailing
	case time.Since(status.At) < shareActiveFor+h.Interval()*2:
		return shareActive
	}
	return shareIdle
}

// shareIcon returns the u00.io mark in the color of the state, nil for none
func (c *LeftWidget) shareIcon(state shareState) image.Image {
	var col color.RGBA
	switch state {
	case shareActive:
		col = colorTrend.get()
	case shareFailing:
		col = colorFailed.get()
	case shareIdle:
		col = colorMuted.get()
	default:
		return nil
	}
	if icon, ok := c.shareIcons[col]; ok {
		return icon
	}
	base := loadIcon("share-16")
	if base == nil {
		return nil
	}
	icon := tintIcon(base, col)
	c.shareIcons[col] = icon
	return icon
}

// trendPoint is one point of the trend chart: the slowest reply of its time
// slot, and whether a ping of it was lost
type trendPoint struct {
	hasValue bool
	value    float64 // ms
	lost     bool
}

// trendCache keeps the computed chart of a host for a moment: the table is
// repainted much more often than the chart changes
type trendCache struct {
	at     time.Time
	points []trendPoint
}

// trendKey is a chart of a host with so many points: the table, the tiles
// and the lanes draw charts of different widths
type trendKey struct {
	hostID string
	n      int
}

// trendPoints returns n points of the trend chart of the host, computed
// at most once a trendCacheLife
func (c *LeftWidget) trendPoints(hostID string, n int) []trendPoint {
	n = max(n, 1)
	key := trendKey{hostID, n}
	if cached, ok := c.trends[key]; ok && time.Since(cached.at) < trendCacheLife {
		return cached.points
	}
	points := trendPointsOf(system.Get().GetHostHistory(hostID), time.Now(), n)
	c.trends[key] = trendCache{at: time.Now(), points: points}
	return points
}

// trendPointsOf splits the statsPeriod before now into n slots.
// The slots are fixed in time, not counted from now: a ping stays in its slot
// while the chart moves on, so the points don't jump between redraws.
func trendPointsOf(history *system.HostHistory, now time.Time, n int) []trendPoint {
	points := make([]trendPoint, n)
	if history == nil {
		return points
	}
	slot := statsPeriod / time.Duration(n)
	// The last slot is the one now is in
	to := now.Truncate(slot).Add(slot)
	from := to.Add(-slot * time.Duration(n))
	history.Visit(from, to, func(s system.HistorySample) {
		if s.Gap || s.DT.Before(from) || !s.DT.Before(to) {
			return
		}
		p := &points[int(s.DT.Sub(from)/slot)]
		if !s.OK {
			p.lost = true
			return
		}
		ms := float64(s.PingTime) / float64(time.Millisecond)
		if !p.hasValue || ms > p.value {
			p.value = ms
		}
		p.hasValue = true
	})
	return points
}

// drawTrendCell draws the replies of the last minutes in the kind of the
// style: a slow host in its state color, except in the classic themes
func (c *LeftWidget) drawTrendCell(cnv *ui.Canvas, width, height int, hostID string, r hostRow) {
	st := currentStyle
	chartW := width - cellPaddingX*2
	top, bottom := trendPaddingY, height-trendPaddingY
	step := trendStep(st)
	if chartW < step || bottom <= top {
		return
	}
	points := c.trendPoints(hostID, chartW/step)
	col := st.Trend
	if !st.classic() && r.isSlow() {
		col = st.Slow
	}
	drawTrend(cnv, st, cellPaddingX, top, chartW, bottom-top, points, step, col)
}

func withAlpha(col color.RGBA, a uint8) color.RGBA {
	// RGBA is premultiplied: the color is scaled with the alpha
	return color.RGBA{
		R: uint8(uint16(col.R) * uint16(a) / 255),
		G: uint8(uint16(col.G) * uint16(a) / 255),
		B: uint8(uint16(col.B) * uint16(a) / 255),
		A: a,
	}
}

func colorToRGBA(col color.Color) color.RGBA {
	return color.RGBAModel.Convert(col).(color.RGBA)
}
