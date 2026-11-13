package forms

import (
	"fmt"
	"image/color"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

// kpiStrip is the row of cards with the totals over the hosts, shown above
// them in the styles with Style.KPIs: how many hosts and in which state,
// the median ping time and how available they were in the last day
type kpiStrip struct {
	ui.Widget

	// The availability over a day is costly to count: it is kept a while
	uptime   float64 // percent, NaN - unknown
	uptimeAt time.Time
}

const (
	kpiUptimePeriod = 24 * time.Hour
	kpiUptimeLife   = 10 * time.Second
)

func newKPIStrip() *kpiStrip {
	k := &kpiStrip{uptime: math.NaN()}
	k.InitWidget()
	k.SetXExpandable(true)
	k.SetOnPaint(k.paint)
	k.applyTheme()
	ui.OnThemeChanged.Connect(func(bool) { k.applyTheme() })
	return k
}

// applyTheme sizes the cards for the font of the style
func (k *kpiStrip) applyTheme() {
	h := int(ui.ThemeFontSize()*5.8) + viewGap
	k.SetMinHeight(h)
	k.SetMaxHeight(h)
}

// kpi is one card: what is counted, the value and a word under it
type kpi struct {
	label, value, note string
	col                color.RGBA // of the value and the bar on the left; zero - the text
}

func (k *kpiStrip) cards() []kpi {
	st := currentStyle
	hosts := config.Get().Hosts
	ok, slow, down := 0, 0, 0
	var firstSlow, firstDown string
	var times []float64
	for _, h := range hosts {
		r := hostRowOf(h)
		switch {
		case !r.processed:
		case r.failed:
			down++
			if firstDown == "" {
				firstDown = r.name
			}
		case r.isSlow():
			slow++
			if firstSlow == "" {
				firstSlow = r.name
			}
		default:
			ok++
		}
		if r.stats.Sent > r.stats.Lost {
			times = append(times, float64(r.stats.Avg)/float64(time.Millisecond))
		}
	}
	median := "-"
	if len(times) > 0 {
		sort.Float64s(times)
		m := times[len(times)/2]
		if len(times)%2 == 0 {
			m = (times[len(times)/2-1] + m) / 2
		}
		median = formatMs(time.Duration(m*float64(time.Millisecond))) + " " + T().Ms
	}
	k.countUptime(hosts)
	uptime := "-"
	if !math.IsNaN(k.uptime) {
		uptime = fmt.Sprintf("%.1f%%", k.uptime)
		if k.uptime >= 99.95 {
			uptime = "100%"
		}
	}
	percent := ""
	if len(hosts) > 0 {
		percent = fmt.Sprintf("%d%%", ok*100/len(hosts))
	}
	return []kpi{
		{T().StatusHosts, fmt.Sprint(len(hosts)), "", color.RGBA{}},
		{T().StatusOK, fmt.Sprint(ok), percent, st.OK},
		{T().StatusSlow, fmt.Sprint(slow), firstSlow, st.Slow},
		{T().StatusDown, fmt.Sprint(down), firstDown, st.Down},
		{T().KPIMedian, median, "", color.RGBA{}},
		{T().KPIUptime, uptime, "", st.OK},
	}
}

// countUptime counts the share of the pings of the last day that got a
// reply, over all the hosts
func (k *kpiStrip) countUptime(hosts []*config.ConfigHost) {
	if time.Since(k.uptimeAt) < kpiUptimeLife {
		return
	}
	k.uptimeAt = time.Now()
	sent, lost := 0, 0
	from := time.Now().Add(-kpiUptimePeriod)
	for _, h := range hosts {
		history := system.Get().GetHostHistory(h.ID)
		if history == nil {
			continue
		}
		s := history.Stats(from)
		sent += s.Sent
		lost += s.Lost
	}
	k.uptime = math.NaN()
	if sent > 0 {
		k.uptime = 100 * float64(sent-lost) / float64(sent)
	}
}

func (k *kpiStrip) paint(cnv *ui.Canvas) {
	st := currentStyle
	p := st.Palette
	cnv.FillRect(0, 0, k.Width(), k.Height(), p.Window)
	cards := k.cards()
	n := len(cards)
	w := (k.Width() - viewGap*(n-1)) / n
	h := k.Height() - viewGap
	fs := ui.ThemeFontSize()
	family := ui.ThemeFontFamily()
	for i, c := range cards {
		x := i * (w + viewGap)
		fillCard(cnv, st, x, 0, w, h, p.Base, p.Border, 1)
		if c.col.A > 0 {
			cnv.FillRect(x+1, 1+st.TileRadius/2, 3, h-2-st.TileRadius, c.col)
		}
		cnv.Save()
		cnv.TranslateAndClip(x+16, 0, w-24, h)
		cnv.SetFontFamily(family)
		cnv.SetHAlign(ui.HAlignLeft)
		cnv.SetVAlign(ui.VAlignTop)
		cnv.SetFontSize(fs - 3)
		cnv.SetColor(st.Muted)
		cnv.DrawText(0, int(fs*0.7), w-24, int(fs*1.4), fitText(family, fs-3, strings.ToUpper(c.label), w-24))
		cnv.SetFontSize(fs * 1.8)
		cnv.SetFontBold(true)
		col := p.Text
		if c.col.A > 0 && c.value != "0" {
			col = c.col
		}
		cnv.SetColor(col)
		cnv.DrawText(0, int(fs*1.75), w-24, int(fs*2.4), c.value)
		cnv.SetFontBold(false)
		cnv.SetFontSize(fs - 3)
		cnv.SetColor(st.Muted)
		cnv.DrawText(0, int(fs*4.0), w-24, int(fs*1.4), fitText(family, fs-3, c.note, w-24))
		cnv.Restore()
	}
}
