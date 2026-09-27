package forms

import (
	"github.com/ipoluianov/altping/config"
	"github.com/u00io/nuiforms/ui"
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
}

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

	c.SetPanelPadding(3)
	return &c
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
	c.Form().SetTitle(title)
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
