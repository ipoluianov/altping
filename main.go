package main

import (
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/forms"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

func main() {
	config.Init()
	ui.SetAppIcon(appIcon())
	form := ui.NewForm()
	mainForm := forms.NewMainForm()
	form.Panel().AddWidget(0, 0, mainForm)
	mainForm.UpdateTitle()
	maximized := mainForm.RestoreWindowState(form)
	form.OnClose = func() bool {
		mainForm.SaveWindowState()
		return true
	}
	system.Get().Start()
	form.SetOnGlobalKeyDown(forms.OnGlobalKeyDown)
	form.Show()
	if maximized {
		form.Maximize()
	}
	mainForm.Activate()
	form.Exec()
	system.Get().Stop()
}
