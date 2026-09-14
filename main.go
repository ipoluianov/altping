package main

import (
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/forms"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

func main() {
	config.Init()
	form := ui.NewForm()
	form.SetTitle("Alt Ping")
	form.SetSize(1100, 800)
	mainForm := forms.NewMainForm()
	form.Panel().AddWidget(mainForm, 0, 0)
	system.Get().Start()
	//ui.MainForm.SetOnGlobalKeyDown(forms.OnGlobalKeyDown)
	form.Show()
	form.Exec()
	system.Get().Stop()
	config.SaveLastConfigId()
}
