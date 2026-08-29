package forms

import (
	"github.com/ipoluianov/altping/system"
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
	c.AddWidgetOnGrid(c.lblStatus, 0, 0)
	c.AddWidgetOnGrid(ui.NewHSpacer(), 0, 1)
	c.AddWidgetOnGrid(ui.NewLabel("AltBins"), 0, 2)

	c.AddTimer(1000, c.timerUpdate)
	return &c
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
