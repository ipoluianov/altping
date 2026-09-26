package forms

import (
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

type TopWidget struct {
	ui.Widget

	btnNew    *ToolButton
	btnOpen   *ToolButton
	btnSaveAs *ToolButton

	btnAddItem    *ToolButton
	btnEditItem   *ToolButton
	btnRemoveItem *ToolButton

	btnDetails *ToolButton
	btnStart   *ToolButton
	btnStop    *ToolButton
}

var lastCreatedTopWidget *TopWidget

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()
	c.SetPanelPadding(6)

	c.btnNew = NewToolButton("new", "New configuration (N)", c.onBtnNew)
	c.btnOpen = NewToolButton("open", "Open configuration (O)", c.onBtnOpen)
	c.btnSaveAs = NewToolButton("saveas", "Save configuration as (S)", c.onBtnSaveAs)

	c.btnAddItem = NewToolButton("add", "Add host (A)", c.onBtnAddItem)
	c.btnEditItem = NewToolButton("edit", "Edit host (E)", c.onBtnEditItem)
	c.btnRemoveItem = NewToolButton("remove", "Remove selected hosts (Del)", c.onBtnRemoveItem)

	c.btnDetails = NewToolButton("details", "Details (D)", c.onBtnDetails)
	c.btnStart = NewToolButton("start", "Start pinging", c.onBtnStart)
	c.btnStop = NewToolButton("stop", "Stop pinging", c.onBtnStop)

	c.AddWidget(0, 0, c.btnNew)
	c.AddWidget(0, 1, c.btnOpen)
	c.AddWidget(0, 2, c.btnSaveAs)

	groupSpace := ui.NewSpace()
	groupSpace.SetSize(16, 0)
	c.AddWidget(0, 3, groupSpace)

	c.AddWidget(0, 4, c.btnAddItem)
	c.AddWidget(0, 5, c.btnEditItem)
	c.AddWidget(0, 6, c.btnRemoveItem)

	c.AddWidget(0, 10, ui.NewHSpacer())

	c.AddWidget(0, 11, c.btnDetails)
	c.AddWidget(0, 12, c.btnStart)
	c.AddWidget(0, 13, c.btnStop)

	c.AddTimer(200, c.timerUpdate)

	lastCreatedTopWidget = &c

	return &c
}

func (c *TopWidget) timerUpdate() {
	if system.Get().IsRunning() {
		c.btnStart.SetEnabled(false)
		c.btnStop.SetEnabled(true)
	} else {
		c.btnStart.SetEnabled(true)
		c.btnStop.SetEnabled(false)
	}
}

func (c *TopWidget) onBtnNew() {
	dialog := NewCreateConfigDialog(func(configName string) {
		if configName != "" {
			cfg, err := config.CreateNewConfig(configName)
			if err != nil {
				ui.ShowMessageBox(c, "Error", err.Error())
				return
			}
			system.Get().Stop()
			err = config.LoadConfig(cfg.ID)
			if err != nil {
				ui.ShowMessageBox(c, "Error", err.Error())
				system.Get().Start()
				return
			}
			lastCreatedLeftWidget.FullRestart()
			system.Get().Start()

			lastCreatedLeftWidget.FocusTable()

		}
	}, func() {
		// onCancel callback
	})

	c.ShowDialog(dialog)
}

func (c *TopWidget) onBtnOpen() {
	currentConfigId := config.Get().ID

	dialogContent := NewOpenConfigDialog(currentConfigId, func(selectedConfigId string) {
		if selectedConfigId != "" {
			system.Get().Stop()
			err := config.LoadConfig(selectedConfigId)
			if err != nil {
				ui.ShowMessageBox(c, "Error", err.Error())
				system.Get().Start()
				return
			}
			lastCreatedLeftWidget.FullRestart()
			system.Get().Start()

			lastCreatedLeftWidget.FocusTable()
		}
	}, func() {
		lastCreatedLeftWidget.FocusTable()
	})

	c.ShowDialog(dialogContent)
}

func (c *TopWidget) onBtnSaveAs() {
}

func (c *TopWidget) onBtnAddItem() {
	dialogContent := NewEditItemDialog(nil, func(hostConfig *config.ConfigHost) {
		if hostConfig != nil {
			hostConfig.ID = config.GenerateRandomID()
			config.Get().AddHost(*hostConfig)
			config.Get().Save()
			lastCreatedLeftWidget.FullRestart()
			lastCreatedLeftWidget.FocusTable()
		}
	}, func() {
		lastCreatedLeftWidget.FocusTable()
	})

	c.ShowDialog(dialogContent)
}

func (c *TopWidget) onBtnEditItem() {
	selectedHosts := lastCreatedLeftWidget.GetSelectedHostConfigs()
	if len(selectedHosts) != 1 {
		return
	}
	selectedHost := selectedHosts[0]
	dialogContent := NewEditItemDialog(selectedHost, func(hostConfig *config.ConfigHost) {
		if hostConfig != nil {
			selectedHost.DisplayName = hostConfig.DisplayName
			selectedHost.Hostname = hostConfig.Hostname
			config.Get().Save()
			lastCreatedLeftWidget.FullRestart()
			lastCreatedLeftWidget.FocusTable()
		}
	}, func() {
		lastCreatedLeftWidget.FocusTable()
	})

	c.ShowDialog(dialogContent)

	/*selectedHost := lastCreatedLeftWidget.GetSelectedHostConfig()
	if selectedHost == nil {
		return
	}

	dialog := ui.NewDialog("Edit Item", 480, 160)
	dialogContent := NewEditItemDialog(selectedHost, dialog.Accept, func() {
		dialog.Reject()
	})
	dialog.ContentPanel().AddWidget(dialogContent, 0, 0)
	dialog.OnAccept = func() {
		hostConfig := dialogContent.GetHostConfig()
		if hostConfig != nil {
			selectedHost.DisplayName = hostConfig.DisplayName
			selectedHost.Hostname = hostConfig.Hostname
			config.Get().Save()
			lastCreatedLeftWidget.FullRestart()

			lastCreatedLeftWidget.FocusTable()
		}
	}
	dialog.ShowDialog()*/
}

func (c *TopWidget) onBtnRemoveItem() {
	selectedHosts := lastCreatedLeftWidget.GetSelectedHostConfigs()
	if len(selectedHosts) == 0 {
		return
	}

	ui.ShowQuestionMessageBoxYesNo(c, "Remove item?", "Remove Selected Items?", func() {
		config := config.Get()
		for _, selectedHost := range selectedHosts {
			config.RemoveHost(selectedHost.ID)
		}
		config.Save()
		lastCreatedLeftWidget.FullRestart()
		lastCreatedLeftWidget.FocusTable()
	}, func() {
		lastCreatedLeftWidget.FocusTable()
	})
}

func (c *TopWidget) onBtnDetails() {
	lastCreatedMainWidget.ToggleDetails()
	lastCreatedLeftWidget.FocusTable()
}

func (c *TopWidget) onBtnStart() {
	system.Get().Start()
}

func (c *TopWidget) onBtnStop() {
	system.Get().Stop()
}

func (c *TopWidget) onCreateConfigDialogAccept() {

}
