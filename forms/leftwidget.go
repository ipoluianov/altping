package forms

import (
	"image/color"
	"strconv"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

type LeftWidget struct {
	ui.Widget

	lvItems *ui.Table
}

func NewLeftWidget(onModeChanged func(mode string)) *LeftWidget {
	var c LeftWidget
	c.InitWidget()
	c.lvItems = ui.NewTable()
	c.AddWidgetOnGrid(c.lvItems, 0, 0)
	//c.SetMinWidth(700)
	//c.SetMaxWidth(700)

	c.lvItems.SetSelectingRow(true)
	c.lvItems.SetSelectingCell(false)
	c.lvItems.SetColumnCount(6)
	c.lvItems.SetColumnWidth(0, 150)
	c.lvItems.SetColumnWidth(1, 160)
	c.lvItems.SetColumnWidth(2, 100)
	c.lvItems.SetColumnWidth(3, 50)
	c.lvItems.SetColumnWidth(4, 50)
	c.lvItems.SetColumnWidth(5, 250)
	c.lvItems.SetColumnName(0, "Name")
	c.lvItems.SetColumnName(1, "IP")
	c.lvItems.SetColumnName(2, "Time, ms")
	c.lvItems.SetColumnName(3, "OK")
	c.lvItems.SetColumnName(4, "ERR")
	c.lvItems.SetColumnName(5, "Details")

	c.loadHosts()

	c.AddTimer(100, c.timerUpdate)

	return &c
}

func (c *LeftWidget) loadHosts() {
	config := config.Get()
	c.lvItems.SetRowCount(len(config.Hosts))
	for i, host := range config.Hosts {

		displayName := host.DisplayName
		if displayName == "" {
			displayName = host.Hostname
		}

		c.lvItems.SetCellData2(i, 0, host)
		c.lvItems.SetCellText2(i, 0, displayName)
	}
}

func (c *LeftWidget) timerUpdate() {
	for row := 0; row < c.lvItems.RowCount(); row++ {
		hostConfig := c.lvItems.GetCellData2(row, 0).(*config.ConfigHost)
		if hostConfig == nil {
			continue
		}

		lastState := system.Get().GetHostLastState(hostConfig.ID)
		isProcessed := false
		if lastState.StatOK > 0 || lastState.StatERR > 0 {
			isProcessed = true
		}

		ipStr := ""
		timeStr := ""
		details := ""
		if lastState.LastError != nil {
			details = lastState.LastError.Error()

			switch lastState.LastError.Error() {
			case "timeout":
				details = "TIMEOUT"
			case "cannot resolve hostname":
				details = "CANNOT RESOLVE HOSTNAME"
			}

			timeStr = "-"
			ipStr = "-"
		} else {
			details = "OK"
			timeStr = strconv.FormatInt(int64(lastState.PingTime.Milliseconds()), 10)
			ipStr = lastState.StatIP
		}

		if !isProcessed {
			details = "-"
			timeStr = "-"
			ipStr = "-"
		}

		c.lvItems.SetCellText2(row, 1, ipStr)
		c.lvItems.SetCellText2(row, 2, timeStr)
		c.lvItems.SetCellText2(row, 3, strconv.FormatInt(int64(lastState.StatOK), 10))
		c.lvItems.SetCellText2(row, 4, strconv.FormatInt(int64(lastState.StatERR), 10))
		c.lvItems.SetCellText2(row, 5, details)

		var col color.Color
		col = ui.ColorFromHex("#888888")
		if isProcessed {
			col = ui.ColorFromHex("#1ebd1e")
		}
		if lastState.LastError != nil {
			col = ui.ColorFromHex("#e6660a")
		}
		c.lvItems.SetCellColor(row, 0, col)
		c.lvItems.SetCellColor(row, 1, col)
		c.lvItems.SetCellColor(row, 2, col)
		c.lvItems.SetCellColor(row, 3, col)
		c.lvItems.SetCellColor(row, 4, col)
		c.lvItems.SetCellColor(row, 5, col)
	}
}
