package forms

import (
	"fmt"
	"slices"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

// DowntimeDialog lists the periods of the last day when a host did not reply
type DowntimeDialog struct {
	ui.DialogContent

	btnClose *ui.Button
}

func NewDowntimeDialog(host *config.ConfigHost) *DowntimeDialog {
	var c DowntimeDialog
	c.InitWidget()

	var outages []system.Outage
	if history := system.Get().GetHostHistory(host.ID); history != nil {
		outages = history.Outages(time.Now().Add(-24 * time.Hour))
	}
	slices.Reverse(outages) // the latest first

	if len(outages) == 0 {
		lbl := ui.NewLabel("No downtime in the last 24 hours")
		lbl.SetTextAlign(ui.HAlignCenter)
		lbl.SetXExpandable(true)
		c.AddWidget(0, 0, ui.NewVSpacer())
		c.AddWidget(1, 0, lbl)
		c.AddWidget(2, 0, ui.NewVSpacer())
	} else {
		c.AddWidget(0, 0, ui.NewLabel(downTimesText(len(outages))))
		table := ui.NewTable()
		table.SetSelectingRows(true)
		table.SetColumnCount(3)
		for i, col := range []struct {
			name  string
			width int
		}{{"Started", 180}, {"Ended", 180}, {"Duration", 120}} {
			table.SetColumnName(i, col.name)
			table.SetColumnWidth(i, col.width)
		}
		table.SetRowCount(len(outages))
		for row, o := range outages {
			table.SetCellText2(row, 0, formatDowntimeTime(o.Start))
			if o.Ongoing {
				table.SetCellText2(row, 1, "still down")
				table.SetCellText2(row, 2, formatSince(time.Since(o.Start)))
			} else {
				table.SetCellText2(row, 1, formatDowntimeTime(o.End))
				table.SetCellText2(row, 2, formatSince(o.End.Sub(o.Start)))
			}
		}
		c.AddWidget(1, 0, table)
	}

	buttons := ui.NewPanel()
	c.AddWidget(3, 0, buttons)
	c.btnClose = ui.NewButton("Close")
	c.btnClose.SetOnClick(func() { c.Form().Close() })
	buttons.AddWidget(0, 0, ui.NewHSpacer())
	buttons.AddWidget(0, 1, c.btnClose)

	c.OnDialogShow = func() {
		c.Form().SetTitle("Downtime - " + hostDisplayName(host))
		c.Form().SetSize(540, 400)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnClose)
		c.Form().SetCancelButton(c.btnClose)
	}
	return &c
}

// formatDowntimeTime shows the time, with the date when it is not today
func formatDowntimeTime(t time.Time) string {
	if y, m, d := t.Date(); y == time.Now().Year() && m == time.Now().Month() && d == time.Now().Day() {
		return t.Format("15:04:05")
	}
	return t.Format("Jan 2 15:04:05")
}

// downTimesText says how many times the host was down
func downTimesText(n int) string {
	if n == 1 {
		return "Down once in the last 24 hours:"
	}
	return fmt.Sprintf("Down %d times in the last 24 hours:", n)
}
