package forms

import (
	"image/color"
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

// Hosts shown on the chart at once, one area per host
const detailsMaxHosts = 4

// detailsPeriods are the time ranges the chart can show, named by Strings.Periods; the history keeps a day
var detailsPeriods = [len(Strings{}.Periods)]time.Duration{5 * time.Minute, time.Hour, 24 * time.Hour}

var lastCreatedDetailsWidget *DetailsWidget

// DetailsWidget shows the ping history of the hosts selected in the table
type DetailsWidget struct {
	ui.Widget

	lblTitle   *ui.Label
	btnPeriods []*ui.Button
	// The numbers of the host on the chart, above it
	stats *detailsStats
	chart *ui.TimeChart

	// The time range shown, up to now
	period time.Duration

	// IDs of the hosts on the chart, to rebuild it when the selection changes
	shownHostIDs string
	// Names of the hosts on the chart, for the title
	shownHostNames []string
}

func NewDetailsWidget() *DetailsWidget {
	var c DetailsWidget
	c.InitWidget()
	c.SetPanelPadding(0)
	c.SetMinWidth(400)

	header := ui.NewPanel()
	header.SetPanelPadding(0)
	c.AddWidget(0, 0, header)
	c.lblTitle = header.AddLabel(0, 0, "")
	// Differs from any selection, so the first check builds the chart (and the title for no selection)
	c.shownHostIDs = "-"
	c.lblTitle.SetXExpandable(true)
	for i, period := range detailsPeriods {
		btn := ui.NewButton("")
		btn.SetTextFunc(func() string { return T().Periods[i] })
		btn.SetMinWidth(48)
		btn.OnClick.Connect(func(ui.Widgeter) { c.SetPeriod(period) })
		header.AddWidget(0, i+1, btn)
		c.btnPeriods = append(c.btnPeriods, btn)
	}

	c.stats = newDetailsStats()
	c.AddWidget(1, 0, c.stats)
	c.chart = ui.NewTimeChart()
	c.AddWidget(2, 0, c.chart)
	c.applyTheme()
	ui.OnThemeChanged.Connect(func(bool) { c.applyTheme() })

	c.AddTimer(500, c.timerUpdate)

	lastCreatedDetailsWidget = &c
	c.SetPeriod(detailsPeriods[0])
	return &c
}

// SetPeriod shows the last period of time on the chart
func (c *DetailsWidget) SetPeriod(period time.Duration) {
	c.period = period
	for i, btn := range c.btnPeriods {
		// The active period is highlighted
		if detailsPeriods[i] == period {
			btn.SetRole("primary")
		} else {
			btn.SetRole("")
		}
	}
	if c.Form() != nil {
		c.updateTimeRange()
		// Leave a zoomed-in view: the period button shows the whole period
		c.chart.ResetZoom()
	}
}

func (c *DetailsWidget) Period() time.Duration {
	return c.period
}

// Refresh shows the hosts selected in the table and the current time range right away
func (c *DetailsWidget) Refresh() {
	c.checkSelection()
	c.timerUpdate()
}

// checkSelection rebuilds the chart when other hosts are selected in the table
func (c *DetailsWidget) checkSelection() {
	if !c.IsVisible() {
		return
	}

	hosts := lastCreatedLeftWidget.GetSelectedHostConfigs()
	if len(hosts) > detailsMaxHosts {
		hosts = hosts[:detailsMaxHosts]
	}
	ids := make([]string, 0, len(hosts))
	for _, h := range hosts {
		ids = append(ids, h.ID)
	}
	idsStr := strings.Join(ids, ",")
	if idsStr != c.shownHostIDs {
		c.shownHostIDs = idsStr
		c.rebuildChart(hosts)
		c.updateTimeRange()
	}
}

func (c *DetailsWidget) timerUpdate() {
	if !c.IsVisible() {
		return
	}
	c.stats.Update()
	c.updateTimeRange()
}

// applyTheme draws the chart in the colors of the style: its lines in the
// trend colors, over a fill in the styles that have one
func (c *DetailsWidget) applyTheme() {
	st := currentStyle
	th := *ui.CurrentChartTheme()
	if !st.classic() {
		th.Palette = append([]color.RGBA{st.Trend, st.Trend2}, th.Palette...)
		th.BadHatch = alpha(st.Down, 110)
		th.Background = st.Palette.Window
		th.PlotBackground = st.Palette.Base
	}
	c.chart.SetTheme(&th)
	for _, area := range c.chart.Areas() {
		for _, series := range area.Series() {
			styleSeries(series)
		}
	}
	c.chart.Update()
}

// styleSeries draws a line of the chart as the style wants it
func styleSeries(series *ui.TimeChartSeries) {
	st := currentStyle
	line, ok := series.Series().(*ui.LineSeries)
	if !ok {
		return
	}
	line.SetLineWidth(st.ChartLine)
	line.SetFill(st.ChartFill > 0, st.ChartFill, true)
}

func (c *DetailsWidget) updateTimeRange() {
	now := time.Now()
	c.chart.SetDefaultTimeRange(now.Add(-c.period), now)
	c.Form().Update()
}

func (c *DetailsWidget) rebuildChart(hosts []*config.ConfigHost) {
	c.chart.RemoveAllAreas()

	c.shownHostNames = c.shownHostNames[:0]
	for _, h := range hosts {
		name := hostDisplayName(h)
		c.shownHostNames = append(c.shownHostNames, name)
		area := c.chart.AddArea()
		series := area.AddSeries(name+", "+T().Ms, &historySource{hostID: h.ID})
		series.SetPaletteColor(0)
		styleSeries(series)
	}
	c.stats.setHost(nil)
	if len(hosts) == 1 {
		c.stats.setHost(hosts[0])
	}
	c.updateTitle()
}

// applyLanguage rebuilds the chart with the texts in the new language
func (c *DetailsWidget) applyLanguage() {
	c.shownHostIDs = "-"
	c.updateTitle()
	c.checkSelection()
}

// updateTitle names the hosts on the chart; the series names are updated when the chart is rebuilt
func (c *DetailsWidget) updateTitle() {
	if len(c.shownHostNames) == 0 {
		c.lblTitle.SetText(T().SelectAHost)
		return
	}
	c.lblTitle.SetText(T().PingHistory + strings.Join(c.shownHostNames, ", "))
}

func hostDisplayName(h *config.ConfigHost) string {
	if h.DisplayName != "" {
		return h.DisplayName
	}
	return h.Address()
}

// historySource feeds the in-memory ping history of a host to the chart.
// The history is looked up on every call: it is created when the host
// is added and dropped when the host is removed from the config.
type historySource struct {
	hostID string
}

func (c *historySource) GetData(from, to time.Time, groupDuration time.Duration) []ui.TimeChartPoint {
	history := system.Get().GetHostHistory(c.hostID)
	if history == nil {
		return nil
	}
	// Up to a day of samples: they are aggregated as they are read, without copying them
	agg := ui.NewTimeChartAggregator(groupDuration, int(to.Sub(from)/max(groupDuration, time.Second))+3)
	history.Visit(from, to, func(s system.HistorySample) {
		switch {
		case s.Gap:
			agg.Add(ui.NewTimeChartGap(s.DT))
		case s.OK:
			agg.Add(ui.NewTimeChartValue(s.DT, float64(s.PingTime.Microseconds())/1000))
		default:
			agg.Add(ui.NewTimeChartBad(s.DT))
		}
	})
	return agg.Points()
}

// detailsStats shows the numbers of the host on the chart over the last 5
// minutes as cards: the last time, the average, the loss and the jitter.
// It is hidden while the chart shows no host or several.
type detailsStats struct {
	ui.Widget
	host *config.ConfigHost
}

func newDetailsStats() *detailsStats {
	d := &detailsStats{}
	d.InitWidget()
	d.SetXExpandable(true)
	d.SetOnPaint(d.paint)
	d.SetVisible(false)
	d.applyTheme()
	ui.OnThemeChanged.Connect(func(bool) { d.applyTheme() })
	return d
}

func (d *detailsStats) applyTheme() {
	h := int(ui.ThemeFontSize()*3.6) + 8
	d.SetMinHeight(h)
	d.SetMaxHeight(h)
}

func (d *detailsStats) setHost(h *config.ConfigHost) {
	d.host = h
	if visible := h != nil; visible != d.IsVisible() {
		d.SetVisible(visible)
		if f := d.Form(); f != nil {
			f.UpdateLayout()
		}
	}
	d.Update()
}

func (d *detailsStats) paint(cnv *ui.Canvas) {
	if d.host == nil {
		return
	}
	st := currentStyle
	p := st.Palette
	cnv.FillRect(0, 0, d.Width(), d.Height(), p.Window)
	r := hostRowOf(d.host)
	t := r.texts()
	cols := T().Columns
	cards := []struct {
		label, value string
		col          color.RGBA
	}{
		{cols.Time, t[colTime], p.Text},
		{cols.Avg, t[colAvg], p.Text},
		{cols.Loss, t[colLoss], p.Text},
		{cols.Jitter, t[colJitter], p.Text},
	}
	if r.failed {
		cards[0].value, cards[0].col = strings.ToUpper(stateText(r)), st.Down
	}
	if t[colLoss] != "0" && t[colLoss] != "-" {
		cards[2].col = st.Slow
		if r.stats.Lost == r.stats.Sent {
			cards[2].col = st.Down
		}
	}
	if r.isSlow() {
		cards[1].col = st.Slow
	}
	const gap = 8
	fs := ui.ThemeFontSize()
	family := ui.ThemeFontFamily()
	n := len(cards)
	w := (d.Width() - gap*(n-1)) / n
	h := d.Height() - 8
	for i, c := range cards {
		x := i * (w + gap)
		fillCard(cnv, st, x, 4, w, h, p.Base, p.Border, 1)
		cnv.Save()
		cnv.TranslateAndClip(x+10, 4, w-16, h)
		cnv.SetFontFamily(family)
		cnv.SetHAlign(ui.HAlignLeft)
		cnv.SetVAlign(ui.VAlignTop)
		cnv.SetFontSize(fs - 4)
		cnv.SetColor(st.Muted)
		cnv.DrawText(0, int(fs*0.45), w-16, int(fs*1.3), fitText(family, fs-4, strings.ToUpper(c.label), w-16))
		cnv.SetFontSize(fs * 1.35)
		cnv.SetFontBold(true)
		cnv.SetColor(c.col)
		cnv.DrawText(0, int(fs*1.45), w-16, int(fs*1.9), fitText(family, fs*1.35, c.value, w-16))
		cnv.SetFontBold(false)
		cnv.Restore()
	}
}
