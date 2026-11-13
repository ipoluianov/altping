package forms

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/ipoluianov/altping/config"
	"github.com/ipoluianov/altping/system"
	"github.com/ipoluianov/nui/ui"
)

// The hosts can be shown as tiles or as heartbeat lanes instead of the
// table (Style.View). The table stays the model: its rows give the order
// and the selection, so sorting, the toolbar, the menu and the chart work
// the same in every view.

// viewHost is a host as a view draws it, refreshed by LeftWidget.timerUpdate
type viewHost struct {
	host *config.ConfigHost
	row  hostRow
}

// hostView draws the hosts of the table as tiles or lanes
type hostView struct {
	ui.Widget
	owner *LeftWidget
	kind  viewKind

	// The rectangles of the hosts by table row, in the content coordinates
	rects []image.Rectangle
	// How far the content is scrolled, and its height
	scroll, contentHeight int
	hover                 int
	// The row a Shift+click selects from
	anchor int
}

const (
	viewGap = 12 // between the tiles and around them
	// The lanes show the pings of this long
	laneWindow = 30 * time.Second
)

func newHostView(owner *LeftWidget, kind viewKind) *hostView {
	v := &hostView{owner: owner, kind: kind, hover: -1, anchor: -1}
	v.InitWidget()
	v.SetXExpandable(true)
	v.SetYExpandable(true)
	v.SetCanBeFocused(true)
	v.SetOnPaint(v.paint)
	v.SetOnMouseDown(v.mouseDown)
	v.SetOnMouseDblClick(func(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
		if button == ui.MouseButtonLeft && v.rowAt(x, y) >= 0 {
			lastCreatedTopWidget.onBtnEditItem()
			return true
		}
		return false
	})
	v.SetOnMouseMove(func(x, y int, mods ui.KeyModifiers) bool {
		if row := v.rowAt(x, y); row != v.hover {
			v.hover = row
			v.Update()
		}
		return false
	})
	v.SetOnMouseLeave(func() {
		v.hover = -1
		v.Update()
	})
	v.SetOnMouseWheel(func(deltaX, deltaY int) bool {
		v.scrollTo(v.scroll - deltaY*48)
		return true
	})
	v.OnKeyDown.Connect(v.keyDown)
	v.OnResize.Connect(func(ui.Size) { v.scrollTo(v.scroll) })
	if kind == viewLanes {
		// The lanes run like a monitor: they move on between the updates
		v.AddTimer(100, func() {
			if v.IsVisible() {
				v.Update()
			}
		})
	}
	return v
}

func (v *hostView) hosts() []viewHost {
	return v.owner.viewHosts
}

// rowAt returns the table row of the host at the point, -1 for none
func (v *hostView) rowAt(x, y int) int {
	pt := image.Pt(x, y+v.scroll)
	for row, r := range v.rects {
		if pt.In(r) {
			return row
		}
	}
	return -1
}

func (v *hostView) scrollTo(scroll int) {
	scroll = min(scroll, v.contentHeight-v.Height())
	scroll = max(scroll, 0)
	if scroll != v.scroll {
		v.scroll = scroll
		v.Update()
	}
}

// scrollToRow scrolls so the host of the row is seen whole
func (v *hostView) scrollToRow(row int) {
	if row < 0 || row >= len(v.rects) {
		return
	}
	r := v.rects[row]
	if r.Min.Y-viewGap < v.scroll {
		v.scrollTo(r.Min.Y - viewGap)
	} else if r.Max.Y+viewGap > v.scroll+v.Height() {
		v.scrollTo(r.Max.Y + viewGap - v.Height())
	}
}

// selectRow selects the host of the row as a click with the modifiers does:
// alone, added or removed with Ctrl, a range from the last click with Shift
func (v *hostView) selectRow(row int, mods ui.KeyModifiers) {
	table := v.owner.lvItems
	switch {
	case mods.Shift && v.anchor >= 0:
		var rows []int
		for r := min(v.anchor, row); r <= max(v.anchor, row); r++ {
			rows = append(rows, r)
		}
		table.SetSelectedRows(rows)
	case mods.Ctrl || mods.Cmd:
		var rows []int
		removed := false
		for _, r := range table.SelectedRows() {
			if r == row {
				removed = true
				continue
			}
			rows = append(rows, r)
		}
		if !removed {
			rows = append(rows, row)
		}
		if len(rows) > 0 {
			table.SetSelectedRows(rows)
		}
		v.anchor = row
	default:
		table.SetSelectedRows([]int{row})
		v.anchor = row
	}
	v.scrollToRow(row)
	v.Update()
}

func (v *hostView) mouseDown(button ui.MouseButton, x, y int, mods ui.KeyModifiers) bool {
	v.Focus()
	row := v.rowAt(x, y)
	if row < 0 {
		return true
	}
	switch button {
	case ui.MouseButtonLeft:
		v.selectRow(row, mods)
	case ui.MouseButtonRight:
		// As in the table: the menu is for the host under the mouse,
		// unless it is one of the selected
		if !v.owner.lvItems.IsRowSelected(row) {
			v.selectRow(row, ui.KeyModifiers{})
		}
		fx, fy := v.RectClientAreaOnWindow()
		v.Form().ShowContextMenu(v.owner.menu, fx+x, fy+y)
	}
	return true
}

// keyDown moves the selection with the arrows: across the tiles by their
// columns, along the lanes up and down
func (v *hostView) keyDown(e *ui.KeyEvent) {
	n := len(v.hosts())
	if n == 0 {
		return
	}
	cur := v.owner.lvItems.CurrentRow()
	cols := v.columns()
	next := cur
	switch e.Key {
	case ui.KeyArrowLeft:
		next = cur - 1
	case ui.KeyArrowRight:
		next = cur + 1
	case ui.KeyArrowUp:
		next = cur - cols
	case ui.KeyArrowDown:
		next = cur + cols
	case ui.KeyHome:
		next = 0
	case ui.KeyEnd:
		next = n - 1
	case ui.KeyA:
		if e.Mods.Ctrl || e.Mods.Cmd {
			rows := make([]int, n)
			for i := range rows {
				rows[i] = i
			}
			v.owner.lvItems.SetSelectedRows(rows)
			e.Handled = true
		}
		return
	default:
		return
	}
	if cur < 0 {
		next = 0
	}
	next = min(max(next, 0), n-1)
	v.selectRow(next, e.Mods)
	e.Handled = true
}

// columns is how many hosts are side by side
func (v *hostView) columns() int {
	if v.kind != viewTiles {
		return 1
	}
	minW := max(currentStyle.TileMinWidth, 180)
	return max(1, (v.Width()-viewGap)/(minW+viewGap))
}

func (v *hostView) paint(cnv *ui.Canvas) {
	st := currentStyle
	cnv.FillRect(0, 0, v.Width(), v.Height(), st.Palette.Window)
	switch v.kind {
	case viewTiles:
		v.layoutTiles()
	default:
		v.layoutLanes()
	}
	v.scroll = max(min(v.scroll, v.contentHeight-v.Height()), 0)
	for row, h := range v.hosts() {
		if row >= len(v.rects) {
			break
		}
		r := v.rects[row].Sub(image.Pt(0, v.scroll))
		if r.Max.Y < 0 || r.Min.Y > v.Height() {
			continue
		}
		selected := v.owner.lvItems.IsRowSelected(row)
		if v.kind == viewTiles {
			v.drawTile(cnv, r, h, selected, row == v.hover)
		} else {
			v.drawLane(cnv, r, h, selected, row == v.hover)
		}
	}
	// A thin bar tells there is more to scroll
	if v.contentHeight > v.Height() {
		h := v.Height() * v.Height() / v.contentHeight
		y := (v.Height() - h) * v.scroll / max(v.contentHeight-v.Height(), 1)
		cnv.FillRoundedRectAA(v.Width()-6, y+2, 4, h-4, 2, alpha(st.Palette.Text, 70))
	}
}

// Tiles

// tileMetrics are the sizes inside a tile for the font of the style
type tileMetrics struct {
	pad, line, small, big float64
	spark                 int
}

func tileSizes() tileMetrics {
	fs := ui.ThemeFontSize()
	m := tileMetrics{pad: 14, line: fs + 1, small: math.Max(fs-3, 9), big: fs * 2.1, spark: 34}
	if currentStyle.TileSolid {
		m.big, m.spark, m.pad = fs*3.2, 52, 18
	}
	return m
}

func (m tileMetrics) height() int {
	return int(m.pad*2 + m.line*1.4 + m.small*1.5 + m.big*1.3 + float64(m.spark) + m.small*1.8 + 8)
}

func (v *hostView) layoutTiles() {
	n := len(v.hosts())
	cols := v.columns()
	w := (v.Width() - viewGap*(cols+1)) / cols
	m := tileSizes()
	h := m.height()
	rows := (n + cols - 1) / cols
	// Few hosts fill the height, up to a half again as tall
	if rows > 0 {
		fit := (v.Height() - viewGap*(rows+1)) / rows
		if fit > h {
			h = min(fit, h*3/2)
		}
	}
	v.rects = v.rects[:0]
	for i := 0; i < n; i++ {
		x := viewGap + (i%cols)*(w+viewGap)
		y := viewGap + (i/cols)*(h+viewGap)
		v.rects = append(v.rects, image.Rect(x, y, x+w, y+h))
	}
	v.contentHeight = viewGap + rows*(h+viewGap)
}

// stateFace marks the state in the friendly styles
func stateFace(r hostRow) string {
	switch {
	case r.failed:
		return "(×_×)"
	case !r.processed:
		return "(-_-)"
	case r.isSlow():
		return "(◔_◔)"
	}
	return "(◕‿◕)"
}

func (v *hostView) drawTile(cnv *ui.Canvas, rc image.Rectangle, h viewHost, selected, hovered bool) {
	st := currentStyle
	p := st.Palette
	r := h.row
	m := tileSizes()
	state := st.stateColor(r)
	family := ui.ThemeFontFamily()
	x, y, w, ht := rc.Min.X, rc.Min.Y, rc.Dx(), rc.Dy()

	fill := p.Base
	text := p.Text
	muted := st.Muted
	border := alpha(p.Border, 0)
	borderWidth := 1.0
	solidDown := st.TileSolid && r.failed
	switch {
	case solidDown:
		fill = state
		text, muted = color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 200}
	case st.TileSolid && r.isSlow():
		fill = ui.MixColors(p.Base, state, 0.22)
	case st.TileCut && r.failed:
		fill = ui.MixColors(p.Base, state, 0.14)
	case st.TileFaces && r.failed:
		fill = ui.MixColors(p.Base, state, 0.07)
	}
	switch {
	case selected:
		border, borderWidth = p.Highlight, 2
	case hovered:
		border = ui.MixColors(p.Border, p.Highlight, 0.5)
	case !st.TileCut && !solidDown:
		border = p.Border
	}
	if st.TileFaces && !selected && !hovered {
		border = alpha(st.Palette.Border, 255)
	}
	if solidDown && !selected {
		border = alpha(border, 0)
	}
	fillCard(cnv, st, x, y, w, ht, fill, border, borderWidth)
	if st.TileCut {
		// A bar of the state color down the left side
		cnv.FillRect(x+1, y+1, 3, ht-16, state)
	}

	cnv.Save()
	cnv.TranslateAndClip(x, y, w, ht)
	defer cnv.Restore()

	pad := int(m.pad)
	cy := pad + int(m.line*0.7)
	tx := pad
	// The state: a face or a mark, then the name and the bell
	cnv.SetFontFamily(family)
	if st.TileFaces {
		cnv.SetFontFamily(ui.FontFamilySans)
		cnv.SetFontSize(m.line - 1)
		cnv.SetFontBold(true)
		cnv.SetColor(state)
		cnv.SetHAlign(ui.HAlignLeft)
		cnv.SetVAlign(ui.VAlignCenter)
		face := stateFace(r)
		fw, _ := cnv.MeasureTextF(face)
		cnv.DrawText(tx, cy-int(m.line), int(fw)+2, int(m.line*2), face)
		cnv.SetFontBold(false)
		tx += int(fw) + 8
	} else {
		col := state
		if solidDown {
			col = text
		}
		drawStateMark(cnv, st, float64(tx)+4, float64(cy), stateMarkSize, col)
		tx += stateMarkSize + 9
	}
	name := r.name
	if st.UpperHeaders && st.TileCut {
		name = strings.ToUpper(name)
	}
	cnv.SetFontFamily(family)
	cnv.SetFontSize(m.line)
	cnv.SetFontBold(true)
	cnv.SetColor(text)
	cnv.SetHAlign(ui.HAlignLeft)
	cnv.SetVAlign(ui.VAlignCenter)
	nameW := w - tx - pad - 30
	if h.host.Notify {
		nameW -= notifyIconSize + 4
	}
	name = fitText(family, m.line, name, nameW)
	cnv.DrawText(tx, cy-int(m.line), nameW, int(m.line*2), name)
	nw, _ := cnv.MeasureTextF(name)
	cnv.SetFontBold(false)
	if h.host.Notify && v.owner.iconNotify != nil {
		cnv.DrawImage(tx+int(nw)+6, cy-notifyIconSize/2, v.owner.iconNotify)
	}
	sig := state
	if solidDown || r.failed {
		sig = text
	}
	drawSignal(cnv, st, w-pad-18, cy+7, signalBars(r), sig)

	// The address
	y1 := cy + int(m.line*0.7)
	cnv.SetFontSize(m.small)
	cnv.SetColor(muted)
	addr := h.host.Address()
	if r.processed && !r.failed && r.ip != "" && r.ip != addr {
		addr += " · " + r.ip
	}
	cnv.DrawText(tx, y1, w-tx-pad, int(m.small*1.5), fitText(family, m.small, addr, w-tx-pad))

	// The big time, and the state plate
	footerH := int(m.small * 1.8)
	sparkTop := ht - pad - footerH - m.spark
	bigY := y1 + int(m.small*1.5)
	bigH := sparkTop - bigY
	texts := r.texts()
	cnv.SetFontSize(m.big)
	cnv.SetFontBold(true)
	valCol := text
	if !solidDown {
		switch {
		case r.failed:
			valCol = state
		case r.isSlow():
			valCol = state
		}
	}
	cnv.SetColor(valCol)
	value := texts[colTime]
	unit := T().Ms
	if r.failed {
		value, unit = strings.ToUpper(stateText(r)), ""
		if st.TileFaces {
			value = "zzz…"
		}
	} else if !r.processed {
		unit = ""
	}
	value = fitText(family, m.big, value, w-pad*2)
	cnv.DrawText(pad, bigY, w-pad*2, bigH, value)
	vw, _ := cnv.MeasureTextF(value)
	cnv.SetFontBold(false)
	if unit != "" {
		cnv.SetFontSize(m.small)
		cnv.SetColor(muted)
		cnv.DrawText(pad+int(vw)+4, bigY+int(m.small*0.4), 40, bigH, unit)
	}
	if !st.TileSolid && !st.TileCut && r.processed {
		pw := pillWidth(st, stateText(r), family, ui.ThemeFontSize())
		drawPill(cnv, st, w-pad-pw, bigY+bigH/2, stateText(r), state, family, ui.ThemeFontSize())
	}

	// The last minutes
	col := st.Trend
	if st.TrendByState || r.isSlow() {
		col = state
	}
	if solidDown {
		col = text
	}
	step := trendStep(st)
	sw := w - pad*2
	points := v.owner.trendPoints(h.host.ID, sw/step)
	drawTrend(cnv, st, pad, sparkTop, sw, m.spark, points, step, col)

	// Loss, jitter and since along the bottom
	fy := ht - pad - footerH
	cnv.SetFontSize(m.small)
	cnv.SetColor(muted)
	cnv.SetVAlign(ui.VAlignBottom)
	third := sw / 3
	loss := texts[colLoss]
	if loss != "-" {
		loss += "%"
	}
	cols := T().Columns
	cnv.SetHAlign(ui.HAlignLeft)
	cnv.DrawText(pad, fy, third, footerH, fitText(family, m.small, shortColumn(cols.Loss)+" "+loss, third))
	cnv.SetHAlign(ui.HAlignCenter)
	cnv.DrawText(pad+third, fy, third, footerH, fitText(family, m.small, shortColumn(cols.Jitter)+" "+texts[colJitter], third))
	cnv.SetHAlign(ui.HAlignRight)
	cnv.DrawText(pad+third*2, fy, sw-third*2, footerH, texts[colSince])
	cnv.SetVAlign(ui.VAlignCenter)
}

// shortColumn is the name of a column without its unit: "Loss %" is "loss"
func shortColumn(name string) string {
	name = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(name, "%"), T().Ms))
	return strings.ToLower(name)
}

// pillWidth is the width drawPill takes for the text
func pillWidth(s *Style, text, family string, fontSize float64) int {
	size := math.Max(fontSize-3, 9)
	w, _, err := ui.MeasureText(family, size, strings.ToUpper(text))
	if err != nil {
		return 60
	}
	dot := 0
	if s.StatusKind == statusPill {
		dot = 12
	}
	// The plate writes in bold: a little wider than measured
	return w + w/12 + int(size*0.75)*2 + dot
}

// Lanes

func (v *hostView) layoutLanes() {
	n := len(v.hosts())
	h := 64
	if n > 0 {
		fit := v.Height() / n
		h = min(max(fit, 56), 100)
	}
	v.rects = v.rects[:0]
	for i := 0; i < n; i++ {
		v.rects = append(v.rects, image.Rect(0, i*h, v.Width(), (i+1)*h))
	}
	v.contentHeight = n * h
}

func (v *hostView) drawLane(cnv *ui.Canvas, rc image.Rectangle, h viewHost, selected, hovered bool) {
	st := currentStyle
	p := st.Palette
	r := h.row
	state := st.stateColor(r)
	family := ui.ThemeFontFamily()
	fs := ui.ThemeFontSize()
	x, y, w, ht := rc.Min.X, rc.Min.Y, rc.Dx(), rc.Dy()

	bg := p.Base
	switch {
	case r.failed:
		bg = ui.MixColors(p.Base, state, 0.09)
	case selected:
		bg = p.Selection
	case hovered:
		bg = ui.MixColors(p.Base, p.Text, 0.04)
	}
	cnv.FillRect(x, y, w, ht, bg)
	cnv.FillRect(x, y+ht-1, w, 1, p.Divider)
	if selected {
		cnv.FillRect(x, y, 3, ht, p.Highlight)
	}

	leftW := min(max(w/5, 170), 260)
	rightW := int(fs * 11)
	traceX, traceW := x+leftW, w-leftW-rightW-16
	cnv.Save()
	cnv.TranslateAndClip(x, y, w, ht)
	defer cnv.Restore()
	traceX -= x

	// The grid paper of a monitor under the trace
	grid := alpha(st.Trend, 14)
	for gx := traceX; gx < traceX+traceW; gx += 20 {
		cnv.FillRect(gx, 0, 1, ht-1, grid)
	}
	for gy := 10; gy < ht-1; gy += 20 {
		cnv.FillRect(traceX, gy, traceW, 1, grid)
	}

	// The name, the bell and the address
	pad := 16
	cnv.SetFontFamily(family)
	cnv.SetFontSize(fs + 1)
	cnv.SetFontBold(true)
	cnv.SetColor(p.Text)
	cnv.SetHAlign(ui.HAlignLeft)
	cnv.SetVAlign(ui.VAlignBottom)
	nameW := leftW - pad - 8
	if h.host.Notify {
		nameW -= notifyIconSize + 4
	}
	name := fitText(family, fs+1, r.name, nameW)
	cnv.DrawText(pad, 0, nameW, ht/2+2, name)
	nw, _ := cnv.MeasureTextF(name)
	cnv.SetFontBold(false)
	if h.host.Notify && v.owner.iconNotify != nil {
		cnv.DrawImage(pad+int(nw)+6, ht/2-notifyIconSize, v.owner.iconNotify)
	}
	cnv.SetVAlign(ui.VAlignTop)
	cnv.SetFontSize(fs - 2)
	cnv.SetColor(st.Muted)
	cnv.DrawText(pad, ht/2+4, leftW-pad-8, ht/2, fitText(family, fs-2, h.host.Address(), leftW-pad-8))

	v.drawHeartbeat(cnv, traceX, 0, traceW, ht, h, state)

	// The big time and the loss and jitter under it
	texts := r.texts()
	rx := w - rightW - pad
	cnv.SetHAlign(ui.HAlignRight)
	cnv.SetVAlign(ui.VAlignBottom)
	big := math.Min(fs*2.3, float64(ht)*0.42)
	cnv.SetFontSize(big)
	cnv.SetFontBold(true)
	cnv.SetColor(state)
	value, unit := texts[colTime], T().Ms
	if r.failed {
		value, unit = "- - -", ""
	} else if !r.processed {
		unit = ""
	}
	unitW := 0
	if unit != "" {
		uw, _, _ := ui.MeasureText(family, fs-3, unit)
		unitW = uw + 4
	}
	cnv.DrawText(rx, 0, rightW-unitW, ht*62/100, value)
	cnv.SetFontBold(false)
	if unit != "" {
		cnv.SetFontSize(fs - 3)
		cnv.DrawText(rx, 0, rightW, ht*62/100-3, unit)
	}
	cnv.SetVAlign(ui.VAlignTop)
	cnv.SetFontSize(fs - 4)
	cnv.SetColor(alpha(state, 190))
	var sub string
	if r.failed {
		sub = strings.ToUpper(stateText(r)) + " · " + texts[colSince]
	} else {
		cols := T().Columns
		sub = strings.ToUpper(fmt.Sprintf("%s %s%% · %s %s", shortColumn(cols.Loss), texts[colLoss], shortColumn(cols.Jitter), texts[colJitter]))
	}
	cnv.DrawText(rx-40, ht*62/100+3, rightW+40, ht/3, fitText(family, fs-4, sub, rightW+40))
}

// drawHeartbeat draws the pings of the last laneWindow like the trace of a
// heart monitor: a beat at each reply, taller for a slower one; a lost
// ping is a flat stretch with a mark. The trace fades to the past.
func (v *hostView) drawHeartbeat(cnv *ui.Canvas, x, y, w, h int, vh viewHost, col color.RGBA) {
	st := currentStyle
	history := system.Get().GetHostHistory(vh.host.ID)
	if history == nil || w < 40 {
		return
	}
	now := time.Now()
	from := now.Add(-laneWindow)
	type beat struct {
		x    float64
		ms   float64
		lost bool
	}
	var beats []beat
	maxMs := 0.0
	history.Visit(from, now, func(s system.HistorySample) {
		if s.Gap || s.DT.Before(from) {
			return
		}
		bx := float64(x) + float64(w)*float64(s.DT.Sub(from))/float64(laneWindow)
		ms := float64(s.PingTime) / float64(time.Millisecond)
		beats = append(beats, beat{bx, ms, !s.OK})
		if s.OK {
			maxMs = math.Max(maxMs, ms)
		}
	})
	if len(beats) == 0 {
		return
	}
	maxMs = math.Max(maxMs, 0.5)
	base := float64(y) + float64(h)*0.62
	amp := float64(h) * 0.5

	p := ui.NewPath()
	first := math.Max(beats[0].x-12, float64(x))
	p.MoveTo(first, base)
	for _, b := range beats {
		if b.lost {
			p.LineTo(b.x, base)
			continue
		}
		a := amp * (0.35 + 0.65*b.ms/maxMs)
		// A small wave, the spike of the beat and a wider wave after it
		p.LineTo(b.x-9, base)
		p.QuadTo(b.x-7, base-4, b.x-5, base)
		p.LineTo(b.x-2, base)
		p.LineTo(b.x-1, base+3)
		p.LineTo(b.x+1.5, base-a)
		p.LineTo(b.x+4, base+6)
		p.LineTo(b.x+5.5, base)
		p.LineTo(b.x+9, base)
		p.QuadTo(b.x+12, base-7, b.x+15, base)
	}
	last := beats[len(beats)-1].x + 16
	p.LineTo(math.Min(last, float64(x+w)), base)

	fade := ui.LinearGradient{X1: float64(x), Y1: 0, X2: float64(x + w), Y2: 0,
		Stops: []ui.GradientStop{{Offset: 0, Color: alpha(col, 30)}, {Offset: 0.8, Color: alpha(col, 230)}, {Offset: 1, Color: col}}}
	if st.Glow {
		glow := fade
		glow.Stops = []ui.GradientStop{{Offset: 0, Color: alpha(col, 0)}, {Offset: 1, Color: alpha(col, 70)}}
		cnv.StrokePathPaint(p, ui.StrokeStyle{Width: 5, Join: ui.LineJoinRound, Cap: ui.LineCapRound}, glow)
	}
	cnv.StrokePathPaint(p, ui.StrokeStyle{Width: 1.7, Join: ui.LineJoinRound, Cap: ui.LineCapRound}, fade)

	for _, b := range beats {
		if b.lost {
			cnv.FillRect(int(b.x)-1, int(base)-6, 3, 12, st.Down)
		}
	}
	// The head of the trace
	hx := math.Min(last, float64(x+w)) - 1
	cnv.FillRoundedRectAA(int(hx)-4, int(base)-4, 8, 8, 4, alpha(col, 90))
	cnv.FillRoundedRectAA(int(hx)-2, int(base)-2, 4, 4, 2, color.RGBA{255, 255, 255, 255})
}
