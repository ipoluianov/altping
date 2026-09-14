package forms

import (
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

type TopWidget struct {
	ui.Widget

	btnNew    *ui.Button
	btnOpen   *ui.Button
	btnSaveAs *ui.Button

	btnAddItem    *ui.Button
	btnEditItem   *ui.Button
	btnRemoveItem *ui.Button

	btnDetails *ui.Button
	btnStart   *ui.Button
	btnStop    *ui.Button
}

var lastCreatedTopWidget *TopWidget

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()
	c.SetPanelPadding(6)

	c.btnNew = ui.NewButton("New")
	c.btnNew.SetOnClick(c.onBtnNew)
	c.btnOpen = ui.NewButton("Open")
	c.btnOpen.SetOnClick(c.onBtnOpen)
	c.btnSaveAs = ui.NewButton("Save As")
	c.btnSaveAs.SetOnClick(c.onBtnSaveAs)

	c.btnAddItem = ui.NewButton("Add")
	c.btnAddItem.SetOnClick(c.onBtnAddItem)
	c.btnEditItem = ui.NewButton("Edit")
	c.btnEditItem.SetOnClick(c.onBtnEditItem)
	c.btnRemoveItem = ui.NewButton("Remove")
	c.btnRemoveItem.SetOnClick(c.onBtnRemoveItem)

	c.btnDetails = ui.NewButton("Details")
	c.btnDetails.SetOnClick(c.onBtnDetails)
	c.btnStart = ui.NewButton("Start")
	c.btnStart.SetOnClick(c.onBtnStart)
	c.btnStop = ui.NewButton("Stop")
	c.btnStop.SetOnClick(c.onBtnStop)

	c.AddWidget(c.btnNew, 0, 0)
	c.AddWidget(c.btnOpen, 0, 1)
	c.AddWidget(c.btnSaveAs, 0, 2)

	c.AddWidget(c.btnAddItem, 0, 4)
	c.AddWidget(c.btnEditItem, 0, 5)
	c.AddWidget(c.btnRemoveItem, 0, 6)

	c.AddWidget(ui.NewHSpacer(), 0, 10)

	c.AddWidget(c.btnDetails, 0, 11)
	c.AddWidget(c.btnStart, 0, 12)
	c.AddWidget(c.btnStop, 0, 13)

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
	//dialog := ui.NewDialog("New Config", 480, 160)

	ShowCreateConfigDialog(c.Form(), func(configName string) {
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

	})
}

func (c *TopWidget) onBtnOpen() {
	currentConfigId := config.Get().ID

	dialog := ui.NewDialog("Open Config", 640, 480)
	dialogContent := NewOpenConfigDialog(currentConfigId, dialog.Accept, func() {
		dialog.Reject()
		lastCreatedLeftWidget.FocusTable()
	})
	dialog.ContentPanel().AddWidget(dialogContent, 0, 0)
	dialog.OnAccept = func() {
		selectedConfig := dialogContent.GetSelectedConfig()
		if selectedConfig != nil {
			system.Get().Stop()
			err := config.LoadConfig(selectedConfig.ID)
			if err != nil {
				ui.ShowMessageBox(c, "Error", err.Error())
				system.Get().Start()
				return
			}
			lastCreatedLeftWidget.FullRestart()
			system.Get().Start()

			lastCreatedLeftWidget.FocusTable()
		}
	}
	dialog.ShowDialog()
}

func (c *TopWidget) onBtnSaveAs() {
}

func (c *TopWidget) onBtnAddItem() {
	dialog := ui.NewDialog("Edit Item", 480, 160)
	dialogContent := NewEditItemDialog(nil, dialog.Accept, func() {
		dialog.Reject()
	})
	dialog.ContentPanel().AddWidget(dialogContent, 0, 0)
	dialog.OnAccept = func() {
		hostConfig := dialogContent.GetHostConfig()
		if hostConfig != nil {
			hostConfig.ID = config.GenerateRandomID()
			config.Get().AddHost(*hostConfig)
			config.Get().Save()
			lastCreatedLeftWidget.FullRestart()

			lastCreatedLeftWidget.FocusTable()
		}
	}
	dialog.ShowDialog()
}

func (c *TopWidget) onBtnEditItem() {
	selectedHost := lastCreatedLeftWidget.GetSelectedHostConfig()
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
	dialog.ShowDialog()
}

func (c *TopWidget) onBtnRemoveItem() {
	selectedHost := lastCreatedLeftWidget.GetSelectedHostConfig()
	if selectedHost == nil {
		return
	}

	ui.ShowQuestionMessageBoxYesNo(c, "Remove item?", "Remove Selected Item?", func() {
		config := config.Get()
		config.RemoveHost(selectedHost.ID)
		config.Save()
		lastCreatedLeftWidget.FullRestart()

		lastCreatedLeftWidget.FocusTable()
	}, func() {
	})
}

func (c *TopWidget) onBtnDetails() {
}

func (c *TopWidget) onBtnStart() {
	system.Get().Start()
}

func (c *TopWidget) onBtnStop() {
	system.Get().Stop()
}

func (c *TopWidget) onCreateConfigDialogAccept() {

}
