package forms

import (
	"strings"

	"github.com/ipoluianov/altping/app"
	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/nui/ui"
)

// shareOpenMax is how many pages are opened at once, so a large selection
// does not flood the browser with tabs
const shareOpenMax = 10

// sharedHosts returns the hosts that are shared on u00.io
func sharedHosts(hosts []*config.ConfigHost) []*config.ConfigHost {
	var res []*config.ConfigHost
	for _, h := range hosts {
		if h.Share && h.ShareURL() != "" {
			res = append(res, h)
		}
	}
	return res
}

// openShareURL opens the public u00.io page of the host in the browser
func openShareURL(w ui.Widgeter, host config.ConfigHost) {
	if err := app.OpenURL(host.ShareURL()); err != nil {
		ui.ShowMessageBox(w, T().Error, err.Error())
	}
}

// copyShareURL puts the link to the public u00.io page of the host in the clipboard
func copyShareURL(w ui.Widgeter, host config.ConfigHost) {
	ui.ClipboardSetText(host.ShareURL())
	ui.ShowToast(w, T().LinkCopied, ui.ToastSuccess)
}

// shortShareURL shows the link like u00.io does: "u00.io/ch/1a2b3c...9f8e7d"
func shortShareURL(host config.ConfigHost) string {
	s := strings.TrimPrefix(host.ShareURL(), "https://")
	if len(s) > 25 {
		s = s[:12] + "..." + s[len(s)-10:]
	}
	return s
}
