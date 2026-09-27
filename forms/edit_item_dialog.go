package forms

import (
	"strconv"
	"strings"

	"github.com/ipoluianov/altping/config"
	"github.com/u00io/nuiforms/ui"
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

	c.panelContent = ui.NewPanel()
	c.AddWidget(0, 0, c.panelContent)
	c.AddWidget(1, 0, ui.NewVSpacer())
	c.panelButtons = ui.NewPanel()
	c.AddWidget(2, 0, c.panelButtons)

	c.btnOK = ui.NewButton("OK")
	c.onAccept = onAccept
	c.onCancel = onCancel
	c.btnOK.SetOnClick(c.Accept)
	c.btnCancel = ui.NewButton("Cancel")
	c.btnCancel.SetOnClick(c.Reject)

	c.panelButtons.AddWidget(0, 0, ui.NewHSpacer())
	c.panelButtons.AddWidget(0, 1, c.btnOK)
	c.panelButtons.AddWidget(0, 2, c.btnCancel)

	c.lblName = ui.NewLabel("Name:")
	c.txtName = ui.NewTextBox()
	c.txtName.SetText(c.hostConfig.DisplayName)
	c.panelContent.AddWidget(0, 0, c.lblName)
	c.panelContent.AddWidget(0, 1, c.txtName)

	// An old "example.com:443" in Hostname is shown as the host and the port
	host, port := c.hostConfig.Target()

	c.lblHost = ui.NewLabel("Host:")
	c.txtHost = ui.NewTextBox()
	c.txtHost.SetText(host)
	c.panelContent.AddWidget(1, 0, c.lblHost)
	c.panelContent.AddWidget(1, 1, c.txtHost)

	// Check a TCP port instead of ping: the checkbox enables the port field
	c.chkPort = ui.NewCheckbox("Check instead of ping")
	c.lblPort = ui.NewLabel("Port:")
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
	c.panelContent.AddWidget(2, 0, ui.NewLabel("TCP port:"))
	c.panelContent.AddWidget(2, 1, c.chkPort)
	c.panelContent.AddWidget(3, 0, c.lblPort)
	c.panelContent.AddWidget(3, 1, c.numPort)

	c.numInterval = newMsBox(c.hostConfig.Interval().Milliseconds(), config.MinIntervalMs, config.MaxIntervalMs)
	lblInterval := ui.NewLabel("Ping every, ms:")
	c.panelContent.AddWidget(4, 0, lblInterval)
	c.panelContent.AddWidget(4, 1, c.numInterval)

	c.numTimeout = newMsBox(c.hostConfig.Timeout().Milliseconds(), config.MinTimeoutMs, config.MaxTimeoutMs)
	c.panelContent.AddWidget(5, 0, ui.NewLabel("Timeout, ms:"))
	c.panelContent.AddWidget(5, 1, c.numTimeout)

	c.numSlow = newMsBox(int64(c.hostConfig.SlowMs), 0, config.MaxSlowMs)
	lblSlow := ui.NewLabel("Slow above, ms:")
	lblSlow.SetTooltip("Show the host in yellow when its average ping time is above it. 0 - off")
	c.numSlow.SetTooltip("Show the host in yellow when its average ping time is above it. 0 - off")
	c.panelContent.AddWidget(6, 0, lblSlow)
	c.panelContent.AddWidget(6, 1, c.numSlow)

	c.chkNotify = ui.NewCheckbox("Beep when down or back")
	c.chkNotify.SetChecked(c.hostConfig.Notify)
	c.chkNotify.SetTooltip("Down means 3 failed pings in a row. The window title shows how many such hosts are down.")
	c.panelContent.AddWidget(7, 0, ui.NewLabel("Notify:"))
	c.panelContent.AddWidget(7, 1, c.chkNotify)

	c.OnDialogShow = func() {
		c.Form().SetSize(480, 400)
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
