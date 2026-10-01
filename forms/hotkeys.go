package forms

import "github.com/ipoluianov/nui/ui"

// OnGlobalKeyDown handles the application hotkeys. A key with modifiers
// that is not a hotkey is passed on to the focused widget (e.g. Ctrl+A to the table).
func OnGlobalKeyDown(key ui.Key, mods ui.KeyModifiers) bool {
	noMods := !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd
	altOnly := !mods.Shift && mods.Alt && !mods.Ctrl && !mods.Cmd
	ctrlOnly := !mods.Shift && !mods.Alt && mods.Ctrl && !mods.Cmd

	switch {
	case key == ui.KeyO && noMods:
		lastCreatedTopWidget.onBtnOpen()
	case key == ui.KeyA && noMods:
		lastCreatedTopWidget.onBtnAddItem()
	case (key == ui.KeyE || key == ui.KeyEnter) && noMods:
		lastCreatedTopWidget.onBtnEditItem()
	case key == ui.KeyF1 && noMods:
		openDocs(lastCreatedMainWidget, "help_f1")
	case key == ui.KeyD && noMods:
		lastCreatedTopWidget.onBtnDetails()
	case (key == ui.KeyDelete || key == ui.KeyBackspace) && noMods:
		lastCreatedTopWidget.onBtnRemoveItem()
	case key == ui.KeyM && ctrlOnly:
		form := lastCreatedMainWidget.Form()
		if form.IsMaximized() {
			form.Restore()
		} else {
			form.Maximize()
		}
	case key == ui.KeyT && noMods:
		lastCreatedTopWidget.onBtnTray()
	case key == ui.KeyH && ctrlOnly:
		lastCreatedMainWidget.Form().Minimize()
	case key == ui.KeyX && altOnly:
		// Closing from code does not call Form.OnClose, so the layout is saved here
		lastCreatedMainWidget.SaveWindowState()
		lastCreatedMainWidget.Form().Close()
	default:
		return false
	}
	return true
}
