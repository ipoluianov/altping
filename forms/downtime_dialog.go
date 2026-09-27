package forms

import (
	"slices"
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

// DowntimeDialog lists the periods of the last day when the hosts did not
// reply; for several hosts the outages of all of them together, with the
// host of each
type DowntimeDialog struct {
	ui.DialogContent

	btnClose *ui.Button
}

// hostOutage is an outage of the host
type hostOutage struct {
	host *config.ConfigHost
	system.Outage
}

// downtimeTitleHosts is how many host names the title lists
const downtimeTitleHosts = 3

func NewDowntimeDialog(hosts []*config.ConfigHost) *DowntimeDialog {
	var c DowntimeDialog
	c.InitWidget()

	var outages []hostOutage
	for _, h := range hosts {
		if history := system.Get().GetHostHistory(h.ID); history != nil {
			for _, o := range history.Outages(time.Now().Add(-24 * time.Hour)) {
				outages = append(outages, hostOutage{host: h, Outage: o})
			}
		}
	}
	// The latest first
	slices.SortStableFunc(outages, func(a, b hostOutage) int { return b.Start.Compare(a.Start) })
	group := len(hosts) > 1

	if len(outages) == 0 {
		lbl := ui.NewLabel(T().NoDowntime)
		lbl.SetTextAlign(ui.HAlignCenter)
		lbl.SetXExpandable(true)
		c.AddWidget(0, 0, ui.NewVSpacer())
		c.AddWidget(1, 0, lbl)
		c.AddWidget(2, 0, ui.NewVSpacer())
	} else {
		summary := T().DownTimes(len(outages))
		if group {
			summary = T().GroupDownTimes(len(outages))
		}
		c.AddWidget(0, 0, ui.NewLabel(summary))

		columns := []struct {
			name  string
			width int
		}{{T().Started, 180}, {T().Ended, 180}, {T().Duration, 120}}
		if group {
			columns = slices.Insert(columns, 0, struct {
				name  string
				width int
			}{T().Columns.Name, 160})
		}
		table := ui.NewTable()
		table.SetSelectingRows(true)
		table.SetColumnCount(len(columns))
		for i, col := range columns {
			table.SetColumnName(i, col.name)
			table.SetColumnWidth(i, col.width)
		}
		table.SetRowCount(len(outages))
		for row, o := range outages {
			cells := []string{formatDowntimeTime(o.Start)}
			if o.Ongoing {
				cells = append(cells, T().StillDown, formatSince(time.Since(o.Start)))
			} else {
				cells = append(cells, formatDowntimeTime(o.End), formatSince(o.End.Sub(o.Start)))
			}
			if group {
				cells = slices.Insert(cells, 0, hostDisplayName(o.host))
			}
			for col, text := range cells {
				table.SetCellText2(row, col, text)
			}
		}
		c.AddWidget(1, 0, table)
	}

	buttons := ui.NewPanel()
	c.AddWidget(3, 0, buttons)
	c.btnClose = ui.NewButton(T().Close)
	c.btnClose.SetOnClick(func() { c.Form().Close() })
	buttons.AddWidget(0, 0, ui.NewHSpacer())
	buttons.AddWidget(0, 1, c.btnClose)

	c.OnDialogShow = func() {
		c.Form().SetTitle(T().DowntimeTitle + " - " + downtimeTitleNames(hosts))
		width := 540
		if group {
			width += 160
		}
		c.Form().SetSize(width, 400)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnClose)
		c.Form().SetCancelButton(c.btnClose)
	}
	return &c
}

// downtimeTitleNames lists the first hosts for the title
func downtimeTitleNames(hosts []*config.ConfigHost) string {
	names := make([]string, 0, downtimeTitleHosts)
	for i, h := range hosts {
		if i == downtimeTitleHosts {
			names = append(names, "…")
			break
		}
		names = append(names, hostDisplayName(h))
	}
	return strings.Join(names, ", ")
}

// formatDowntimeTime shows the time, with the date when it is not today
func formatDowntimeTime(t time.Time) string {
	if y, m, d := t.Date(); y == time.Now().Year() && m == time.Now().Month() && d == time.Now().Day() {
		return t.Format("15:04:05")
	}
	return t.Format(T().DateTimeLayout)
}
