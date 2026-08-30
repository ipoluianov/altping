package forms

import (
	"fmt"
	"strings"

	"github.com/ipoluianov/altping/config"
	"github.com/u00io/nuiforms/ui"
)

type OpenConfigDialog struct {
	ui.Widget

	lvConfigs *ui.Table

	btnOK     *ui.Button
	btnCancel *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel

	onAccept func()
	onCancel func()
}

func NewOpenConfigDialog(currentConfigId string, onAccept func(), onCancel func()) *OpenConfigDialog {
	var c OpenConfigDialog
	c.InitWidget()
	c.panelContent = ui.NewPanel()
	c.AddWidgetOnGrid(c.panelContent, 0, 0)
	c.panelButtons = ui.NewPanel()
	c.AddWidgetOnGrid(c.panelButtons, 1, 0)

	c.btnOK = ui.NewButton("OK")
	c.onAccept = onAccept
	c.onCancel = onCancel
	c.btnOK.SetOnClick(func() {
		if c.onAccept != nil {
			c.onAccept()
		}
	})
	c.btnCancel = ui.NewButton("Cancel")
	c.btnCancel.SetOnClick(func() {
		if c.onCancel != nil {
			c.onCancel()
		}
	})

	c.panelButtons.AddWidgetOnGrid(c.btnOK, 0, 0)
	c.panelButtons.AddWidgetOnGrid(c.btnCancel, 0, 1)

	c.lvConfigs = ui.NewTable()
	c.lvConfigs.SetSelectingCell(false)
	c.panelContent.AddWidgetOnGrid(ui.NewLabel("Configs:"), 0, 0)
	c.panelContent.AddWidgetOnGrid(c.lvConfigs, 1, 0)

	c.lvConfigs.SetColumnCount(3)
	c.lvConfigs.SetColumnWidth(0, 200)
	c.lvConfigs.SetColumnWidth(1, 100)
	c.lvConfigs.SetColumnWidth(2, 100)
	c.lvConfigs.SetColumnName(0, "Name")
	c.lvConfigs.SetColumnName(1, "Hosts Count")
	c.lvConfigs.SetColumnName(2, "Hosts")

	configs := config.Configs()
	c.lvConfigs.SetRowCount(len(configs))
	for i, cfg := range configs {
		c.lvConfigs.SetCellData2(i, 0, cfg)
		c.lvConfigs.SetCellText2(i, 0, cfg.Name)
		c.lvConfigs.SetCellText2(i, 1, fmt.Sprintf("%d", len(cfg.Hosts)))
		hostnames := make([]string, len(cfg.Hosts))
		for j, host := range cfg.Hosts {
			hostnames[j] = host.Hostname
		}
		c.lvConfigs.SetCellText2(i, 2, strings.Join(hostnames, ", "))
	}

	// Select the current config
	for i, cfg := range configs {
		if cfg.ID == currentConfigId {
			c.lvConfigs.SetCurrentCell2(i, 0)
			c.lvConfigs.ScrollToCell2(i, 0)
			break
		}
	}

	c.lvConfigs.Focus()

	return &c
}

func (c *OpenConfigDialog) GetSelectedConfig() *config.Config {
	row := c.lvConfigs.CurrentRow()
	if row < 0 {
		return nil
	}
	config := c.lvConfigs.GetCellData2(row, 0).(*config.Config)
	return config
}
