package forms

import (
	"image"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

// The host table draws two cells itself: the name, after a dot in the color of
// the host state, and the trend, a small chart of the last minutes. The other
// cells are plain text; only the values that mean something are colored.

var (
	// colorMuted is for the text that is there but does not matter much: "-", "OK"
	colorMuted = themeColors{ui.ColorFromHex("#7d7d7d"), ui.ColorFromHex("#8c8c8c")}
	// colorTrend is the line of the trend chart
	colorTrend = themeColors{ui.ColorFromHex("#4c9aff"), ui.ColorFromHex("#1d6fb5")}
)

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

// drawNameCell draws the dot of the host state, the bell of a host that beeps,
// and the name
func (c *LeftWidget) drawNameCell(cnv *ui.Canvas, width, height int, r hostRow, notify bool) {
	x := cellPaddingX
	cnv.DrawImage(x, (height-statusDotSize)/2, statusDot(colorToRGBA(r.color())))
	x += statusDotSize + statusDotGap
	if notify && c.iconNotify != nil {
		cnv.DrawImage(x, (height-notifyIconSize)/2, c.iconNotify)
	}
	x += notifyIconSize + notifyIconGap

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
	width  int
	points []trendPoint
}

// trendPoints returns the trend chart of the host, computed at most once a trendCacheLife
func (c *LeftWidget) trendPoints(hostID string, width int) []trendPoint {
	if cached, ok := c.trends[hostID]; ok && cached.width == width && time.Since(cached.at) < trendCacheLife {
		return cached.points
	}
	n := max(width/trendBarWidth, 1)
	points := trendPointsOf(system.Get().GetHostHistory(hostID), time.Now(), n)
	c.trends[hostID] = trendCache{at: time.Now(), width: width, points: points}
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

// drawTrendCell draws the replies of the last minutes as a line, scaled to
// the slowest one; lost pings are ticks along the bottom
func (c *LeftWidget) drawTrendCell(cnv *ui.Canvas, width, height int, hostID string) {
	chartW := width - cellPaddingX*2
	top, bottom := trendPaddingY, height-trendPaddingY
	if chartW < trendBarWidth || bottom <= top {
		return
	}
	points := c.trendPoints(hostID, chartW)

	cnv.FillRect(cellPaddingX, bottom, chartW, 1, withAlpha(colorMuted.get(), 90))

	maxValue := 0.0
	for _, p := range points {
		if p.hasValue {
			maxValue = math.Max(maxValue, p.value)
		}
	}
	// Some room above the top, and no huge swings for a host that replies in no time
	maxValue = math.Max(maxValue*1.15, 1)
	yOf := func(v float64) float64 {
		return float64(bottom) - v/maxValue*float64(bottom-top)
	}

	lineColor := colorTrend.get()
	lostColor := colorFailed.get()
	prevX, prevY, prevOK := 0.0, 0.0, false
	for i, p := range points {
		x := cellPaddingX + i*trendBarWidth
		if p.lost {
			cnv.FillRect(x, bottom-trendLostHeight+1, trendBarWidth, trendLostHeight, lostColor)
		}
		if !p.hasValue {
			prevOK = false
			continue
		}
		y := yOf(p.value)
		cx := float64(x) + float64(trendBarWidth)/2
		if prevOK {
			cnv.DrawLineF(prevX, prevY, cx, y, lineColor)
		} else if i+1 == len(points) || !points[i+1].hasValue {
			// A reply between gaps or losses has no line to it: a short dash
			cnv.FillRect(x, int(math.Round(y)), trendBarWidth, 1, lineColor)
		}
		prevX, prevY, prevOK = cx, y, true
	}
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
