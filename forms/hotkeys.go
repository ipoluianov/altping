package forms

import (
	"os"

	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nuiforms/ui"
)

func OnGlobalKeyDown(key nuikey.Key, mods nuikey.KeyModifiers) bool {
	if key == nuikey.KeyN {
		if !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd { // N
			lastCreatedTopWidget.onBtnNew()
		}
		return true
	}

	if key == nuikey.KeyO {
		if !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd { // O
			lastCreatedTopWidget.onBtnOpen()
		}
		return true
	}

	if key == nuikey.KeyS {
		if !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd { // S
			lastCreatedTopWidget.onBtnSaveAs()
		}
		return true
	}

	if key == nuikey.KeyA {
		if !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd { // N
			lastCreatedTopWidget.onBtnAddItem()
		}
		return true
	}

	if key == nuikey.KeyE {
		if !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd { // E
			lastCreatedTopWidget.onBtnEditItem()
		}
		return true
	}

	if key == nuikey.KeyDelete || key == nuikey.KeyBackspace {
		if !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd { // Delete or Backspace
			lastCreatedTopWidget.onBtnRemoveItem()
		}
		return true
	}

	if key == nuikey.KeyX {
		if !mods.Shift && mods.Alt && !mods.Ctrl && !mods.Cmd { // Alt+X
			lastCreatedMainWidget.Form().Close()
		}
		return true
	}

	if key == nuikey.KeyEsc {
		if !mods.Shift && !mods.Alt && !mods.Ctrl && !mods.Cmd { // Esc
			if lastCreatedMainWidget.Form().TopPopupWidget() == nil {
				ui.ShowQuestionMessageBoxOKCancel(lastCreatedMainWidget, "Closing", "Close the application?", func() {
					os.Exit(0)
				}, func() {
				})
				return true
			}
		}
	}

	return false
}
