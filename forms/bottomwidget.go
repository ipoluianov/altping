package forms

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/ipoluianov/altping/app"
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/install"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

// BottomWidget is the status bar: how the hosts are doing, the ping mode and the links
type BottomWidget struct {
	ui.Widget

	// The dot in the color of the worst host state, and the counts of the hosts in each state
	statusDot      *ui.Space
	statusDotColor color.RGBA
	lblStatus      *ui.Label
	lblMode        *ui.Label

	// What the install link does: install, update or uninstall
	installStatus install.Status
}

func NewBottomWidget() *BottomWidget {
	var c BottomWidget
	c.InitWidget()
	c.SetPanelPadding(6)
	c.statusDot = ui.NewSpace()
	c.statusDot.SetSize(statusDotSize, ui.ThemeControlHeight())
	c.statusDot.SetOnPaint(func(cnv *ui.Canvas) {
		cnv.DrawImage(0, (c.statusDot.Height()-statusDotSize)/2, statusDot(c.statusDotColor))
	})
	c.AddWidget(0, 0, c.statusDot)
	dotSpace := ui.NewSpace()
	dotSpace.SetSize(statusDotGap, 0)
	c.AddWidget(0, 1, dotSpace)
	c.lblStatus = ui.NewLabel("")
	c.AddWidget(0, 2, c.lblStatus)
	modeSpace := ui.NewSpace()
	modeSpace.SetSize(24, 0)
	c.AddWidget(0, 3, modeSpace)
	c.lblMode = ui.NewLabel("")
	c.AddWidget(0, 4, c.lblMode)
	c.AddWidget(0, 5, ui.NewHSpacer())
	links := []*ui.Label{
		newLinkLabel(func() string { return T().Settings }, func() { lastCreatedMainWidget.ShowSettings() }),
		newLinkLabel(func() string { return T().Help }, func() { openDocs(&c, "help") }),
	}
	// A downloaded copy offers to install or update itself, the installed one to be removed
	c.installStatus = install.CurrentStatus()
	if c.installStatus != install.StatusNone {
		links = append(links, newLinkLabel(c.installText, c.onInstallLink))
	}
	links = append(links, newLinkLabel(func() string { return T().About }, c.onAbout))
	for i, lbl := range links {
		if i > 0 {
			space := ui.NewSpace()
			space.SetSize(12, 0)
			c.AddWidget(0, 5+i*2, space)
		}
		c.AddWidget(0, 6+i*2, lbl)
	}

	c.AddTimer(500, c.timerUpdate)
	c.timerUpdate()

	return &c
}

// newLinkLabel creates a hyperlink-style label with the text from text() that calls onClick on left click.
// It is underlined only under the mouse, so a row of links stays quiet.
func newLinkLabel(text func() string, onClick func()) *ui.Label {
	lbl := ui.NewLabel("")
	lbl.SetTextFunc(text)
	lbl.SetForegroundColor(colorLink.get())
	lbl.SetOnMouseEnter(func() { lbl.SetUnderline(true) })
	lbl.SetOnMouseLeave(func() { lbl.SetUnderline(false) })
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

// installText is the text of the install link for what it does now
func (c *BottomWidget) installText() string {
	switch c.installStatus {
	case install.StatusUpdate:
		return T().Update
	case install.StatusUninstall:
		return T().Uninstall
	}
	return T().Install
}

func (c *BottomWidget) onInstallLink() {
	if c.installStatus == install.StatusUninstall {
		c.onUninstall()
		return
	}
	c.onInstall()
}

// onInstall copies the application to ~/.altbins and registers it, then
// quits for the installed copy to start (see main). Over an older installed
// version it is the same: that one is replaced.
func (c *BottomWidget) onInstall() {
	ui.ShowQuestionMessageBoxOKCancel(c, c.installText(), T().InstallAsk(install.Dir()), func() {
		if err := install.Install(); err != nil {
			ui.ShowMessageBox(c, T().Error, T().InstallFailed(err.Error()))
			return
		}
		install.RelaunchAfterExit()
		lastCreatedMainWidget.SaveWindowState()
		lastCreatedMainWidget.Form().Close()
	}, nil)
}

// onUninstall removes the installed copy; the settings and host lists stay.
// When that is this copy, it quits and its binary is deleted once it has.
func (c *BottomWidget) onUninstall() {
	ui.ShowQuestionMessageBoxOKCancel(c, T().Uninstall, T().UninstallAsk(config.ConfigDirectory()), func() {
		installed := install.IsInstalledCopy()
		if err := install.Uninstall(); err != nil {
			ui.ShowMessageBox(c, T().Error, err.Error())
			return
		}
		if installed {
			lastCreatedMainWidget.SaveWindowState()
			lastCreatedMainWidget.Form().Close()
			return
		}
		c.installStatus = install.CurrentStatus()
		ui.ShowToast(c, T().Uninstalled, ui.ToastSuccess)
	}, nil)
}

func (c *BottomWidget) timerUpdate() {
	mode := system.Get().PingServerMode()
	if mode == "udp" {
		mode = T().ModeUDP
	}
	if mode == "icmp" {
		mode = T().ModeICMP
	}
	c.lblMode.SetText(mode)
	c.lblMode.SetForegroundColor(colorMuted.get())

	text, col := hostsSummary()
	c.lblStatus.SetText(text)
	if col != c.statusDotColor {
		c.statusDotColor = col
		if f := c.Form(); f != nil {
			f.Update()
		}
	}
}

// hostsSummary counts the hosts in each state, e.g. "Hosts: 12 · OK: 10 · Down: 2",
// and returns the color of the worst state for the dot
func hostsSummary() (string, color.RGBA) {
	hosts := config.Get().Hosts
	if !system.Get().IsRunning() {
		return T().StatusStopped, colorNotPinged.get()
	}
	ok, slow, down := 0, 0, 0
	for _, h := range hosts {
		r := hostRowOf(h)
		switch {
		case !r.processed:
		case r.failed:
			down++
		case r.isSlow():
			slow++
		default:
			ok++
		}
	}
	parts := []string{fmt.Sprintf("%s: %d", T().StatusHosts, len(hosts))}
	for _, part := range []struct {
		name  string
		count int
	}{{T().StatusOK, ok}, {T().StatusSlow, slow}, {T().StatusDown, down}} {
		if part.count > 0 {
			parts = append(parts, fmt.Sprintf("%s: %d", part.name, part.count))
		}
	}
	col := colorNotPinged.get()
	switch {
	case down > 0:
		col = colorFailed.get()
	case slow > 0:
		col = colorSlow.get()
	case ok > 0:
		col = colorOK.get()
	}
	return strings.Join(parts, "  ·  "), col
}
