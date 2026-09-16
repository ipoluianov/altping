package forms

import (
	"github.com/ipoluianov/altping/config"
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

type EditItemDialog struct {
	ui.Widget

	hostConfig config.ConfigHost

	lblName *ui.Label
	txtName *ui.TextBox

	lblHost *ui.Label
	txtHost *ui.TextBox

	btnOK     *ui.Button
	btnCancel *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel

	onAccept func()
	onCancel func()
}

func NewEditItemDialog(hostConfig *config.ConfigHost, onAccept func(), onCancel func()) *EditItemDialog {
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

	if c.hostConfig.ID == "" {
		c.txtHost.SetText("127.0.0.1")
		c.txtHost.MoveCursorToEnd()
		c.txtHost.SelectAllText()
		c.txtHost.Focus()
	}

	c.txtName.SetOnKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool {
		if key == nuikey.KeyEnter {
			if c.onAccept != nil {
				c.onAccept()
			}
			return true
		}
		if key == nuikey.KeyEsc {
			if c.onCancel != nil {
				c.onCancel()
			}
			return true
		}
		return c.txtName.KeyDown(key, mods)
	})

	c.txtHost.SetOnKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool {
		if key == nuikey.KeyEnter {
			if c.onAccept != nil {
				c.onAccept()
			}
			return true
		}
		if key == nuikey.KeyEsc {
			if c.onCancel != nil {
				c.onCancel()
			}
			return true
		}
		return c.txtHost.KeyDown(key, mods)
	})

	c.txtHost.Focus()

	return &c
}

func (c *EditItemDialog) GetHostConfig() *config.ConfigHost {
	c.hostConfig.DisplayName = c.txtName.Text()
	c.hostConfig.Hostname = c.txtHost.Text()
	return &c.hostConfig
}
