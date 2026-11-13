package forms

import (
	"image/color"

	"github.com/ipoluianov/nui/ui"
)

// A style is the whole look of the application: the palette of the widgets,
// the colors of the host states and how the hosts are drawn. It is picked
// in the toolbar (stylepicker.go) and kept in config.Settings.Theme by ID.

// viewKind is how the list of hosts is shown
type viewKind int

const (
	viewTable viewKind = iota // the table with all the columns
	viewTiles                 // a grid of cards with a big ping time
	viewLanes                 // a heartbeat lane per host
)

// trendKind is how the small charts of the last minutes are drawn
type trendKind int

const (
	trendLine trendKind = iota
	trendArea           // a line over a fading fill
	trendGlow           // an area with a glowing line
	trendBars           // a bar per point
)

// statusKind is how the state of a host is written in the table
type statusKind int

const (
	statusText  statusKind = iota // a colored word
	statusPill                    // a word on a tinted rounded plate
	statusSolid                   // a dark word on a plate of the state color
)

// dotKind is the shape of the mark of the host state before its name
type dotKind int

const (
	dotCircle dotKind = iota
	dotSquare
	dotDiamond
)

type Style struct {
	ID string
	// name is the name in the style picker; nil - Name
	name func() string
	Name string

	Dark    bool
	Palette ui.ThemePalette

	// The host states: replies, replies slowly, does not reply, not pinged yet
	OK, Slow, Down, Idle color.RGBA
	// Muted is the text that is there but does not matter much: "-", "OK"
	Muted color.RGBA
	// Trend is the line of the small charts and of the history chart;
	// Trend2 the other end of the history chart's gradient
	Trend, Trend2 color.RGBA
	// Icon is the color of the toolbar icons; zero - as they are drawn
	Icon color.RGBA

	// Mono shows the whole interface in the monospaced font
	Mono bool
	// FontSize of the interface; 0 - the default
	FontSize float64
	// RowExtra makes the table rows taller than the theme's
	RowExtra int
	// UpperHeaders writes the column names in capitals
	UpperHeaders bool
	// HeatTime colors the ping time by how fast it is
	HeatTime bool

	TrendKind  trendKind
	StatusKind statusKind
	DotKind    dotKind
	// Glow surrounds the state marks with a soft light
	Glow bool

	View viewKind
	// TileMinWidth is the narrowest a tile gets: the wider, the fewer columns
	TileMinWidth int
	// TileCut cuts two corners of the tiles, TileSolid fills a failed host's
	// tile with the state color, TileFaces marks the state with a face
	TileCut, TileSolid, TileFaces bool
	// TileRadius rounds the tiles and the cards
	TileRadius int
	// TrendByState draws the small charts in the color of the host state
	TrendByState bool

	// KPIs shows the cards with the totals above the hosts
	KPIs bool
	// ChartFill is the opacity of the fill under the history chart; 0 - none
	ChartFill float64
	// ChartLine is the width of the history chart's line
	ChartLine float64
}

// DisplayName is the name of the style in the current language
func (s *Style) DisplayName() string {
	if s.name != nil {
		return s.name()
	}
	return s.Name
}

func hex(s string) color.RGBA { return ui.ColorFromHex(s) }

// palette makes the palette of a style from the few colors it picks: the
// window, the background of lists and inputs, the text, the accent with
// the text on it, and the frames. The rest is mixed from them.
func palette(dark bool, window, base, text, accent, accentText, border string) ui.ThemePalette {
	w, b, t, a := hex(window), hex(base), hex(text), hex(accent)
	p := ui.ThemePalette{
		Window:          w,
		WindowText:      t,
		Base:            b,
		Text:            t,
		Button:          ui.MixColors(w, t, 0.08),
		ButtonText:      t,
		Border:          hex(border),
		Divider:         ui.MixColors(b, t, 0.09),
		Highlight:       a,
		HighlightedText: hex(accentText),
		Selection:       ui.MixColors(b, a, 0.22),
		DisabledText:    ui.MixColors(t, w, 0.55),
		PlaceholderText: ui.MixColors(t, w, 0.45),
		PopupBase:       b,
		ToolTipBase:     b,
		ToolTipText:     t,
		Link:            a,
	}
	if dark {
		p.PopupBase = ui.MixColors(w, t, 0.06)
		p.ToolTipBase = p.PopupBase
	}
	return p
}

// withStates sets the state colors of the style and of its palette
func (s Style) withStates(ok, slow, down, idle, muted string) Style {
	s.OK, s.Slow, s.Down, s.Idle, s.Muted = hex(ok), hex(slow), hex(down), hex(idle), hex(muted)
	s.Palette.Success, s.Palette.Warning, s.Palette.Error = s.OK, s.Slow, s.Down
	return s
}

// The IDs of the classic themes are those the settings had before the styles
const (
	styleClassicDark  = ""
	styleClassicLight = "light"
)

// styles are offered in this order
var styles = []*Style{
	classicDark(),
	classicLight(),
	ptr(Style{ID: "neon", Name: "Neon Midnight", Dark: true,
		Palette: palette(true, "#0b1022", "#10172e", "#e6e9ff", "#22d3ee", "#0b1022", "#26305a"),
		Trend:   hex("#22d3ee"), Trend2: hex("#a78bfa"), Icon: hex("#a9b3e0"),
		TrendKind: trendGlow, StatusKind: statusPill, Glow: true, HeatTime: true, UpperHeaders: true,
		RowExtra: 14, TileRadius: 10, ChartFill: 0.4, ChartLine: 1.4,
	}.withStates("#34f5a4", "#fbbf24", "#ff4d7d", "#5d6690", "#7d88b5")),
	ptr(Style{ID: "heartbeat", Name: "Heartbeat", Dark: true,
		Palette: palette(true, "#04090a", "#060d0e", "#d8fff9", "#00ffc8", "#04090a", "#12302f"),
		Trend:   hex("#19ff9c"), Trend2: hex("#00b3ff"), Icon: hex("#8fd9d0"),
		TrendKind: trendGlow, StatusKind: statusPill, Glow: true, View: viewLanes,
		RowExtra: 12, TileRadius: 8, ChartFill: 0.3, ChartLine: 1.4,
	}.withStates("#19ff9c", "#ffd400", "#ff3355", "#3a5a58", "#4d8b88")),
	ptr(Style{ID: "hud", Name: "Gamer HUD", Dark: true,
		Palette: palette(true, "#0c0c12", "#14141d", "#f2f2ff", "#f7ff00", "#0c0c12", "#2a2a40"),
		Trend:   hex("#00f0ff"), Trend2: hex("#f7ff00"), Icon: hex("#c9c9e8"),
		TrendKind: trendGlow, StatusKind: statusPill, Glow: true, UpperHeaders: true, View: viewTiles,
		TileMinWidth: 230, TileCut: true, RowExtra: 12, ChartFill: 0.3, ChartLine: 1.4,
	}.withStates("#39ff88", "#ffb300", "#ff2a6d", "#55556e", "#7f7f9c")),
	ptr(Style{ID: "noc", Name: "NOC Wall", Dark: true,
		Palette: palette(true, "#0a0d12", "#10141b", "#f4fff8", "#ffffff", "#0a0d12", "#252d39"),
		Trend:   hex("#4ade80"), Trend2: hex("#4ade80"), Icon: hex("#c7d0db"),
		TrendKind: trendArea, StatusKind: statusPill, View: viewTiles, TrendByState: true,
		TileMinWidth: 260, TileSolid: true, TileRadius: 10, RowExtra: 12, ChartFill: 0.35, ChartLine: 1.4,
	}.withStates("#4ade80", "#fbbf24", "#e11d48", "#5b6472", "#9aa3ad")),
	ptr(Style{ID: "command", Name: "Command Center", Dark: true,
		Palette: palette(true, "#101312", "#161a19", "#e8efec", "#a3e635", "#101312", "#2a312e"),
		Trend:   hex("#a3e635"), Trend2: hex("#22d3ee"), Icon: hex("#b9c4c0"),
		TrendKind: trendArea, StatusKind: statusPill, Glow: true, HeatTime: true, UpperHeaders: true,
		KPIs: true, RowExtra: 12, TileRadius: 10, ChartFill: 0.35, ChartLine: 1.4,
	}.withStates("#4ade80", "#facc15", "#f87171", "#59625f", "#7c8783")),
	ptr(Style{ID: "arctic", Name: "Arctic", Dark: true,
		Palette: palette(true, "#2e3440", "#3b4252", "#eceff4", "#88c0d0", "#2e3440", "#4c566a"),
		Trend:   hex("#88c0d0"), Trend2: hex("#b48ead"), Icon: hex("#d8dee9"),
		TrendKind: trendArea, StatusKind: statusPill, KPIs: true, UpperHeaders: true,
		RowExtra: 12, TileRadius: 10, ChartFill: 0.35, ChartLine: 1.4,
	}.withStates("#a3be8c", "#ebcb8b", "#bf616a", "#7b869c", "#a7b1c5")),
	ptr(Style{ID: "mocha", Name: "Pastel Mocha", Dark: true,
		Palette: palette(true, "#1e1e2e", "#24273a", "#cdd6f4", "#cba6f7", "#1e1e2e", "#45475a"),
		Trend:   hex("#cba6f7"), Trend2: hex("#89b4fa"), Icon: hex("#bac2de"),
		TrendKind: trendArea, StatusKind: statusPill, HeatTime: true,
		RowExtra: 12, TileRadius: 14, ChartFill: 0.4, ChartLine: 1.5,
	}.withStates("#a6e3a1", "#f9e2af", "#f38ba8", "#6c7086", "#9399b2")),
	ptr(Style{ID: "synthwave", Name: "Synthwave", Dark: true,
		Palette: palette(true, "#160a26", "#1d0d33", "#ffe9fb", "#ff3cac", "#1a0a2a", "#3d1d5c"),
		Trend:   hex("#ff3cac"), Trend2: hex("#ffb340"), Icon: hex("#f3b8ff"),
		TrendKind: trendGlow, StatusKind: statusPill, Glow: true, UpperHeaders: true,
		RowExtra: 12, TileRadius: 10, ChartFill: 0.5, ChartLine: 1.5,
	}.withStates("#2bffd5", "#ffd23f", "#ff2e63", "#7a5a8c", "#b783c9")),
	ptr(Style{ID: "phosphor", Name: "Phosphor", Dark: true, Mono: true,
		Palette: palette(true, "#041008", "#05130a", "#8dffb3", "#4dff88", "#041008", "#1a4a2a"),
		Trend:   hex("#4dff88"), Trend2: hex("#b6ff4d"), Icon: hex("#6cff9c"),
		TrendKind: trendBars, StatusKind: statusText, DotKind: dotSquare, Glow: true, UpperHeaders: true,
		RowExtra: 10, ChartFill: 0.45, ChartLine: 1.2,
	}.withStates("#4dff88", "#ffcc33", "#ff3b3b", "#2f6e44", "#3f9a5c")),
	trading(),
	ptr(Style{ID: "warm", Name: "Warm Retro", Dark: true, Mono: true,
		Palette: palette(true, "#282828", "#32302f", "#ebdbb2", "#fabd2f", "#282828", "#504945"),
		Trend:   hex("#83a598"), Trend2: hex("#fabd2f"), Icon: hex("#d5c4a1"),
		TrendKind: trendBars, StatusKind: statusPill, HeatTime: true,
		RowExtra: 10, TileRadius: 6, ChartFill: 0.4, ChartLine: 1.4,
	}.withStates("#b8bb26", "#fabd2f", "#fb4934", "#7c6f64", "#a89984")),
	ptr(Style{ID: "blueprint", Name: "Blueprint", Dark: true, Mono: true,
		Palette: palette(true, "#0f3f85", "#0d3a7c", "#ffffff", "#ffffff", "#0f3f85", "#6f93c8"),
		Trend:   hex("#ffffff"), Trend2: hex("#9cd3ff"), Icon: hex("#ffffff"),
		TrendKind: trendLine, StatusKind: statusText, DotKind: dotDiamond, UpperHeaders: true,
		RowExtra: 10, ChartFill: 0.12, ChartLine: 1.1,
	}.withStates("#8dffc0", "#ffe27a", "#ff8f8f", "#7fa8dc", "#9cc3f5")),
	ptr(Style{ID: "clean", Name: "Clean Light",
		Palette: palette(false, "#f7f8fb", "#ffffff", "#171a2b", "#4f46e5", "#ffffff", "#dfe2ec"),
		Trend:   hex("#4f46e5"), Trend2: hex("#06b6d4"), Icon: hex("#4b5170"),
		TrendKind: trendArea, StatusKind: statusPill, UpperHeaders: true,
		RowExtra: 14, TileRadius: 12, ChartFill: 0.28, ChartLine: 1.6,
	}.withStates("#10b981", "#f59e0b", "#ef4444", "#a3a8ba", "#7b8199")),
	ptr(Style{ID: "native", Name: "Desktop",
		Palette: palette(false, "#f6f5f7", "#ffffff", "#1d1d1f", "#0a7aff", "#ffffff", "#d9d8de"),
		Trend:   hex("#0a7aff"), Trend2: hex("#5ac8fa"), Icon: hex("#515154"),
		TrendKind: trendArea, StatusKind: statusPill,
		RowExtra: 10, TileRadius: 8, ChartFill: 0.25, ChartLine: 1.5,
	}.withStates("#28a745", "#f29d00", "#ff3b30", "#a1a1a6", "#86868b")),
	ptr(Style{ID: "bubblegum", Name: "Bubblegum",
		Palette: palette(false, "#fff7fb", "#ffffff", "#4a3350", "#ff6fae", "#ffffff", "#ffd6e8"),
		Trend:   hex("#ff6fae"), Trend2: hex("#a58bff"), Icon: hex("#8d6b96"),
		TrendKind: trendArea, StatusKind: statusPill, View: viewTiles, TileFaces: true, TrendByState: true,
		TileMinWidth: 220, TileRadius: 18, RowExtra: 14, ChartFill: 0.45, ChartLine: 2,
	}.withStates("#2fcf9a", "#ffa63d", "#ff5d8f", "#c9b3c9", "#a98ba8")),
	highContrast(),
}

func ptr(s Style) *Style { return &s }

func classicDark() *Style {
	s := Style{ID: styleClassicDark, name: func() string { return T().ThemeDark }, Dark: true,
		Palette: ui.DarkPalette, Trend: hex("#4c9aff"), Trend2: hex("#4c9aff"),
		RowExtra: 8, TileRadius: 6, ChartLine: 1.5,
	}.withStates("#1ebd1e", "#d9b72b", "#e6660a", "#888888", "#7d7d7d")
	s.Palette.Link = hex("#3d8bf2")
	return &s
}

func classicLight() *Style {
	s := Style{ID: styleClassicLight, name: func() string { return T().ThemeLight },
		Palette: ui.LightPalette, Trend: hex("#1d6fb5"), Trend2: hex("#1d6fb5"), Icon: iconColorLight,
		RowExtra: 8, TileRadius: 6, ChartLine: 1.5,
	}.withStates("#178a17", "#a88400", "#d05500", "#7a7a7a", "#8c8c8c")
	s.Palette.Link = hex("#1d6fb5")
	return &s
}

// trading looks like a market terminal: amber on black
func trading() *Style {
	s := Style{ID: "terminal", Name: "Trading Terminal", Dark: true, Mono: true,
		Palette: palette(true, "#000000", "#000000", "#ffb84d", "#ff9e1b", "#000000", "#3a2a0c"),
		Trend:   hex("#ff9e1b"), Trend2: hex("#ffd27a"), Icon: hex("#ff9e1b"),
		TrendKind: trendBars, StatusKind: statusText, HeatTime: true, UpperHeaders: true,
		RowExtra: 6, ChartFill: 0.8, ChartLine: 1.2,
	}.withStates("#00e676", "#ffea00", "#ff1744", "#6b4a12", "#a36a14")
	s.Palette.Selection = hex("#2b1d05")
	return &s
}

// highContrast is for poor sight and projectors: big type, pure colors,
// the states told apart by the shape of their mark too
func highContrast() *Style {
	s := Style{ID: "contrast", Name: "High Contrast", Dark: true,
		Palette: palette(true, "#000000", "#000000", "#ffffff", "#ffff00", "#000000", "#ffffff"),
		Trend:   hex("#ffff00"), Trend2: hex("#ffff00"), Icon: hex("#ffffff"),
		TrendKind: trendLine, StatusKind: statusSolid, FontSize: 17, RowExtra: 14,
		TileRadius: 0, ChartLine: 2,
	}.withStates("#00ff66", "#ffd000", "#ff5c5c", "#bdbdbd", "#d0d0d0")
	s.Palette.Selection = hex("#003b8e")
	s.Palette.Divider = hex("#5a5a5a")
	return &s
}

// currentStyle is the style applied last
var currentStyle = styles[0]

// styleByID returns the style with the ID, the classic dark one for an unknown ID
func styleByID(id string) *Style {
	for _, s := range styles {
		if s.ID == id {
			return s
		}
	}
	return styles[0]
}

// ApplyTheme switches the application to the style of the settings and
// repaints the open windows. The colors of the table follow on its next
// update; the colors and icons the application sets itself follow
// ui.OnThemeChanged.
func ApplyTheme(id string) {
	s := styleByID(id)
	currentStyle = s
	family := ui.FontFamilySans
	if s.Mono {
		family = ui.FontFamilyMono
	}
	size := s.FontSize
	if size == 0 {
		size = ui.DefaultFontSize
	}
	if ui.ThemeFontFamily() != family || ui.ThemeFontSize() != size {
		ui.Theme["fontFamily"] = family
		ui.ApplyBaseFontSize(size)
	}
	ui.ApplyPalette(s.Palette, s.Dark)
}

// styleColor is a color of the application taken from the current style
type styleColor func(s *Style) color.RGBA

func (f styleColor) get() color.RGBA { return f(currentStyle) }

var (
	colorLink          = styleColor(func(s *Style) color.RGBA { return s.Palette.Link })
	colorSecondaryText = styleColor(func(s *Style) color.RGBA { return s.Muted })
	colorMuted         = styleColor(func(s *Style) color.RGBA { return s.Muted })
	colorTrend         = styleColor(func(s *Style) color.RGBA { return s.Trend })
	colorNotPinged     = styleColor(func(s *Style) color.RGBA { return s.Idle })
	colorOK            = styleColor(func(s *Style) color.RGBA { return s.OK })
	colorSlow          = styleColor(func(s *Style) color.RGBA { return s.Slow })
	colorFailed        = styleColor(func(s *Style) color.RGBA { return s.Down })
)
