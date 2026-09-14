package forms

import "github.com/u00io/nuiforms/ui"

type MainForm struct {
	ui.Widget

	panelCenter *ui.Panel

	topWidget    *TopWidget
	leftWidget   *LeftWidget
	centerWidget *CenterWidget
	bottomWidget *BottomWidget
}

func NewMainForm() *MainForm {
	var c MainForm
	c.InitWidget()
	c.topWidget = NewTopWidget()
	c.leftWidget = NewLeftWidget(c.SetMode)
	c.centerWidget = NewCenterWidget()
	c.bottomWidget = NewBottomWidget()

	c.AddWidget(c.topWidget, 0, 0)
	c.panelCenter = ui.NewPanel()
	c.panelCenter.AddWidget(c.leftWidget, 0, 0)
	c.panelCenter.AddWidget(c.centerWidget, 0, 1)
	c.AddWidget(c.panelCenter, 1, 0)
	c.AddWidget(c.bottomWidget, 2, 0)
	return &c
}

func (c *MainForm) SetMode(mode string) {
	//c.centerWidget.SetMode(mode)
}
