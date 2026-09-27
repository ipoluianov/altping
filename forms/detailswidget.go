package forms

import (
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

// Hosts shown on the chart at once, one area per host
const detailsMaxHosts = 4

// detailsPeriods are the time ranges the chart can show; the history keeps a day
var detailsPeriods = []struct {
	name     string
	duration time.Duration
}{
	{"5m", 5 * time.Minute},
	{"1h", time.Hour},
	{"24h", 24 * time.Hour},
}

var lastCreatedDetailsWidget *DetailsWidget

// DetailsWidget shows the ping history of the hosts selected in the table
type DetailsWidget struct {
	ui.Widget

	lblTitle   *ui.Label
	btnPeriods []*ui.Button
	chart      *ui.TimeChart

	// The time range shown, up to now
	period time.Duration

	// IDs of the hosts on the chart, to rebuild it when the selection changes
	shownHostIDs string
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
	for i, p := range detailsPeriods {
		btn := ui.NewButton(p.name)
		btn.SetMinWidth(48)
		btn.SetMaxWidth(48)
		btn.SetOnClick(func() { c.SetPeriod(p.duration) })
		header.AddWidget(0, i+1, btn)
		c.btnPeriods = append(c.btnPeriods, btn)
	}

	c.chart = ui.NewTimeChart()
	c.AddWidget(1, 0, c.chart)

	// The selection is checked often so the chart follows the table without a visible delay
	c.AddTimer(50, c.checkSelection)
	c.AddTimer(500, c.timerUpdate)

	lastCreatedDetailsWidget = &c
	c.SetPeriod(detailsPeriods[0].duration)
	return &c
}

// SetPeriod shows the last period of time on the chart
func (c *DetailsWidget) SetPeriod(period time.Duration) {
	c.period = period
	for i, btn := range c.btnPeriods {
		// The active period is highlighted
		if detailsPeriods[i].duration == period {
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
	c.updateTimeRange()
}

func (c *DetailsWidget) updateTimeRange() {
	now := time.Now()
	c.chart.SetDefaultTimeRange(now.Add(-c.period), now)
	c.Form().Update()
}

func (c *DetailsWidget) rebuildChart(hosts []*config.ConfigHost) {
	c.chart.RemoveAllAreas()

	if len(hosts) == 0 {
		c.lblTitle.SetText("Select a host")
		return
	}

	names := make([]string, 0, len(hosts))
	for _, h := range hosts {
		name := hostDisplayName(h)
		names = append(names, name)
		area := c.chart.AddArea()
		area.AddSeries(name+", ms", &historySource{hostID: h.ID}).SetPaletteColor(0)
	}
	c.lblTitle.SetText("Ping history: " + strings.Join(names, ", "))
}

func hostDisplayName(h *config.ConfigHost) string {
	if h.DisplayName != "" {
		return h.DisplayName
	}
	return h.Hostname
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
