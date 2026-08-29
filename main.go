package main

import (
	"github.com/ipoluianov/altping/forms"
	"github.com/ipoluianov/altping/localstorage"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

func main() {
	localstorage.Init("altping")
	form := ui.NewForm()
	form.SetTitle("Alt Ping")
	form.SetSize(1100, 800)
	mainForm := forms.NewMainForm()
	form.Panel().AddWidgetOnGrid(mainForm, 0, 0)
	system.Get().Start()
	form.Exec()
	system.Get().Stop()
}
