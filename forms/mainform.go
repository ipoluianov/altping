package forms

import (
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

type MainForm struct {
	ui.Widget

	splitter *ui.Splitter

	topWidget    *TopWidget
	leftWidget   *LeftWidget
	details      *DetailsWidget
	bottomWidget *BottomWidget

	// The width the user gave the details, kept while they are hidden
	detailsWidth int

	// Hosts known to be down (true) or up (false), for the notifications
	hostDown  map[string]bool
	downCount int
}

// A host is down after this many failed pings in a row, so a single lost ping does not beep
const downAfterFailures = 3

const (
	appTitle            = "Alt Ping"
	defaultWindowWidth  = 1100
	defaultWindowHeight = 800
	detailsInitialWidth = 500
	// When the window gets narrow the details shrink first, down to this table width
	tableMinWidth = 400
)

var lastCreatedMainWidget *MainForm

func NewMainForm() *MainForm {
	var c MainForm
	c.InitWidget()
	c.topWidget = NewTopWidget()
	c.leftWidget = NewLeftWidget()
	c.bottomWidget = NewBottomWidget()

	c.AddWidget(0, 0, c.topWidget)
	c.details = NewDetailsWidget()
	c.details.SetVisible(false)
	// The details keep their width when the window is resized
	c.leftWidget.SetMinWidth(tableMinWidth)
	c.splitter = ui.NewHSplitter()
	c.splitter.SetWidgets(c.leftWidget, c.details)
	c.detailsWidth = detailsInitialWidth
	c.splitter.SetSecondSize(c.detailsWidth)
	c.splitter.SetOnSplitChanged(func() {
		c.detailsWidth = c.splitter.SecondSize()
	})
	c.AddWidget(1, 0, c.splitter)
	c.AddWidget(2, 0, c.bottomWidget)
	lastCreatedMainWidget = &c

	c.hostDown = make(map[string]bool)
	c.AddTimer(500, c.checkHostsDown)

	c.SetPanelPadding(3)
	return &c
}

// ApplySettings saves the settings and applies what they change in the window
func (c *MainForm) ApplySettings(s config.Settings) {
	languageChanged := s.Language != config.GetSettings().Language
	themeChanged := s.Theme != config.GetSettings().Theme
	if err := config.SetSettings(s); err != nil {
		ui.ShowMessageBox(c, T().Error, err.Error())
	}
	if languageChanged {
		SetLanguage(s.Language)
	}
	if themeChanged {
		ApplyTheme(s.Theme)
	}
	c.leftWidget.ApplyColumns()
	c.Form().SetAlwaysOnTop(s.AlwaysOnTop)
	c.checkHostsDown()
}

// ShowSettings opens the settings dialog
func (c *MainForm) ShowSettings() {
	c.ShowDialog(NewSettingsDialog(config.GetSettings(), c.ApplySettings))
}

// checkHostsDown beeps and marks the window in the taskbar when a host with
// notifications on goes down or comes back. The first state seen is not reported.
func (c *MainForm) checkHostsDown() {
	if !system.Get().IsRunning() {
		clear(c.hostDown)
		c.setDownCount(0)
		return
	}

	changed := false
	downCount := 0
	seen := make(map[string]bool)
	for _, h := range config.Get().Hosts {
		history := system.Get().GetHostHistory(h.ID)
		if !h.Notify || history == nil {
			continue
		}
		seen[h.ID] = true
		streak := history.FailStreak()
		wasDown, known := c.hostDown[h.ID]
		down := streak >= downAfterFailures
		if !down && wasDown && streak > 0 {
			down = true // still failing: not back yet
		}
		if down {
			downCount++
		}
		// No ping result after a start: nothing to compare with yet
		if history.LastChange().IsZero() {
			continue
		}
		if known && down != wasDown {
			changed = true
		}
		c.hostDown[h.ID] = down
	}
	for id := range c.hostDown {
		if !seen[id] {
			delete(c.hostDown, id)
		}
	}

	c.setDownCount(downCount)
	if changed {
		c.Form().Beep()
		c.Form().RequestAttention()
	}
}

func (c *MainForm) setDownCount(n int) {
	if c.downCount == n {
		return
	}
	c.downCount = n
	c.UpdateTitle()
}

func (c *MainForm) ToggleDetails() {
	c.details.SetVisible(!c.details.IsVisible())
	c.details.Refresh()
	c.Form().UpdateLayout()
	c.Form().Update()
}

// UpdateTitle shows the name of the opened config in the window title
func (c *MainForm) UpdateTitle() {
	title := appTitle
	if cfg := config.Get(); cfg != nil && cfg.Name != "" {
		title += " - " + cfg.Name
	}
	// Seen in the taskbar while the window is in the background
	if c.downCount > 0 {
		title = T().DownCount(c.downCount) + " " + title
	}
	c.Form().SetTitle(title)
}

// ApplyLanguage updates the texts that do not follow the language by themselves
func (c *MainForm) ApplyLanguage() {
	c.UpdateTitle()
	c.leftWidget.updateColumnNames()
	c.leftWidget.timerUpdate()
	c.details.applyLanguage()
	c.bottomWidget.timerUpdate()
}

// RestoreWindowState applies the saved window layout before the form is shown.
// Returns whether the window should be maximized once shown.
func (c *MainForm) RestoreWindowState(form *ui.Form) (maximized bool) {
	form.SetSize(defaultWindowWidth, defaultWindowHeight)
	state, ok := config.LoadWindowState()
	if !ok {
		return false
	}
	form.SetSize(state.Width, state.Height)
	// 0,0 means the position was not known (the window manager did not report it)
	if state.X != 0 || state.Y != 0 {
		form.Move(state.X, state.Y)
	}
	if state.DetailsWidth > 0 {
		c.detailsWidth = state.DetailsWidth
		c.splitter.SetSecondSize(c.detailsWidth)
	}
	if state.DetailsVisible {
		c.setDetailsVisible(true)
	}
	return state.Maximized
}

// SaveWindowState remembers the window layout for the next start
func (c *MainForm) SaveWindowState() {
	form := c.Form()
	state, _ := config.LoadWindowState()
	state.Maximized = form.IsMaximized()
	// The size of a maximized window is the screen size: keep the normal one
	if !state.Maximized {
		state.X, state.Y = form.Position()
		state.Width, state.Height = form.Size()
	}
	state.DetailsVisible = c.details.IsVisible()
	state.DetailsWidth = c.detailsWidth
	config.SaveWindowState(state)
}

// setDetailsVisible shows or hides the details the same way as the toolbar button
func (c *MainForm) setDetailsVisible(visible bool) {
	if c.details.IsVisible() != visible {
		c.topWidget.onBtnDetails()
	}
}

func (c *MainForm) IsDetailsVisible() bool {
	return c.details.IsVisible()
}

func (c *MainForm) Activate() {
	c.leftWidget.FocusTable()
}
