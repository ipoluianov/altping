package forms

import "github.com/u00io/nuiforms/ui"

type MainForm struct {
	ui.Widget

	panelCenter *ui.Panel

	topWidget    *TopWidget
	leftWidget   *LeftWidget
	centerWidget *CenterWidget
	details      *DetailsWidget
	bottomWidget *BottomWidget
}

var lastCreatedMainWidget *MainForm

func NewMainForm() *MainForm {
	var c MainForm
	c.InitWidget()
	c.topWidget = NewTopWidget()
	c.leftWidget = NewLeftWidget()
	c.centerWidget = NewCenterWidget()
	c.bottomWidget = NewBottomWidget()

	c.AddWidget(0, 0, c.topWidget)
	c.panelCenter = ui.NewPanel()
	c.panelCenter.SetPanelPadding(0)
	c.panelCenter.AddWidget(0, 0, c.leftWidget)
	c.panelCenter.AddWidget(0, 1, c.centerWidget)
	c.details = NewDetailsWidget()
	c.details.SetVisible(false)
	c.panelCenter.AddWidget(0, 2, c.details)
	c.AddWidget(1, 0, c.panelCenter)
	c.AddWidget(2, 0, c.bottomWidget)
	lastCreatedMainWidget = &c

	c.SetPanelPadding(3)
	return &c
}

func (c *MainForm) ToggleDetails() {
	c.details.SetVisible(!c.details.IsVisible())
	c.Form().UpdateLayout()
	c.Form().Update()
}

func (c *MainForm) Activate() {
	c.leftWidget.FocusTable()
}
