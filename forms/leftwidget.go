package forms

import (
	"image/color"
	"strconv"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

var lastCreatedLeftWidget *LeftWidget

type LeftWidget struct {
	ui.Widget

	lvItems *ui.Table
}

func NewLeftWidget() *LeftWidget {
	var c LeftWidget
	c.InitWidget()
	c.lvItems = ui.NewTable()
	c.AddWidget(0, 0, c.lvItems)
	//c.SetMinWidth(700)
	//c.SetMaxWidth(700)

	c.lvItems.SetSelectingRows(true)
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

	c.lvItems.SetMultiselect(true)

	// Right click selects the row under the mouse (unless it is already selected) and shows the menu
	menu := ui.NewContextMenu(c.lvItems)
	menu.AddItem("Edit... (E)", func() { lastCreatedTopWidget.onBtnEditItem() }).SetImage(loadIcon("edit-16"))
	menu.AddItem("Remove (Del)", func() { lastCreatedTopWidget.onBtnRemoveItem() }).SetImage(loadIcon("remove-16"))
	c.lvItems.SetContextMenu(menu)

	c.loadHosts()

	c.AddTimer(200, c.timerUpdate)

	lastCreatedLeftWidget = &c

	c.lvItems.Focus()

	c.SetPanelPadding(0)

	return &c
}

func (c *LeftWidget) FocusTable() {
	c.lvItems.Focus()
}

// FullRestart restarts all the hosts, e.g. after another config is opened
func (c *LeftWidget) FullRestart() {
	system.Get().Stop()
	c.loadHosts()
	system.Get().Start()
}

// ApplyHostsChange shows the hosts of the config after they were added, edited or removed.
// Only the changed hosts are started, stopped or restarted, the others keep running.
// The row is selected afterwards (clamped to the table).
func (c *LeftWidget) ApplyHostsChange(selectRow int) {
	c.loadHosts()
	system.Get().SyncHosts()
	c.timerUpdate()
	selectRow = min(selectRow, c.lvItems.RowCount()-1)
	if selectRow >= 0 {
		c.lvItems.SetCurrentCell2(selectRow, 0)
	}
}

// RowOfHost returns the table row of the host, -1 if there is none
func (c *LeftWidget) RowOfHost(id string) int {
	for row := 0; row < c.lvItems.RowCount(); row++ {
		if host, ok := c.lvItems.GetCellData2(row, 0).(*config.ConfigHost); ok && host.ID == id {
			return row
		}
	}
	return -1
}

// loadHosts fills the table with the hosts of the config;
// the state columns are filled by timerUpdate
func (c *LeftWidget) loadHosts() {
	config := config.Get()
	c.lvItems.SetRowCount(len(config.Hosts))
	for i, host := range config.Hosts {
		c.lvItems.SetCellData2(i, 0, host)
		c.lvItems.SetCellText2(i, 0, hostDisplayName(host))
	}
}

func (c *LeftWidget) GetSelectedHostConfigs() []*config.ConfigHost {
	selectedRows := c.lvItems.SelectedRows()
	if len(selectedRows) == 0 {
		return nil
	}
	hosts := make([]*config.ConfigHost, 0, len(selectedRows))
	for _, row := range selectedRows {
		// The selection may still hold rows removed from the table
		if host, ok := c.lvItems.GetCellData2(row, 0).(*config.ConfigHost); ok && host != nil {
			hosts = append(hosts, host)
		}
	}
	return hosts
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
