package forms

import (
	"image"

	"github.com/u00io/nuiforms/ui"
)

const toolButtonSize = 48

// ToolButton is a square image button for the toolbar.
// ui.ButtonImage has no disabled state, so it is handled here:
// a grayed-out image is shown and clicks are ignored.
type ToolButton struct {
	*ui.ButtonImage

	img         image.Image
	imgDisabled image.Image
	enabled     bool
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
