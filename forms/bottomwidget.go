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
	c.AddWidget(0, 0, c.lblStatus)
	c.AddWidget(0, 1, ui.NewHSpacer())
	c.AddWidget(0, 2, ui.NewLabel("AltBins"))

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
