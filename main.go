package main

import (
	"github.com/ipoluianov/altping/app"
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/forms"
	"github.com/ipoluianov/altping/install"
	"github.com/ipoluianov/altping/instance"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

func main() {
	install.SetIcon(iconPNG)
	// Started from "Installed apps" to remove it
	if install.HasArg(install.UninstallArg) {
		uninstall(install.HasArg(install.QuietArg))
		return
	}

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
		if install.HasArg(install.InstalledArg) {
			mainForm.ShowInstalled()
		}
	})
	form.Exec()
	forms.CloseTray()
	system.Get().Stop()

	// Installed from this copy: the installed one takes over, so the lock goes first
	if install.RelaunchPending() {
		inst.Close()
		if err := install.StartInstalled(); err != nil {
			install.Inform(app.DisplayName, err.Error())
		}
	}
}
