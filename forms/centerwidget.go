package forms

import (
	"github.com/u00io/nuiforms/ui"
)

type CenterWidget struct {
	ui.Widget
}

func NewCenterWidget() *CenterWidget {
	var c CenterWidget
	c.InitWidget()
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	//c.SetMode("common")

	return &c
}
