package forms

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nui/nui"
	"github.com/u00io/nuiforms/ui"
)

// exportHistory saves the last day of pings of the hosts to a CSV file the user chooses
func exportHistory(hosts []*config.ConfigHost) {
	if len(hosts) == 0 {
		return
	}
	name := "altping-history"
	if len(hosts) == 1 {
		name = "altping-" + fileNamePart(hostDisplayName(hosts[0]))
	}
	name += "-" + time.Now().Format("20060102-1504") + ".csv"

	form := lastCreatedMainWidget.Form()
	opts := nui.SaveFileDialogOptions{
		Title:           "Export history",
		DefaultFileName: name,
		Filters:         []nui.FileDialogFilter{{DisplayName: "CSV files", Patterns: []string{"*.csv"}}},
	}
	// The system dialog is a separate window: it would open below a window kept on top
	onTop := config.GetSettings().AlwaysOnTop
	if onTop {
		form.SetAlwaysOnTop(false)
	}
	form.ShowSaveFileDialog(opts, func(path string, err error) {
		if onTop {
			form.SetAlwaysOnTop(true)
		}
		// No system dialog (e.g. Linux without zenity or kdialog): save to the home folder
		noDialog := errors.Is(err, nui.ErrNoFileDialog)
		if err != nil && !noDialog {
			ui.ShowMessageBox(lastCreatedMainWidget, "Error", err.Error())
			return
		}
		if noDialog {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, name)
		}
		if path == "" {
			return // cancelled
		}
		if err := writeHistoryCSV(path, hosts); err != nil {
			ui.ShowMessageBox(lastCreatedMainWidget, "Error", err.Error())
			return
		}
		if noDialog {
			ui.ShowMessageBox(lastCreatedMainWidget, "Export history", "Saved to "+path)
		}
	})
}

// writeHistoryCSV writes one line per ping: time, host, address, result (ok/lost), ms
func writeHistoryCSV(path string, hosts []*config.ConfigHost) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.Write([]string{"time", "host", "address", "result", "ms"})
	to := time.Now()
	from := to.Add(-24 * time.Hour)
	for _, h := range hosts {
		history := system.Get().GetHostHistory(h.ID)
		if history == nil {
			continue
		}
		name := hostDisplayName(h)
		history.Visit(from, to, func(s system.HistorySample) {
			// Visit also gives the nearest samples outside the range
			if s.Gap || s.DT.Before(from) || s.DT.After(to) {
				return
			}
			result, ms := "lost", ""
			if s.OK {
				result = "ok"
				ms = fmt.Sprintf("%.3f", float64(s.PingTime)/float64(time.Millisecond))
			}
			w.Write([]string{s.DT.Format("2006-01-02 15:04:05.000"), name, h.Address(), result, ms})
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// fileNamePart keeps letters, digits, dots and dashes of a host name for a file name
func fileNamePart(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '.' || r == '-' || r == '_' || ('0' <= r && r <= '9') || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			return r
		}
		return '_'
	}, s)
}
