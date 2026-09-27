package forms

import (
	"os"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

// OnGlobalKeyDown handles the application hotkeys. A key with modifiers
// that is not a hotkey is passed on to the focused widget (e.g. Ctrl+A to the table).
func OnGlobalKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	noMods := !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd
	altOnly := !mods.Shift && mods.Alt && !mods.Ctrl && !mods.Cmd

	switch {
	case key == nuikey.KeyO && noMods:
		lastCreatedTopWidget.onBtnOpen()
	case key == nuikey.KeyA && noMods:
		lastCreatedTopWidget.onBtnAddItem()
	case key == nuikey.KeyE && noMods:
		lastCreatedTopWidget.onBtnEditItem()
	case key == nuikey.KeyD && noMods:
		lastCreatedTopWidget.onBtnDetails()
	case (key == nuikey.KeyDelete || key == nuikey.KeyBackspace) && noMods:
		lastCreatedTopWidget.onBtnRemoveItem()
	case key == nuikey.KeyX && altOnly:
		lastCreatedMainWidget.Form().Close()
	case key == nuikey.KeyEsc && noMods && lastCreatedMainWidget.Form().TopPopupWidget() == nil:
		ui.ShowQuestionMessageBoxOKCancel(lastCreatedMainWidget, "Closing", "Close the application?", func() {
			os.Exit(0)
		}, func() {
		})
	default:
		return false
	}
	return true
}
