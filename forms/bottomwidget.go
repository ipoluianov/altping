package forms

import (
	"github.com/ipoluianov/altping/app"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
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
		newLinkLabel(func() string { return T().Settings }, func() { lastCreatedMainWidget.ShowSettings() }),
		newLinkLabel(func() string { return T().Help }, func() { openDocs(&c, "help") }),
		newLinkLabel(func() string { return T().About }, c.onAbout),
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

// newLinkLabel creates a hyperlink-style label with the text from text() that calls onClick on left click
func newLinkLabel(text func() string, onClick func()) *ui.Label {
	lbl := ui.NewLabel("")
	lbl.SetTextFunc(text)
	lbl.SetUnderline(true)
	lbl.SetForegroundColor(colorLink.get())
	linkLabels = append(linkLabels, lbl)
	lbl.SetMouseCursor(ui.MouseCursorPointer)
	lbl.SetOnMouseDown(func(button ui.MouseButton, x int, y int, mods ui.KeyModifiers) bool {
		if button != ui.MouseButtonLeft {
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
		ui.ShowMessageBox(parent, T().Error, err.Error())
	}
}

func (c *BottomWidget) onAbout() {
	c.ShowDialog(NewAboutDialog())
}

func (c *BottomWidget) timerUpdate() {
	mode := system.Get().PingServerMode()
	if mode == "udp" {
		mode = T().ModeUDP
	}
	if mode == "icmp" {
		mode = T().ModeICMP
	}
	c.lblStatus.SetText(mode)
}
