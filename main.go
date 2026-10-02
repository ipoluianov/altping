package main

import (
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/forms"
	"github.com/ipoluianov/altping/instance"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

func main() {
	// Taken before anything is read or written: the files belong to one copy
	inst, ok := instance.Acquire(config.ConfigDirectory())
	if !ok {
		return // the running copy shows its window instead
	}
	defer inst.Close()

	config.Init()
	forms.SetLanguage(config.GetSettings().Language)
	forms.ApplyTheme(config.GetSettings().Theme)
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
	form.SetOnLanguageChanged(mainForm.ApplyLanguage)
	form.Show()
	inst.Serve(func() { form.Invoke(mainForm.BringToFront) })
	// The form is handled by its own goroutine once shown
	form.Invoke(func() {
		if maximized {
			form.Maximize()
		}
		form.SetAlwaysOnTop(config.GetSettings().AlwaysOnTop)
		mainForm.Activate()
		mainForm.ShowStartupErrors()
	})
	form.Exec()
	forms.CloseTray()
	system.Get().Stop()
}
