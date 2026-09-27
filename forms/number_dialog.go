package forms

import (
	"github.com/ipoluianov/nui/ui"
)

// NumberDialog asks for one whole number, e.g. a port or a time in ms
type NumberDialog struct {
	ui.DialogContent

	title     string
	num       *ui.NumBox
	btnOK     *ui.Button
	btnCancel *ui.Button

	onAccept func(value int)
}

func NewNumberDialog(title, label string, value, minValue, maxValue, step int, onAccept func(value int)) *NumberDialog {
	var c NumberDialog
	c.InitWidget()
	c.title = title
	c.onAccept = onAccept

	content := ui.NewPanel()
	c.AddWidget(0, 0, content)
	c.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	c.AddWidget(2, 0, buttons)

	c.num = newMsBox(int64(value), minValue, maxValue)
	c.num.SetStep(float64(step))
	content.AddWidget(0, 0, ui.NewLabel(label))
	content.AddWidget(0, 1, c.num)

	c.btnOK = ui.NewButton("OK")
	c.btnOK.SetOnClick(c.Accept)
	c.btnCancel = ui.NewButton("Cancel")
	c.btnCancel.SetOnClick(func() { c.Form().Close() })
	buttons.AddWidget(0, 0, ui.NewHSpacer())
	buttons.AddWidget(0, 1, c.btnOK)
	buttons.AddWidget(0, 2, c.btnCancel)

	c.OnDialogShow = func() {
		c.Form().SetTitle(c.title)
		c.Form().SetSize(400, 150)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnOK)
		c.Form().SetCancelButton(c.btnCancel)
		c.num.Focus()
	}
	return &c
}

func (c *NumberDialog) Accept() {
	// The callback changes the window below: it runs on its goroutine
	if c.onAccept != nil {
		value := int(c.num.Value())
		c.RunInParent(func() { c.onAccept(value) })
	}
	c.Form().Close()
}
