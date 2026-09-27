package forms

import (
	"fmt"
	"image/color"
	"math"
	"net/netip"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/u00io/nuiforms/ui"
)

var lastCreatedLeftWidget *LeftWidget

// statsPeriod is the period the Loss, Avg and Max columns are computed over
const statsPeriod = 5 * time.Minute

const (
	colName = iota
	colIP
	colTime
	colLoss
	colAvg
	colMax
	colDetails
	colCount
)

var hostColumns = [colCount]struct {
	name  string
	width int
}{
	{"Name", 150},
	{"IP", 140},
	{"Time ms", 90},
	{"Loss %", 90},
	{"Avg ms", 90},
	{"Max ms", 90},
	{"Details", 200},
}

type LeftWidget struct {
	ui.Widget

	lvItems *ui.Table
	// Shown instead of the table while there are no hosts
	emptyHint *ui.Panel

	// Column the rows are sorted by (-1 - config order). The rows are sorted
	// when a header is clicked, not on every update, so they do not jump around.
	sortColumn int
	sortDesc   bool
}

func NewLeftWidget() *LeftWidget {
	var c LeftWidget
	c.InitWidget()
	c.lvItems = ui.NewTable()
	c.AddWidget(0, 0, c.lvItems)
	c.emptyHint = newEmptyHint()
	c.emptyHint.SetVisible(false)
	c.AddWidget(1, 0, c.emptyHint)
	//c.SetMinWidth(700)
	//c.SetMaxWidth(700)

	c.sortColumn = -1
	c.lvItems.SetSelectingRows(true)
	c.lvItems.SetColumnCount(colCount)
	for i, col := range hostColumns {
		c.lvItems.SetColumnWidth(i, col.width)
	}
	c.updateColumnNames()
	c.lvItems.SetOnColumnClick(c.onColumnClick)

	c.lvItems.SetMultiselect(true)

	// Right click selects the row under the mouse (unless it is already selected) and shows the menu
	menu := ui.NewContextMenu(c.lvItems)
	menu.AddItem("Edit... (E)", func() { lastCreatedTopWidget.onBtnEditItem() }).SetImage(loadIcon("edit-16"))
	menu.AddItem("Remove (Del)", func() { lastCreatedTopWidget.onBtnRemoveItem() }).SetImage(loadIcon("remove-16"))
	c.lvItems.SetContextMenu(menu)

	// Double click edits the host, like Enter and E
	c.lvItems.SetOnCellMouseDblClick(func() {
		if ev, ok := ui.CurrentEvent().Parameter.(*ui.EventTableCellMouseDblClick); ok {
			ev.Processed = true
		}
		lastCreatedTopWidget.onBtnEditItem()
	})

	c.loadHosts()

	c.AddTimer(200, c.timerUpdate)

	lastCreatedLeftWidget = &c

	c.lvItems.Focus()

	c.SetPanelPadding(0)

	return &c
}

func (c *LeftWidget) FocusTable() {
	c.lvItems.Focus()
}

// FullRestart restarts all the hosts, e.g. after another config is opened
func (c *LeftWidget) FullRestart() {
	system.Get().Stop()
	c.loadHosts()
	system.Get().Start()
}

// ApplyHostsChange shows the hosts of the config after they were added, edited or removed.
// Only the changed hosts are started, stopped or restarted, the others keep running.
// Afterwards the host selectID is selected, or if there is none, the row selectRow (clamped to the table).
func (c *LeftWidget) ApplyHostsChange(selectID string, selectRow int) {
	c.loadHosts()
	system.Get().SyncHosts()
	c.timerUpdate()
	if row := c.RowOfHost(selectID); row >= 0 {
		selectRow = row
	}
	selectRow = min(selectRow, c.lvItems.RowCount()-1)
	if selectRow >= 0 {
		c.lvItems.SetCurrentCell2(selectRow, 0)
	}
}

// RowOfHost returns the table row of the host, -1 if there is none
func (c *LeftWidget) RowOfHost(id string) int {
	for row := 0; row < c.lvItems.RowCount(); row++ {
		if host, ok := c.lvItems.GetCellData2(row, 0).(*config.ConfigHost); ok && host.ID == id {
			return row
		}
	}
	return -1
}

// loadHosts fills the table with the hosts of the config in the sort order;
// the state columns are filled by timerUpdate
func (c *LeftWidget) loadHosts() {
	hosts := append([]*config.ConfigHost(nil), config.Get().Hosts...)
	if c.sortColumn >= 0 {
		rows := make(map[string]hostRow, len(hosts))
		for _, h := range hosts {
			rows[h.ID] = hostRowOf(h)
		}
		sort.SliceStable(hosts, func(i, j int) bool {
			return lessHostRows(rows[hosts[i].ID], rows[hosts[j].ID], c.sortColumn, c.sortDesc)
		})
	}

	c.lvItems.SetRowCount(len(hosts))
	for i, host := range hosts {
		c.lvItems.SetCellData2(i, 0, host)
		c.lvItems.SetCellText2(i, 0, hostDisplayName(host))
	}

	empty := len(hosts) == 0
	if c.emptyHint.IsVisible() != empty {
		c.lvItems.SetVisible(!empty)
		c.emptyHint.SetVisible(empty)
		if c.Form() != nil {
			c.Form().UpdateLayout()
		}
	}
}

// newEmptyHint tells what to do with an empty list and links to the docs
func newEmptyHint() *ui.Panel {
	p := ui.NewPanel()
	// The same darker area as the table, so the empty list does not look like a gap
	p.SetAutoFillBackground(true)
	p.SetElevation(-3)
	p.AddWidget(0, 0, ui.NewVSpacer())
	for i, text := range []string{"No hosts yet", "Press A to add a host"} {
		lbl := ui.NewLabel(text)
		lbl.SetTextAlign(ui.HAlignCenter)
		lbl.SetXExpandable(true)
		p.AddWidget(i+1, 0, lbl)
	}
	linkRow := ui.NewPanel()
	linkRow.AddWidget(0, 0, ui.NewHSpacer())
	linkRow.AddWidget(0, 1, newLinkLabel("How it works", func() { openDocs(p, "empty_list") }))
	linkRow.AddWidget(0, 2, ui.NewHSpacer())
	p.AddWidget(3, 0, linkRow)
	p.AddWidget(4, 0, ui.NewVSpacer())
	return p
}

// onColumnClick sorts by the column; a second click reverses the order
func (c *LeftWidget) onColumnClick(col int) {
	if col == c.sortColumn {
		c.sortDesc = !c.sortDesc
	} else {
		c.sortColumn = col
		c.sortDesc = false
	}
	c.updateColumnNames()

	// Keep the current host selected
	currentID := ""
	if host, ok := c.lvItems.GetCellData2(c.lvItems.CurrentRow(), 0).(*config.ConfigHost); ok {
		currentID = host.ID
	}
	c.loadHosts()
	c.timerUpdate()
	if row := c.RowOfHost(currentID); row >= 0 {
		c.lvItems.SetCurrentCell2(row, 0)
	}
}

// updateColumnNames marks the sort column with an arrow
func (c *LeftWidget) updateColumnNames() {
	for i, col := range hostColumns {
		name := col.name
		if i == c.sortColumn {
			if c.sortDesc {
				name += " ▼"
			} else {
				name += " ▲"
			}
		}
		c.lvItems.SetColumnName(i, name)
	}
}

func (c *LeftWidget) GetSelectedHostConfigs() []*config.ConfigHost {
	selectedRows := c.lvItems.SelectedRows()
	if len(selectedRows) == 0 {
		return nil
	}
	hosts := make([]*config.ConfigHost, 0, len(selectedRows))
	for _, row := range selectedRows {
		// The selection may still hold rows removed from the table
		if host, ok := c.lvItems.GetCellData2(row, 0).(*config.ConfigHost); ok && host != nil {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

// hostRow is what a table row shows about a host
type hostRow struct {
	name      string
	processed bool // at least one ping was done
	failed    bool // the last ping failed
	ip        string
	time      time.Duration // of the last ping
	stats     system.HistoryStats
	details   string
}

func hostRowOf(h *config.ConfigHost) hostRow {
	state := system.Get().GetHostLastState(h.ID)
	r := hostRow{
		name:      hostDisplayName(h),
		processed: state.StatOK > 0 || state.StatERR > 0,
		failed:    state.LastError != nil,
		ip:        state.StatIP,
		time:      state.PingTime,
		details:   "OK",
	}
	if history := system.Get().GetHostHistory(h.ID); history != nil {
		r.stats = history.Stats(time.Now().Add(-statsPeriod))
	}
	if r.failed {
		switch state.LastError.Error() {
		case "timeout":
			r.details = "TIMEOUT"
		case "cannot resolve hostname":
			r.details = "CANNOT RESOLVE HOSTNAME"
		default:
			r.details = state.LastError.Error()
		}
	}
	return r
}

// texts returns the cell texts; "-" where there is nothing to show yet
func (r hostRow) texts() [colCount]string {
	t := [colCount]string{colName: r.name}
	for i := colIP; i < colCount; i++ {
		t[i] = "-"
	}
	if r.processed {
		t[colDetails] = r.details
		if !r.failed {
			t[colIP] = r.ip
			t[colTime] = formatMs(r.time)
		}
	}
	switch {
	case r.stats.Sent == 0:
	case r.stats.Lost == 0:
		t[colLoss] = "0"
	case r.stats.Lost == r.stats.Sent:
		t[colLoss] = "100"
	default:
		t[colLoss] = fmt.Sprintf("%.1f", r.stats.LossPercent())
	}
	if r.stats.Sent > r.stats.Lost {
		t[colAvg] = formatMs(r.stats.Avg)
		t[colMax] = formatMs(r.stats.Max)
	}
	return t
}

// formatMs shows short times with a decimal, e.g. 0.4 for a LAN host
func formatMs(d time.Duration) string {
	ms := float64(d) / float64(time.Millisecond)
	if ms < 10 {
		return fmt.Sprintf("%.1f", ms)
	}
	return strconv.Itoa(int(math.Round(ms)))
}

// lessHostRows orders the rows by the column; rows without a value go last in both orders
func lessHostRows(a, b hostRow, col int, desc bool) bool {
	var va, vb float64
	var sa, sb string
	numeric := true
	switch col {
	case colName:
		sa, sb, numeric = strings.ToLower(a.name), strings.ToLower(b.name), false
	case colIP:
		return lessIPs(a.texts()[colIP], b.texts()[colIP], desc)
	case colDetails:
		sa, sb, numeric = a.texts()[colDetails], b.texts()[colDetails], false
	case colTime:
		va, vb = rowValue(a.processed && !a.failed, a.time), rowValue(b.processed && !b.failed, b.time)
	case colLoss:
		va, vb = math.NaN(), math.NaN()
		if a.stats.Sent > 0 {
			va = a.stats.LossPercent()
		}
		if b.stats.Sent > 0 {
			vb = b.stats.LossPercent()
		}
	case colAvg:
		va, vb = rowValue(a.stats.Sent > a.stats.Lost, a.stats.Avg), rowValue(b.stats.Sent > b.stats.Lost, b.stats.Avg)
	case colMax:
		va, vb = rowValue(a.stats.Sent > a.stats.Lost, a.stats.Max), rowValue(b.stats.Sent > b.stats.Lost, b.stats.Max)
	}

	if !numeric {
		if (sa == "-") != (sb == "-") {
			return sb == "-"
		}
		if desc {
			return sa > sb
		}
		return sa < sb
	}
	if math.IsNaN(va) || math.IsNaN(vb) {
		return !math.IsNaN(va) && math.IsNaN(vb)
	}
	if desc {
		return va > vb
	}
	return va < vb
}

// lessIPs compares the addresses as numbers, so 8.8.8.8 comes before 123.123.123.123;
// a text that is not an address ("-") goes last
func lessIPs(a, b string, desc bool) bool {
	ipA, errA := netip.ParseAddr(a)
	ipB, errB := netip.ParseAddr(b)
	if (errA == nil) != (errB == nil) {
		return errA == nil
	}
	if errA != nil {
		return false
	}
	if desc {
		return ipA.Compare(ipB) > 0
	}
	return ipA.Compare(ipB) < 0
}

// rowValue returns the duration for sorting, NaN when the row has no value
func rowValue(has bool, d time.Duration) float64 {
	if !has {
		return math.NaN()
	}
	return float64(d)
}

func (c *LeftWidget) timerUpdate() {
	for row := 0; row < c.lvItems.RowCount(); row++ {
		hostConfig, ok := c.lvItems.GetCellData2(row, 0).(*config.ConfigHost)
		if !ok || hostConfig == nil {
			continue
		}
		r := hostRowOf(hostConfig)
		for i, text := range r.texts() {
			c.lvItems.SetCellText2(row, i, text)
		}

		var col color.Color = ui.ColorFromHex("#888888")
		if r.processed {
			col = ui.ColorFromHex("#1ebd1e")
		}
		if r.failed {
			col = ui.ColorFromHex("#e6660a")
		}
		for i := 0; i < colCount; i++ {
			c.lvItems.SetCellColor(row, i, col)
		}
	}
}
