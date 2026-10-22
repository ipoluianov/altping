package forms

import (
	"strings"
	"unicode"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
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

	btnTray *ui.ToolButton

	btnShareOpen *ui.Button
}

var lastCreatedTopWidget *TopWidget

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()
	c.SetPanelPadding(0)

	c.btnOpen = ui.NewToolButton(nil, "", c.onBtnOpen)
	setIcon("open", c.btnOpen.SetImage)
	c.btnOpen.SetTooltipFunc(func() string { return T().ToolConfigs })

	c.btnAddItem = ui.NewToolButton(nil, "", c.onBtnAddItem)
	setIcon("add", c.btnAddItem.SetImage)
	c.btnAddItem.SetTooltipFunc(func() string { return T().ToolAddHost })
	c.btnEditItem = ui.NewToolButton(nil, "", c.onBtnEditItem)
	setIcon("edit", c.btnEditItem.SetImage)
	c.btnEditItem.SetTooltipFunc(func() string { return T().ToolEditHost })
	c.btnRemoveItem = ui.NewToolButton(nil, "", c.onBtnRemoveItem)
	setIcon("remove", c.btnRemoveItem.SetImage)
	c.btnRemoveItem.SetTooltipFunc(func() string { return T().ToolRemoveHosts })

	c.btnDetails = ui.NewToolButton(nil, "", c.onBtnDetails)
	setIcon("details", c.btnDetails.SetImage)
	c.btnDetails.SetTooltipFunc(func() string { return T().ToolDetails })
	c.btnStart = ui.NewToolButton(nil, "", c.onBtnStart)
	setIcon("start", c.btnStart.SetImage)
	c.btnStart.SetTooltipFunc(func() string { return T().ToolStart })
	c.btnStop = ui.NewToolButton(nil, "", c.onBtnStop)
	setIcon("stop", c.btnStop.SetImage)
	c.btnStop.SetTooltipFunc(func() string { return T().ToolStop })

	// The u00.io page stands out among the icons: a button with its name, in the middle
	c.btnShareOpen = ui.NewButton("")
	c.btnShareOpen.SetTextFunc(func() string { return T().MenuShareOpen })
	c.btnShareOpen.SetOnClick(c.onBtnShareOpen)
	c.btnShareOpen.SetEnabled(false) // until a shared host is selected
	// Like the icons, it leaves the keys to the table: U opens the pages
	c.btnShareOpen.SetCanBeFocused(false)
	c.btnShareOpen.SetTooltipFunc(func() string { return T().ToolShareOpen })

	c.btnTray = ui.NewToolButton(nil, "", c.onBtnTray)
	setIcon("tray", c.btnTray.SetImage)
	c.btnTray.SetTooltipFunc(func() string { return T().ToolTray })

	// Flat icons: the toolbar stays light, a button shows its shape under the mouse
	for _, btn := range []*ui.ToolButton{c.btnAddItem, c.btnEditItem, c.btnRemoveItem,
		c.btnDetails, c.btnStart, c.btnStop, c.btnOpen, c.btnTray} {
		btn.SetFlat(true)
	}

	left := ui.NewPanel()
	left.SetPanelPadding(0)
	left.AddWidget(0, 0, c.btnAddItem)
	left.AddWidget(0, 1, c.btnEditItem)
	left.AddWidget(0, 2, c.btnRemoveItem)
	left.AddWidget(0, 3, ui.NewHSpacer())

	right := ui.NewPanel()
	right.SetPanelPadding(0)
	right.AddWidget(0, 0, ui.NewHSpacer())
	right.AddWidget(0, 1, c.btnDetails)
	right.AddWidget(0, 2, c.btnStart)
	right.AddWidget(0, 3, c.btnStop)
	groupSpace := ui.NewSpace()
	groupSpace.SetSize(16, 0)
	right.AddWidget(0, 4, groupSpace)
	right.AddWidget(0, 5, c.btnOpen)
	right.AddWidget(0, 6, c.btnTray)

	// The sides are as wide as each other, so the button between them is in the middle
	sideWidth := max(left.MinWidth(), right.MinWidth())
	left.SetMinWidth(sideWidth)
	right.SetMinWidth(sideWidth)
	left.ClearLayoutCache()
	right.ClearLayoutCache()

	c.AddWidget(0, 0, left)
	// The button keeps its own height, in the middle of the toolbar's
	// (spacers would let the toolbar grow)
	middle := ui.NewPanel()
	middle.SetPanelPadding(0)
	middle.SetCellPadding(0)
	top := ui.NewSpace()
	c.btnShareOpen.SetMaxHeight(c.btnShareOpen.MinHeight())
	top.SetSize(0, (ui.ToolButtonDefaultSize-c.btnShareOpen.MinHeight())/2)
	middle.AddWidget(0, 0, top)
	middle.AddWidget(1, 0, c.btnShareOpen)
	c.AddWidget(0, 1, middle)
	c.AddWidget(0, 2, right)

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

	// Only one host is edited at a time; pages are opened for the shared hosts
	if lastCreatedLeftWidget != nil {
		c.btnEditItem.SetEnabled(len(lastCreatedLeftWidget.GetSelectedHostConfigs()) == 1)
		shareOpen := len(sharedHosts(lastCreatedLeftWidget.GetSelectedHostConfigs())) > 0
		// Set on a change only: a property of a Button lays the window out again
		if c.btnShareOpen.Enabled() != shareOpen {
			c.btnShareOpen.SetEnabled(shareOpen)
		}
	}

	// An empty config: point to the first thing to do
	if len(config.Get().Hosts) == 0 {
		c.btnAddItem.SetHighlight(addHighlightColor)
	} else {
		c.btnAddItem.SetHighlight(nil)
	}
}

var addHighlightColor = ui.ColorFromHex("#3fb950")

func (c *TopWidget) onBtnOpen() {
	currentConfigId := config.Get().ID

	dialogContent := NewOpenConfigDialog(currentConfigId, func(selectedConfigId string) {
		if selectedConfigId != "" {
			system.Get().Stop()
			err := config.LoadConfig(selectedConfigId)
			if err != nil {
				ui.ShowMessageBox(c, T().Error, err.Error())
				system.Get().Start()
				return
			}
			lastCreatedLeftWidget.FullRestart()
			system.Get().Start()
			lastCreatedMainWidget.UpdateTitle()

			lastCreatedLeftWidget.FocusTable()
		}
	}, func() {
		lastCreatedLeftWidget.FocusTable()
	})

	c.ShowDialog(dialogContent)
}

func (c *TopWidget) onBtnAddItem() {
	dialogContent := NewEditItemDialog(nil, func(hostConfig *config.ConfigHost) {
		if hostConfig == nil {
			return
		}
		// Several hosts can be pasted at once, separated by commas, semicolons or spaces;
		// then the name is not used, each host is shown by its address
		addresses := strings.FieldsFunc(hostConfig.Hostname, func(r rune) bool {
			return r == ',' || r == ';' || unicode.IsSpace(r)
		})
		if len(addresses) == 0 {
			return
		}
		firstID := ""
		for _, address := range addresses {
			host := config.ConfigHost{
				ID:         config.GenerateRandomID(),
				Hostname:   address,
				Port:       hostConfig.Port,
				IntervalMs: hostConfig.IntervalMs,
				TimeoutMs:  hostConfig.TimeoutMs,
				SlowMs:     hostConfig.SlowMs,
				Notify:     hostConfig.Notify,
				Share:      hostConfig.Share,
			}
			if len(addresses) == 1 {
				host.DisplayName = hostConfig.DisplayName
			}
			// The first host gets the key the dialog showed, the others new ones (AddHost)
			if firstID == "" {
				host.ShareKey = hostConfig.ShareKey
			}
			if firstID == "" {
				firstID = host.ID
			}
			config.Get().AddHost(host)
		}
		config.Get().Save()
		lastCreatedLeftWidget.ApplyHostsChange(firstID, -1)
		lastCreatedLeftWidget.FocusTable()
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
			selectedHost.Port = hostConfig.Port
			selectedHost.IntervalMs = hostConfig.IntervalMs
			selectedHost.TimeoutMs = hostConfig.TimeoutMs
			selectedHost.SlowMs = hostConfig.SlowMs
			selectedHost.Notify = hostConfig.Notify
			selectedHost.Share = hostConfig.Share
			selectedHost.ShareKey = hostConfig.ShareKey
			config.Get().Save()
			lastCreatedLeftWidget.ApplyHostsChange(selectedHost.ID, -1)
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

	ui.ShowQuestionMessageBoxYesNo(c, T().RemoveHostsTitle, removeHostsQuestion(selectedHosts), func() {
		// After removing select the row that followed the first removed one
		firstRow := lastCreatedLeftWidget.RowOfHost(selectedHosts[0].ID)
		config := config.Get()
		for _, selectedHost := range selectedHosts {
			config.RemoveHost(selectedHost.ID)
		}
		config.Save()
		lastCreatedLeftWidget.ApplyHostsChange("", firstRow)
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

// onBtnShareOpen opens the u00.io pages of the selected hosts that are shared
func (c *TopWidget) onBtnShareOpen() {
	hosts := sharedHosts(lastCreatedLeftWidget.GetSelectedHostConfigs())
	for i, host := range hosts {
		if i == shareOpenMax {
			break
		}
		openShareURL(c, *host)
	}
}

func (c *TopWidget) onBtnTray() {
	lastCreatedMainWidget.HideToTray()
}

func (c *TopWidget) onCreateConfigDialogAccept() {

}

// removeHostsListMax is how many hosts the remove question lists; the message box fits 20 lines
const removeHostsListMax = 10

// removeHostsQuestion asks to remove the hosts and lists them
func removeHostsQuestion(hosts []*config.ConfigHost) string {
	var sb strings.Builder
	if len(hosts) == 1 {
		sb.WriteString(T().RemoveHost + "\n")
	} else {
		sb.WriteString(T().RemoveHosts(len(hosts)) + "\n")
	}
	for i, h := range hosts {
		if i == removeHostsListMax {
			sb.WriteString("\n" + T().AndMore(len(hosts)-removeHostsListMax))
			break
		}
		sb.WriteString("\n" + hostDisplayName(h))
		if h.DisplayName != "" && h.DisplayName != h.Address() {
			sb.WriteString(" (" + h.Address() + ")")
		}
	}
	return sb.String()
}
