package forms

import (
	"net/url"

	"github.com/ipoluianov/altping/app"
	"github.com/u00io/nuiforms/ui"
)

type AboutDialog struct {
	ui.DialogContent

	btnWebsite *ui.Button
	btnClose   *ui.Button

	panelContent *ui.Panel
	panelButtons *ui.Panel
}

func NewAboutDialog() *AboutDialog {
	var c AboutDialog
	c.InitWidget()

	c.panelContent = ui.NewPanel()
	c.AddWidget(0, 0, c.panelContent)
	c.AddWidget(1, 0, ui.NewVSpacer())
	c.panelButtons = ui.NewPanel()
	c.AddWidget(2, 0, c.panelButtons)

	lblName := ui.NewLabel(app.DisplayName)
	lblName.SetFontSize(24)
	lblName.SetTextAlign(ui.HAlignCenter)
	lblVersion := ui.NewLabel("Version " + app.Version)
	lblVersion.SetTextAlign(ui.HAlignCenter)
	lblAuthor := ui.NewLabel("Author: " + app.Author)
	lblAuthor.SetTextAlign(ui.HAlignCenter)
	lblCopyright := ui.NewLabel(app.Copyright())
	lblCopyright.SetTextAlign(ui.HAlignCenter)
	lblLicense := ui.NewLabel("License: " + app.License)
	lblLicense.SetTextAlign(ui.HAlignCenter)
	lblWebsite := ui.NewLabel(app.Website)
	lblWebsite.SetTextAlign(ui.HAlignCenter)

	c.panelContent.AddWidget(0, 0, lblName)
	c.panelContent.AddWidget(1, 0, lblVersion)
	c.panelContent.AddWidget(2, 0, lblAuthor)
	c.panelContent.AddWidget(3, 0, lblCopyright)
	c.panelContent.AddWidget(4, 0, lblLicense)
	c.panelContent.AddWidget(5, 0, lblWebsite)

	c.btnWebsite = ui.NewButton("Visit Website")
	c.btnWebsite.SetOnClick(c.VisitWebsite)
	c.btnClose = ui.NewButton("Close")
	c.btnClose.SetOnClick(c.Close)

	c.panelButtons.AddWidget(0, 0, ui.NewHSpacer())
	c.panelButtons.AddWidget(0, 1, c.btnWebsite)
	c.panelButtons.AddWidget(0, 2, c.btnClose)

	c.OnDialogShow = func() {
		c.Form().SetTitle("About " + app.DisplayName)
		c.Form().SetSize(400, 300)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnClose)
		c.Form().SetCancelButton(c.btnClose)
	}

	return &c
}

func (c *AboutDialog) VisitWebsite() {
	// UTM tags mark visits that came from the About dialog
	q := url.Values{}
	q.Set("utm_source", "altping")
	q.Set("utm_medium", "app")
	q.Set("utm_campaign", "about_dialog")
	q.Set("utm_content", app.Version)
	if err := app.OpenURL(app.Website + "?" + q.Encode()); err != nil {
		ui.ShowMessageBox(c, "Error", err.Error())
	}
}

func (c *AboutDialog) Close() {
	c.Form().Close()
}
