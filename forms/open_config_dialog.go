package forms

import (
	"fmt"
	"strings"

	"github.com/ipoluianov/altping/config"
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

type OpenConfigDialog struct {
	ui.DialogContent

	lvConfigs *ui.Table

	btnRemove *ui.Button
	btnOK     *ui.Button
	btnCancel *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel

	openedConfigId string

	onAccept func(configId string)
	onCancel func()
}

func NewOpenConfigDialog(openedConfigId string, onAccept func(configId string), onCancel func()) *OpenConfigDialog {
	var c OpenConfigDialog
	c.InitWidget()

	c.openedConfigId = openedConfigId

	c.panelContent = ui.NewPanel()
	c.AddWidget(0, 0, c.panelContent)
	c.panelButtons = ui.NewPanel()
	c.AddWidget(1, 0, c.panelButtons)

	c.btnRemove = ui.NewButton("Remove")
	c.btnRemove.SetOnClick(c.RemoveSelectedConfig)

	c.btnOK = ui.NewButton("OK")
	c.onAccept = onAccept
	c.onCancel = onCancel
	c.btnOK.SetOnClick(func() {
		c.Accept()
	})
	c.btnCancel = ui.NewButton("Cancel")
	c.btnCancel.SetOnClick(func() {
		c.Cancel()
	})

	c.panelButtons.AddWidget(0, 0, c.btnRemove)
	c.panelButtons.AddWidget(0, 1, ui.NewHSpacer())
	c.panelButtons.AddWidget(0, 2, c.btnOK)
	c.panelButtons.AddWidget(0, 3, c.btnCancel)

	c.lvConfigs = ui.NewTable()
	c.lvConfigs.SetSelectingCell(false)
	c.panelContent.AddWidget(0, 0, ui.NewLabel("Configs:"))
	c.panelContent.AddWidget(1, 0, c.lvConfigs)

	c.lvConfigs.SetColumnCount(3)
	c.lvConfigs.SetColumnWidth(0, 200)
	c.lvConfigs.SetColumnWidth(1, 100)
	c.lvConfigs.SetColumnWidth(2, 500)
	c.lvConfigs.SetColumnName(0, "Name")
	c.lvConfigs.SetColumnName(1, "Hosts Count")
	c.lvConfigs.SetColumnName(2, "Hosts")

	c.lvConfigs.SetOnCellMouseDblClick(c.onConfigDoubleClick)

	c.lvConfigs.SetOnKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool {
		if key == nuikey.KeyDelete {
			c.RemoveSelectedConfig()
			return true
		}
		return false
	})

	c.OnDialogShow = func() {
		c.Form().SetTitle("Open Config")
		c.lvConfigs.Focus()
		c.Form().SetAcceptButton(c.btnOK)
		c.Form().SetCancelButton(c.btnCancel)

		c.LoadTable(-1)
	}

	return &c
}

func (c *OpenConfigDialog) RemoveSelectedConfig() {
	ui.ShowQuestionMessageBoxOKCancel(c, "Remove config", "Remote selected config?", func() {
		selectedConfigIndex := c.lvConfigs.CurrentRow()
		selectedConfig := c.GetSelectedConfig()

		if selectedConfig.ID == c.openedConfigId {
			ui.ShowMessageBox(c, "Error", "Cannot remove the currently opened config.")
			return
		}

		if selectedConfig != nil {
			config.RemoveConfig(selectedConfig.ID)
			c.LoadTable(selectedConfigIndex)
			c.lvConfigs.Focus()
		}
	}, func() {
		c.lvConfigs.Focus()
	})
}

func (c *OpenConfigDialog) LoadTable(selectedConfigIndex int) {
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

		c.lvConfigs.SetCellColor(i, 2, ui.ColorFromHex("#555555"))
	}

	if selectedConfigIndex >= 0 && selectedConfigIndex < len(configs) {
		c.lvConfigs.SetCurrentCell2(selectedConfigIndex, 0)
		c.lvConfigs.ScrollToCell2(selectedConfigIndex, 0)
	} else {
		// Select the current config
		for i, cfg := range configs {
			if cfg.ID == c.openedConfigId {
				c.lvConfigs.SetCurrentCell2(i, 0)
				c.lvConfigs.ScrollToCell2(i, 0)
				break
			}
		}
	}
}

func (c *OpenConfigDialog) GetSelectedConfig() *config.Config {
	row := c.lvConfigs.CurrentRow()
	if row < 0 {
		return nil
	}
	config := c.lvConfigs.GetCellData2(row, 0).(*config.Config)
	return config
}

func (c *OpenConfigDialog) onConfigDoubleClick() {
	c.Accept()
}

func (c *OpenConfigDialog) Accept() {
	if c.onAccept != nil {
		selectedConfig := c.GetSelectedConfig()
		if selectedConfig != nil {
			c.onAccept(selectedConfig.ID)
		}
	}
	c.Form().Close()
}

func (c *OpenConfigDialog) Cancel() {
	if c.onCancel != nil {
		c.onCancel()
	}
	c.Form().Close()
}
