package forms

import (
	"fmt"
	"strings"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

type TopWidget struct {
	ui.Widget

	btnOpen *ui.ToolButton

	btnAddItem    *ui.ToolButton
	btnEditItem   *ui.ToolButton
	btnRemoveItem *ui.ToolButton

	btnDetails *ui.ToolButton
	btnStart   *ui.ToolButton
	btnStop    *ui.ToolButton
}

var lastCreatedTopWidget *TopWidget

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()
	c.SetPanelPadding(0)

	c.btnOpen = ui.NewToolButton(loadIcon("open"), "Configurations (O)", c.onBtnOpen)

	c.btnAddItem = ui.NewToolButton(loadIcon("add"), "Add host (A)", c.onBtnAddItem)
	// The main action: twice as wide as the other buttons
	c.btnAddItem.SetButtonSize(ui.ToolButtonDefaultSize*2, ui.ToolButtonDefaultSize)
	c.btnEditItem = ui.NewToolButton(loadIcon("edit"), "Edit host (E)", c.onBtnEditItem)
	c.btnRemoveItem = ui.NewToolButton(loadIcon("remove"), "Remove selected hosts (Del)", c.onBtnRemoveItem)

	c.btnDetails = ui.NewToolButton(loadIcon("details"), "Details (D)", c.onBtnDetails)
	c.btnStart = ui.NewToolButton(loadIcon("start"), "Start pinging", c.onBtnStart)
	c.btnStop = ui.NewToolButton(loadIcon("stop"), "Stop pinging", c.onBtnStop)

	c.AddWidget(0, 4, c.btnAddItem)
	c.AddWidget(0, 5, c.btnEditItem)
	c.AddWidget(0, 6, c.btnRemoveItem)

	c.AddWidget(0, 10, ui.NewHSpacer())

	c.AddWidget(0, 11, c.btnDetails)
	c.AddWidget(0, 12, c.btnStart)
	c.AddWidget(0, 13, c.btnStop)

	groupSpace := ui.NewSpace()
	groupSpace.SetSize(16, 0)
	c.AddWidget(0, 14, groupSpace)

	c.AddWidget(0, 15, c.btnOpen)

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

func (c *TopWidget) onBtnAddItem() {
	dialogContent := NewEditItemDialog(nil, func(hostConfig *config.ConfigHost) {
		if hostConfig != nil {
			hostConfig.ID = config.GenerateRandomID()
			config.Get().AddHost(*hostConfig)
			config.Get().Save()
			lastCreatedLeftWidget.ApplyHostsChange(len(config.Get().Hosts) - 1)
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
			lastCreatedLeftWidget.ApplyHostsChange(lastCreatedLeftWidget.RowOfHost(selectedHost.ID))
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

	ui.ShowQuestionMessageBoxYesNo(c, "Remove hosts", removeHostsQuestion(selectedHosts), func() {
		// After removing select the row that followed the first removed one
		firstRow := lastCreatedLeftWidget.RowOfHost(selectedHosts[0].ID)
		config := config.Get()
		for _, selectedHost := range selectedHosts {
			config.RemoveHost(selectedHost.ID)
		}
		config.Save()
		lastCreatedLeftWidget.ApplyHostsChange(firstRow)
		lastCreatedLeftWidget.FocusTable()
	}, func() {
		lastCreatedLeftWidget.FocusTable()
	})
}

func (c *TopWidget) onBtnDetails() {
	lastCreatedMainWidget.ToggleDetails()
	c.btnDetails.SetChecked(lastCreatedMainWidget.IsDetailsVisible())
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

// removeHostsListMax is how many hosts the remove question lists; the message box fits 20 lines
const removeHostsListMax = 10

// removeHostsQuestion asks to remove the hosts and lists them
func removeHostsQuestion(hosts []*config.ConfigHost) string {
	var sb strings.Builder
	if len(hosts) == 1 {
		sb.WriteString("Remove the host?\n")
	} else {
		fmt.Fprintf(&sb, "Remove %d hosts?\n", len(hosts))
	}
	for i, h := range hosts {
		if i == removeHostsListMax {
			fmt.Fprintf(&sb, "\n... and %d more", len(hosts)-removeHostsListMax)
			break
		}
		sb.WriteString("\n" + hostDisplayName(h))
		if h.DisplayName != "" && h.DisplayName != h.Hostname {
			sb.WriteString(" (" + h.Hostname + ")")
		}
	}
	return sb.String()
}
