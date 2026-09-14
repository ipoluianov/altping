package forms

import (
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

type CreateConfigDialog struct {
	ui.Widget

	txtName *ui.TextBox

	btnOK     *ui.Button
	btnCancel *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel

	onAccept func(name string)
	onCancel func()
}

func ShowCreateConfigDialog(parent *ui.Form, onAccept func(name string), onCancel func()) {
	form := ui.NewForm()
	var c CreateConfigDialog
	c.InitWidget()

	c.panelContent = ui.NewPanel()
	c.AddWidget(c.panelContent, 0, 0)
	c.AddWidget(ui.NewVSpacer(), 1, 0)
	c.panelButtons = ui.NewPanel()
	c.AddWidget(c.panelButtons, 2, 0)

	c.btnOK = ui.NewButton("OK")
	c.onAccept = onAccept
	c.onCancel = onCancel
	c.btnOK.SetOnClick(func() {
		if c.onAccept != nil {
			c.onAccept(c.txtName.Text())
		}
		form.Close()
	})
	c.btnCancel = ui.NewButton("Cancel")
	c.btnCancel.SetOnClick(func() {
		if c.onCancel != nil {
			c.onCancel()
		}
		form.Close()
	})

	c.panelButtons.AddWidget(c.btnOK, 0, 0)
	c.panelButtons.AddWidget(c.btnCancel, 0, 1)

	c.txtName = ui.NewTextBox()
	c.panelContent.AddWidget(ui.NewLabel("Name:"), 0, 0)
	c.panelContent.AddWidget(c.txtName, 0, 1)

	c.txtName.Focus()
	c.txtName.SetOnKeyDown(func(key nuikey.Key, mods nuikey.KeyModifiers) bool {
		if key == nuikey.KeyEnter {
			if c.onAccept != nil {
				c.onAccept(c.txtName.Text())
			}
			form.Close()
			return true
		}
		if key == nuikey.KeyEsc {
			if c.onCancel != nil {
				c.onCancel()
			}
			form.Close()
			return true
		}
		return c.txtName.KeyDown(key, mods)
	})

	form.Panel().AddWidget(&c, 0, 0)
	form.ShowModal(parent)
}
