package forms

import "github.com/ipoluianov/nui/ui"

// trayIcon is created the first time the window goes to the tray and stays
// until the application quits
var trayIcon *ui.TrayIcon

// HideToTray hides the window, leaving the icon in the system tray to bring
// it back. Where there is no tray the window is just minimized.
func (c *MainForm) HideToTray() {
	if trayIcon == nil {
		icon, err := ui.NewTrayIcon()
		if err != nil {
			c.Form().Minimize()
			return
		}
		trayIcon = icon
		trayIcon.SetOnClick(c.ShowFromTray)
		c.updateTrayMenu()
		c.UpdateTitle()
	}
	c.Form().Hide()
}

// ShowFromTray brings the hidden window back
func (c *MainForm) ShowFromTray() {
	c.Form().Show()
	c.Activate()
}

// quitFromTray closes the application from the tray menu
func (c *MainForm) quitFromTray() {
	// Closing from code does not call Form.OnClose, so the layout is saved here
	c.SaveWindowState()
	c.Form().Close()
}

// updateTrayMenu sets the menu again, in the current language
func (c *MainForm) updateTrayMenu() {
	if trayIcon == nil {
		return
	}
	trayIcon.SetMenu(
		ui.TrayMenuItem{Text: T().TrayShow, OnClick: c.ShowFromTray},
		ui.TrayMenuItem{Separator: true},
		ui.TrayMenuItem{Text: T().TrayQuit, OnClick: c.quitFromTray},
	)
}

// updateTrayTooltip shows the window title, with the hosts down, on the icon
func updateTrayTooltip(title string) {
	if trayIcon != nil {
		trayIcon.SetTooltip(title)
	}
}

// CloseTray removes the tray icon; call it when the application quits
func CloseTray() {
	if trayIcon != nil {
		trayIcon.Close()
		trayIcon = nil
	}
}
