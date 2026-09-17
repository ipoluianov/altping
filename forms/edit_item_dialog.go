package forms

import (
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
	c.panelButtons = ui.NewPanel()
	c.AddWidget(1, 0, c.panelButtons)

	c.btnOK = ui.NewButton("OK")
	c.onAccept = onAccept
	c.onCancel = onCancel
	c.btnOK.SetOnClick(c.Accept)
	c.btnCancel = ui.NewButton("Cancel")
	c.btnCancel.SetOnClick(c.Reject)

	c.panelButtons.AddWidget(0, 0, c.btnOK)
	c.panelButtons.AddWidget(0, 1, c.btnCancel)

	c.lblName = ui.NewLabel("Name:")
	c.txtName = ui.NewTextBox()
	c.txtName.SetText(c.hostConfig.DisplayName)
	c.panelContent.AddWidget(0, 0, c.lblName)
	c.panelContent.AddWidget(0, 1, c.txtName)

	c.lblHost = ui.NewLabel("Host:")
	c.txtHost = ui.NewTextBox()
	c.txtHost.SetText(c.hostConfig.Hostname)
	c.panelContent.AddWidget(1, 0, c.lblHost)
	c.panelContent.AddWidget(1, 1, c.txtHost)

	c.OnDialogShow = func() {
		c.Form().SetSize(400, 200)
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

func (c *EditItemDialog) GetHostConfig() *config.ConfigHost {
	c.hostConfig.DisplayName = c.txtName.Text()
	c.hostConfig.Hostname = c.txtHost.Text()
	return &c.hostConfig
}

func (c *EditItemDialog) Accept() {
	if c.onAccept != nil {
		c.onAccept(c.GetHostConfig())
	}
	c.Form().Close()
}

func (c *EditItemDialog) Reject() {
	if c.onCancel != nil {
		c.onCancel()
	}
	c.Form().Close()
}
