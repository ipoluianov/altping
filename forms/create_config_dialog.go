package forms

import (
	"github.com/u00io/nuiforms/ui"
)

type CreateConfigDialog struct {
	ui.DialogContent

	txtName *ui.TextBox

	btnOK     *ui.Button
	btnCancel *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel

	onAccept func(name string)
	onCancel func()
}

func NewCreateConfigDialog(onAccept func(name string), onCancel func()) *CreateConfigDialog {
	var c CreateConfigDialog
	c.InitWidget()

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
	c.btnCancel.SetOnClick(c.Cancel)

	c.panelButtons.AddWidget(0, 0, ui.NewHSpacer())
	c.panelButtons.AddWidget(0, 1, c.btnOK)
	c.panelButtons.AddWidget(0, 2, c.btnCancel)

	c.txtName = ui.NewTextBox()
	c.panelContent.AddWidget(0, 0, ui.NewLabel("Name:"))
	c.panelContent.AddWidget(0, 1, c.txtName)

	c.OnDialogShow = func() {
		c.Form().SetTitle("New Config")
		c.Form().SetSize(400, 200)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnOK)
		c.Form().SetCancelButton(c.btnCancel)

		c.txtName.Focus()
	}

	return &c
}

func (c *CreateConfigDialog) Accept() {
	if c.onAccept != nil {
		c.onAccept(c.txtName.Text())
	}
	c.Form().Close()
}

func (c *CreateConfigDialog) Cancel() {
	if c.onCancel != nil {
		c.onCancel()
	}
	c.Form().Close()
}
