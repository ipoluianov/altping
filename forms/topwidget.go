package forms

import (
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

type TopWidget struct {
	ui.Widget

	btnStart *ui.Button
	btnStop  *ui.Button
}

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()
	c.SetPanelPadding(6)

	c.btnStart = ui.NewButton("Start")
	c.btnStart.SetOnClick(c.onBtnStart)
	c.btnStop = ui.NewButton("Stop")
	c.btnStop.SetOnClick(c.onBtnStop)

	c.AddWidgetOnGrid(c.btnStart, 0, 0)
	c.AddWidgetOnGrid(c.btnStop, 0, 1)
	c.AddWidgetOnGrid(ui.NewHSpacer(), 0, 10)

	c.AddTimer(200, c.timerUpdate)

	return &c
}

func (c *TopWidget) timerUpdate() {
	if system.Get().IsRunning() {
		c.btnStart.SetEnabled(false)
		c.btnStop.SetEnabled(true)
	} else {
		c.btnStart.SetEnabled(true)
		c.btnStop.SetEnabled(false)
	}
}

func (c *TopWidget) onBtnStart() {
	system.Get().Start()
}

func (c *TopWidget) onBtnStop() {
	system.Get().Stop()
}
