package forms

import (
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/nui/ui"
)

// SettingsDialog edits the application options (config.Settings)
type SettingsDialog struct {
	ui.DialogContent

	chkMin    *ui.Checkbox
	chkJitter *ui.Checkbox
	chkSince  *ui.Checkbox

	chkOnTop *ui.Checkbox

	btnOK     *ui.Button
	btnCancel *ui.Button

	onAccept func(config.Settings)
}

func NewSettingsDialog(settings config.Settings, onAccept func(config.Settings)) *SettingsDialog {
	var c SettingsDialog
	c.InitWidget()
	c.onAccept = onAccept

	checks := ui.NewPanel()
	c.AddWidget(0, 0, checks)
	c.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	c.AddWidget(2, 0, buttons)

	row := 0
	addCheck := func(text string, checked bool) *ui.Checkbox {
		chk := ui.NewCheckbox(text)
		chk.SetChecked(checked)
		chk.SetXExpandable(true)
		checks.AddWidget(row, 0, chk)
		row++
		return chk
	}
	checks.AddWidget(row, 0, ui.NewLabel("Extra columns:"))
	row++
	c.chkMin = addCheck("Min", settings.ShowMin)
	c.chkJitter = addCheck("Jitter", settings.ShowJitter)
	c.chkSince = addCheck("Since (last change)", settings.ShowSince)
	c.chkOnTop = addCheck("Always on top", settings.AlwaysOnTop)

	c.btnOK = ui.NewButton("OK")
	c.btnOK.SetOnClick(c.Accept)
	c.btnCancel = ui.NewButton("Cancel")
	c.btnCancel.SetOnClick(c.Cancel)
	buttons.AddWidget(0, 0, ui.NewHSpacer())
	buttons.AddWidget(0, 1, c.btnOK)
	buttons.AddWidget(0, 2, c.btnCancel)

	c.OnDialogShow = func() {
		c.Form().SetTitle("Settings")
		c.Form().SetSize(420, 290)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnOK)
		c.Form().SetCancelButton(c.btnCancel)
	}
	return &c
}

func (c *SettingsDialog) Accept() {
	s := config.GetSettings()
	s.ShowMin = c.chkMin.Checked()
	s.ShowJitter = c.chkJitter.Checked()
	s.ShowSince = c.chkSince.Checked()
	s.AlwaysOnTop = c.chkOnTop.Checked()
	// The callback changes the window below: it runs on its goroutine
	if c.onAccept != nil {
		c.RunInParent(func() { c.onAccept(s) })
	}
	c.Form().Close()
}

func (c *SettingsDialog) Cancel() {
	c.Form().Close()
}
