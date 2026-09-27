package forms

import (
	"image"

	"github.com/u00io/nuiforms/ui"
)

const (
	toolButtonSize           = 48
	toolButtonCheckedLift    = 3
	toolButtonCheckMarkWidth = 3
	// Every button has a lighter bottom edge so it stands out from the toolbar
	toolButtonEdgeWidth = 2
	toolButtonEdgeLift  = 8
)

// ToolButton is a square image button for the toolbar.
// ui.ButtonImage has no disabled state, so it is handled here:
// a grayed-out image is shown and clicks are ignored.
type ToolButton struct {
	*ui.ButtonImage

	img         image.Image
	imgDisabled image.Image
	enabled     bool
	checked     bool
}

func NewToolButton(iconName string, tooltip string, onClick func()) *ToolButton {
	var c ToolButton
	c.img = loadIcon(iconName)
	c.imgDisabled = disabledIcon(c.img)
	c.enabled = true
	c.ButtonImage = ui.NewButtonImage(c.img)
	c.SetMinSize(toolButtonSize, toolButtonSize)
	c.SetMaxSize(toolButtonSize, toolButtonSize)
	c.SetTooltip(tooltip)
	c.SetOnPostPaint(c.drawBottomEdge)
	c.SetOnButtonClick(func(btn *ui.ButtonImage) {
		if c.enabled && onClick != nil {
			onClick()
		}
	})
	return &c
}

func (c *ToolButton) SetEnabled(enabled bool) {
	if c.enabled == enabled {
		return
	}
	c.enabled = enabled
	c.ButtonImage.SetEnabled(enabled)
	if enabled {
		c.SetImage(c.img)
	} else {
		c.SetImage(c.imgDisabled)
	}
}

func (c *ToolButton) IsChecked() bool {
	return c.checked
}

// SetChecked shows the button as toggled on: a lighter background
// and a bar along the bottom edge.
func (c *ToolButton) SetChecked(checked bool) {
	if c.checked == checked {
		return
	}
	c.checked = checked
	if checked {
		c.SetElevation(toolButtonCheckedLift)
	} else {
		c.SetElevation(0)
	}
	c.Form().Update()
}

// drawBottomEdge draws the lighter bottom edge, or the bright bar when the button is checked
func (c *ToolButton) drawBottomEdge(cnv *ui.Canvas) {
	if c.checked {
		cnv.FillRect(0, c.Height()-toolButtonCheckMarkWidth, c.Width(), toolButtonCheckMarkWidth, c.ForegroundColor())
		return
	}
	cnv.FillRect(0, c.Height()-toolButtonEdgeWidth, c.Width(), toolButtonEdgeWidth, c.BackgroundColorWithAddElevation(toolButtonEdgeLift))
}
