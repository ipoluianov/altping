package forms

import (
	"github.com/ipoluianov/nui/ui"
)

type CreateConfigDialog struct {
	ui.DialogContent

	txtName *ui.TextBox
	title   string

	btnOK     *ui.Button
	btnCancel *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel

	onAccept func(name string)
	onCancel func()
}

// NewCreateConfigDialog asks for a config name, initially set to name
func NewCreateConfigDialog(title string, name string, onAccept func(name string), onCancel func()) *CreateConfigDialog {
	var c CreateConfigDialog
	c.InitWidget()
	c.title = title

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
	c.txtName.SetText(name)
	c.panelContent.AddWidget(0, 0, ui.NewLabel("Name:"))
	c.panelContent.AddWidget(0, 1, c.txtName)

	c.OnDialogShow = func() {
		c.Form().SetTitle(c.title)
		c.Form().SetSize(400, 200)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnOK)
		c.Form().SetCancelButton(c.btnCancel)

		c.txtName.Focus()
		c.txtName.SelectAllText()
	}

	return &c
}

func (c *CreateConfigDialog) Accept() {
	// The callbacks change the window below: they run on its goroutine
	if c.onAccept != nil {
		name := c.txtName.Text()
		c.RunInParent(func() { c.onAccept(name) })
	}
	c.Form().Close()
}

func (c *CreateConfigDialog) Cancel() {
	c.RunInParent(c.onCancel)
	c.Form().Close()
}
