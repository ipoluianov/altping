package forms

import (
	"strconv"
	"strings"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/nui/ui"
)

type EditItemDialog struct {
	ui.DialogContent

	hostConfig config.ConfigHost

	lblName *ui.Label
	txtName *ui.TextBox

	lblHost *ui.Label
	txtHost *ui.TextBox

	chkPort *ui.Checkbox
	lblPort *ui.Label
	numPort *ui.NumBox

	numInterval *ui.NumBox
	numTimeout  *ui.NumBox
	numSlow     *ui.NumBox
	chkNotify   *ui.Checkbox

	chkShare     *ui.Checkbox
	btnShareOpen *ui.Button
	btnShareCopy *ui.Button
	lblShareURL  *ui.Label

	btnOK     *ui.Button
	btnCancel *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel

	onAccept func(*config.ConfigHost)
	onCancel func()
}

func NewEditItemDialog(hostConfig *config.ConfigHost, onAccept func(hostConfig *config.ConfigHost), onCancel func()) *EditItemDialog {
	var c EditItemDialog
	c.InitWidget()

	if hostConfig != nil {
		c.hostConfig = *hostConfig
	}
	// A new host gets its key now, so its link can be copied before it is saved
	if config.ShareAddress(c.hostConfig.ShareKey) == "" {
		c.hostConfig.ShareKey = config.NewShareKey()
	}

	c.panelContent = ui.NewPanel()
	c.AddWidget(0, 0, c.panelContent)
	c.AddWidget(1, 0, ui.NewVSpacer())
	c.panelButtons = ui.NewPanel()
	c.AddWidget(2, 0, c.panelButtons)

	c.btnOK = ui.NewButton(ui.UIText().OK)
	c.onAccept = onAccept
	c.onCancel = onCancel
	c.btnOK.SetOnClick(c.Accept)
	c.btnCancel = ui.NewButton(ui.UIText().Cancel)
	c.btnCancel.SetOnClick(c.Reject)

	c.panelButtons.AddWidget(0, 0, ui.NewHSpacer())
	c.panelButtons.AddWidget(0, 1, c.btnOK)
	c.panelButtons.AddWidget(0, 2, c.btnCancel)

	c.lblName = ui.NewLabel(T().Name)
	c.txtName = ui.NewTextBox()
	c.txtName.SetText(c.hostConfig.DisplayName)
	c.panelContent.AddWidget(0, 0, c.lblName)
	c.panelContent.AddWidget(0, 1, c.txtName)

	// An old "example.com:443" in Hostname is shown as the host and the port
	host, port := c.hostConfig.Target()

	c.lblHost = ui.NewLabel(T().Host)
	c.txtHost = ui.NewTextBox()
	c.txtHost.SetText(host)
	c.panelContent.AddWidget(1, 0, c.lblHost)
	c.panelContent.AddWidget(1, 1, c.txtHost)

	// Check a TCP port instead of ping: the checkbox enables the port field
	c.chkPort = ui.NewCheckbox(T().CheckInstead)
	c.lblPort = ui.NewLabel(T().Port)
	c.numPort = ui.NewNumBox()
	c.numPort.SetDecimals(0)
	c.numPort.SetMin(1)
	c.numPort.SetMax(65535)
	c.numPort.SetStep(1)
	c.numPort.SetValue(defaultTCPPort)
	if p, err := strconv.Atoi(port); err == nil {
		c.numPort.SetValue(float64(p))
		c.chkPort.SetChecked(true)
	}
	c.showPortRow()
	c.chkPort.SetOnStateChanged(c.showPortRow)
	c.panelContent.AddWidget(2, 0, ui.NewLabel(T().TCPPort))
	c.panelContent.AddWidget(2, 1, c.chkPort)
	c.panelContent.AddWidget(3, 0, c.lblPort)
	c.panelContent.AddWidget(3, 1, c.numPort)

	c.numInterval = newMsBox(c.hostConfig.Interval().Milliseconds(), config.MinIntervalMs, config.MaxIntervalMs)
	lblInterval := ui.NewLabel(T().PingEveryMs)
	c.panelContent.AddWidget(4, 0, lblInterval)
	c.panelContent.AddWidget(4, 1, c.numInterval)

	c.numTimeout = newMsBox(c.hostConfig.Timeout().Milliseconds(), config.MinTimeoutMs, config.MaxTimeoutMs)
	c.panelContent.AddWidget(5, 0, ui.NewLabel(T().TimeoutMs))
	c.panelContent.AddWidget(5, 1, c.numTimeout)

	c.numSlow = newMsBox(int64(c.hostConfig.SlowMs), 0, config.MaxSlowMs)
	lblSlow := ui.NewLabel(T().SlowAboveMs)
	lblSlow.SetTooltip(T().SlowTooltip)
	c.numSlow.SetTooltip(T().SlowTooltip)
	c.panelContent.AddWidget(6, 0, lblSlow)
	c.panelContent.AddWidget(6, 1, c.numSlow)

	c.chkNotify = ui.NewCheckbox(T().BeepDownBack)
	c.chkNotify.SetChecked(c.hostConfig.Notify)
	c.chkNotify.SetTooltip(T().NotifyTooltip)
	c.panelContent.AddWidget(7, 0, ui.NewLabel(T().Notify))
	c.panelContent.AddWidget(7, 1, c.chkNotify)

	// Sharing on u00.io: the page and its link are there before sharing is
	// turned on, the buttons work once it is on
	c.chkShare = ui.NewCheckbox(T().ShareOn)
	c.chkShare.SetChecked(c.hostConfig.Share)
	c.chkShare.SetTooltip(T().ShareTooltip)
	c.chkShare.SetOnStateChanged(c.showShareRow)
	c.panelContent.AddWidget(8, 0, ui.NewLabel(T().Share))
	c.panelContent.AddWidget(8, 1, c.chkShare)

	shareRow := ui.NewPanel()
	shareRow.SetPanelPadding(0)
	c.btnShareOpen = ui.NewButton(T().ShareOpen)
	c.btnShareOpen.SetOnClick(func() { openShareURL(&c, c.hostConfig) })
	c.btnShareCopy = ui.NewButton(T().ShareCopy)
	c.btnShareCopy.SetOnClick(func() { copyShareURL(&c, c.hostConfig) })
	c.lblShareURL = ui.NewLabel(shortShareURL(c.hostConfig))
	c.lblShareURL.SetTooltip(c.hostConfig.ShareURL())
	shareRow.AddWidget(0, 0, c.btnShareOpen)
	shareRow.AddWidget(0, 1, c.btnShareCopy)
	shareRow.AddWidget(0, 2, c.lblShareURL)
	shareRow.AddWidget(0, 3, ui.NewHSpacer())
	c.panelContent.AddWidget(9, 1, shareRow)
	c.showShareRow()

	c.OnDialogShow = func() {
		c.Form().SetSize(520, 470)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnOK)
		c.Form().SetCancelButton(c.btnCancel)

		if c.hostConfig.ID == "" {
			c.txtHost.SetText("127.0.0.1")
			c.txtHost.MoveCursorToEnd()
			c.txtHost.SelectAllText()
		}
		c.txtHost.Focus()

	}

	return &c
}

// showShareRow enables the buttons of the page while sharing is on
func (c *EditItemDialog) showShareRow() {
	on := c.chkShare.Checked()
	c.btnShareOpen.SetEnabled(on)
	c.btnShareCopy.SetEnabled(on)
	if on {
		c.lblShareURL.SetForegroundColor(colorMuted.get())
	} else {
		c.lblShareURL.SetForegroundColor(ui.ThemeForegroundColorDisabled())
	}
}

// showPortRow grays out the port field while the port check is off
func (c *EditItemDialog) showPortRow() {
	on := c.chkPort.Checked()
	c.numPort.SetEnabled(on)
	if on {
		c.lblPort.SetForegroundColor(nil) // the theme color
	} else {
		c.lblPort.SetForegroundColor(ui.ThemeForegroundColorDisabled())
	}
}

func (c *EditItemDialog) GetHostConfig() *config.ConfigHost {
	c.hostConfig.DisplayName = c.txtName.Text()
	c.hostConfig.Hostname = strings.TrimSpace(c.txtHost.Text())
	c.hostConfig.Port = 0
	if c.chkPort.Checked() {
		c.hostConfig.Port = int(c.numPort.Value())
	}
	c.hostConfig.IntervalMs = int(c.numInterval.Value())
	if c.hostConfig.IntervalMs == config.DefaultIntervalMs {
		c.hostConfig.IntervalMs = 0 // the default, not stored
	}
	c.hostConfig.TimeoutMs = int(c.numTimeout.Value())
	if c.hostConfig.TimeoutMs == config.DefaultTimeoutMs {
		c.hostConfig.TimeoutMs = 0 // the default, not stored
	}
	c.hostConfig.SlowMs = int(c.numSlow.Value())
	c.hostConfig.Notify = c.chkNotify.Checked()
	c.hostConfig.Share = c.chkShare.Checked()
	return &c.hostConfig
}

func (c *EditItemDialog) Accept() {
	// The callbacks change the window below: they run on its goroutine
	if c.onAccept != nil {
		hostConfig := c.GetHostConfig()
		c.RunInParent(func() { c.onAccept(hostConfig) })
	}
	c.Form().Close()
}

func (c *EditItemDialog) Reject() {
	c.RunInParent(c.onCancel)
	c.Form().Close()
}

// newMsBox is a field for whole milliseconds
func newMsBox(value int64, minValue, maxValue int) *ui.NumBox {
	num := ui.NewNumBox()
	num.SetDecimals(0)
	num.SetMin(float64(minValue))
	num.SetMax(float64(maxValue))
	num.SetStep(100)
	num.SetValue(float64(value))
	return num
}

// defaultTCPPort is offered when the port check is turned on: HTTPS, the most common one
const defaultTCPPort = 443
