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

	btnNew    *ui.Button
	btnSaveAs *ui.Button
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

	c.btnNew = ui.NewButton("New...")
	c.btnNew.SetOnClick(c.CreateConfig)
	c.btnSaveAs = ui.NewButton("Save As...")
	c.btnSaveAs.SetOnClick(c.SaveOpenedConfigAs)
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

	c.panelButtons.AddWidget(0, 0, c.btnNew)
	c.panelButtons.AddWidget(0, 1, c.btnSaveAs)
	c.panelButtons.AddWidget(0, 2, c.btnRemove)
	c.panelButtons.AddWidget(0, 3, ui.NewHSpacer())
	c.panelButtons.AddWidget(0, 4, c.btnOK)
	c.panelButtons.AddWidget(0, 5, c.btnCancel)

	c.lvConfigs = ui.NewTable()
	c.lvConfigs.SetSelectingRows(true)
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

		c.LoadTableSelectID(c.openedConfigId)
	}

	return &c
}

// CreateConfig creates an empty config and selects it in the table
func (c *OpenConfigDialog) CreateConfig() {
	c.ShowDialog(NewCreateConfigDialog("New Config", "", func(name string) {
		if name == "" {
			return
		}
		cfg, err := config.CreateNewConfig(name)
		if err != nil {
			ui.ShowMessageBox(c, "Error", err.Error())
			return
		}
		c.LoadTableSelectID(cfg.ID)
		c.lvConfigs.Focus()
	}, func() {
		c.lvConfigs.Focus()
	}))
}

// SaveOpenedConfigAs saves a copy of the opened config under a new name and selects it in the table
func (c *OpenConfigDialog) SaveOpenedConfigAs() {
	opened := config.Get()
	c.ShowDialog(NewCreateConfigDialog("Save Config As", opened.Name+" copy", func(name string) {
		if name == "" {
			return
		}
		cfg, err := config.CopyConfig(opened, name)
		if err != nil {
			ui.ShowMessageBox(c, "Error", err.Error())
			return
		}
		c.LoadTableSelectID(cfg.ID)
		c.lvConfigs.Focus()
	}, func() {
		c.lvConfigs.Focus()
	}))
}

func (c *OpenConfigDialog) RemoveSelectedConfig() {
	ui.ShowQuestionMessageBoxOKCancel(c, "Remove config", "Remove the selected config?", func() {
		selectedConfigIndex := c.lvConfigs.CurrentRow()
		selectedConfig := c.GetSelectedConfig()

		if selectedConfig != nil && selectedConfig.ID == c.openedConfigId {
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

	// After removing the last row select the new last one
	selectedConfigIndex = min(selectedConfigIndex, len(configs)-1)
	if selectedConfigIndex >= 0 {
		c.lvConfigs.SetCurrentCell2(selectedConfigIndex, 0)
		c.lvConfigs.ScrollToCell2(selectedConfigIndex, 0)
	}
}

// LoadTableSelectID fills the table and selects the config with the ID
func (c *OpenConfigDialog) LoadTableSelectID(id string) {
	index := -1
	for i, cfg := range config.Configs() {
		if cfg.ID == id {
			index = i
			break
		}
	}
	c.LoadTable(index)
}

func (c *OpenConfigDialog) GetSelectedConfig() *config.Config {
	cfg, _ := c.lvConfigs.GetCellData2(c.lvConfigs.CurrentRow(), 0).(*config.Config)
	return cfg
}

func (c *OpenConfigDialog) onConfigDoubleClick() {
	c.Accept()
}

func (c *OpenConfigDialog) Accept() {
	// The callbacks change the window below: they run on its goroutine
	if selectedConfig := c.GetSelectedConfig(); selectedConfig != nil && c.onAccept != nil {
		id := selectedConfig.ID
		c.RunInParent(func() { c.onAccept(id) })
	}
	c.Form().Close()
}

func (c *OpenConfigDialog) Cancel() {
	c.RunInParent(c.onCancel)
	c.Form().Close()
}
