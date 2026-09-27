package forms

import (
	"fmt"
	"image/color"
	"math"
	"net/netip"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

var lastCreatedLeftWidget *LeftWidget

// statsPeriod is the period the Loss, Min, Avg, Max and Jitter columns are computed over
const statsPeriod = 5 * time.Minute

// The fields a host row can show; the table shows some of them (LeftWidget.columns)
const (
	colName = iota
	colIP
	colTime
	colLoss
	colMin
	colAvg
	colMax
	colJitter
	colSince
	colDetails
	colCount
)

// hostColumns describes the fields in the order of the table columns;
// the optional ones are shown only when turned on in the settings
var hostColumns = [colCount]struct {
	name  string
	width int
	shown func(s config.Settings) bool // nil - always
}{
	colName:    {"Name", 150, nil},
	colIP:      {"IP", 140, nil},
	colTime:    {"Time ms", 90, nil},
	colLoss:    {"Loss %", 90, nil},
	colMin:     {"Min ms", 90, func(s config.Settings) bool { return s.ShowMin }},
	colAvg:     {"Avg ms", 90, nil},
	colMax:     {"Max ms", 90, nil},
	colJitter:  {"Jitter ms", 110, func(s config.Settings) bool { return s.ShowJitter }},
	colSince:   {"Since", 90, func(s config.Settings) bool { return s.ShowSince }},
	colDetails: {"Details", 200, nil},
}

var (
	colorNotPinged = ui.ColorFromHex("#888888")
	colorOK        = ui.ColorFromHex("#1ebd1e")
	colorSlow      = ui.ColorFromHex("#d9b72b")
	colorFailed    = ui.ColorFromHex("#e6660a")
)

type LeftWidget struct {
	ui.Widget

	lvItems *ui.Table
	// Shown instead of the table while there are no hosts
	emptyHint *ui.Panel

	// The fields shown, one per table column
	columns []int

	// Field the rows are sorted by (-1 - config order). The rows are sorted
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
	c.ApplyColumns()
	c.lvItems.SetOnColumnClick(c.onColumnClick)

	c.lvItems.SetMultiselect(true)

	// Right click selects the row under the mouse (unless it is already selected) and shows the menu
	menu := ui.NewContextMenu(c.lvItems)
	menu.AddItem("Edit... (E)", func() { lastCreatedTopWidget.onBtnEditItem() }).SetImage(loadIcon("edit-16"))
	menu.AddItem("Remove (Del)", func() { lastCreatedTopWidget.onBtnRemoveItem() }).SetImage(loadIcon("remove-16"))
	menu.AddSeparator()
	menu.AddItem("Downtime...", func() {
		if hosts := c.GetSelectedHostConfigs(); len(hosts) > 0 {
			c.ShowDialog(NewDowntimeDialog(hosts[0]))
		}
	})
	menu.AddItem("Export history...", func() { exportHistory(c.GetSelectedHostConfigs()) })
	menu.AddSeparator()
	menu.AddItemWithSubmenu("Change selected", c.newChangeMenu())
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

// newChangeMenu changes an option of all the selected hosts at once
func (c *LeftWidget) newChangeMenu() *ui.ContextMenu {
	menu := ui.NewContextMenu(c.lvItems)
	menu.AddItem("Beep on", func() { c.changeSelected(func(h *config.ConfigHost) { h.Notify = true }) })
	menu.AddItem("Beep off", func() { c.changeSelected(func(h *config.ConfigHost) { h.Notify = false }) })
	menu.AddSeparator()
	menu.AddItem("Check TCP port...", func() {
		c.askSelected("Check TCP port", "Port:", func(h *config.ConfigHost) int {
			if _, port := h.Target(); port != "" {
				p, _ := strconv.Atoi(port)
				return p
			}
			return defaultTCPPort
		}, 1, 65535, 1, func(h *config.ConfigHost, v int) {
			h.Hostname, _ = h.Target() // an old "host:port" becomes the host
			h.Port = v
		})
	})
	menu.AddItem("Use ping", func() {
		c.changeSelected(func(h *config.ConfigHost) {
			h.Hostname, _ = h.Target()
			h.Port = 0
		})
	})
	menu.AddSeparator()
	menu.AddItem("Ping every...", func() {
		c.askSelected("Ping every", "Ping every, ms:", func(h *config.ConfigHost) int { return int(h.Interval().Milliseconds()) },
			config.MinIntervalMs, config.MaxIntervalMs, 100, func(h *config.ConfigHost, v int) { h.IntervalMs = storedMs(v, config.DefaultIntervalMs) })
	})
	menu.AddItem("Timeout...", func() {
		c.askSelected("Timeout", "Timeout, ms:", func(h *config.ConfigHost) int { return int(h.Timeout().Milliseconds()) },
			config.MinTimeoutMs, config.MaxTimeoutMs, 100, func(h *config.ConfigHost, v int) { h.TimeoutMs = storedMs(v, config.DefaultTimeoutMs) })
	})
	menu.AddItem("Slow above...", func() {
		c.askSelected("Slow above", "Slow above, ms (0 - off):", func(h *config.ConfigHost) int { return h.SlowMs },
			0, config.MaxSlowMs, 10, func(h *config.ConfigHost, v int) { h.SlowMs = v })
	})
	return menu
}

// storedMs keeps the default out of the config file
func storedMs(v, def int) int {
	if v == def {
		return 0
	}
	return v
}

// askSelected asks for a number, starting from the value of the first selected host,
// and sets it for all the selected hosts
func (c *LeftWidget) askSelected(title, label string, current func(h *config.ConfigHost) int, minValue, maxValue, step int, set func(h *config.ConfigHost, v int)) {
	hosts := c.GetSelectedHostConfigs()
	if len(hosts) == 0 {
		return
	}
	c.ShowDialog(NewNumberDialog(title, label, current(hosts[0]), minValue, maxValue, step, func(v int) {
		c.changeSelected(func(h *config.ConfigHost) { set(h, v) })
	}))
}

// changeSelected changes all the selected hosts, saves the config and restarts
// only the hosts that need it; the selection is kept
func (c *LeftWidget) changeSelected(change func(h *config.ConfigHost)) {
	hosts := c.GetSelectedHostConfigs()
	if len(hosts) == 0 {
		return
	}
	for _, h := range hosts {
		change(h)
	}
	if err := config.Get().Save(); err != nil {
		ui.ShowMessageBox(c, "Error", err.Error())
	}
	c.loadHosts()
	system.Get().SyncHosts()
	c.timerUpdate()
	rows := make([]int, 0, len(hosts))
	for _, h := range hosts {
		rows = append(rows, c.RowOfHost(h.ID))
	}
	c.lvItems.SetSelectedRows(rows)
	c.FocusTable()
}

// ApplyColumns shows the columns turned on in the settings
func (c *LeftWidget) ApplyColumns() {
	settings := config.GetSettings()
	c.columns = c.columns[:0]
	for field, col := range hostColumns {
		if col.shown == nil || col.shown(settings) {
			c.columns = append(c.columns, field)
		}
	}
	if !slices.Contains(c.columns, c.sortColumn) {
		c.sortColumn = -1
	}
	c.lvItems.SetColumnCount(len(c.columns))
	for i, field := range c.columns {
		c.lvItems.SetColumnWidth(i, hostColumns[field].width)
	}
	c.updateColumnNames()
	c.timerUpdate()
}

// onColumnClick sorts by the column; a second click reverses the order
func (c *LeftWidget) onColumnClick(index int) {
	if index < 0 || index >= len(c.columns) {
		return
	}
	col := c.columns[index]
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
	for i, field := range c.columns {
		name := hostColumns[field].name
		if field == c.sortColumn {
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
	since     time.Time // when the result last changed; zero - unknown
	slowMs    int       // the host's limit for "slow"; 0 - off
	details   string
}

func hostRowOf(h *config.ConfigHost) hostRow {
	state := system.Get().GetHostLastState(h.ID)
	r := hostRow{
		name:      hostDisplayName(h),
		slowMs:    h.SlowMs,
		processed: state.StatOK > 0 || state.StatERR > 0,
		failed:    state.LastError != nil,
		ip:        state.StatIP,
		time:      state.PingTime,
		details:   "OK",
	}
	if history := system.Get().GetHostHistory(h.ID); history != nil {
		r.stats = history.Stats(time.Now().Add(-statsPeriod))
		r.since = history.LastChange()
	}
	if r.failed {
		switch state.LastError.Error() {
		case "timeout":
			r.details = "TIMEOUT"
		case "cannot resolve hostname":
			r.details = "CANNOT RESOLVE HOSTNAME"
		case "port closed":
			r.details = "PORT CLOSED"
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
		t[colMin] = formatMs(r.stats.Min)
		t[colAvg] = formatMs(r.stats.Avg)
		t[colMax] = formatMs(r.stats.Max)
		t[colJitter] = formatMs(r.stats.Jitter)
	}
	if !r.since.IsZero() {
		t[colSince] = formatSince(time.Since(r.since))
	}
	return t
}

// formatSince shows how long ago in the largest two units: 45s, 12m 5s, 3h 20m, 2d 4h
func formatSince(d time.Duration) string {
	sec := int(d.Seconds())
	switch {
	case sec < 60:
		return fmt.Sprintf("%ds", sec)
	case sec < 3600:
		return fmt.Sprintf("%dm %ds", sec/60, sec%60)
	case sec < 86400:
		return fmt.Sprintf("%dh %dm", sec/3600, sec%3600/60)
	}
	return fmt.Sprintf("%dd %dh", sec/86400, sec%86400/3600)
}

// color shows the state of the host: not pinged yet, replies, replies slowly, failed
func (r hostRow) color(settings config.Settings) color.Color {
	switch {
	case r.failed:
		return colorFailed
	case !r.processed:
		return colorNotPinged
	case r.slowMs > 0 && r.stats.Sent > r.stats.Lost && r.stats.Avg > time.Duration(r.slowMs)*time.Millisecond:
		return colorSlow
	}
	return colorOK
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
	case colMin:
		va, vb = rowValue(a.stats.Sent > a.stats.Lost, a.stats.Min), rowValue(b.stats.Sent > b.stats.Lost, b.stats.Min)
	case colMax:
		va, vb = rowValue(a.stats.Sent > a.stats.Lost, a.stats.Max), rowValue(b.stats.Sent > b.stats.Lost, b.stats.Max)
	case colJitter:
		va, vb = rowValue(a.stats.Sent > a.stats.Lost, a.stats.Jitter), rowValue(b.stats.Sent > b.stats.Lost, b.stats.Jitter)
	case colSince:
		va, vb = rowValue(!a.since.IsZero(), time.Since(a.since)), rowValue(!b.since.IsZero(), time.Since(b.since))
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
	settings := config.GetSettings()
	for row := 0; row < c.lvItems.RowCount(); row++ {
		hostConfig, ok := c.lvItems.GetCellData2(row, 0).(*config.ConfigHost)
		if !ok || hostConfig == nil {
			continue
		}
		r := hostRowOf(hostConfig)
		texts := r.texts()
		col := r.color(settings)
		for i, field := range c.columns {
			c.lvItems.SetCellText2(row, i, texts[field])
			c.lvItems.SetCellColor(row, i, col)
		}
	}
}
