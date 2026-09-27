package forms

import (
	"github.com/u00io/nui/nuikey"
)

// OnGlobalKeyDown handles the application hotkeys. A key with modifiers
// that is not a hotkey is passed on to the focused widget (e.g. Ctrl+A to the table).
func OnGlobalKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	noMods := !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd
	altOnly := !mods.Shift && mods.Alt && !mods.Ctrl && !mods.Cmd
	ctrlOnly := !mods.Shift && !mods.Alt && mods.Ctrl && !mods.Cmd

	switch {
	case key == nuikey.KeyO && noMods:
		lastCreatedTopWidget.onBtnOpen()
	case key == nuikey.KeyA && noMods:
		lastCreatedTopWidget.onBtnAddItem()
	case (key == nuikey.KeyE || key == nuikey.KeyEnter) && noMods:
		lastCreatedTopWidget.onBtnEditItem()
	case key == nuikey.KeyD && noMods:
		lastCreatedTopWidget.onBtnDetails()
	case (key == nuikey.KeyDelete || key == nuikey.KeyBackspace) && noMods:
		lastCreatedTopWidget.onBtnRemoveItem()
	case key == nuikey.KeyM && ctrlOnly:
		form := lastCreatedMainWidget.Form()
		if form.IsMaximized() {
			form.Restore()
		} else {
			form.Maximize()
		}
	case key == nuikey.KeyH && ctrlOnly:
		lastCreatedMainWidget.Form().Minimize()
	case key == nuikey.KeyX && altOnly:
		// Closing from code does not call Form.OnClose, so the layout is saved here
		lastCreatedMainWidget.SaveWindowState()
		lastCreatedMainWidget.Form().Close()
	default:
		return false
	}
	return true
}
