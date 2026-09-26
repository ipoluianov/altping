package forms

import (
	"image/color"

	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nui/nuikey"
	"github.com/u00io/nui/nuimouse"
	"github.com/u00io/nuiforms/ui"
)

type BottomWidget struct {
	ui.Widget

	lblStatus *ui.Label
	lblAbout  *ui.Label
}

func NewBottomWidget() *BottomWidget {
	var c BottomWidget
	c.InitWidget()
	c.SetPanelPadding(6)
	c.lblStatus = ui.NewLabel("---")
	c.AddWidget(0, 0, c.lblStatus)
	c.AddWidget(0, 1, ui.NewHSpacer())
	c.lblAbout = newLinkLabel("About", c.onAbout)
	c.AddWidget(0, 2, c.lblAbout)

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
