package forms

import "github.com/u00io/nuiforms/ui"

type CreateConfigDialog struct {
	ui.Widget

	txtName *ui.TextBox

	btnOK     *ui.Button
	btnCancel *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel

	onAccept func()
	onCancel func()
}

func NewCreateConfigDialog(onAccept func(), onCancel func()) *CreateConfigDialog {
	var c CreateConfigDialog
	c.InitWidget()
	c.panelContent = ui.NewPanel()
	c.AddWidgetOnGrid(c.panelContent, 0, 0)
	c.AddWidgetOnGrid(ui.NewVSpacer(), 1, 0)
	c.panelButtons = ui.NewPanel()
	c.AddWidgetOnGrid(c.panelButtons, 2, 0)

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

	c.panelButtons.AddWidgetOnGrid(c.btnOK, 0, 0)
	c.panelButtons.AddWidgetOnGrid(c.btnCancel, 0, 1)

	c.txtName = ui.NewTextBox()
	c.panelContent.AddWidgetOnGrid(ui.NewLabel("Name:"), 0, 0)
	c.panelContent.AddWidgetOnGrid(c.txtName, 0, 1)

	return &c
}

func (c *CreateConfigDialog) GetConfigName() string {
	return c.txtName.Text()
}
