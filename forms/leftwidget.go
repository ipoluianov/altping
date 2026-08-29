package forms

import (
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
	c.SetMinWidth(700)
	c.SetMaxWidth(700)

	c.lvItems.SetSelectingRow(true)
	c.lvItems.SetSelectingCell(false)
	c.lvItems.SetColumnCount(7)
	c.lvItems.SetColumnWidth(0, 30)
	c.lvItems.SetColumnWidth(1, 150)
	c.lvItems.SetColumnWidth(2, 150)
	c.lvItems.SetColumnWidth(3, 100)
	c.lvItems.SetColumnWidth(4, 100)
	c.lvItems.SetColumnWidth(5, 100)
	c.lvItems.SetColumnWidth(6, 100)
	c.lvItems.SetColumnName(0, "*")
	c.lvItems.SetColumnName(1, "Hostname")
	c.lvItems.SetColumnName(2, "Resp. Time, ms")
	c.lvItems.SetColumnName(3, "OK")
	c.lvItems.SetColumnName(4, "ERR")
	c.lvItems.SetColumnName(5, "IP")
	c.lvItems.SetColumnName(6, "ID")

	c.loadHosts()

	c.AddTimer(100, c.timerUpdate)

	return &c
}

func (c *LeftWidget) loadHosts() {
	config := config.Get()
	c.lvItems.SetRowCount(len(config.Hosts))
	for i, host := range config.Hosts {
		c.lvItems.SetCellData2(i, 0, host)
		c.lvItems.SetCellText2(i, 0, "-")
		c.lvItems.SetCellText2(i, 1, host.Hostname)
		c.lvItems.SetCellText2(i, 2, "")
		c.lvItems.SetCellText2(i, 3, "")
		c.lvItems.SetCellText2(i, 4, "")
		c.lvItems.SetCellText2(i, 5, "")
		c.lvItems.SetCellText2(i, 6, host.ID)

	}
}

func (c *LeftWidget) timerUpdate() {
	for row := 0; row < c.lvItems.RowCount(); row++ {
		hostConfig := c.lvItems.GetCellData2(row, 0).(*config.ConfigHost)
		if hostConfig != nil {
			lastState := system.Get().GetHostLastState(hostConfig.ID)
			c.lvItems.SetCellText2(row, 2, strconv.FormatInt(int64(lastState.PingTime.Milliseconds()), 10))
			c.lvItems.SetCellText2(row, 3, strconv.FormatInt(int64(lastState.StatOK), 10))
			c.lvItems.SetCellText2(row, 4, strconv.FormatInt(int64(lastState.StatERR), 10))
			c.lvItems.SetCellText2(row, 5, lastState.StatIP)
			c.lvItems.SetCellText2(row, 6, lastState.ConfigHost.ID)
		}
	}
}
