package forms

import (
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

const (
	detailsLiveWindow = 5 * time.Minute
	// Hosts shown on the chart at once, one area per host
	detailsMaxHosts = 4
)

var lastCreatedDetailsWidget *DetailsWidget

// DetailsWidget shows the ping history of the hosts selected in the table
type DetailsWidget struct {
	ui.Widget

	lblTitle *ui.Label
	chart    *ui.TimeChart

	// IDs of the hosts on the chart, to rebuild it when the selection changes
	shownHostIDs string
}

func NewDetailsWidget() *DetailsWidget {
	var c DetailsWidget
	c.InitWidget()
	c.SetPanelPadding(0)
	c.SetMinWidth(400)

	c.lblTitle = c.AddLabel(0, 0, "")
	c.lblTitle.SetXExpandable(true)

	c.chart = ui.NewTimeChart()
	c.AddWidget(1, 0, c.chart)

	c.AddTimer(500, c.timerUpdate)

	lastCreatedDetailsWidget = &c
	return &c
}

func (c *DetailsWidget) timerUpdate() {
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
	}

	now := time.Now()
	c.chart.SetDefaultTimeRange(now.Add(-detailsLiveWindow), now)
	c.Form().Update()
}

func (c *DetailsWidget) rebuildChart(hosts []*config.ConfigHost) {
	c.chart.RemoveAllAreas()

	if len(hosts) == 0 {
		c.lblTitle.SetText("Select a host to see its ping history")
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
	samples := history.Range(from, to)
	points := make([]ui.TimeChartPoint, len(samples))
	for i, s := range samples {
		switch {
		case s.Gap:
			points[i] = ui.NewTimeChartGap(s.DT)
		case s.OK:
			points[i] = ui.NewTimeChartValue(s.DT, float64(s.PingTime.Microseconds())/1000)
		default:
			points[i] = ui.NewTimeChartBad(s.DT)
		}
	}
	return ui.TimeChartDownsample(points, from, to, groupDuration)
}
