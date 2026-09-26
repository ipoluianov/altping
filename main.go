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
	form.SetTitle("Alt Ping")
	form.SetSize(1100, 800)
	mainForm := forms.NewMainForm()
	form.Panel().AddWidget(0, 0, mainForm)
	system.Get().Start()
	form.SetOnGlobalKeyDown(forms.OnGlobalKeyDown)
	form.Show()
	mainForm.Activate()
	form.Exec()
	system.Get().Stop()
	config.SaveLastConfigId()
}
