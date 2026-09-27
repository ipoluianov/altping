package forms

import (
	"image/color"

	"github.com/ipoluianov/altping/app"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type BottomWidget struct {
	ui.Widget

	lblStatus *ui.Label
}

func NewBottomWidget() *BottomWidget {
	var c BottomWidget
	c.InitWidget()
	c.SetPanelPadding(6)
	c.lblStatus = ui.NewLabel("---")
	c.AddWidget(0, 0, c.lblStatus)
	c.AddWidget(0, 1, ui.NewHSpacer())
	links := []*ui.Label{
		newLinkLabel("Settings", func() { lastCreatedMainWidget.ShowSettings() }),
		newLinkLabel("Help", func() { openDocs(&c, "help") }),
		newLinkLabel("About", c.onAbout),
	}
	for i, lbl := range links {
		if i > 0 {
			space := ui.NewSpace()
			space.SetSize(12, 0)
			c.AddWidget(0, 1+i*2, space)
		}
		c.AddWidget(0, 2+i*2, lbl)
	}

	c.AddTimer(1000, c.timerUpdate)

	return &c
}

// newLinkLabel creates a hyperlink-style label that calls onClick on left click
func newLinkLabel(text string, onClick func()) *ui.Label {
	lbl := ui.NewLabel(text)
	lbl.SetUnderline(true)
	lbl.SetForegroundColor(color.RGBA{0x3D, 0x8B, 0xF2, 0xFF})
	lbl.SetMouseCursor(nuimouse.MouseCursorPointer)
	lbl.SetOnMouseDown(func(button nuimouse.MouseButton, x int, y int, mods nuikey.KeyModifiers) bool {
		if button != nuimouse.MouseButtonLeft {
			return false
		}
		onClick()
		return true
	})
	return lbl
}

// openDocs opens the docs on the site; campaign tells which place in the app the visit came from
func openDocs(parent ui.Widgeter, campaign string) {
	if err := app.OpenSiteURL(app.DocsURL, campaign); err != nil {
		ui.ShowMessageBox(parent, "Error", err.Error())
	}
}

func (c *BottomWidget) onAbout() {
	c.ShowDialog(NewAboutDialog())
}

func (c *BottomWidget) timerUpdate() {
	mode := system.Get().PingServerMode()
	if mode == "udp" {
		mode = "UDP Mode"
	}
	if mode == "icmp" {
		mode = "ICMP Mode"
	}
	c.lblStatus.SetText(mode)
}
