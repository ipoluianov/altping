package forms

import "github.com/u00io/nuiforms/ui"

type MainForm struct {
	ui.Widget

	splitter *ui.Splitter

	topWidget    *TopWidget
	leftWidget   *LeftWidget
	details      *DetailsWidget
	bottomWidget *BottomWidget
}

const (
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
	c.splitter.SetSecondSize(detailsInitialWidth)
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

func (c *MainForm) IsDetailsVisible() bool {
	return c.details.IsVisible()
}

func (c *MainForm) Activate() {
	c.leftWidget.FocusTable()
}
