package forms

import (
	"fmt"
	"strings"

	"github.com/ipoluianov/nui/ui"
	"github.com/ipoluianov/nui/ui/i18n"
)

// Strings are all the texts of the application. A misspelled field is a
// compile error; a field a language leaves empty falls back to English
// (strings_test.go checks that none is left).
type Strings struct {
	Error string

	// Main window
	DownCount       func(n int) string // in the title, seen in the taskbar
	ToolConfigs     string
	ToolAddHost     string
	ToolEditHost    string
	ToolRemoveHosts string
	ToolDetails     string
	ToolStart       string
	ToolStop        string
	ToolTray        string
	ToolShareOpen   string
	ToolStyle       string
	// The cards with the totals of some styles
	KPIMedian string
	KPIUptime string
	TrayShow  string
	TrayQuit  string
	Settings  string
	Help      string
	About     string
	ModeUDP   string
	ModeICMP  string
	// The status bar: how many hosts are in which state
	StatusHosts   string
	StatusOK      string
	StatusSlow    string
	StatusDown    string
	StatusStopped string

	// Table of hosts
	Columns         ColumnStrings
	StateOK         string
	StateTimeout    string
	StateNoResolve  string
	StatePortClosed string
	NoHosts         string
	PressAToAdd     string
	HowItWorks      string

	// Units of the durations like 12m 5s
	Units UnitStrings

	// Context menu of the table
	MenuEdit       string
	MenuRemove     string
	MenuDowntime   string
	MenuExport     string
	MenuChange     string
	BeepOn         string
	BeepOff        string
	CheckTCPPort   string
	UsePing        string
	PingEvery      string
	Timeout        string
	SlowAbove      string
	PortLabel      string
	PingEveryMs    string
	TimeoutMs      string
	SlowAboveMs    string
	SlowAboveMsOff string

	// Remove hosts
	RemoveHostsTitle string
	RemoveHost       string
	RemoveHosts      func(n int) string
	AndMore          func(n int) string

	// Details
	Periods     [3]string
	SelectAHost string
	PingHistory string
	Ms          string

	// Host dialog
	Name          string
	Host          string
	TCPPort       string
	CheckInstead  string
	Port          string
	SlowTooltip   string
	Notify        string
	BeepDownBack  string
	NotifyTooltip string
	// Sharing the ping times on u00.io
	Share         string
	ShareOn       string
	ShareTooltip  string
	ShareOpen     string
	ShareCopy     string
	LinkCopied    string
	MenuShareOpen string
	MenuShareCopy string
	ShareStart    string
	ShareStop     string

	// Downtime dialog
	DowntimeTitle string
	NoDowntime    string
	DownTimes     func(n int) string
	// The number of outages of several hosts
	GroupDownTimes func(n int) string
	Started        string
	Ended          string
	Duration       string
	StillDown      string
	// Go time layout of a time that is not today
	DateTimeLayout string

	// Configs dialog
	ConfigsTitle       string
	Configs            string
	NewConfig          string
	SaveConfigAs       string
	RemoveConfig       string
	HostsCount         string
	Hosts              string
	NewConfigTitle     string
	SaveConfigAsTitle  string
	CopySuffix         string
	RemoveConfigTitle  string
	RemoveConfigAsk    string
	CannotRemoveOpened string

	// Settings dialog
	ExtraColumns   string
	ShowMin        string
	ShowJitter     string
	ShowSince      string
	AlwaysOnTop    string
	Language       string
	LanguageSystem string
	Theme          string
	ThemeDark      string
	ThemeLight     string

	// About dialog
	AboutTitle   func(name string) string
	Version      string
	Author       string
	License      string
	VisitWebsite string
	Close        string

	// Export
	ExportTitle string
	CSVFiles    string
	SavedTo     func(path string) string
	// At start, when the last opened config could not be read
	ConfigLoadFailed func(err string) string

	// Installing a downloaded copy to ~/.altbins (Windows) and removing it
	Install          string
	Update           string // an older version is installed
	Uninstall        string
	InstallAsk       func(dir string) string
	InstallFailed    func(err string) string
	Installed        string // a toast in the installed copy
	UninstallAsk     func(dir string) string
	UninstallRunning string
	Uninstalled      string
}

// ColumnStrings are the names of the table columns
type ColumnStrings struct {
	Name    string
	IP      string
	Time    string
	Trend   string // the small chart of the last minutes
	Loss    string
	Min     string
	Avg     string
	Max     string
	Jitter  string
	Since   string
	Details string
}

type UnitStrings struct {
	Second string
	Minute string
	Hour   string
	Day    string
}

// languages are offered in the settings; the names stay in their own language
var languages = []struct{ tag, name string }{
	{"en", "English"},
	{"ru", "Русский"},
	{"pl", "Polski"},
	{"sr", "Српски"},
	{"de", "Deutsch"},
	{"fr", "Français"},
	{"es", "Español"},
	{"it", "Italiano"},
	{"pt", "Português"},
	{"zh", "中文"},
	{"ja", "日本語"},
	{"ko", "한국어"},
	{"nl", "Nederlands"},
	{"sv", "Svenska"},
	{"da", "Dansk"},
	{"nb", "Norsk bokmål"},
	{"fi", "Suomi"},
	{"cs", "Čeština"},
	{"sk", "Slovenčina"},
	{"sl", "Slovenščina"},
	{"hr", "Hrvatski"},
	{"hu", "Magyar"},
	{"ro", "Română"},
	{"bg", "Български"},
	{"el", "Ελληνικά"},
	{"lt", "Lietuvių"},
	{"lv", "Latviešu"},
	{"et", "Eesti"},
	{"tr", "Türkçe"},
	{"id", "Bahasa Indonesia"},
	{"vi", "Tiếng Việt"},
}

var en = Strings{
	Error: "Error",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d down)", n) },
	ToolConfigs:     "Configurations (O)",
	ToolAddHost:     "Add host (A)",
	ToolEditHost:    "Edit host (E)",
	ToolRemoveHosts: "Remove selected hosts (Del)",
	ToolDetails:     "Details (D)",
	ToolStart:       "Start pinging",
	ToolStop:        "Stop pinging",
	ToolTray:        "Minimize to tray (T)",
	ToolShareOpen:   "Open selected hosts on u00.io (U)",
	ToolStyle:       "Style (S)",
	KPIMedian:       "Median ping",
	KPIUptime:       "Uptime 24 h",
	TrayShow:        "Show",
	TrayQuit:        "Quit",
	Settings:        "Settings",
	Help:            "Help",
	About:           "About",
	ModeUDP:         "UDP Mode",
	ModeICMP:        "ICMP Mode",
	StatusHosts:     "Hosts",
	StatusOK:        "OK",
	StatusSlow:      "Slow",
	StatusDown:      "Down",
	StatusStopped:   "Pinging stopped",

	Columns: ColumnStrings{
		Name: "Name", IP: "IP", Time: "Time ms", Trend: "Last 5 min", Loss: "Loss %", Min: "Min ms", Avg: "Avg ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Since", Details: "Details",
	},
	StateOK:         "OK",
	StateTimeout:    "TIMEOUT",
	StateNoResolve:  "CANNOT RESOLVE HOSTNAME",
	StatePortClosed: "PORT CLOSED",
	NoHosts:         "No hosts yet",
	PressAToAdd:     "Press A to add a host",
	HowItWorks:      "How it works",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "h", Day: "d"},

	MenuEdit:       "Edit... (E)",
	MenuRemove:     "Remove (Del)",
	MenuDowntime:   "Downtime...",
	MenuExport:     "Export history...",
	MenuChange:     "Change selected",
	BeepOn:         "Beep on",
	BeepOff:        "Beep off",
	CheckTCPPort:   "Check TCP port",
	UsePing:        "Use ping",
	PingEvery:      "Ping every",
	Timeout:        "Timeout",
	SlowAbove:      "Slow above",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping every, ms:",
	TimeoutMs:      "Timeout, ms:",
	SlowAboveMs:    "Slow above, ms:",
	SlowAboveMsOff: "Slow above, ms (0 - off):",

	RemoveHostsTitle: "Remove hosts",
	RemoveHost:       "Remove the host?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("Remove %d hosts?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... and %d more", n) },

	Periods:     [3]string{"5m", "1h", "24h"},
	SelectAHost: "Select a host",
	PingHistory: "Ping history: ",
	Ms:          "ms",

	Name:          "Name:",
	Host:          "Host:",
	TCPPort:       "TCP port:",
	CheckInstead:  "Check instead of ping",
	Port:          "Port:",
	SlowTooltip:   "Show the host in yellow when its average ping time is above it. 0 - off",
	Notify:        "Notify:",
	BeepDownBack:  "Beep when down or back",
	NotifyTooltip: "Down means 3 failed pings in a row. The window title shows how many such hosts are down.",
	Share:         "Share:",
	ShareOn:       "Publish ping time on u00.io",
	ShareTooltip:  "The ping times are sent to a public page. Anyone with the link can see them.",
	ShareOpen:     "Open",
	ShareCopy:     "Copy link",
	LinkCopied:    "Link copied",
	MenuShareOpen: "Open on u00.io",
	MenuShareCopy: "Copy u00.io link",
	ShareStart:    "Publish on u00.io",
	ShareStop:     "Stop publishing",

	DowntimeTitle: "Downtime",
	NoDowntime:    "No downtime in the last 24 hours",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Down once in the last 24 hours:"
		}
		return fmt.Sprintf("Down %d times in the last 24 hours:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Outages in the last 24 hours: %d", n) },
	Started:        "Started",
	Ended:          "Ended",
	Duration:       "Duration",
	StillDown:      "still down",
	DateTimeLayout: "Jan 2 15:04:05",

	ConfigsTitle:       "Open Config",
	Configs:            "Configs:",
	NewConfig:          "New...",
	SaveConfigAs:       "Save As...",
	RemoveConfig:       "Remove",
	HostsCount:         "Hosts Count",
	Hosts:              "Hosts",
	NewConfigTitle:     "New Config",
	SaveConfigAsTitle:  "Save Config As",
	CopySuffix:         " copy",
	RemoveConfigTitle:  "Remove config",
	RemoveConfigAsk:    "Remove the selected config?",
	CannotRemoveOpened: "Cannot remove the currently opened config.",

	ExtraColumns:   "Extra columns:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Since (last change)",
	AlwaysOnTop:    "Always on top",
	Language:       "Language:",
	LanguageSystem: "As in the system",
	Theme:          "Theme:",
	ThemeDark:      "Dark",
	ThemeLight:     "Light",

	AboutTitle:   func(name string) string { return "About " + name },
	Version:      "Version",
	Author:       "Author:",
	License:      "License:",
	VisitWebsite: "Visit Website",
	Close:        "Close",

	ExportTitle: "Export history",
	CSVFiles:    "CSV files",
	SavedTo:     func(path string) string { return "Saved to " + path },
	ConfigLoadFailed: func(err string) string {
		return "The host list could not be read, a new empty one was opened. The file was left as is:\n\n" + err
	},

	Install:   "Install",
	Update:    "Update",
	Uninstall: "Uninstall",
	InstallAsk: func(dir string) string {
		return "AltPing will be copied to " + dir + " and added to the Start menu, the desktop and the list of installed apps. If it is already installed there, it is updated. Settings and host lists stay as they are.\n\nAltPing will then restart from there."
	},
	InstallFailed: func(err string) string { return "AltPing could not be installed:\n\n" + err },
	Installed:     "AltPing is installed",
	UninstallAsk: func(dir string) string {
		return "Remove AltPing from this computer?\n\nSettings and host lists stay in " + dir
	},
	UninstallRunning: "AltPing is running. Close it and try again.",
	Uninstalled:      "AltPing has been removed.",
}

var ru = Strings{
	Error: "Ошибка",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d недоступно)", n) },
	ToolConfigs:     "Конфигурации (O)",
	ToolAddHost:     "Добавить хост (A)",
	ToolEditHost:    "Изменить хост (E)",
	ToolRemoveHosts: "Удалить выбранные хосты (Del)",
	ToolDetails:     "Подробности (D)",
	ToolStart:       "Начать пинг",
	ToolStop:        "Остановить пинг",
	ToolTray:        "Свернуть в трей (T)",
	ToolShareOpen:   "Открыть выбранные на u00.io (U)",
	ToolStyle:       "Стиль (S)",
	KPIMedian:       "Медиана пинга",
	KPIUptime:       "Доступность 24 ч",
	TrayShow:        "Показать",
	TrayQuit:        "Выход",
	Settings:        "Настройки",
	Help:            "Справка",
	About:           "О программе",
	ModeUDP:         "Режим UDP",
	ModeICMP:        "Режим ICMP",
	StatusHosts:     "Хосты",
	StatusOK:        "В норме",
	StatusSlow:      "Медленные",
	StatusDown:      "Недоступны",
	StatusStopped:   "Пинг остановлен",

	Columns: ColumnStrings{
		Name: "Имя", IP: "IP", Time: "Время мс", Trend: "За 5 мин", Loss: "Потери %", Min: "Мин мс", Avg: "Сред мс",
		Max: "Макс мс", Jitter: "Джиттер мс", Since: "С тех пор", Details: "Состояние",
	},
	StateOK:         "OK",
	StateTimeout:    "ТАЙМАУТ",
	StateNoResolve:  "ИМЯ НЕ НАЙДЕНО",
	StatePortClosed: "ПОРТ ЗАКРЫТ",
	NoHosts:         "Хостов пока нет",
	PressAToAdd:     "Нажмите A, чтобы добавить хост",
	HowItWorks:      "Как это работает",

	Units: UnitStrings{Second: "с", Minute: "м", Hour: "ч", Day: "д"},

	MenuEdit:       "Изменить... (E)",
	MenuRemove:     "Удалить (Del)",
	MenuDowntime:   "Простои...",
	MenuExport:     "Экспорт истории...",
	MenuChange:     "Изменить выбранные",
	BeepOn:         "Включить сигнал",
	BeepOff:        "Выключить сигнал",
	CheckTCPPort:   "Проверять TCP-порт",
	UsePing:        "Использовать пинг",
	PingEvery:      "Пинг каждые",
	Timeout:        "Таймаут",
	SlowAbove:      "Медленно при",
	PortLabel:      "Порт:",
	PingEveryMs:    "Пинг каждые, мс:",
	TimeoutMs:      "Таймаут, мс:",
	SlowAboveMs:    "Медленно выше, мс:",
	SlowAboveMsOff: "Медленно выше, мс (0 - выкл.):",

	RemoveHostsTitle: "Удаление хостов",
	RemoveHost:       "Удалить хост?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Удалить %d %s?", n, i18n.Plural("ru", n, "хост", "хоста", "хостов"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... и ещё %d", n) },

	Periods:     [3]string{"5 мин", "1 ч", "24 ч"},
	SelectAHost: "Выберите хост",
	PingHistory: "История пинга: ",
	Ms:          "мс",

	Name:          "Имя:",
	Host:          "Хост:",
	TCPPort:       "TCP-порт:",
	CheckInstead:  "Проверять вместо пинга",
	Port:          "Порт:",
	SlowTooltip:   "Хост показывается жёлтым, когда среднее время пинга выше. 0 - выкл.",
	Notify:        "Оповещение:",
	BeepDownBack:  "Сигнал при пропадании и возврате",
	NotifyTooltip: "Пропал - 3 неудачных пинга подряд. Заголовок окна показывает, сколько таких хостов недоступно.",
	Share:         "Поделиться:",
	ShareOn:       "Публиковать пинг на u00.io",
	ShareTooltip:  "Время пинга отправляется на публичную страницу. Её может открыть любой, у кого есть ссылка.",
	ShareOpen:     "Открыть",
	ShareCopy:     "Копировать ссылку",
	LinkCopied:    "Ссылка скопирована",
	MenuShareOpen: "Открыть на u00.io",
	MenuShareCopy: "Копировать ссылку u00.io",
	ShareStart:    "Публиковать на u00.io",
	ShareStop:     "Не публиковать",

	DowntimeTitle: "Простои",
	NoDowntime:    "Простоев за последние 24 часа не было",
	DownTimes: func(n int) string {
		return fmt.Sprintf("Недоступен %d %s за последние 24 часа:", n, i18n.Plural("ru", n, "раз", "раза", "раз"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Простои за последние 24 часа: %d", n) },
	Started:        "Начало",
	Ended:          "Конец",
	Duration:       "Длительность",
	StillDown:      "всё ещё недоступен",
	DateTimeLayout: "02.01 15:04:05",

	ConfigsTitle:       "Открыть конфигурацию",
	Configs:            "Конфигурации:",
	NewConfig:          "Новая...",
	SaveConfigAs:       "Сохранить как...",
	RemoveConfig:       "Удалить",
	HostsCount:         "Хостов",
	Hosts:              "Хосты",
	NewConfigTitle:     "Новая конфигурация",
	SaveConfigAsTitle:  "Сохранить конфигурацию как",
	CopySuffix:         " копия",
	RemoveConfigTitle:  "Удаление конфигурации",
	RemoveConfigAsk:    "Удалить выбранную конфигурацию?",
	CannotRemoveOpened: "Нельзя удалить открытую сейчас конфигурацию.",

	ExtraColumns:   "Дополнительные колонки:",
	ShowMin:        "Мин",
	ShowJitter:     "Джиттер",
	ShowSince:      "С тех пор (последнее изменение)",
	AlwaysOnTop:    "Поверх всех окон",
	Language:       "Язык:",
	LanguageSystem: "Как в системе",
	Theme:          "Тема:",
	ThemeDark:      "Тёмная",
	ThemeLight:     "Светлая",

	AboutTitle:   func(name string) string { return "О программе " + name },
	Version:      "Версия",
	Author:       "Автор:",
	License:      "Лицензия:",
	VisitWebsite: "Открыть сайт",
	Close:        "Закрыть",

	ExportTitle: "Экспорт истории",
	CSVFiles:    "Файлы CSV",
	SavedTo:     func(path string) string { return "Сохранено в " + path },
	ConfigLoadFailed: func(err string) string {
		return "Не удалось прочитать список хостов, открыт новый пустой. Файл оставлен как есть:\n\n" + err
	},

	Install:   "Установить",
	Update:    "Обновить",
	Uninstall: "Удалить",
	InstallAsk: func(dir string) string {
		return "AltPing будет скопирован в " + dir + " и добавлен в меню «Пуск», на рабочий стол и в список установленных приложений. Если он там уже установлен, он обновится. Настройки и списки хостов останутся как есть.\n\nЗатем AltPing перезапустится оттуда."
	},
	InstallFailed: func(err string) string { return "Не удалось установить AltPing:\n\n" + err },
	Installed:     "AltPing установлен",
	UninstallAsk: func(dir string) string {
		return "Удалить AltPing с этого компьютера?\n\nНастройки и списки хостов останутся в " + dir
	},
	UninstallRunning: "AltPing запущен. Закройте его и попробуйте снова.",
	Uninstalled:      "AltPing удалён.",
}

var pl = Strings{
	Error: "Błąd",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d niedostępnych)", n) },
	ToolConfigs:     "Konfiguracje (O)",
	ToolAddHost:     "Dodaj host (A)",
	ToolEditHost:    "Edytuj host (E)",
	ToolRemoveHosts: "Usuń zaznaczone hosty (Del)",
	ToolDetails:     "Szczegóły (D)",
	ToolStart:       "Rozpocznij pingowanie",
	ToolStop:        "Zatrzymaj pingowanie",
	ToolTray:        "Minimalizuj do zasobnika (T)",
	ToolShareOpen:   "Otwórz zaznaczone na u00.io (U)",
	ToolStyle:       "Styl (S)",
	KPIMedian:       "Mediana pingu",
	KPIUptime:       "Dostępność 24 h",
	TrayShow:        "Pokaż",
	TrayQuit:        "Zakończ",
	Settings:        "Ustawienia",
	Help:            "Pomoc",
	About:           "O programie",
	ModeUDP:         "Tryb UDP",
	ModeICMP:        "Tryb ICMP",
	StatusHosts:     "Hosty",
	StatusOK:        "OK",
	StatusSlow:      "Wolne",
	StatusDown:      "Niedostępne",
	StatusStopped:   "Pingowanie zatrzymane",

	Columns: ColumnStrings{
		Name: "Nazwa", IP: "IP", Time: "Czas ms", Trend: "Ostatnie 5 min", Loss: "Straty %", Min: "Min ms", Avg: "Śr. ms",
		Max: "Maks ms", Jitter: "Jitter ms", Since: "Od", Details: "Stan",
	},
	StateOK:         "OK",
	StateTimeout:    "PRZEKROCZONO CZAS",
	StateNoResolve:  "NIE MOŻNA ROZWIĄZAĆ NAZWY",
	StatePortClosed: "PORT ZAMKNIĘTY",
	NoHosts:         "Brak hostów",
	PressAToAdd:     "Naciśnij A, aby dodać host",
	HowItWorks:      "Jak to działa",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "d"},

	MenuEdit:       "Edytuj... (E)",
	MenuRemove:     "Usuń (Del)",
	MenuDowntime:   "Przestoje...",
	MenuExport:     "Eksportuj historię...",
	MenuChange:     "Zmień zaznaczone",
	BeepOn:         "Włącz sygnał dźwiękowy",
	BeepOff:        "Wyłącz sygnał dźwiękowy",
	CheckTCPPort:   "Sprawdzaj port TCP",
	UsePing:        "Używaj pingu",
	PingEvery:      "Ping co",
	Timeout:        "Limit czasu",
	SlowAbove:      "Wolno powyżej",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping co, ms:",
	TimeoutMs:      "Limit czasu, ms:",
	SlowAboveMs:    "Wolno powyżej, ms:",
	SlowAboveMsOff: "Wolno powyżej, ms (0 - wył.):",

	RemoveHostsTitle: "Usuwanie hostów",
	RemoveHost:       "Usunąć host?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Usunąć %d %s?", n, i18n.Plural("pl", n, "host", "hosty", "hostów"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... i %d więcej", n) },

	Periods:     [3]string{"5 min", "1 godz", "24 godz"},
	SelectAHost: "Wybierz host",
	PingHistory: "Historia pingu: ",
	Ms:          "ms",

	Name:          "Nazwa:",
	Host:          "Host:",
	TCPPort:       "Port TCP:",
	CheckInstead:  "Sprawdzaj zamiast pingu",
	Port:          "Port:",
	SlowTooltip:   "Host jest pokazywany na żółto, gdy jego średni czas pingu przekracza tę wartość. 0 - wył.",
	Notify:        "Powiadamianie:",
	BeepDownBack:  "Sygnał przy awarii i powrocie",
	NotifyTooltip: "Awaria oznacza 3 nieudane pingi z rzędu. Tytuł okna pokazuje, ile takich hostów jest niedostępnych.",
	Share:         "Udostępnij:",
	ShareOn:       "Publikuj ping na u00.io",
	ShareTooltip:  "Czasy pingu są wysyłane na publiczną stronę. Każdy, kto ma link, może je zobaczyć.",
	ShareOpen:     "Otwórz",
	ShareCopy:     "Kopiuj link",
	LinkCopied:    "Skopiowano link",
	MenuShareOpen: "Otwórz na u00.io",
	MenuShareCopy: "Kopiuj link u00.io",
	ShareStart:    "Publikuj na u00.io",
	ShareStop:     "Przestań publikować",

	DowntimeTitle: "Przestoje",
	NoDowntime:    "Brak przestojów w ciągu ostatnich 24 godzin",
	DownTimes: func(n int) string {
		return fmt.Sprintf("Niedostępny %d %s w ciągu ostatnich 24 godzin:", n, i18n.Plural("pl", n, "raz", "razy", "razy"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Przestoje w ciągu ostatnich 24 godzin: %d", n) },
	Started:        "Początek",
	Ended:          "Koniec",
	Duration:       "Czas trwania",
	StillDown:      "nadal niedostępny",
	DateTimeLayout: "02.01 15:04:05",

	ConfigsTitle:       "Otwórz konfigurację",
	Configs:            "Konfiguracje:",
	NewConfig:          "Nowa...",
	SaveConfigAs:       "Zapisz jako...",
	RemoveConfig:       "Usuń",
	HostsCount:         "Liczba hostów",
	Hosts:              "Hosty",
	NewConfigTitle:     "Nowa konfiguracja",
	SaveConfigAsTitle:  "Zapisz konfigurację jako",
	CopySuffix:         " kopia",
	RemoveConfigTitle:  "Usuwanie konfiguracji",
	RemoveConfigAsk:    "Usunąć zaznaczoną konfigurację?",
	CannotRemoveOpened: "Nie można usunąć aktualnie otwartej konfiguracji.",

	ExtraColumns:   "Dodatkowe kolumny:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Od (ostatnia zmiana)",
	AlwaysOnTop:    "Zawsze na wierzchu",
	Language:       "Język:",
	LanguageSystem: "Jak w systemie",
	Theme:          "Motyw:",
	ThemeDark:      "Ciemny",
	ThemeLight:     "Jasny",

	AboutTitle:   func(name string) string { return "O programie " + name },
	Version:      "Wersja",
	Author:       "Autor:",
	License:      "Licencja:",
	VisitWebsite: "Odwiedź stronę",
	Close:        "Zamknij",

	ExportTitle: "Eksport historii",
	CSVFiles:    "Pliki CSV",
	SavedTo:     func(path string) string { return "Zapisano w " + path },
	ConfigLoadFailed: func(err string) string {
		return "Nie udało się odczytać listy hostów, otwarto nową pustą. Plik pozostawiono bez zmian:\n\n" + err
	},

	Install:   "Zainstaluj",
	Update:    "Zaktualizuj",
	Uninstall: "Odinstaluj",
	InstallAsk: func(dir string) string {
		return "AltPing zostanie skopiowany do " + dir + " i dodany do menu Start, na pulpit i do listy zainstalowanych aplikacji. Jeśli jest już tam zainstalowany, zostanie zaktualizowany. Ustawienia i listy hostów pozostaną bez zmian.\n\nNastępnie AltPing uruchomi się ponownie z tego miejsca."
	},
	InstallFailed: func(err string) string { return "Nie udało się zainstalować AltPing:\n\n" + err },
	Installed:     "AltPing został zainstalowany",
	UninstallAsk: func(dir string) string {
		return "Usunąć AltPing z tego komputera?\n\nUstawienia i listy hostów pozostaną w " + dir
	},
	UninstallRunning: "AltPing jest uruchomiony. Zamknij go i spróbuj ponownie.",
	Uninstalled:      "AltPing został usunięty.",
}

// Serbian in Cyrillic, the script of the "sr" tag
var sr = Strings{
	Error: "Грешка",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d недоступно)", n) },
	ToolConfigs:     "Конфигурације (O)",
	ToolAddHost:     "Додај хост (A)",
	ToolEditHost:    "Измени хост (E)",
	ToolRemoveHosts: "Уклони изабране хостове (Del)",
	ToolDetails:     "Детаљи (D)",
	ToolStart:       "Покрени пинговање",
	ToolStop:        "Заустави пинговање",
	ToolTray:        "Минимизуј у системску палету (T)",
	ToolShareOpen:   "Отвори изабране на u00.io (U)",
	ToolStyle:       "Стил (S)",
	KPIMedian:       "Медијана пинга",
	KPIUptime:       "Доступност 24 ч",
	TrayShow:        "Прикажи",
	TrayQuit:        "Изађи",
	Settings:        "Подешавања",
	Help:            "Помоћ",
	About:           "О програму",
	ModeUDP:         "UDP режим",
	ModeICMP:        "ICMP режим",
	StatusHosts:     "Хостови",
	StatusOK:        "У реду",
	StatusSlow:      "Спори",
	StatusDown:      "Недоступни",
	StatusStopped:   "Пинговање заустављено",

	Columns: ColumnStrings{
		Name: "Име", IP: "IP", Time: "Време мс", Trend: "Последњих 5 мин", Loss: "Губитак %", Min: "Мин мс", Avg: "Прос мс",
		Max: "Макс мс", Jitter: "Џитер мс", Since: "Од", Details: "Стање",
	},
	StateOK:         "OK",
	StateTimeout:    "ИСТЕКЛО ВРЕМЕ",
	StateNoResolve:  "ИМЕ НИЈЕ ПРОНАЂЕНО",
	StatePortClosed: "ПОРТ ЗАТВОРЕН",
	NoHosts:         "Још нема хостова",
	PressAToAdd:     "Притисните A да додате хост",
	HowItWorks:      "Како ради",

	Units: UnitStrings{Second: "с", Minute: "м", Hour: "ч", Day: "д"},

	MenuEdit:       "Измени... (E)",
	MenuRemove:     "Уклони (Del)",
	MenuDowntime:   "Прекиди...",
	MenuExport:     "Извези историју...",
	MenuChange:     "Измени изабране",
	BeepOn:         "Укључи звучни сигнал",
	BeepOff:        "Искључи звучни сигнал",
	CheckTCPPort:   "Проверавај TCP порт",
	UsePing:        "Користи пинг",
	PingEvery:      "Пинг сваких",
	Timeout:        "Време чекања",
	SlowAbove:      "Споро изнад",
	PortLabel:      "Порт:",
	PingEveryMs:    "Пинг сваких, мс:",
	TimeoutMs:      "Време чекања, мс:",
	SlowAboveMs:    "Споро изнад, мс:",
	SlowAboveMsOff: "Споро изнад, мс (0 - искључено):",

	RemoveHostsTitle: "Уклањање хостова",
	RemoveHost:       "Уклонити хост?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Уклонити %d %s?", n, i18n.Plural("sr", n, "хост", "хоста", "хостова"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... и још %d", n) },

	Periods:     [3]string{"5 мин", "1 ч", "24 ч"},
	SelectAHost: "Изаберите хост",
	PingHistory: "Историја пинга: ",
	Ms:          "мс",

	Name:          "Име:",
	Host:          "Хост:",
	TCPPort:       "TCP порт:",
	CheckInstead:  "Проверавај уместо пинга",
	Port:          "Порт:",
	SlowTooltip:   "Хост се приказује жутом бојом када је његово просечно време пинга изнад ове вредности. 0 - искључено",
	Notify:        "Обавештење:",
	BeepDownBack:  "Звучни сигнал при паду и повратку",
	NotifyTooltip: "Пад значи 3 неуспела пинга заредом. Наслов прозора показује колико таквих хостова није доступно.",
	Share:         "Подели:",
	ShareOn:       "Објављуј пинг на u00.io",
	ShareTooltip:  "Времена пинга шаљу се на јавну страницу. Свако ко има везу може да их види.",
	ShareOpen:     "Отвори",
	ShareCopy:     "Копирај везу",
	LinkCopied:    "Веза је копирана",
	MenuShareOpen: "Отвори на u00.io",
	MenuShareCopy: "Копирај u00.io везу",
	ShareStart:    "Објављуј на u00.io",
	ShareStop:     "Престани да објављујеш",

	DowntimeTitle: "Прекиди",
	NoDowntime:    "Нема прекида у последња 24 сата",
	DownTimes: func(n int) string {
		return fmt.Sprintf("Недоступан %d %s у последња 24 сата:", n, i18n.Plural("sr", n, "пут", "пута", "пута"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Прекиди у последња 24 сата: %d", n) },
	Started:        "Почетак",
	Ended:          "Крај",
	Duration:       "Трајање",
	StillDown:      "још увек недоступан",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Отвори конфигурацију",
	Configs:            "Конфигурације:",
	NewConfig:          "Нова...",
	SaveConfigAs:       "Сачувај као...",
	RemoveConfig:       "Уклони",
	HostsCount:         "Број хостова",
	Hosts:              "Хостови",
	NewConfigTitle:     "Нова конфигурација",
	SaveConfigAsTitle:  "Сачувај конфигурацију као",
	CopySuffix:         " копија",
	RemoveConfigTitle:  "Уклањање конфигурације",
	RemoveConfigAsk:    "Уклонити изабрану конфигурацију?",
	CannotRemoveOpened: "Није могуће уклонити тренутно отворену конфигурацију.",

	ExtraColumns:   "Додатне колоне:",
	ShowMin:        "Мин",
	ShowJitter:     "Џитер",
	ShowSince:      "Од (последња промена)",
	AlwaysOnTop:    "Увек на врху",
	Language:       "Језик:",
	LanguageSystem: "Као у систему",
	Theme:          "Тема:",
	ThemeDark:      "Тамна",
	ThemeLight:     "Светла",

	AboutTitle:   func(name string) string { return "О програму " + name },
	Version:      "Верзија",
	Author:       "Аутор:",
	License:      "Лиценца:",
	VisitWebsite: "Посети сајт",
	Close:        "Затвори",

	ExportTitle: "Извоз историје",
	CSVFiles:    "CSV датотеке",
	SavedTo:     func(path string) string { return "Сачувано у " + path },
	ConfigLoadFailed: func(err string) string {
		return "Листа хостова није могла да се прочита, отворена је нова празна. Датотека је остављена каква јесте:\n\n" + err
	},

	Install:   "Инсталирај",
	Update:    "Ажурирај",
	Uninstall: "Деинсталирај",
	InstallAsk: func(dir string) string {
		return "AltPing ће бити копиран у " + dir + " и додат у мени Старт, на радну површину и у листу инсталираних апликација. Ако је тамо већ инсталиран, биће ажуриран. Подешавања и листе хостова остају непромењени.\n\nЗатим ће се AltPing поново покренути одатле."
	},
	InstallFailed: func(err string) string {
		return "Инсталирање програма AltPing није успело:\n\n" + err
	},
	Installed: "AltPing је инсталиран",
	UninstallAsk: func(dir string) string {
		return "Уклонити AltPing са овог рачунара?\n\nПодешавања и листе хостова остају у " + dir
	},
	UninstallRunning: "AltPing је покренут. Затворите га и покушајте поново.",
	Uninstalled:      "AltPing је уклоњен.",
}

var de = Strings{
	Error: "Fehler",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d offline)", n) },
	ToolConfigs:     "Konfigurationen (O)",
	ToolAddHost:     "Host hinzufügen (A)",
	ToolEditHost:    "Host bearbeiten (E)",
	ToolRemoveHosts: "Ausgewählte Hosts entfernen (Del)",
	ToolDetails:     "Details (D)",
	ToolStart:       "Ping starten",
	ToolStop:        "Ping stoppen",
	ToolTray:        "In den Infobereich minimieren (T)",
	ToolShareOpen:   "Ausgewählte auf u00.io öffnen (U)",
	ToolStyle:       "Stil (S)",
	KPIMedian:       "Median-Ping",
	KPIUptime:       "Verfügbarkeit 24 h",
	TrayShow:        "Anzeigen",
	TrayQuit:        "Beenden",
	Settings:        "Einstellungen",
	Help:            "Hilfe",
	About:           "Über",
	ModeUDP:         "UDP-Modus",
	ModeICMP:        "ICMP-Modus",
	StatusHosts:     "Hosts",
	StatusOK:        "OK",
	StatusSlow:      "Langsam",
	StatusDown:      "Nicht erreichbar",
	StatusStopped:   "Ping gestoppt",

	Columns: ColumnStrings{
		Name: "Name", IP: "IP", Time: "Zeit ms", Trend: "Letzte 5 Min.", Loss: "Verlust %", Min: "Min ms", Avg: "Mittel ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Seit", Details: "Details",
	},
	StateOK:         "OK",
	StateTimeout:    "ZEITÜBERSCHREITUNG",
	StateNoResolve:  "HOSTNAME NICHT AUFLÖSBAR",
	StatePortClosed: "PORT GESCHLOSSEN",
	NoHosts:         "Noch keine Hosts",
	PressAToAdd:     "Drücken Sie A, um einen Host hinzuzufügen",
	HowItWorks:      "So funktioniert es",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "h", Day: "T"},

	MenuEdit:       "Bearbeiten... (E)",
	MenuRemove:     "Entfernen (Del)",
	MenuDowntime:   "Ausfallzeiten...",
	MenuExport:     "Verlauf exportieren...",
	MenuChange:     "Ausgewählte ändern",
	BeepOn:         "Signalton ein",
	BeepOff:        "Signalton aus",
	CheckTCPPort:   "TCP-Port prüfen",
	UsePing:        "Ping verwenden",
	PingEvery:      "Ping alle",
	Timeout:        "Zeitlimit",
	SlowAbove:      "Langsam ab",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping alle, ms:",
	TimeoutMs:      "Zeitlimit, ms:",
	SlowAboveMs:    "Langsam ab, ms:",
	SlowAboveMsOff: "Langsam ab, ms (0 - aus):",

	RemoveHostsTitle: "Hosts entfernen",
	RemoveHost:       "Den Host entfernen?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("%d Hosts entfernen?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... und %d weitere", n) },

	Periods:     [3]string{"5 Min", "1 Std", "24 Std"},
	SelectAHost: "Wählen Sie einen Host",
	PingHistory: "Ping-Verlauf: ",
	Ms:          "ms",

	Name:          "Name:",
	Host:          "Host:",
	TCPPort:       "TCP-Port:",
	CheckInstead:  "Statt Ping prüfen",
	Port:          "Port:",
	SlowTooltip:   "Den Host gelb anzeigen, wenn seine mittlere Ping-Zeit darüber liegt. 0 - aus",
	Notify:        "Benachrichtigen:",
	BeepDownBack:  "Signalton bei Ausfall und Rückkehr",
	NotifyTooltip: "Ausgefallen heißt 3 fehlgeschlagene Pings in Folge. Der Fenstertitel zeigt, wie viele solche Hosts ausgefallen sind.",
	Share:         "Teilen:",
	ShareOn:       "Ping auf u00.io veröffentlichen",
	ShareTooltip:  "Die Ping-Zeiten werden an eine öffentliche Seite gesendet. Jeder mit dem Link kann sie sehen.",
	ShareOpen:     "Öffnen",
	ShareCopy:     "Link kopieren",
	LinkCopied:    "Link kopiert",
	MenuShareOpen: "Auf u00.io öffnen",
	MenuShareCopy: "u00.io-Link kopieren",
	ShareStart:    "Auf u00.io veröffentlichen",
	ShareStop:     "Nicht mehr veröffentlichen",

	DowntimeTitle: "Ausfallzeiten",
	NoDowntime:    "Keine Ausfälle in den letzten 24 Stunden",
	DownTimes: func(n int) string {
		if n == 1 {
			return "In den letzten 24 Stunden einmal ausgefallen:"
		}
		return fmt.Sprintf("In den letzten 24 Stunden %d-mal ausgefallen:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Ausfälle in den letzten 24 Stunden: %d", n) },
	Started:        "Beginn",
	Ended:          "Ende",
	Duration:       "Dauer",
	StillDown:      "noch ausgefallen",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Konfiguration öffnen",
	Configs:            "Konfigurationen:",
	NewConfig:          "Neu...",
	SaveConfigAs:       "Speichern unter...",
	RemoveConfig:       "Entfernen",
	HostsCount:         "Anzahl Hosts",
	Hosts:              "Hosts",
	NewConfigTitle:     "Neue Konfiguration",
	SaveConfigAsTitle:  "Konfiguration speichern unter",
	CopySuffix:         " Kopie",
	RemoveConfigTitle:  "Konfiguration entfernen",
	RemoveConfigAsk:    "Die ausgewählte Konfiguration entfernen?",
	CannotRemoveOpened: "Die gerade geöffnete Konfiguration kann nicht entfernt werden.",

	ExtraColumns:   "Zusätzliche Spalten:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Seit (letzte Änderung)",
	AlwaysOnTop:    "Immer im Vordergrund",
	Language:       "Sprache:",
	LanguageSystem: "Wie im System",
	Theme:          "Design:",
	ThemeDark:      "Dunkel",
	ThemeLight:     "Hell",

	AboutTitle:   func(name string) string { return "Über " + name },
	Version:      "Version",
	Author:       "Autor:",
	License:      "Lizenz:",
	VisitWebsite: "Website besuchen",
	Close:        "Schließen",

	ExportTitle: "Verlauf exportieren",
	CSVFiles:    "CSV-Dateien",
	SavedTo:     func(path string) string { return "Gespeichert unter " + path },
	ConfigLoadFailed: func(err string) string {
		return "Die Hostliste konnte nicht gelesen werden, eine neue leere wurde geöffnet. Die Datei wurde unverändert gelassen:\n\n" + err
	},

	Install:   "Installieren",
	Update:    "Aktualisieren",
	Uninstall: "Deinstallieren",
	InstallAsk: func(dir string) string {
		return "AltPing wird nach " + dir + " kopiert und zum Startmenü, zum Desktop und zur Liste der installierten Apps hinzugefügt. Ist es dort bereits installiert, wird es aktualisiert. Einstellungen und Hostlisten bleiben erhalten.\n\nAnschließend startet AltPing von dort neu."
	},
	InstallFailed: func(err string) string { return "AltPing konnte nicht installiert werden:\n\n" + err },
	Installed:     "AltPing ist installiert",
	UninstallAsk: func(dir string) string {
		return "AltPing von diesem Computer entfernen?\n\nEinstellungen und Hostlisten bleiben in " + dir
	},
	UninstallRunning: "AltPing läuft. Schließen Sie es und versuchen Sie es erneut.",
	Uninstalled:      "AltPing wurde entfernt.",
}

var fr = Strings{
	Error: "Erreur",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d hors ligne)", n) },
	ToolConfigs:     "Configurations (O)",
	ToolAddHost:     "Ajouter un hôte (A)",
	ToolEditHost:    "Modifier l'hôte (E)",
	ToolRemoveHosts: "Supprimer les hôtes sélectionnés (Del)",
	ToolDetails:     "Détails (D)",
	ToolStart:       "Démarrer le ping",
	ToolStop:        "Arrêter le ping",
	ToolTray:        "Réduire dans la zone de notification (T)",
	ToolShareOpen:   "Ouvrir la sélection sur u00.io (U)",
	ToolStyle:       "Style (S)",
	KPIMedian:       "Ping médian",
	KPIUptime:       "Disponibilité 24 h",
	TrayShow:        "Afficher",
	TrayQuit:        "Quitter",
	Settings:        "Paramètres",
	Help:            "Aide",
	About:           "À propos",
	ModeUDP:         "Mode UDP",
	ModeICMP:        "Mode ICMP",
	StatusHosts:     "Hôtes",
	StatusOK:        "OK",
	StatusSlow:      "Lents",
	StatusDown:      "Injoignables",
	StatusStopped:   "Ping arrêté",

	Columns: ColumnStrings{
		Name: "Nom", IP: "IP", Time: "Temps ms", Trend: "5 dernières min", Loss: "Perte %", Min: "Min ms", Avg: "Moy ms",
		Max: "Max ms", Jitter: "Gigue ms", Since: "Depuis", Details: "Détails",
	},
	StateOK:         "OK",
	StateTimeout:    "DÉLAI DÉPASSÉ",
	StateNoResolve:  "NOM D'HÔTE INTROUVABLE",
	StatePortClosed: "PORT FERMÉ",
	NoHosts:         "Aucun hôte pour l'instant",
	PressAToAdd:     "Appuyez sur A pour ajouter un hôte",
	HowItWorks:      "Comment ça marche",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "j"},

	MenuEdit:       "Modifier... (E)",
	MenuRemove:     "Supprimer (Del)",
	MenuDowntime:   "Interruptions...",
	MenuExport:     "Exporter l'historique...",
	MenuChange:     "Modifier la sélection",
	BeepOn:         "Bip activé",
	BeepOff:        "Bip désactivé",
	CheckTCPPort:   "Vérifier un port TCP",
	UsePing:        "Utiliser le ping",
	PingEvery:      "Ping toutes les",
	Timeout:        "Délai",
	SlowAbove:      "Lent au-delà de",
	PortLabel:      "Port :",
	PingEveryMs:    "Ping toutes les, ms :",
	TimeoutMs:      "Délai, ms :",
	SlowAboveMs:    "Lent au-delà de, ms :",
	SlowAboveMsOff: "Lent au-delà de, ms (0 - désactivé) :",

	RemoveHostsTitle: "Supprimer des hôtes",
	RemoveHost:       "Supprimer l'hôte ?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("Supprimer %d hôtes ?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... et %d de plus", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Sélectionnez un hôte",
	PingHistory: "Historique du ping : ",
	Ms:          "ms",

	Name:          "Nom :",
	Host:          "Hôte :",
	TCPPort:       "Port TCP :",
	CheckInstead:  "Vérifier au lieu du ping",
	Port:          "Port :",
	SlowTooltip:   "Afficher l'hôte en jaune quand son temps de ping moyen dépasse cette valeur. 0 - désactivé",
	Notify:        "Notifier :",
	BeepDownBack:  "Bip en cas de panne ou de retour",
	NotifyTooltip: "En panne signifie 3 pings échoués d'affilée. Le titre de la fenêtre indique combien de ces hôtes sont en panne.",
	Share:         "Partager :",
	ShareOn:       "Publier le ping sur u00.io",
	ShareTooltip:  "Les temps de ping sont envoyés sur une page publique. Toute personne ayant le lien peut les voir.",
	ShareOpen:     "Ouvrir",
	ShareCopy:     "Copier le lien",
	LinkCopied:    "Lien copié",
	MenuShareOpen: "Ouvrir sur u00.io",
	MenuShareCopy: "Copier le lien u00.io",
	ShareStart:    "Publier sur u00.io",
	ShareStop:     "Arrêter la publication",

	DowntimeTitle: "Interruptions",
	NoDowntime:    "Aucune interruption au cours des dernières 24 heures",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Hors ligne une fois au cours des dernières 24 heures :"
		}
		return fmt.Sprintf("Hors ligne %d fois au cours des dernières 24 heures :", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Interruptions au cours des dernières 24 heures : %d", n) },
	Started:        "Début",
	Ended:          "Fin",
	Duration:       "Durée",
	StillDown:      "toujours hors ligne",
	DateTimeLayout: "02/01 15:04:05",

	ConfigsTitle:       "Ouvrir une configuration",
	Configs:            "Configurations :",
	NewConfig:          "Nouvelle...",
	SaveConfigAs:       "Enregistrer sous...",
	RemoveConfig:       "Supprimer",
	HostsCount:         "Nb d'hôtes",
	Hosts:              "Hôtes",
	NewConfigTitle:     "Nouvelle configuration",
	SaveConfigAsTitle:  "Enregistrer la configuration sous",
	CopySuffix:         " copie",
	RemoveConfigTitle:  "Supprimer la configuration",
	RemoveConfigAsk:    "Supprimer la configuration sélectionnée ?",
	CannotRemoveOpened: "Impossible de supprimer la configuration actuellement ouverte.",

	ExtraColumns:   "Colonnes supplémentaires :",
	ShowMin:        "Min",
	ShowJitter:     "Gigue",
	ShowSince:      "Depuis (dernier changement)",
	AlwaysOnTop:    "Toujours au premier plan",
	Language:       "Langue :",
	LanguageSystem: "Comme dans le système",
	Theme:          "Thème :",
	ThemeDark:      "Sombre",
	ThemeLight:     "Clair",

	AboutTitle:   func(name string) string { return "À propos de " + name },
	Version:      "Version",
	Author:       "Auteur :",
	License:      "Licence :",
	VisitWebsite: "Visiter le site",
	Close:        "Fermer",

	ExportTitle: "Exporter l'historique",
	CSVFiles:    "Fichiers CSV",
	SavedTo:     func(path string) string { return "Enregistré dans " + path },
	ConfigLoadFailed: func(err string) string {
		return "Impossible de lire la liste des hôtes, une nouvelle liste vide a été ouverte. Le fichier a été laissé tel quel :\n\n" + err
	},

	Install:   "Installer",
	Update:    "Mettre à jour",
	Uninstall: "Désinstaller",
	InstallAsk: func(dir string) string {
		return "AltPing sera copié dans " + dir + " et ajouté au menu Démarrer, au bureau et à la liste des applications installées. S'il y est déjà installé, il sera mis à jour. Les paramètres et les listes d'hôtes restent inchangés.\n\nAltPing redémarrera ensuite depuis cet emplacement."
	},
	InstallFailed: func(err string) string { return "Impossible d'installer AltPing :\n\n" + err },
	Installed:     "AltPing est installé",
	UninstallAsk: func(dir string) string {
		return "Supprimer AltPing de cet ordinateur ?\n\nLes paramètres et les listes d'hôtes restent dans " + dir
	},
	UninstallRunning: "AltPing est en cours d'exécution. Fermez-le et réessayez.",
	Uninstalled:      "AltPing a été supprimé.",
}

var es = Strings{
	Error: "Error",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d caídos)", n) },
	ToolConfigs:     "Configuraciones (O)",
	ToolAddHost:     "Añadir host (A)",
	ToolEditHost:    "Editar host (E)",
	ToolRemoveHosts: "Eliminar los hosts seleccionados (Del)",
	ToolDetails:     "Detalles (D)",
	ToolStart:       "Iniciar ping",
	ToolStop:        "Detener ping",
	ToolTray:        "Minimizar a la bandeja (T)",
	ToolShareOpen:   "Abrir los seleccionados en u00.io (U)",
	ToolStyle:       "Estilo (S)",
	KPIMedian:       "Ping mediano",
	KPIUptime:       "Disponibilidad 24 h",
	TrayShow:        "Mostrar",
	TrayQuit:        "Salir",
	Settings:        "Ajustes",
	Help:            "Ayuda",
	About:           "Acerca de",
	ModeUDP:         "Modo UDP",
	ModeICMP:        "Modo ICMP",
	StatusHosts:     "Hosts",
	StatusOK:        "OK",
	StatusSlow:      "Lentos",
	StatusDown:      "Caídos",
	StatusStopped:   "Ping detenido",

	Columns: ColumnStrings{
		Name: "Nombre", IP: "IP", Time: "Tiempo ms", Trend: "Últimos 5 min", Loss: "Pérdida %", Min: "Mín ms", Avg: "Media ms",
		Max: "Máx ms", Jitter: "Jitter ms", Since: "Desde", Details: "Detalles",
	},
	StateOK:         "OK",
	StateTimeout:    "TIEMPO AGOTADO",
	StateNoResolve:  "NO SE PUEDE RESOLVER EL NOMBRE",
	StatePortClosed: "PUERTO CERRADO",
	NoHosts:         "Aún no hay hosts",
	PressAToAdd:     "Pulse A para añadir un host",
	HowItWorks:      "Cómo funciona",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "h", Day: "d"},

	MenuEdit:       "Editar... (E)",
	MenuRemove:     "Eliminar (Del)",
	MenuDowntime:   "Caídas...",
	MenuExport:     "Exportar historial...",
	MenuChange:     "Cambiar seleccionados",
	BeepOn:         "Pitido activado",
	BeepOff:        "Pitido desactivado",
	CheckTCPPort:   "Comprobar puerto TCP",
	UsePing:        "Usar ping",
	PingEvery:      "Ping cada",
	Timeout:        "Tiempo de espera",
	SlowAbove:      "Lento a partir de",
	PortLabel:      "Puerto:",
	PingEveryMs:    "Ping cada, ms:",
	TimeoutMs:      "Tiempo de espera, ms:",
	SlowAboveMs:    "Lento a partir de, ms:",
	SlowAboveMsOff: "Lento a partir de, ms (0 - desactivado):",

	RemoveHostsTitle: "Eliminar hosts",
	RemoveHost:       "¿Eliminar el host?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("¿Eliminar %d hosts?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... y %d más", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Seleccione un host",
	PingHistory: "Historial de ping: ",
	Ms:          "ms",

	Name:          "Nombre:",
	Host:          "Host:",
	TCPPort:       "Puerto TCP:",
	CheckInstead:  "Comprobar en lugar de ping",
	Port:          "Puerto:",
	SlowTooltip:   "Mostrar el host en amarillo cuando su tiempo medio de ping lo supere. 0 - desactivado",
	Notify:        "Avisar:",
	BeepDownBack:  "Pitido al caer o volver",
	NotifyTooltip: "Caído significa 3 pings fallidos seguidos. El título de la ventana muestra cuántos de esos hosts están caídos.",
	Share:         "Compartir:",
	ShareOn:       "Publicar el ping en u00.io",
	ShareTooltip:  "Los tiempos de ping se envían a una página pública. Cualquiera con el enlace puede verlos.",
	ShareOpen:     "Abrir",
	ShareCopy:     "Copiar enlace",
	LinkCopied:    "Enlace copiado",
	MenuShareOpen: "Abrir en u00.io",
	MenuShareCopy: "Copiar enlace de u00.io",
	ShareStart:    "Publicar en u00.io",
	ShareStop:     "Dejar de publicar",

	DowntimeTitle: "Caídas",
	NoDowntime:    "Sin caídas en las últimas 24 horas",
	DownTimes: func(n int) string {
		return fmt.Sprintf("Caído %d %s en las últimas 24 horas:", n, i18n.Plural("es", n, "vez", "veces"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Caídas en las últimas 24 horas: %d", n) },
	Started:        "Inicio",
	Ended:          "Fin",
	Duration:       "Duración",
	StillDown:      "sigue caído",
	DateTimeLayout: "02/01 15:04:05",

	ConfigsTitle:       "Abrir configuración",
	Configs:            "Configuraciones:",
	NewConfig:          "Nueva...",
	SaveConfigAs:       "Guardar como...",
	RemoveConfig:       "Eliminar",
	HostsCount:         "N.º de hosts",
	Hosts:              "Hosts",
	NewConfigTitle:     "Nueva configuración",
	SaveConfigAsTitle:  "Guardar configuración como",
	CopySuffix:         " copia",
	RemoveConfigTitle:  "Eliminar configuración",
	RemoveConfigAsk:    "¿Eliminar la configuración seleccionada?",
	CannotRemoveOpened: "No se puede eliminar la configuración abierta actualmente.",

	ExtraColumns:   "Columnas adicionales:",
	ShowMin:        "Mín",
	ShowJitter:     "Jitter",
	ShowSince:      "Desde (último cambio)",
	AlwaysOnTop:    "Siempre visible",
	Language:       "Idioma:",
	LanguageSystem: "Como en el sistema",
	Theme:          "Tema:",
	ThemeDark:      "Oscuro",
	ThemeLight:     "Claro",

	AboutTitle:   func(name string) string { return "Acerca de " + name },
	Version:      "Versión",
	Author:       "Autor:",
	License:      "Licencia:",
	VisitWebsite: "Visitar el sitio web",
	Close:        "Cerrar",

	ExportTitle: "Exportar historial",
	CSVFiles:    "Archivos CSV",
	SavedTo:     func(path string) string { return "Guardado en " + path },
	ConfigLoadFailed: func(err string) string {
		return "No se pudo leer la lista de hosts; se abrió una nueva vacía. El archivo se dejó tal cual:\n\n" + err
	},

	Install:   "Instalar",
	Update:    "Actualizar",
	Uninstall: "Desinstalar",
	InstallAsk: func(dir string) string {
		return "AltPing se copiará en " + dir + " y se añadirá al menú Inicio, al escritorio y a la lista de aplicaciones instaladas. Si ya está instalado allí, se actualizará. La configuración y las listas de hosts se mantienen.\n\nDespués, AltPing se reiniciará desde allí."
	},
	InstallFailed: func(err string) string { return "No se pudo instalar AltPing:\n\n" + err },
	Installed:     "AltPing está instalado",
	UninstallAsk: func(dir string) string {
		return "¿Eliminar AltPing de este equipo?\n\nLa configuración y las listas de hosts se quedan en " + dir
	},
	UninstallRunning: "AltPing se está ejecutando. Ciérrelo e inténtelo de nuevo.",
	Uninstalled:      "AltPing se ha eliminado.",
}

var it = Strings{
	Error: "Errore",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d non raggiungibili)", n) },
	ToolConfigs:     "Configurazioni (O)",
	ToolAddHost:     "Aggiungi host (A)",
	ToolEditHost:    "Modifica host (E)",
	ToolRemoveHosts: "Rimuovi gli host selezionati (Del)",
	ToolDetails:     "Dettagli (D)",
	ToolStart:       "Avvia ping",
	ToolStop:        "Ferma ping",
	ToolTray:        "Riduci nell'area di notifica (T)",
	ToolShareOpen:   "Apri i selezionati su u00.io (U)",
	ToolStyle:       "Stile (S)",
	KPIMedian:       "Ping mediano",
	KPIUptime:       "Disponibilità 24 h",
	TrayShow:        "Mostra",
	TrayQuit:        "Esci",
	Settings:        "Impostazioni",
	Help:            "Aiuto",
	About:           "Informazioni",
	ModeUDP:         "Modalità UDP",
	ModeICMP:        "Modalità ICMP",
	StatusHosts:     "Host",
	StatusOK:        "OK",
	StatusSlow:      "Lenti",
	StatusDown:      "Non raggiungibili",
	StatusStopped:   "Ping fermato",

	Columns: ColumnStrings{
		Name: "Nome", IP: "IP", Time: "Tempo ms", Trend: "Ultimi 5 min", Loss: "Perdita %", Min: "Min ms", Avg: "Media ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Da", Details: "Dettagli",
	},
	StateOK:         "OK",
	StateTimeout:    "TIMEOUT",
	StateNoResolve:  "NOME HOST NON RISOLTO",
	StatePortClosed: "PORTA CHIUSA",
	NoHosts:         "Nessun host",
	PressAToAdd:     "Premi A per aggiungere un host",
	HowItWorks:      "Come funziona",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "h", Day: "g"},

	MenuEdit:       "Modifica... (E)",
	MenuRemove:     "Rimuovi (Del)",
	MenuDowntime:   "Interruzioni...",
	MenuExport:     "Esporta cronologia...",
	MenuChange:     "Modifica selezionati",
	BeepOn:         "Segnale acustico attivo",
	BeepOff:        "Segnale acustico disattivo",
	CheckTCPPort:   "Controlla porta TCP",
	UsePing:        "Usa ping",
	PingEvery:      "Ping ogni",
	Timeout:        "Timeout",
	SlowAbove:      "Lento oltre",
	PortLabel:      "Porta:",
	PingEveryMs:    "Ping ogni, ms:",
	TimeoutMs:      "Timeout, ms:",
	SlowAboveMs:    "Lento oltre, ms:",
	SlowAboveMsOff: "Lento oltre, ms (0 - disattivato):",

	RemoveHostsTitle: "Rimuovi host",
	RemoveHost:       "Rimuovere l'host?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("Rimuovere %d host?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... e altri %d", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Seleziona un host",
	PingHistory: "Cronologia ping: ",
	Ms:          "ms",

	Name:          "Nome:",
	Host:          "Host:",
	TCPPort:       "Porta TCP:",
	CheckInstead:  "Controlla invece del ping",
	Port:          "Porta:",
	SlowTooltip:   "Mostra l'host in giallo quando il tempo medio di ping lo supera. 0 - disattivato",
	Notify:        "Notifica:",
	BeepDownBack:  "Segnale quando cade o torna",
	NotifyTooltip: "Caduto significa 3 ping falliti di fila. Il titolo della finestra mostra quanti di questi host sono caduti.",
	Share:         "Condividi:",
	ShareOn:       "Pubblica il ping su u00.io",
	ShareTooltip:  "I tempi di ping vengono inviati a una pagina pubblica. Chiunque abbia il link può vederli.",
	ShareOpen:     "Apri",
	ShareCopy:     "Copia link",
	LinkCopied:    "Link copiato",
	MenuShareOpen: "Apri su u00.io",
	MenuShareCopy: "Copia link u00.io",
	ShareStart:    "Pubblica su u00.io",
	ShareStop:     "Interrompi pubblicazione",

	DowntimeTitle: "Interruzioni",
	NoDowntime:    "Nessuna interruzione nelle ultime 24 ore",
	DownTimes: func(n int) string {
		return fmt.Sprintf("Non raggiungibile %d %s nelle ultime 24 ore:", n, i18n.Plural("it", n, "volta", "volte"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Interruzioni nelle ultime 24 ore: %d", n) },
	Started:        "Inizio",
	Ended:          "Fine",
	Duration:       "Durata",
	StillDown:      "ancora non raggiungibile",
	DateTimeLayout: "02/01 15:04:05",

	ConfigsTitle:       "Apri configurazione",
	Configs:            "Configurazioni:",
	NewConfig:          "Nuova...",
	SaveConfigAs:       "Salva come...",
	RemoveConfig:       "Rimuovi",
	HostsCount:         "Numero host",
	Hosts:              "Host",
	NewConfigTitle:     "Nuova configurazione",
	SaveConfigAsTitle:  "Salva configurazione come",
	CopySuffix:         " copia",
	RemoveConfigTitle:  "Rimuovi configurazione",
	RemoveConfigAsk:    "Rimuovere la configurazione selezionata?",
	CannotRemoveOpened: "Impossibile rimuovere la configurazione attualmente aperta.",

	ExtraColumns:   "Colonne aggiuntive:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Da (ultimo cambiamento)",
	AlwaysOnTop:    "Sempre in primo piano",
	Language:       "Lingua:",
	LanguageSystem: "Come nel sistema",
	Theme:          "Tema:",
	ThemeDark:      "Scuro",
	ThemeLight:     "Chiaro",

	AboutTitle:   func(name string) string { return "Informazioni su " + name },
	Version:      "Versione",
	Author:       "Autore:",
	License:      "Licenza:",
	VisitWebsite: "Visita il sito",
	Close:        "Chiudi",

	ExportTitle: "Esporta cronologia",
	CSVFiles:    "File CSV",
	SavedTo:     func(path string) string { return "Salvato in " + path },
	ConfigLoadFailed: func(err string) string {
		return "Impossibile leggere l'elenco degli host, ne è stato aperto uno nuovo vuoto. Il file è stato lasciato invariato:\n\n" + err
	},

	Install:   "Installa",
	Update:    "Aggiorna",
	Uninstall: "Disinstalla",
	InstallAsk: func(dir string) string {
		return "AltPing verrà copiato in " + dir + " e aggiunto al menu Start, al desktop e all'elenco delle app installate. Se è già installato lì, verrà aggiornato. Le impostazioni e gli elenchi degli host restano invariati.\n\nPoi AltPing verrà riavviato da lì."
	},
	InstallFailed: func(err string) string { return "Impossibile installare AltPing:\n\n" + err },
	Installed:     "AltPing è installato",
	UninstallAsk: func(dir string) string {
		return "Rimuovere AltPing da questo computer?\n\nLe impostazioni e gli elenchi degli host restano in " + dir
	},
	UninstallRunning: "AltPing è in esecuzione. Chiudilo e riprova.",
	Uninstalled:      "AltPing è stato rimosso.",
}

var pt = Strings{
	Error: "Erro",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d fora do ar)", n) },
	ToolConfigs:     "Configurações (O)",
	ToolAddHost:     "Adicionar host (A)",
	ToolEditHost:    "Editar host (E)",
	ToolRemoveHosts: "Remover os hosts selecionados (Del)",
	ToolDetails:     "Detalhes (D)",
	ToolStart:       "Iniciar ping",
	ToolStop:        "Parar ping",
	ToolTray:        "Minimizar para a bandeja (T)",
	ToolShareOpen:   "Abrir os selecionados no u00.io (U)",
	ToolStyle:       "Estilo (S)",
	KPIMedian:       "Ping mediano",
	KPIUptime:       "Disponibilidade 24 h",
	TrayShow:        "Mostrar",
	TrayQuit:        "Sair",
	Settings:        "Configurações",
	Help:            "Ajuda",
	About:           "Sobre",
	ModeUDP:         "Modo UDP",
	ModeICMP:        "Modo ICMP",
	StatusHosts:     "Hosts",
	StatusOK:        "OK",
	StatusSlow:      "Lentos",
	StatusDown:      "Fora do ar",
	StatusStopped:   "Ping parado",

	Columns: ColumnStrings{
		Name: "Nome", IP: "IP", Time: "Tempo ms", Trend: "Últimos 5 min", Loss: "Perda %", Min: "Mín ms", Avg: "Média ms",
		Max: "Máx ms", Jitter: "Jitter ms", Since: "Desde", Details: "Detalhes",
	},
	StateOK:         "OK",
	StateTimeout:    "TEMPO ESGOTADO",
	StateNoResolve:  "NÃO FOI POSSÍVEL RESOLVER O NOME",
	StatePortClosed: "PORTA FECHADA",
	NoHosts:         "Ainda não há hosts",
	PressAToAdd:     "Pressione A para adicionar um host",
	HowItWorks:      "Como funciona",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "h", Day: "d"},

	MenuEdit:       "Editar... (E)",
	MenuRemove:     "Remover (Del)",
	MenuDowntime:   "Quedas...",
	MenuExport:     "Exportar histórico...",
	MenuChange:     "Alterar selecionados",
	BeepOn:         "Bipe ligado",
	BeepOff:        "Bipe desligado",
	CheckTCPPort:   "Verificar porta TCP",
	UsePing:        "Usar ping",
	PingEvery:      "Ping a cada",
	Timeout:        "Tempo limite",
	SlowAbove:      "Lento acima de",
	PortLabel:      "Porta:",
	PingEveryMs:    "Ping a cada, ms:",
	TimeoutMs:      "Tempo limite, ms:",
	SlowAboveMs:    "Lento acima de, ms:",
	SlowAboveMsOff: "Lento acima de, ms (0 - desligado):",

	RemoveHostsTitle: "Remover hosts",
	RemoveHost:       "Remover o host?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("Remover %d hosts?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... e mais %d", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Selecione um host",
	PingHistory: "Histórico de ping: ",
	Ms:          "ms",

	Name:          "Nome:",
	Host:          "Host:",
	TCPPort:       "Porta TCP:",
	CheckInstead:  "Verificar em vez do ping",
	Port:          "Porta:",
	SlowTooltip:   "Mostrar o host em amarelo quando o tempo médio de ping passar deste valor. 0 - desligado",
	Notify:        "Notificar:",
	BeepDownBack:  "Bipe ao cair ou voltar",
	NotifyTooltip: "Fora do ar significa 3 pings falhos seguidos. O título da janela mostra quantos desses hosts estão fora do ar.",
	Share:         "Compartilhar:",
	ShareOn:       "Publicar o ping no u00.io",
	ShareTooltip:  "Os tempos de ping são enviados para uma página pública. Qualquer pessoa com o link pode vê-los.",
	ShareOpen:     "Abrir",
	ShareCopy:     "Copiar link",
	LinkCopied:    "Link copiado",
	MenuShareOpen: "Abrir no u00.io",
	MenuShareCopy: "Copiar link do u00.io",
	ShareStart:    "Publicar no u00.io",
	ShareStop:     "Parar de publicar",

	DowntimeTitle: "Quedas",
	NoDowntime:    "Nenhuma queda nas últimas 24 horas",
	DownTimes: func(n int) string {
		return fmt.Sprintf("Fora do ar %d %s nas últimas 24 horas:", n, i18n.Plural("pt", n, "vez", "vezes"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Quedas nas últimas 24 horas: %d", n) },
	Started:        "Início",
	Ended:          "Fim",
	Duration:       "Duração",
	StillDown:      "ainda fora do ar",
	DateTimeLayout: "02/01 15:04:05",

	ConfigsTitle:       "Abrir configuração",
	Configs:            "Configurações:",
	NewConfig:          "Nova...",
	SaveConfigAs:       "Salvar como...",
	RemoveConfig:       "Remover",
	HostsCount:         "Nº de hosts",
	Hosts:              "Hosts",
	NewConfigTitle:     "Nova configuração",
	SaveConfigAsTitle:  "Salvar configuração como",
	CopySuffix:         " cópia",
	RemoveConfigTitle:  "Remover configuração",
	RemoveConfigAsk:    "Remover a configuração selecionada?",
	CannotRemoveOpened: "Não é possível remover a configuração aberta no momento.",

	ExtraColumns:   "Colunas extras:",
	ShowMin:        "Mín",
	ShowJitter:     "Jitter",
	ShowSince:      "Desde (última mudança)",
	AlwaysOnTop:    "Sempre visível",
	Language:       "Idioma:",
	LanguageSystem: "Como no sistema",
	Theme:          "Tema:",
	ThemeDark:      "Escuro",
	ThemeLight:     "Claro",

	AboutTitle:   func(name string) string { return "Sobre o " + name },
	Version:      "Versão",
	Author:       "Autor:",
	License:      "Licença:",
	VisitWebsite: "Visitar o site",
	Close:        "Fechar",

	ExportTitle: "Exportar histórico",
	CSVFiles:    "Arquivos CSV",
	SavedTo:     func(path string) string { return "Salvo em " + path },
	ConfigLoadFailed: func(err string) string {
		return "Não foi possível ler a lista de hosts; uma nova lista vazia foi aberta. O arquivo foi mantido como está:\n\n" + err
	},

	Install:   "Instalar",
	Update:    "Atualizar",
	Uninstall: "Desinstalar",
	InstallAsk: func(dir string) string {
		return "O AltPing será copiado para " + dir + " e adicionado ao menu Iniciar, à área de trabalho e à lista de aplicativos instalados. Se já estiver instalado lá, será atualizado. As configurações e as listas de hosts permanecem como estão.\n\nEm seguida, o AltPing será reiniciado a partir de lá."
	},
	InstallFailed: func(err string) string { return "Não foi possível instalar o AltPing:\n\n" + err },
	Installed:     "O AltPing está instalado",
	UninstallAsk: func(dir string) string {
		return "Remover o AltPing deste computador?\n\nAs configurações e as listas de hosts permanecem em " + dir
	},
	UninstallRunning: "O AltPing está em execução. Feche-o e tente novamente.",
	Uninstalled:      "O AltPing foi removido.",
}

var zh = Strings{
	Error: "错误",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d 个中断)", n) },
	ToolConfigs:     "配置 (O)",
	ToolAddHost:     "添加主机 (A)",
	ToolEditHost:    "编辑主机 (E)",
	ToolRemoveHosts: "删除所选主机 (Del)",
	ToolDetails:     "详情 (D)",
	ToolStart:       "开始 ping",
	ToolStop:        "停止 ping",
	ToolTray:        "最小化到托盘 (T)",
	ToolShareOpen:   "在 u00.io 上打开所选主机 (U)",
	ToolStyle:       "样式 (S)",
	KPIMedian:       "延迟中位数",
	KPIUptime:       "24 小时可用率",
	TrayShow:        "显示",
	TrayQuit:        "退出",
	Settings:        "设置",
	Help:            "帮助",
	About:           "关于",
	ModeUDP:         "UDP 模式",
	ModeICMP:        "ICMP 模式",
	StatusHosts:     "主机",
	StatusOK:        "正常",
	StatusSlow:      "缓慢",
	StatusDown:      "不可达",
	StatusStopped:   "已停止 ping",

	Columns: ColumnStrings{
		Name: "名称", IP: "IP", Time: "时间 ms", Trend: "最近 5 分钟", Loss: "丢包 %", Min: "最小 ms", Avg: "平均 ms",
		Max: "最大 ms", Jitter: "抖动 ms", Since: "持续", Details: "状态",
	},
	StateOK:         "正常",
	StateTimeout:    "超时",
	StateNoResolve:  "无法解析主机名",
	StatePortClosed: "端口关闭",
	NoHosts:         "还没有主机",
	PressAToAdd:     "按 A 添加主机",
	HowItWorks:      "使用说明",

	Units: UnitStrings{Second: "秒", Minute: "分", Hour: "时", Day: "天"},

	MenuEdit:       "编辑... (E)",
	MenuRemove:     "删除 (Del)",
	MenuDowntime:   "中断记录...",
	MenuExport:     "导出历史...",
	MenuChange:     "更改所选",
	BeepOn:         "开启提示音",
	BeepOff:        "关闭提示音",
	CheckTCPPort:   "检查 TCP 端口",
	UsePing:        "使用 ping",
	PingEvery:      "ping 间隔",
	Timeout:        "超时",
	SlowAbove:      "慢速阈值",
	PortLabel:      "端口：",
	PingEveryMs:    "ping 间隔，毫秒：",
	TimeoutMs:      "超时，毫秒：",
	SlowAboveMs:    "慢速阈值，毫秒：",
	SlowAboveMsOff: "慢速阈值，毫秒（0 - 关闭）：",

	RemoveHostsTitle: "删除主机",
	RemoveHost:       "删除该主机？",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("删除 %d 个主机？", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... 以及另外 %d 个", n) },

	Periods:     [3]string{"5 分", "1 时", "24 时"},
	SelectAHost: "请选择主机",
	PingHistory: "ping 历史：",
	Ms:          "ms",

	Name:          "名称：",
	Host:          "主机：",
	TCPPort:       "TCP 端口：",
	CheckInstead:  "检查端口代替 ping",
	Port:          "端口：",
	SlowTooltip:   "平均 ping 时间超过该值时主机显示为黄色。0 - 关闭",
	Notify:        "通知：",
	BeepDownBack:  "中断或恢复时发出提示音",
	NotifyTooltip: "中断指连续 3 次 ping 失败。窗口标题显示有多少这样的主机中断。",
	Share:         "分享：",
	ShareOn:       "在 u00.io 上发布 ping",
	ShareTooltip:  "ping 时间会发送到一个公开页面，任何拥有链接的人都可以查看。",
	ShareOpen:     "打开",
	ShareCopy:     "复制链接",
	LinkCopied:    "链接已复制",
	MenuShareOpen: "在 u00.io 上打开",
	MenuShareCopy: "复制 u00.io 链接",
	ShareStart:    "在 u00.io 上发布",
	ShareStop:     "停止发布",

	DowntimeTitle:  "中断记录",
	NoDowntime:     "过去 24 小时内没有中断",
	DownTimes:      func(n int) string { return fmt.Sprintf("过去 24 小时内中断 %d 次：", n) },
	GroupDownTimes: func(n int) string { return fmt.Sprintf("过去 24 小时内的中断：%d", n) },
	Started:        "开始",
	Ended:          "结束",
	Duration:       "时长",
	StillDown:      "仍在中断",
	DateTimeLayout: "01-02 15:04:05",

	ConfigsTitle:       "打开配置",
	Configs:            "配置：",
	NewConfig:          "新建...",
	SaveConfigAs:       "另存为...",
	RemoveConfig:       "删除",
	HostsCount:         "主机数",
	Hosts:              "主机",
	NewConfigTitle:     "新建配置",
	SaveConfigAsTitle:  "配置另存为",
	CopySuffix:         " 副本",
	RemoveConfigTitle:  "删除配置",
	RemoveConfigAsk:    "删除所选配置？",
	CannotRemoveOpened: "无法删除当前打开的配置。",

	ExtraColumns:   "附加列：",
	ShowMin:        "最小",
	ShowJitter:     "抖动",
	ShowSince:      "持续（上次变化）",
	AlwaysOnTop:    "窗口置顶",
	Language:       "语言：",
	LanguageSystem: "跟随系统",
	Theme:          "主题：",
	ThemeDark:      "深色",
	ThemeLight:     "浅色",

	AboutTitle:   func(name string) string { return "关于 " + name },
	Version:      "版本",
	Author:       "作者：",
	License:      "许可证：",
	VisitWebsite: "访问网站",
	Close:        "关闭",

	ExportTitle: "导出历史",
	CSVFiles:    "CSV 文件",
	SavedTo:     func(path string) string { return "已保存到 " + path },
	ConfigLoadFailed: func(err string) string {
		return "无法读取主机列表，已打开一个新的空列表。原文件保持不变：\n\n" + err
	},

	Install:   "安装",
	Update:    "更新",
	Uninstall: "卸载",
	InstallAsk: func(dir string) string {
		return "AltPing 将被复制到 " + dir + "，并添加到“开始”菜单、桌面和已安装应用列表中。如果已在该位置安装，则会进行更新。设置和主机列表保持不变。\n\n随后 AltPing 将从该位置重新启动。"
	},
	InstallFailed: func(err string) string { return "无法安装 AltPing：\n\n" + err },
	Installed:     "AltPing 已安装",
	UninstallAsk: func(dir string) string {
		return "要从此计算机中删除 AltPing 吗？\n\n设置和主机列表将保留在 " + dir
	},
	UninstallRunning: "AltPing 正在运行。请关闭它后重试。",
	Uninstalled:      "AltPing 已删除。",
}

var ja = Strings{
	Error: "エラー",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d 件ダウン)", n) },
	ToolConfigs:     "構成 (O)",
	ToolAddHost:     "ホストを追加 (A)",
	ToolEditHost:    "ホストを編集 (E)",
	ToolRemoveHosts: "選択したホストを削除 (Del)",
	ToolDetails:     "詳細 (D)",
	ToolStart:       "ping を開始",
	ToolStop:        "ping を停止",
	ToolTray:        "トレイに最小化 (T)",
	ToolShareOpen:   "選択したホストを u00.io で開く (U)",
	ToolStyle:       "スタイル (S)",
	KPIMedian:       "Ping 中央値",
	KPIUptime:       "24 時間の稼働率",
	TrayShow:        "表示",
	TrayQuit:        "終了",
	Settings:        "設定",
	Help:            "ヘルプ",
	About:           "バージョン情報",
	ModeUDP:         "UDP モード",
	ModeICMP:        "ICMP モード",
	StatusHosts:     "ホスト",
	StatusOK:        "正常",
	StatusSlow:      "遅い",
	StatusDown:      "応答なし",
	StatusStopped:   "ping 停止中",

	Columns: ColumnStrings{
		Name: "名前", IP: "IP", Time: "時間 ms", Trend: "直近 5 分", Loss: "損失 %", Min: "最小 ms", Avg: "平均 ms",
		Max: "最大 ms", Jitter: "ジッター ms", Since: "経過", Details: "状態",
	},
	StateOK:         "OK",
	StateTimeout:    "タイムアウト",
	StateNoResolve:  "ホスト名を解決できません",
	StatePortClosed: "ポートが閉じています",
	NoHosts:         "ホストがありません",
	PressAToAdd:     "A を押してホストを追加",
	HowItWorks:      "使い方",

	Units: UnitStrings{Second: "秒", Minute: "分", Hour: "時間", Day: "日"},

	MenuEdit:       "編集... (E)",
	MenuRemove:     "削除 (Del)",
	MenuDowntime:   "ダウンタイム...",
	MenuExport:     "履歴をエクスポート...",
	MenuChange:     "選択項目を変更",
	BeepOn:         "ビープ音オン",
	BeepOff:        "ビープ音オフ",
	CheckTCPPort:   "TCP ポートを確認",
	UsePing:        "ping を使用",
	PingEvery:      "ping 間隔",
	Timeout:        "タイムアウト",
	SlowAbove:      "低速のしきい値",
	PortLabel:      "ポート:",
	PingEveryMs:    "ping 間隔、ミリ秒:",
	TimeoutMs:      "タイムアウト、ミリ秒:",
	SlowAboveMs:    "低速のしきい値、ミリ秒:",
	SlowAboveMsOff: "低速のしきい値、ミリ秒 (0 - オフ):",

	RemoveHostsTitle: "ホストの削除",
	RemoveHost:       "ホストを削除しますか？",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("%d 件のホストを削除しますか？", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... 他 %d 件", n) },

	Periods:     [3]string{"5 分", "1 時間", "24 時間"},
	SelectAHost: "ホストを選択してください",
	PingHistory: "ping 履歴: ",
	Ms:          "ms",

	Name:          "名前:",
	Host:          "ホスト:",
	TCPPort:       "TCP ポート:",
	CheckInstead:  "ping の代わりに確認",
	Port:          "ポート:",
	SlowTooltip:   "平均 ping 時間がこの値を超えるとホストを黄色で表示します。0 - オフ",
	Notify:        "通知:",
	BeepDownBack:  "ダウン時と復旧時にビープ音",
	NotifyTooltip: "ダウンとは ping が 3 回連続で失敗した状態です。ウィンドウのタイトルにダウン中のホスト数が表示されます。",
	Share:         "共有:",
	ShareOn:       "u00.io で ping を公開",
	ShareTooltip:  "ping の時間は公開ページに送信されます。リンクを知っている人は誰でも見られます。",
	ShareOpen:     "開く",
	ShareCopy:     "リンクをコピー",
	LinkCopied:    "リンクをコピーしました",
	MenuShareOpen: "u00.io で開く",
	MenuShareCopy: "u00.io のリンクをコピー",
	ShareStart:    "u00.io で公開",
	ShareStop:     "公開を停止",

	DowntimeTitle:  "ダウンタイム",
	NoDowntime:     "過去 24 時間にダウンはありません",
	DownTimes:      func(n int) string { return fmt.Sprintf("過去 24 時間に %d 回ダウン:", n) },
	GroupDownTimes: func(n int) string { return fmt.Sprintf("過去 24 時間のダウン: %d", n) },
	Started:        "開始",
	Ended:          "終了",
	Duration:       "期間",
	StillDown:      "ダウン中",
	DateTimeLayout: "01/02 15:04:05",

	ConfigsTitle:       "構成を開く",
	Configs:            "構成:",
	NewConfig:          "新規...",
	SaveConfigAs:       "名前を付けて保存...",
	RemoveConfig:       "削除",
	HostsCount:         "ホスト数",
	Hosts:              "ホスト",
	NewConfigTitle:     "新しい構成",
	SaveConfigAsTitle:  "構成に名前を付けて保存",
	CopySuffix:         " のコピー",
	RemoveConfigTitle:  "構成の削除",
	RemoveConfigAsk:    "選択した構成を削除しますか？",
	CannotRemoveOpened: "現在開いている構成は削除できません。",

	ExtraColumns:   "追加の列:",
	ShowMin:        "最小",
	ShowJitter:     "ジッター",
	ShowSince:      "経過 (最後の変化から)",
	AlwaysOnTop:    "常に手前に表示",
	Language:       "言語:",
	LanguageSystem: "システムに従う",
	Theme:          "テーマ:",
	ThemeDark:      "ダーク",
	ThemeLight:     "ライト",

	AboutTitle:   func(name string) string { return name + " について" },
	Version:      "バージョン",
	Author:       "作者:",
	License:      "ライセンス:",
	VisitWebsite: "Web サイトを開く",
	Close:        "閉じる",

	ExportTitle: "履歴のエクスポート",
	CSVFiles:    "CSV ファイル",
	SavedTo:     func(path string) string { return path + " に保存しました" },
	ConfigLoadFailed: func(err string) string {
		return "ホスト一覧を読み込めなかったため、新しい空の一覧を開きました。ファイルはそのまま残してあります:\n\n" + err
	},

	Install:   "インストール",
	Update:    "更新",
	Uninstall: "アンインストール",
	InstallAsk: func(dir string) string {
		return "AltPing を " + dir + " にコピーし、スタート メニュー、デスクトップ、インストール済みアプリの一覧に追加します。すでにインストールされている場合は更新されます。設定とホストの一覧はそのまま残ります。\n\nその後、AltPing はそこから再起動します。"
	},
	InstallFailed: func(err string) string { return "AltPing をインストールできませんでした:\n\n" + err },
	Installed:     "AltPing をインストールしました",
	UninstallAsk: func(dir string) string {
		return "このコンピューターから AltPing を削除しますか?\n\n設定とホストの一覧は " + dir + " に残ります。"
	},
	UninstallRunning: "AltPing が実行中です。閉じてからもう一度お試しください。",
	Uninstalled:      "AltPing を削除しました。",
}

var ko = Strings{
	Error: "오류",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d개 다운)", n) },
	ToolConfigs:     "구성 (O)",
	ToolAddHost:     "호스트 추가 (A)",
	ToolEditHost:    "호스트 편집 (E)",
	ToolRemoveHosts: "선택한 호스트 삭제 (Del)",
	ToolDetails:     "자세히 (D)",
	ToolStart:       "ping 시작",
	ToolStop:        "ping 중지",
	ToolTray:        "트레이로 최소화 (T)",
	ToolShareOpen:   "선택한 호스트를 u00.io에서 열기 (U)",
	ToolStyle:       "스타일 (S)",
	KPIMedian:       "핑 중앙값",
	KPIUptime:       "24시간 가용성",
	TrayShow:        "표시",
	TrayQuit:        "종료",
	Settings:        "설정",
	Help:            "도움말",
	About:           "정보",
	ModeUDP:         "UDP 모드",
	ModeICMP:        "ICMP 모드",
	StatusHosts:     "호스트",
	StatusOK:        "정상",
	StatusSlow:      "느림",
	StatusDown:      "응답 없음",
	StatusStopped:   "ping 중지됨",

	Columns: ColumnStrings{
		Name: "이름", IP: "IP", Time: "시간 ms", Trend: "최근 5분", Loss: "손실 %", Min: "최소 ms", Avg: "평균 ms",
		Max: "최대 ms", Jitter: "지터 ms", Since: "경과", Details: "상태",
	},
	StateOK:         "정상",
	StateTimeout:    "시간 초과",
	StateNoResolve:  "호스트 이름을 확인할 수 없음",
	StatePortClosed: "포트 닫힘",
	NoHosts:         "아직 호스트가 없습니다",
	PressAToAdd:     "A를 눌러 호스트를 추가하세요",
	HowItWorks:      "사용 방법",

	Units: UnitStrings{Second: "초", Minute: "분", Hour: "시간", Day: "일"},

	MenuEdit:       "편집... (E)",
	MenuRemove:     "삭제 (Del)",
	MenuDowntime:   "다운타임...",
	MenuExport:     "기록 내보내기...",
	MenuChange:     "선택 항목 변경",
	BeepOn:         "알림음 켜기",
	BeepOff:        "알림음 끄기",
	CheckTCPPort:   "TCP 포트 확인",
	UsePing:        "ping 사용",
	PingEvery:      "ping 간격",
	Timeout:        "시간 제한",
	SlowAbove:      "느림 기준",
	PortLabel:      "포트:",
	PingEveryMs:    "ping 간격, ms:",
	TimeoutMs:      "시간 제한, ms:",
	SlowAboveMs:    "느림 기준, ms:",
	SlowAboveMsOff: "느림 기준, ms (0 - 끄기):",

	RemoveHostsTitle: "호스트 삭제",
	RemoveHost:       "호스트를 삭제하시겠습니까?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("호스트 %d개를 삭제하시겠습니까?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... 외 %d개", n) },

	Periods:     [3]string{"5분", "1시간", "24시간"},
	SelectAHost: "호스트를 선택하세요",
	PingHistory: "ping 기록: ",
	Ms:          "ms",

	Name:          "이름:",
	Host:          "호스트:",
	TCPPort:       "TCP 포트:",
	CheckInstead:  "ping 대신 확인",
	Port:          "포트:",
	SlowTooltip:   "평균 ping 시간이 이 값을 넘으면 호스트를 노란색으로 표시합니다. 0 - 끄기",
	Notify:        "알림:",
	BeepDownBack:  "다운 및 복구 시 알림음",
	NotifyTooltip: "다운은 ping이 3번 연속 실패한 상태입니다. 창 제목에 다운된 호스트 수가 표시됩니다.",
	Share:         "공유:",
	ShareOn:       "u00.io에 ping 게시",
	ShareTooltip:  "ping 시간이 공개 페이지로 전송됩니다. 링크가 있는 사람은 누구나 볼 수 있습니다.",
	ShareOpen:     "열기",
	ShareCopy:     "링크 복사",
	LinkCopied:    "링크가 복사되었습니다",
	MenuShareOpen: "u00.io에서 열기",
	MenuShareCopy: "u00.io 링크 복사",
	ShareStart:    "u00.io에 게시",
	ShareStop:     "게시 중지",

	DowntimeTitle:  "다운타임",
	NoDowntime:     "지난 24시간 동안 다운 없음",
	DownTimes:      func(n int) string { return fmt.Sprintf("지난 24시간 동안 %d회 다운:", n) },
	GroupDownTimes: func(n int) string { return fmt.Sprintf("지난 24시간 동안의 다운: %d", n) },
	Started:        "시작",
	Ended:          "종료",
	Duration:       "지속 시간",
	StillDown:      "다운 중",
	DateTimeLayout: "01.02 15:04:05",

	ConfigsTitle:       "구성 열기",
	Configs:            "구성:",
	NewConfig:          "새로 만들기...",
	SaveConfigAs:       "다른 이름으로 저장...",
	RemoveConfig:       "삭제",
	HostsCount:         "호스트 수",
	Hosts:              "호스트",
	NewConfigTitle:     "새 구성",
	SaveConfigAsTitle:  "구성을 다른 이름으로 저장",
	CopySuffix:         " 사본",
	RemoveConfigTitle:  "구성 삭제",
	RemoveConfigAsk:    "선택한 구성을 삭제하시겠습니까?",
	CannotRemoveOpened: "현재 열려 있는 구성은 삭제할 수 없습니다.",

	ExtraColumns:   "추가 열:",
	ShowMin:        "최소",
	ShowJitter:     "지터",
	ShowSince:      "경과 (마지막 변경 후)",
	AlwaysOnTop:    "항상 위에 표시",
	Language:       "언어:",
	LanguageSystem: "시스템 설정 따름",
	Theme:          "테마:",
	ThemeDark:      "어둡게",
	ThemeLight:     "밝게",

	AboutTitle:   func(name string) string { return name + " 정보" },
	Version:      "버전",
	Author:       "만든 이:",
	License:      "라이선스:",
	VisitWebsite: "웹사이트 방문",
	Close:        "닫기",

	ExportTitle: "기록 내보내기",
	CSVFiles:    "CSV 파일",
	SavedTo:     func(path string) string { return path + "에 저장됨" },
	ConfigLoadFailed: func(err string) string {
		return "호스트 목록을 읽을 수 없어 새 빈 목록을 열었습니다. 파일은 그대로 두었습니다:\n\n" + err
	},

	Install:   "설치",
	Update:    "업데이트",
	Uninstall: "제거",
	InstallAsk: func(dir string) string {
		return "AltPing이 " + dir + "에 복사되고 시작 메뉴, 바탕 화면, 설치된 앱 목록에 추가됩니다. 이미 설치되어 있으면 업데이트됩니다. 설정과 호스트 목록은 그대로 유지됩니다.\n\n그런 다음 AltPing이 그 위치에서 다시 시작됩니다."
	},
	InstallFailed: func(err string) string { return "AltPing을 설치할 수 없습니다:\n\n" + err },
	Installed:     "AltPing이 설치되었습니다",
	UninstallAsk: func(dir string) string {
		return "이 컴퓨터에서 AltPing을 제거하시겠습니까?\n\n설정과 호스트 목록은 " + dir + "에 남아 있습니다."
	},
	UninstallRunning: "AltPing이 실행 중입니다. 닫은 후 다시 시도하세요.",
	Uninstalled:      "AltPing이 제거되었습니다.",
}

var nl = Strings{
	Error: "Fout",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d offline)", n) },
	ToolConfigs:     "Configuraties (O)",
	ToolAddHost:     "Host toevoegen (A)",
	ToolEditHost:    "Host bewerken (E)",
	ToolRemoveHosts: "Geselecteerde hosts verwijderen (Del)",
	ToolDetails:     "Details (D)",
	ToolStart:       "Pingen starten",
	ToolStop:        "Pingen stoppen",
	ToolTray:        "Minimaliseren naar systeemvak (T)",
	ToolShareOpen:   "Geselecteerde hosts openen op u00.io (U)",
	ToolStyle:       "Stijl (S)",
	KPIMedian:       "Mediane ping",
	KPIUptime:       "Beschikbaarheid 24 u",
	TrayShow:        "Weergeven",
	TrayQuit:        "Afsluiten",
	Settings:        "Instellingen",
	Help:            "Help",
	About:           "Over",
	ModeUDP:         "UDP-modus",
	ModeICMP:        "ICMP-modus",
	StatusHosts:     "Hosts",
	StatusOK:        "OK",
	StatusSlow:      "Traag",
	StatusDown:      "Offline",
	StatusStopped:   "Pingen gestopt",

	Columns: ColumnStrings{
		Name: "Naam", IP: "IP", Time: "Tijd ms", Trend: "Laatste 5 min", Loss: "Verlies %", Min: "Min ms", Avg: "Gem ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Sinds", Details: "Details",
	},
	StateOK:         "OK",
	StateTimeout:    "TIME-OUT",
	StateNoResolve:  "KAN HOSTNAAM NIET OMZETTEN",
	StatePortClosed: "POORT GESLOTEN",
	NoHosts:         "Nog geen hosts",
	PressAToAdd:     "Druk op A om een host toe te voegen",
	HowItWorks:      "Hoe het werkt",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "u", Day: "d"},

	MenuEdit:       "Bewerken... (E)",
	MenuRemove:     "Verwijderen (Del)",
	MenuDowntime:   "Uitvaltijd...",
	MenuExport:     "Geschiedenis exporteren...",
	MenuChange:     "Selectie wijzigen",
	BeepOn:         "Pieptoon aan",
	BeepOff:        "Pieptoon uit",
	CheckTCPPort:   "TCP-poort controleren",
	UsePing:        "Ping gebruiken",
	PingEvery:      "Ping elke",
	Timeout:        "Time-out",
	SlowAbove:      "Traag boven",
	PortLabel:      "Poort:",
	PingEveryMs:    "Ping elke, ms:",
	TimeoutMs:      "Time-out, ms:",
	SlowAboveMs:    "Traag boven, ms:",
	SlowAboveMsOff: "Traag boven, ms (0 - uit):",

	RemoveHostsTitle: "Hosts verwijderen",
	RemoveHost:       "De host verwijderen?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("%d %s verwijderen?", n, i18n.Plural("nl", n, "host", "hosts"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... en nog %d", n) },

	Periods:     [3]string{"5m", "1u", "24u"},
	SelectAHost: "Selecteer een host",
	PingHistory: "Pinggeschiedenis: ",
	Ms:          "ms",

	Name:          "Naam:",
	Host:          "Host:",
	TCPPort:       "TCP-poort:",
	CheckInstead:  "Controleren in plaats van pingen",
	Port:          "Poort:",
	SlowTooltip:   "Toon de host in geel als de gemiddelde pingtijd hoger is. 0 - uit",
	Notify:        "Melden:",
	BeepDownBack:  "Pieptoon bij uitval of herstel",
	NotifyTooltip: "Offline betekent 3 mislukte pings op rij. De venstertitel toont hoeveel van zulke hosts offline zijn.",
	Share:         "Delen:",
	ShareOn:       "Pingtijd publiceren op u00.io",
	ShareTooltip:  "De pingtijden worden naar een openbare pagina gestuurd. Iedereen met de link kan ze zien.",
	ShareOpen:     "Openen",
	ShareCopy:     "Link kopiëren",
	LinkCopied:    "Link gekopieerd",
	MenuShareOpen: "Openen op u00.io",
	MenuShareCopy: "u00.io-link kopiëren",
	ShareStart:    "Publiceren op u00.io",
	ShareStop:     "Publiceren stoppen",

	DowntimeTitle: "Uitvaltijd",
	NoDowntime:    "Geen uitval in de afgelopen 24 uur",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Eén keer offline in de afgelopen 24 uur:"
		}
		return fmt.Sprintf("%d %s offline in de afgelopen 24 uur:", n, i18n.Plural("nl", n, "keer", "keer"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Storingen in de afgelopen 24 uur: %d", n) },
	Started:        "Begin",
	Ended:          "Einde",
	Duration:       "Duur",
	StillDown:      "nog offline",
	DateTimeLayout: "02-01 15:04:05",

	ConfigsTitle:       "Configuratie openen",
	Configs:            "Configuraties:",
	NewConfig:          "Nieuw...",
	SaveConfigAs:       "Opslaan als...",
	RemoveConfig:       "Verwijderen",
	HostsCount:         "Aantal hosts",
	Hosts:              "Hosts",
	NewConfigTitle:     "Nieuwe configuratie",
	SaveConfigAsTitle:  "Configuratie opslaan als",
	CopySuffix:         " kopie",
	RemoveConfigTitle:  "Configuratie verwijderen",
	RemoveConfigAsk:    "De geselecteerde configuratie verwijderen?",
	CannotRemoveOpened: "De momenteel geopende configuratie kan niet worden verwijderd.",

	ExtraColumns:   "Extra kolommen:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Sinds (laatste wijziging)",
	AlwaysOnTop:    "Altijd bovenaan",
	Language:       "Taal:",
	LanguageSystem: "Zoals in het systeem",
	Theme:          "Thema:",
	ThemeDark:      "Donker",
	ThemeLight:     "Licht",

	AboutTitle:   func(name string) string { return "Over " + name },
	Version:      "Versie",
	Author:       "Auteur:",
	License:      "Licentie:",
	VisitWebsite: "Website bezoeken",
	Close:        "Sluiten",

	ExportTitle: "Geschiedenis exporteren",
	CSVFiles:    "CSV-bestanden",
	SavedTo:     func(path string) string { return "Opgeslagen in " + path },
	ConfigLoadFailed: func(err string) string {
		return "De hostlijst kon niet worden gelezen, er is een nieuwe lege lijst geopend. Het bestand is ongewijzigd gelaten:\n\n" + err
	},

	Install:   "Installeren",
	Update:    "Bijwerken",
	Uninstall: "Verwijderen",
	InstallAsk: func(dir string) string {
		return "AltPing wordt gekopieerd naar " + dir + " en toegevoegd aan het menu Start, het bureaublad en de lijst met geïnstalleerde apps. Als het daar al is geïnstalleerd, wordt het bijgewerkt. Instellingen en hostlijsten blijven behouden.\n\nDaarna wordt AltPing vanaf daar opnieuw gestart."
	},
	InstallFailed: func(err string) string { return "AltPing kon niet worden geïnstalleerd:\n\n" + err },
	Installed:     "AltPing is geïnstalleerd",
	UninstallAsk: func(dir string) string {
		return "AltPing van deze computer verwijderen?\n\nInstellingen en hostlijsten blijven in " + dir
	},
	UninstallRunning: "AltPing is actief. Sluit het en probeer het opnieuw.",
	Uninstalled:      "AltPing is verwijderd.",
}

var sv = Strings{
	Error: "Fel",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d nere)", n) },
	ToolConfigs:     "Konfigurationer (O)",
	ToolAddHost:     "Lägg till värd (A)",
	ToolEditHost:    "Redigera värd (E)",
	ToolRemoveHosts: "Ta bort markerade värdar (Del)",
	ToolDetails:     "Detaljer (D)",
	ToolStart:       "Starta ping",
	ToolStop:        "Stoppa ping",
	ToolTray:        "Minimera till meddelandefältet (T)",
	ToolShareOpen:   "Öppna markerade värdar på u00.io (U)",
	ToolStyle:       "Stil (S)",
	KPIMedian:       "Medianping",
	KPIUptime:       "Tillgänglighet 24 h",
	TrayShow:        "Visa",
	TrayQuit:        "Avsluta",
	Settings:        "Inställningar",
	Help:            "Hjälp",
	About:           "Om",
	ModeUDP:         "UDP-läge",
	ModeICMP:        "ICMP-läge",
	StatusHosts:     "Värdar",
	StatusOK:        "OK",
	StatusSlow:      "Långsam",
	StatusDown:      "Nere",
	StatusStopped:   "Ping stoppad",

	Columns: ColumnStrings{
		Name: "Namn", IP: "IP", Time: "Tid ms", Trend: "Senaste 5 min", Loss: "Förlust %", Min: "Min ms", Avg: "Medel ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Sedan", Details: "Detaljer",
	},
	StateOK:         "OK",
	StateTimeout:    "TIDSGRÄNS",
	StateNoResolve:  "KAN INTE SLÅ UPP VÄRDNAMN",
	StatePortClosed: "PORTEN STÄNGD",
	NoHosts:         "Inga värdar än",
	PressAToAdd:     "Tryck på A för att lägga till en värd",
	HowItWorks:      "Så fungerar det",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "h", Day: "d"},

	MenuEdit:       "Redigera... (E)",
	MenuRemove:     "Ta bort (Del)",
	MenuDowntime:   "Avbrott...",
	MenuExport:     "Exportera historik...",
	MenuChange:     "Ändra markerade",
	BeepOn:         "Pip på",
	BeepOff:        "Pip av",
	CheckTCPPort:   "Kontrollera TCP-port",
	UsePing:        "Använd ping",
	PingEvery:      "Pingintervall",
	Timeout:        "Tidsgräns",
	SlowAbove:      "Långsam över",
	PortLabel:      "Port:",
	PingEveryMs:    "Pingintervall, ms:",
	TimeoutMs:      "Tidsgräns, ms:",
	SlowAboveMs:    "Långsam över, ms:",
	SlowAboveMsOff: "Långsam över, ms (0 - av):",

	RemoveHostsTitle: "Ta bort värdar",
	RemoveHost:       "Ta bort värden?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Ta bort %d %s?", n, i18n.Plural("sv", n, "värd", "värdar"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... och %d till", n) },

	Periods:     [3]string{"5m", "1h", "24h"},
	SelectAHost: "Välj en värd",
	PingHistory: "Pinghistorik: ",
	Ms:          "ms",

	Name:          "Namn:",
	Host:          "Värd:",
	TCPPort:       "TCP-port:",
	CheckInstead:  "Kontrollera i stället för ping",
	Port:          "Port:",
	SlowTooltip:   "Visa värden i gult när dess genomsnittliga pingtid är över värdet. 0 - av",
	Notify:        "Avisera:",
	BeepDownBack:  "Pip när nere eller tillbaka",
	NotifyTooltip: "Nere betyder 3 misslyckade ping i rad. Fönstertiteln visar hur många sådana värdar som är nere.",
	Share:         "Dela:",
	ShareOn:       "Publicera pingtid på u00.io",
	ShareTooltip:  "Pingtiderna skickas till en offentlig sida. Alla som har länken kan se dem.",
	ShareOpen:     "Öppna",
	ShareCopy:     "Kopiera länk",
	LinkCopied:    "Länken har kopierats",
	MenuShareOpen: "Öppna på u00.io",
	MenuShareCopy: "Kopiera u00.io-länk",
	ShareStart:    "Publicera på u00.io",
	ShareStop:     "Sluta publicera",

	DowntimeTitle: "Avbrott",
	NoDowntime:    "Inga avbrott de senaste 24 timmarna",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Nere en gång de senaste 24 timmarna:"
		}
		return fmt.Sprintf("Nere %d %s de senaste 24 timmarna:", n, i18n.Plural("sv", n, "gång", "gånger"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Avbrott de senaste 24 timmarna: %d", n) },
	Started:        "Början",
	Ended:          "Slut",
	Duration:       "Varaktighet",
	StillDown:      "fortfarande nere",
	DateTimeLayout: "01-02 15:04:05",

	ConfigsTitle:       "Öppna konfiguration",
	Configs:            "Konfigurationer:",
	NewConfig:          "Ny...",
	SaveConfigAs:       "Spara som...",
	RemoveConfig:       "Ta bort",
	HostsCount:         "Antal värdar",
	Hosts:              "Värdar",
	NewConfigTitle:     "Ny konfiguration",
	SaveConfigAsTitle:  "Spara konfiguration som",
	CopySuffix:         " kopia",
	RemoveConfigTitle:  "Ta bort konfiguration",
	RemoveConfigAsk:    "Ta bort den markerade konfigurationen?",
	CannotRemoveOpened: "Den konfiguration som är öppen kan inte tas bort.",

	ExtraColumns:   "Extra kolumner:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Sedan (senaste ändring)",
	AlwaysOnTop:    "Alltid överst",
	Language:       "Språk:",
	LanguageSystem: "Som i systemet",
	Theme:          "Tema:",
	ThemeDark:      "Mörkt",
	ThemeLight:     "Ljust",

	AboutTitle:   func(name string) string { return "Om " + name },
	Version:      "Version",
	Author:       "Författare:",
	License:      "Licens:",
	VisitWebsite: "Besök webbplatsen",
	Close:        "Stäng",

	ExportTitle: "Exportera historik",
	CSVFiles:    "CSV-filer",
	SavedTo:     func(path string) string { return "Sparad i " + path },
	ConfigLoadFailed: func(err string) string {
		return "Värdlistan kunde inte läsas, så en ny tom lista öppnades. Filen lämnades orörd:\n\n" + err
	},

	Install:   "Installera",
	Update:    "Uppdatera",
	Uninstall: "Avinstallera",
	InstallAsk: func(dir string) string {
		return "AltPing kopieras till " + dir + " och läggs till i Start-menyn, på skrivbordet och i listan över installerade appar. Om det redan är installerat där uppdateras det. Inställningar och värdlistor behålls.\n\nSedan startas AltPing om därifrån."
	},
	InstallFailed: func(err string) string { return "AltPing kunde inte installeras:\n\n" + err },
	Installed:     "AltPing är installerat",
	UninstallAsk: func(dir string) string {
		return "Ta bort AltPing från den här datorn?\n\nInställningar och värdlistor finns kvar i " + dir
	},
	UninstallRunning: "AltPing körs. Stäng programmet och försök igen.",
	Uninstalled:      "AltPing har tagits bort.",
}

var da = Strings{
	Error: "Fejl",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d nede)", n) },
	ToolConfigs:     "Konfigurationer (O)",
	ToolAddHost:     "Tilføj vært (A)",
	ToolEditHost:    "Rediger vært (E)",
	ToolRemoveHosts: "Fjern valgte værter (Del)",
	ToolDetails:     "Detaljer (D)",
	ToolStart:       "Start ping",
	ToolStop:        "Stop ping",
	ToolTray:        "Minimer til meddelelsesområdet (T)",
	ToolShareOpen:   "Åbn valgte værter på u00.io (U)",
	ToolStyle:       "Stil (S)",
	KPIMedian:       "Median-ping",
	KPIUptime:       "Tilgængelighed 24 t",
	TrayShow:        "Vis",
	TrayQuit:        "Afslut",
	Settings:        "Indstillinger",
	Help:            "Hjælp",
	About:           "Om",
	ModeUDP:         "UDP-tilstand",
	ModeICMP:        "ICMP-tilstand",
	StatusHosts:     "Værter",
	StatusOK:        "OK",
	StatusSlow:      "Langsom",
	StatusDown:      "Nede",
	StatusStopped:   "Ping stoppet",

	Columns: ColumnStrings{
		Name: "Navn", IP: "IP", Time: "Tid ms", Trend: "Seneste 5 min", Loss: "Tab %", Min: "Min ms", Avg: "Gns ms",
		Max: "Maks ms", Jitter: "Jitter ms", Since: "Siden", Details: "Detaljer",
	},
	StateOK:         "OK",
	StateTimeout:    "TIDSUDLØB",
	StateNoResolve:  "KAN IKKE SLÅ VÆRTSNAVN OP",
	StatePortClosed: "PORT LUKKET",
	NoHosts:         "Ingen værter endnu",
	PressAToAdd:     "Tryk på A for at tilføje en vært",
	HowItWorks:      "Sådan virker det",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "t", Day: "d"},

	MenuEdit:       "Rediger... (E)",
	MenuRemove:     "Fjern (Del)",
	MenuDowntime:   "Nedetid...",
	MenuExport:     "Eksportér historik...",
	MenuChange:     "Ændr valgte",
	BeepOn:         "Bip til",
	BeepOff:        "Bip fra",
	CheckTCPPort:   "Kontrollér TCP-port",
	UsePing:        "Brug ping",
	PingEvery:      "Ping hver",
	Timeout:        "Tidsgrænse",
	SlowAbove:      "Langsom over",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping hver, ms:",
	TimeoutMs:      "Tidsgrænse, ms:",
	SlowAboveMs:    "Langsom over, ms:",
	SlowAboveMsOff: "Langsom over, ms (0 - fra):",

	RemoveHostsTitle: "Fjern værter",
	RemoveHost:       "Fjern værten?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Fjern %d %s?", n, i18n.Plural("da", n, "vært", "værter"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... og %d til", n) },

	Periods:     [3]string{"5m", "1t", "24t"},
	SelectAHost: "Vælg en vært",
	PingHistory: "Pinghistorik: ",
	Ms:          "ms",

	Name:          "Navn:",
	Host:          "Vært:",
	TCPPort:       "TCP-port:",
	CheckInstead:  "Kontrollér i stedet for ping",
	Port:          "Port:",
	SlowTooltip:   "Vis værten med gult, når dens gennemsnitlige pingtid er over denne værdi. 0 - fra",
	Notify:        "Giv besked:",
	BeepDownBack:  "Bip når nede eller tilbage",
	NotifyTooltip: "Nede betyder 3 mislykkede ping i træk. Vinduets titel viser, hvor mange sådanne værter der er nede.",
	Share:         "Del:",
	ShareOn:       "Udgiv pingtid på u00.io",
	ShareTooltip:  "Pingtiderne sendes til en offentlig side. Alle med linket kan se dem.",
	ShareOpen:     "Åbn",
	ShareCopy:     "Kopiér link",
	LinkCopied:    "Linket er kopieret",
	MenuShareOpen: "Åbn på u00.io",
	MenuShareCopy: "Kopiér u00.io-link",
	ShareStart:    "Udgiv på u00.io",
	ShareStop:     "Stop udgivelse",

	DowntimeTitle: "Nedetid",
	NoDowntime:    "Ingen nedetid de seneste 24 timer",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Nede én gang de seneste 24 timer:"
		}
		return fmt.Sprintf("Nede %d %s de seneste 24 timer:", n, i18n.Plural("da", n, "gang", "gange"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Nedbrud de seneste 24 timer: %d", n) },
	Started:        "Start",
	Ended:          "Slut",
	Duration:       "Varighed",
	StillDown:      "stadig nede",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Åbn konfiguration",
	Configs:            "Konfigurationer:",
	NewConfig:          "Ny...",
	SaveConfigAs:       "Gem som...",
	RemoveConfig:       "Fjern",
	HostsCount:         "Antal værter",
	Hosts:              "Værter",
	NewConfigTitle:     "Ny konfiguration",
	SaveConfigAsTitle:  "Gem konfiguration som",
	CopySuffix:         " kopi",
	RemoveConfigTitle:  "Fjern konfiguration",
	RemoveConfigAsk:    "Fjern den valgte konfiguration?",
	CannotRemoveOpened: "Den konfiguration, der er åben nu, kan ikke fjernes.",

	ExtraColumns:   "Ekstra kolonner:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Siden (seneste ændring)",
	AlwaysOnTop:    "Altid øverst",
	Language:       "Sprog:",
	LanguageSystem: "Som i systemet",
	Theme:          "Tema:",
	ThemeDark:      "Mørkt",
	ThemeLight:     "Lyst",

	AboutTitle:   func(name string) string { return "Om " + name },
	Version:      "Version",
	Author:       "Forfatter:",
	License:      "Licens:",
	VisitWebsite: "Besøg webstedet",
	Close:        "Luk",

	ExportTitle: "Eksportér historik",
	CSVFiles:    "CSV-filer",
	SavedTo:     func(path string) string { return "Gemt i " + path },
	ConfigLoadFailed: func(err string) string {
		return "Værtslisten kunne ikke læses, så en ny tom liste blev åbnet. Filen er ikke ændret:\n\n" + err
	},

	Install:   "Installer",
	Update:    "Opdater",
	Uninstall: "Afinstaller",
	InstallAsk: func(dir string) string {
		return "AltPing kopieres til " + dir + " og føjes til menuen Start, skrivebordet og listen over installerede apps. Hvis det allerede er installeret der, opdateres det. Indstillinger og værtslister bevares.\n\nDerefter genstartes AltPing derfra."
	},
	InstallFailed: func(err string) string { return "AltPing kunne ikke installeres:\n\n" + err },
	Installed:     "AltPing er installeret",
	UninstallAsk: func(dir string) string {
		return "Fjern AltPing fra denne computer?\n\nIndstillinger og værtslister bliver i " + dir
	},
	UninstallRunning: "AltPing kører. Luk programmet, og prøv igen.",
	Uninstalled:      "AltPing er fjernet.",
}

var nb = Strings{
	Error: "Feil",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d nede)", n) },
	ToolConfigs:     "Konfigurasjoner (O)",
	ToolAddHost:     "Legg til vert (A)",
	ToolEditHost:    "Rediger vert (E)",
	ToolRemoveHosts: "Fjern valgte verter (Del)",
	ToolDetails:     "Detaljer (D)",
	ToolStart:       "Start ping",
	ToolStop:        "Stopp ping",
	ToolTray:        "Minimer til systemstatusfeltet (T)",
	ToolShareOpen:   "Åpne valgte verter på u00.io (U)",
	ToolStyle:       "Stil (S)",
	KPIMedian:       "Median-ping",
	KPIUptime:       "Tilgjengelighet 24 t",
	TrayShow:        "Vis",
	TrayQuit:        "Avslutt",
	Settings:        "Innstillinger",
	Help:            "Hjelp",
	About:           "Om",
	ModeUDP:         "UDP-modus",
	ModeICMP:        "ICMP-modus",
	StatusHosts:     "Verter",
	StatusOK:        "OK",
	StatusSlow:      "Treg",
	StatusDown:      "Nede",
	StatusStopped:   "Ping stoppet",

	Columns: ColumnStrings{
		Name: "Navn", IP: "IP", Time: "Tid ms", Trend: "Siste 5 min", Loss: "Tap %", Min: "Min ms", Avg: "Snitt ms",
		Max: "Maks ms", Jitter: "Jitter ms", Since: "Siden", Details: "Detaljer",
	},
	StateOK:         "OK",
	StateTimeout:    "TIDSAVBRUDD",
	StateNoResolve:  "KAN IKKE SLÅ OPP VERTSNAVN",
	StatePortClosed: "PORT LUKKET",
	NoHosts:         "Ingen verter ennå",
	PressAToAdd:     "Trykk A for å legge til en vert",
	HowItWorks:      "Slik fungerer det",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "t", Day: "d"},

	MenuEdit:       "Rediger... (E)",
	MenuRemove:     "Fjern (Del)",
	MenuDowntime:   "Nedetid...",
	MenuExport:     "Eksporter historikk...",
	MenuChange:     "Endre valgte",
	BeepOn:         "Pip på",
	BeepOff:        "Pip av",
	CheckTCPPort:   "Sjekk TCP-port",
	UsePing:        "Bruk ping",
	PingEvery:      "Ping hvert",
	Timeout:        "Tidsavbrudd",
	SlowAbove:      "Treg over",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping hvert, ms:",
	TimeoutMs:      "Tidsavbrudd, ms:",
	SlowAboveMs:    "Treg over, ms:",
	SlowAboveMsOff: "Treg over, ms (0 - av):",

	RemoveHostsTitle: "Fjern verter",
	RemoveHost:       "Fjerne verten?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Fjerne %d %s?", n, i18n.Plural("nb", n, "vert", "verter"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... og %d til", n) },

	Periods:     [3]string{"5m", "1t", "24t"},
	SelectAHost: "Velg en vert",
	PingHistory: "Pinghistorikk: ",
	Ms:          "ms",

	Name:          "Navn:",
	Host:          "Vert:",
	TCPPort:       "TCP-port:",
	CheckInstead:  "Sjekk i stedet for ping",
	Port:          "Port:",
	SlowTooltip:   "Vis verten i gult når den gjennomsnittlige pingtiden er over denne verdien. 0 - av",
	Notify:        "Varsle:",
	BeepDownBack:  "Pip når nede eller tilbake",
	NotifyTooltip: "Nede betyr 3 mislykkede ping på rad. Vindustittelen viser hvor mange slike verter som er nede.",
	Share:         "Del:",
	ShareOn:       "Publiser pingtid på u00.io",
	ShareTooltip:  "Pingtidene sendes til en offentlig side. Alle med lenken kan se dem.",
	ShareOpen:     "Åpne",
	ShareCopy:     "Kopier lenke",
	LinkCopied:    "Lenken er kopiert",
	MenuShareOpen: "Åpne på u00.io",
	MenuShareCopy: "Kopier u00.io-lenke",
	ShareStart:    "Publiser på u00.io",
	ShareStop:     "Stopp publisering",

	DowntimeTitle: "Nedetid",
	NoDowntime:    "Ingen nedetid de siste 24 timene",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Nede én gang de siste 24 timene:"
		}
		return fmt.Sprintf("Nede %d %s de siste 24 timene:", n, i18n.Plural("nb", n, "gang", "ganger"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Avbrudd de siste 24 timene: %d", n) },
	Started:        "Start",
	Ended:          "Slutt",
	Duration:       "Varighet",
	StillDown:      "fortsatt nede",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Åpne konfigurasjon",
	Configs:            "Konfigurasjoner:",
	NewConfig:          "Ny...",
	SaveConfigAs:       "Lagre som...",
	RemoveConfig:       "Fjern",
	HostsCount:         "Antall verter",
	Hosts:              "Verter",
	NewConfigTitle:     "Ny konfigurasjon",
	SaveConfigAsTitle:  "Lagre konfigurasjon som",
	CopySuffix:         " kopi",
	RemoveConfigTitle:  "Fjern konfigurasjon",
	RemoveConfigAsk:    "Fjerne den valgte konfigurasjonen?",
	CannotRemoveOpened: "Konfigurasjonen som er åpen nå, kan ikke fjernes.",

	ExtraColumns:   "Ekstra kolonner:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Siden (siste endring)",
	AlwaysOnTop:    "Alltid øverst",
	Language:       "Språk:",
	LanguageSystem: "Som i systemet",
	Theme:          "Tema:",
	ThemeDark:      "Mørkt",
	ThemeLight:     "Lyst",

	AboutTitle:   func(name string) string { return "Om " + name },
	Version:      "Versjon",
	Author:       "Forfatter:",
	License:      "Lisens:",
	VisitWebsite: "Besøk nettstedet",
	Close:        "Lukk",

	ExportTitle: "Eksporter historikk",
	CSVFiles:    "CSV-filer",
	SavedTo:     func(path string) string { return "Lagret i " + path },
	ConfigLoadFailed: func(err string) string {
		return "Vertslisten kunne ikke leses, så en ny tom liste ble åpnet. Filen ble ikke endret:\n\n" + err
	},

	Install:   "Installer",
	Update:    "Oppdater",
	Uninstall: "Avinstaller",
	InstallAsk: func(dir string) string {
		return "AltPing kopieres til " + dir + " og legges til på Start-menyen, skrivebordet og i listen over installerte apper. Hvis det allerede er installert der, oppdateres det. Innstillinger og vertslister beholdes.\n\nDeretter startes AltPing på nytt derfra."
	},
	InstallFailed: func(err string) string { return "AltPing kunne ikke installeres:\n\n" + err },
	Installed:     "AltPing er installert",
	UninstallAsk: func(dir string) string {
		return "Fjerne AltPing fra denne datamaskinen?\n\nInnstillinger og vertslister blir liggende i " + dir
	},
	UninstallRunning: "AltPing kjører. Lukk programmet og prøv igjen.",
	Uninstalled:      "AltPing er fjernet.",
}

var fi = Strings{
	Error: "Virhe",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d ei vastaa)", n) },
	ToolConfigs:     "Määritykset (O)",
	ToolAddHost:     "Lisää isäntä (A)",
	ToolEditHost:    "Muokkaa isäntää (E)",
	ToolRemoveHosts: "Poista valitut isännät (Del)",
	ToolDetails:     "Tiedot (D)",
	ToolStart:       "Aloita ping",
	ToolStop:        "Lopeta ping",
	ToolTray:        "Pienennä ilmaisinalueelle (T)",
	ToolShareOpen:   "Avaa valitut isännät u00.io:ssa (U)",
	ToolStyle:       "Tyyli (S)",
	KPIMedian:       "Pingin mediaani",
	KPIUptime:       "Saatavuus 24 h",
	TrayShow:        "Näytä",
	TrayQuit:        "Lopeta",
	Settings:        "Asetukset",
	Help:            "Ohje",
	About:           "Tietoja",
	ModeUDP:         "UDP-tila",
	ModeICMP:        "ICMP-tila",
	StatusHosts:     "Isännät",
	StatusOK:        "OK",
	StatusSlow:      "Hidas",
	StatusDown:      "Ei vastaa",
	StatusStopped:   "Ping pysäytetty",

	Columns: ColumnStrings{
		Name: "Nimi", IP: "IP", Time: "Aika ms", Trend: "Viim. 5 min", Loss: "Hävikki %", Min: "Min ms", Avg: "Ka ms",
		Max: "Maks ms", Jitter: "Värinä ms", Since: "Alkaen", Details: "Tiedot",
	},
	StateOK:         "OK",
	StateTimeout:    "AIKAKATKAISU",
	StateNoResolve:  "ISÄNTÄNIMEÄ EI VOI SELVITTÄÄ",
	StatePortClosed: "PORTTI SULJETTU",
	NoHosts:         "Ei vielä isäntiä",
	PressAToAdd:     "Lisää isäntä painamalla A",
	HowItWorks:      "Miten tämä toimii",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "pv"},

	MenuEdit:       "Muokkaa... (E)",
	MenuRemove:     "Poista (Del)",
	MenuDowntime:   "Katkokset...",
	MenuExport:     "Vie historia...",
	MenuChange:     "Muuta valittuja",
	BeepOn:         "Äänimerkki päälle",
	BeepOff:        "Äänimerkki pois",
	CheckTCPPort:   "Tarkista TCP-portti",
	UsePing:        "Käytä pingiä",
	PingEvery:      "Ping-väli",
	Timeout:        "Aikakatkaisu",
	SlowAbove:      "Hidas yli",
	PortLabel:      "Portti:",
	PingEveryMs:    "Ping-väli, ms:",
	TimeoutMs:      "Aikakatkaisu, ms:",
	SlowAboveMs:    "Hidas yli, ms:",
	SlowAboveMsOff: "Hidas yli, ms (0 - pois):",

	RemoveHostsTitle: "Poista isännät",
	RemoveHost:       "Poistetaanko isäntä?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Poistetaanko %d %s?", n, i18n.Plural("fi", n, "isäntä", "isäntää"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... ja %d muuta", n) },

	Periods:     [3]string{"5min", "1h", "24h"},
	SelectAHost: "Valitse isäntä",
	PingHistory: "Ping-historia: ",
	Ms:          "ms",

	Name:          "Nimi:",
	Host:          "Isäntä:",
	TCPPort:       "TCP-portti:",
	CheckInstead:  "Tarkista pingin sijaan",
	Port:          "Portti:",
	SlowTooltip:   "Näytä isäntä keltaisena, kun sen keskimääräinen ping-aika ylittää tämän arvon. 0 - pois",
	Notify:        "Ilmoita:",
	BeepDownBack:  "Äänimerkki, kun yhteys katkeaa tai palaa",
	NotifyTooltip: "Ei vastaa tarkoittaa 3 epäonnistunutta pingiä peräkkäin. Ikkunan otsikko näyttää, montako tällaista isäntää ei vastaa.",
	Share:         "Jaa:",
	ShareOn:       "Julkaise ping-aika u00.io:ssa",
	ShareTooltip:  "Ping-ajat lähetetään julkiselle sivulle. Kuka tahansa, jolla on linkki, voi nähdä ne.",
	ShareOpen:     "Avaa",
	ShareCopy:     "Kopioi linkki",
	LinkCopied:    "Linkki kopioitu",
	MenuShareOpen: "Avaa u00.io:ssa",
	MenuShareCopy: "Kopioi u00.io-linkki",
	ShareStart:    "Julkaise u00.io:ssa",
	ShareStop:     "Lopeta julkaisu",

	DowntimeTitle: "Katkokset",
	NoDowntime:    "Ei katkoksia viimeisten 24 tunnin aikana",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Ei vastannut kerran viimeisten 24 tunnin aikana:"
		}
		return fmt.Sprintf("Ei vastannut %d %s viimeisten 24 tunnin aikana:", n, i18n.Plural("fi", n, "kerran", "kertaa"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Katkoksia viimeisten 24 tunnin aikana: %d", n) },
	Started:        "Alkoi",
	Ended:          "Päättyi",
	Duration:       "Kesto",
	StillDown:      "katkos jatkuu",
	DateTimeLayout: "2.1. 15:04:05",

	ConfigsTitle:       "Avaa määritys",
	Configs:            "Määritykset:",
	NewConfig:          "Uusi...",
	SaveConfigAs:       "Tallenna nimellä...",
	RemoveConfig:       "Poista",
	HostsCount:         "Isäntien määrä",
	Hosts:              "Isännät",
	NewConfigTitle:     "Uusi määritys",
	SaveConfigAsTitle:  "Tallenna määritys nimellä",
	CopySuffix:         " kopio",
	RemoveConfigTitle:  "Poista määritys",
	RemoveConfigAsk:    "Poistetaanko valittu määritys?",
	CannotRemoveOpened: "Avoinna olevaa määritystä ei voi poistaa.",

	ExtraColumns:   "Lisäsarakkeet:",
	ShowMin:        "Min",
	ShowJitter:     "Värinä",
	ShowSince:      "Alkaen (viimeisin muutos)",
	AlwaysOnTop:    "Aina päällimmäisenä",
	Language:       "Kieli:",
	LanguageSystem: "Järjestelmän mukaan",
	Theme:          "Teema:",
	ThemeDark:      "Tumma",
	ThemeLight:     "Vaalea",

	AboutTitle:   func(name string) string { return "Tietoja: " + name },
	Version:      "Versio",
	Author:       "Tekijä:",
	License:      "Lisenssi:",
	VisitWebsite: "Siirry verkkosivustoon",
	Close:        "Sulje",

	ExportTitle: "Vie historia",
	CSVFiles:    "CSV-tiedostot",
	SavedTo:     func(path string) string { return "Tallennettu: " + path },
	ConfigLoadFailed: func(err string) string {
		return "Isäntäluetteloa ei voitu lukea, joten avattiin uusi tyhjä luettelo. Tiedostoa ei muutettu:\n\n" + err
	},

	Install:   "Asenna",
	Update:    "Päivitä",
	Uninstall: "Poista asennus",
	InstallAsk: func(dir string) string {
		return "AltPing kopioidaan kansioon " + dir + " ja lisätään Käynnistä-valikkoon, työpöydälle ja asennettujen sovellusten luetteloon. Jos se on jo asennettu sinne, se päivitetään. Asetukset ja isäntäluettelot säilyvät ennallaan.\n\nSen jälkeen AltPing käynnistyy uudelleen sieltä."
	},
	InstallFailed: func(err string) string { return "AltPingin asennus epäonnistui:\n\n" + err },
	Installed:     "AltPing on asennettu",
	UninstallAsk: func(dir string) string {
		return "Poistetaanko AltPing tästä tietokoneesta?\n\nAsetukset ja isäntäluettelot jäävät kansioon " + dir
	},
	UninstallRunning: "AltPing on käynnissä. Sulje se ja yritä uudelleen.",
	Uninstalled:      "AltPing on poistettu.",
}

var cs = Strings{
	Error: "Chyba",

	DownCount: func(n int) string {
		return fmt.Sprintf("(%d %s)", n, i18n.Plural("cs", n, "nedostupný", "nedostupné", "nedostupných"))
	},
	ToolConfigs:     "Konfigurace (O)",
	ToolAddHost:     "Přidat hostitele (A)",
	ToolEditHost:    "Upravit hostitele (E)",
	ToolRemoveHosts: "Odebrat vybrané hostitele (Del)",
	ToolDetails:     "Podrobnosti (D)",
	ToolStart:       "Spustit pingování",
	ToolStop:        "Zastavit pingování",
	ToolTray:        "Minimalizovat do oznamovací oblasti (T)",
	ToolShareOpen:   "Otevřít vybrané hostitele na u00.io (U)",
	ToolStyle:       "Styl (S)",
	KPIMedian:       "Medián pingu",
	KPIUptime:       "Dostupnost 24 h",
	TrayShow:        "Zobrazit",
	TrayQuit:        "Ukončit",
	Settings:        "Nastavení",
	Help:            "Nápověda",
	About:           "O aplikaci",
	ModeUDP:         "Režim UDP",
	ModeICMP:        "Režim ICMP",
	StatusHosts:     "Hostitelé",
	StatusOK:        "OK",
	StatusSlow:      "Pomalé",
	StatusDown:      "Nedostupné",
	StatusStopped:   "Pingování zastaveno",

	Columns: ColumnStrings{
		Name: "Název", IP: "IP", Time: "Čas ms", Trend: "Posl. 5 min", Loss: "Ztráty %", Min: "Min ms", Avg: "Prům. ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Od", Details: "Stav",
	},
	StateOK:         "OK",
	StateTimeout:    "ČAS VYPRŠEL",
	StateNoResolve:  "NELZE PŘELOŽIT NÁZEV HOSTITELE",
	StatePortClosed: "PORT ZAVŘEN",
	NoHosts:         "Zatím žádní hostitelé",
	PressAToAdd:     "Stisknutím A přidáte hostitele",
	HowItWorks:      "Jak to funguje",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "d"},

	MenuEdit:       "Upravit... (E)",
	MenuRemove:     "Odebrat (Del)",
	MenuDowntime:   "Výpadky...",
	MenuExport:     "Exportovat historii...",
	MenuChange:     "Změnit vybrané",
	BeepOn:         "Zapnout pípání",
	BeepOff:        "Vypnout pípání",
	CheckTCPPort:   "Kontrolovat port TCP",
	UsePing:        "Používat ping",
	PingEvery:      "Ping každých",
	Timeout:        "Časový limit",
	SlowAbove:      "Pomalé nad",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping každých, ms:",
	TimeoutMs:      "Časový limit, ms:",
	SlowAboveMs:    "Pomalé nad, ms:",
	SlowAboveMsOff: "Pomalé nad, ms (0 - vyp.):",

	RemoveHostsTitle: "Odebrat hostitele",
	RemoveHost:       "Odebrat hostitele?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Odebrat %d %s?", n, i18n.Plural("cs", n, "hostitele", "hostitele", "hostitelů"))
	},
	AndMore: func(n int) string {
		return fmt.Sprintf("... a %d %s", n, i18n.Plural("cs", n, "další", "další", "dalších"))
	},

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Vyberte hostitele",
	PingHistory: "Historie pingu: ",
	Ms:          "ms",

	Name:          "Název:",
	Host:          "Hostitel:",
	TCPPort:       "Port TCP:",
	CheckInstead:  "Kontrolovat místo pingu",
	Port:          "Port:",
	SlowTooltip:   "Hostitel se zobrazí žlutě, když jeho průměrná doba pingu překročí tuto hodnotu. 0 - vyp.",
	Notify:        "Upozornění:",
	BeepDownBack:  "Pípnout při výpadku a obnovení",
	NotifyTooltip: "Výpadek znamená 3 neúspěšné pingy za sebou. Titulek okna ukazuje, kolik takových hostitelů je nedostupných.",
	Share:         "Sdílení:",
	ShareOn:       "Zveřejňovat dobu pingu na u00.io",
	ShareTooltip:  "Doby pingu se odesílají na veřejnou stránku. Uvidí je každý, kdo má odkaz.",
	ShareOpen:     "Otevřít",
	ShareCopy:     "Kopírovat odkaz",
	LinkCopied:    "Odkaz zkopírován",
	MenuShareOpen: "Otevřít na u00.io",
	MenuShareCopy: "Kopírovat odkaz u00.io",
	ShareStart:    "Zveřejnit na u00.io",
	ShareStop:     "Ukončit zveřejňování",

	DowntimeTitle: "Výpadky",
	NoDowntime:    "Žádné výpadky za posledních 24 hodin",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Nedostupný jednou za posledních 24 hodin:"
		}
		return fmt.Sprintf("Nedostupný %dkrát za posledních 24 hodin:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Výpadky za posledních 24 hodin: %d", n) },
	Started:        "Začátek",
	Ended:          "Konec",
	Duration:       "Trvání",
	StillDown:      "stále nedostupný",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Otevřít konfiguraci",
	Configs:            "Konfigurace:",
	NewConfig:          "Nová...",
	SaveConfigAs:       "Uložit jako...",
	RemoveConfig:       "Odebrat",
	HostsCount:         "Počet hostitelů",
	Hosts:              "Hostitelé",
	NewConfigTitle:     "Nová konfigurace",
	SaveConfigAsTitle:  "Uložit konfiguraci jako",
	CopySuffix:         " kopie",
	RemoveConfigTitle:  "Odebrat konfiguraci",
	RemoveConfigAsk:    "Odebrat vybranou konfiguraci?",
	CannotRemoveOpened: "Právě otevřenou konfiguraci nelze odebrat.",

	ExtraColumns:   "Další sloupce:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Od (poslední změna)",
	AlwaysOnTop:    "Vždy navrchu",
	Language:       "Jazyk:",
	LanguageSystem: "Podle systému",
	Theme:          "Motiv:",
	ThemeDark:      "Tmavý",
	ThemeLight:     "Světlý",

	AboutTitle:   func(name string) string { return "O aplikaci " + name },
	Version:      "Verze",
	Author:       "Autor:",
	License:      "Licence:",
	VisitWebsite: "Navštívit web",
	Close:        "Zavřít",

	ExportTitle: "Export historie",
	CSVFiles:    "Soubory CSV",
	SavedTo:     func(path string) string { return "Uloženo do " + path },
	ConfigLoadFailed: func(err string) string {
		return "Seznam hostitelů nelze načíst, byl otevřen nový prázdný. Soubor zůstal beze změny:\n\n" + err
	},

	Install:   "Nainstalovat",
	Update:    "Aktualizovat",
	Uninstall: "Odinstalovat",
	InstallAsk: func(dir string) string {
		return "AltPing bude zkopírován do " + dir + " a přidán do nabídky Start, na plochu a do seznamu nainstalovaných aplikací. Pokud tam už nainstalován je, bude aktualizován. Nastavení a seznamy hostitelů zůstanou beze změny.\n\nAltPing se poté znovu spustí odtamtud."
	},
	InstallFailed: func(err string) string { return "AltPing nelze nainstalovat:\n\n" + err },
	Installed:     "AltPing je nainstalován",
	UninstallAsk: func(dir string) string {
		return "Odebrat AltPing z tohoto počítače?\n\nNastavení a seznamy hostitelů zůstanou v " + dir
	},
	UninstallRunning: "AltPing je spuštěn. Zavřete ho a zkuste to znovu.",
	Uninstalled:      "AltPing byl odebrán.",
}

var sk = Strings{
	Error: "Chyba",

	DownCount: func(n int) string {
		return fmt.Sprintf("(%d %s)", n, i18n.Plural("sk", n, "nedostupný", "nedostupné", "nedostupných"))
	},
	ToolConfigs:     "Konfigurácie (O)",
	ToolAddHost:     "Pridať hostiteľa (A)",
	ToolEditHost:    "Upraviť hostiteľa (E)",
	ToolRemoveHosts: "Odstrániť vybraných hostiteľov (Del)",
	ToolDetails:     "Podrobnosti (D)",
	ToolStart:       "Spustiť pingovanie",
	ToolStop:        "Zastaviť pingovanie",
	ToolTray:        "Minimalizovať do oblasti oznámení (T)",
	ToolShareOpen:   "Otvoriť vybraných hostiteľov na u00.io (U)",
	ToolStyle:       "Štýl (S)",
	KPIMedian:       "Medián pingu",
	KPIUptime:       "Dostupnosť 24 h",
	TrayShow:        "Zobraziť",
	TrayQuit:        "Ukončiť",
	Settings:        "Nastavenia",
	Help:            "Pomocník",
	About:           "O aplikácii",
	ModeUDP:         "Režim UDP",
	ModeICMP:        "Režim ICMP",
	StatusHosts:     "Hostitelia",
	StatusOK:        "OK",
	StatusSlow:      "Pomalé",
	StatusDown:      "Nedostupné",
	StatusStopped:   "Pingovanie zastavené",

	Columns: ColumnStrings{
		Name: "Názov", IP: "IP", Time: "Čas ms", Trend: "Posl. 5 min", Loss: "Straty %", Min: "Min ms", Avg: "Priem. ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Od", Details: "Stav",
	},
	StateOK:         "OK",
	StateTimeout:    "ČAS VYPRŠAL",
	StateNoResolve:  "NEDÁ SA PRELOŽIŤ NÁZOV HOSTITEĽA",
	StatePortClosed: "PORT ZATVORENÝ",
	NoHosts:         "Zatiaľ žiadni hostitelia",
	PressAToAdd:     "Stlačením A pridáte hostiteľa",
	HowItWorks:      "Ako to funguje",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "d"},

	MenuEdit:       "Upraviť... (E)",
	MenuRemove:     "Odstrániť (Del)",
	MenuDowntime:   "Výpadky...",
	MenuExport:     "Exportovať históriu...",
	MenuChange:     "Zmeniť vybrané",
	BeepOn:         "Zapnúť pípanie",
	BeepOff:        "Vypnúť pípanie",
	CheckTCPPort:   "Kontrolovať port TCP",
	UsePing:        "Používať ping",
	PingEvery:      "Ping každých",
	Timeout:        "Časový limit",
	SlowAbove:      "Pomalé nad",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping každých, ms:",
	TimeoutMs:      "Časový limit, ms:",
	SlowAboveMs:    "Pomalé nad, ms:",
	SlowAboveMsOff: "Pomalé nad, ms (0 - vyp.):",

	RemoveHostsTitle: "Odstrániť hostiteľov",
	RemoveHost:       "Odstrániť hostiteľa?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Odstrániť %d %s?", n, i18n.Plural("sk", n, "hostiteľa", "hostiteľov", "hostiteľov"))
	},
	AndMore: func(n int) string {
		return fmt.Sprintf("... a %d %s", n, i18n.Plural("sk", n, "ďalší", "ďalšie", "ďalších"))
	},

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Vyberte hostiteľa",
	PingHistory: "História pingu: ",
	Ms:          "ms",

	Name:          "Názov:",
	Host:          "Hostiteľ:",
	TCPPort:       "Port TCP:",
	CheckInstead:  "Kontrolovať namiesto pingu",
	Port:          "Port:",
	SlowTooltip:   "Hostiteľ sa zobrazí nažlto, keď jeho priemerný čas pingu prekročí túto hodnotu. 0 - vyp.",
	Notify:        "Upozornenie:",
	BeepDownBack:  "Pípnuť pri výpadku a obnovení",
	NotifyTooltip: "Výpadok znamená 3 neúspešné pingy za sebou. Záhlavie okna ukazuje, koľko takých hostiteľov je nedostupných.",
	Share:         "Zdieľanie:",
	ShareOn:       "Zverejňovať čas pingu na u00.io",
	ShareTooltip:  "Časy pingu sa odosielajú na verejnú stránku. Uvidí ich každý, kto má odkaz.",
	ShareOpen:     "Otvoriť",
	ShareCopy:     "Kopírovať odkaz",
	LinkCopied:    "Odkaz skopírovaný",
	MenuShareOpen: "Otvoriť na u00.io",
	MenuShareCopy: "Kopírovať odkaz u00.io",
	ShareStart:    "Zverejniť na u00.io",
	ShareStop:     "Ukončiť zverejňovanie",

	DowntimeTitle: "Výpadky",
	NoDowntime:    "Žiadne výpadky za posledných 24 hodín",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Nedostupný raz za posledných 24 hodín:"
		}
		return fmt.Sprintf("Nedostupný %d-krát za posledných 24 hodín:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Výpadky za posledných 24 hodín: %d", n) },
	Started:        "Začiatok",
	Ended:          "Koniec",
	Duration:       "Trvanie",
	StillDown:      "stále nedostupný",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Otvoriť konfiguráciu",
	Configs:            "Konfigurácie:",
	NewConfig:          "Nová...",
	SaveConfigAs:       "Uložiť ako...",
	RemoveConfig:       "Odstrániť",
	HostsCount:         "Počet hostiteľov",
	Hosts:              "Hostitelia",
	NewConfigTitle:     "Nová konfigurácia",
	SaveConfigAsTitle:  "Uložiť konfiguráciu ako",
	CopySuffix:         " kópia",
	RemoveConfigTitle:  "Odstrániť konfiguráciu",
	RemoveConfigAsk:    "Odstrániť vybranú konfiguráciu?",
	CannotRemoveOpened: "Práve otvorenú konfiguráciu nemožno odstrániť.",

	ExtraColumns:   "Ďalšie stĺpce:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Od (posledná zmena)",
	AlwaysOnTop:    "Vždy navrchu",
	Language:       "Jazyk:",
	LanguageSystem: "Podľa systému",
	Theme:          "Motív:",
	ThemeDark:      "Tmavý",
	ThemeLight:     "Svetlý",

	AboutTitle:   func(name string) string { return "O aplikácii " + name },
	Version:      "Verzia",
	Author:       "Autor:",
	License:      "Licencia:",
	VisitWebsite: "Navštíviť web",
	Close:        "Zavrieť",

	ExportTitle: "Export histórie",
	CSVFiles:    "Súbory CSV",
	SavedTo:     func(path string) string { return "Uložené do " + path },
	ConfigLoadFailed: func(err string) string {
		return "Zoznam hostiteľov sa nepodarilo načítať, otvoril sa nový prázdny. Súbor zostal bez zmeny:\n\n" + err
	},

	Install:   "Nainštalovať",
	Update:    "Aktualizovať",
	Uninstall: "Odinštalovať",
	InstallAsk: func(dir string) string {
		return "AltPing sa skopíruje do " + dir + " a pridá do ponuky Štart, na pracovnú plochu a do zoznamu nainštalovaných aplikácií. Ak je tam už nainštalovaný, aktualizuje sa. Nastavenia a zoznamy hostiteľov zostanú bez zmeny.\n\nAltPing sa potom znova spustí odtiaľ."
	},
	InstallFailed: func(err string) string { return "AltPing sa nepodarilo nainštalovať:\n\n" + err },
	Installed:     "AltPing je nainštalovaný",
	UninstallAsk: func(dir string) string {
		return "Odstrániť AltPing z tohto počítača?\n\nNastavenia a zoznamy hostiteľov zostanú v " + dir
	},
	UninstallRunning: "AltPing je spustený. Zatvorte ho a skúste to znova.",
	Uninstalled:      "AltPing bol odstránený.",
}

var sl = Strings{
	Error: "Napaka",

	DownCount: func(n int) string {
		return fmt.Sprintf("(%d %s)", n, i18n.Plural("sl", n, "nedosegljiv", "nedosegljiva", "nedosegljivi", "nedosegljivih"))
	},
	ToolConfigs:     "Konfiguracije (O)",
	ToolAddHost:     "Dodaj gostitelja (A)",
	ToolEditHost:    "Uredi gostitelja (E)",
	ToolRemoveHosts: "Odstrani izbrane gostitelje (Del)",
	ToolDetails:     "Podrobnosti (D)",
	ToolStart:       "Začni pinganje",
	ToolStop:        "Ustavi pinganje",
	ToolTray:        "Pomanjšaj v območje za obvestila (T)",
	ToolShareOpen:   "Odpri izbrane gostitelje na u00.io (U)",
	ToolStyle:       "Slog (S)",
	KPIMedian:       "Mediana pinga",
	KPIUptime:       "Razpoložljivost 24 h",
	TrayShow:        "Pokaži",
	TrayQuit:        "Izhod",
	Settings:        "Nastavitve",
	Help:            "Pomoč",
	About:           "O programu",
	ModeUDP:         "Način UDP",
	ModeICMP:        "Način ICMP",
	StatusHosts:     "Gostitelji",
	StatusOK:        "OK",
	StatusSlow:      "Počasni",
	StatusDown:      "Nedosegljivi",
	StatusStopped:   "Pinganje ustavljeno",

	Columns: ColumnStrings{
		Name: "Ime", IP: "IP", Time: "Čas ms", Trend: "Zadnjih 5 min", Loss: "Izguba %", Min: "Min ms", Avg: "Povp. ms",
		Max: "Maks ms", Jitter: "Jitter ms", Since: "Od", Details: "Stanje",
	},
	StateOK:         "OK",
	StateTimeout:    "ČAS JE POTEKEL",
	StateNoResolve:  "IMENA GOSTITELJA NI MOGOČE RAZREŠITI",
	StatePortClosed: "VRATA ZAPRTA",
	NoHosts:         "Še ni gostiteljev",
	PressAToAdd:     "Pritisnite A, da dodate gostitelja",
	HowItWorks:      "Kako deluje",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "d"},

	MenuEdit:       "Uredi... (E)",
	MenuRemove:     "Odstrani (Del)",
	MenuDowntime:   "Izpadi...",
	MenuExport:     "Izvozi zgodovino...",
	MenuChange:     "Spremeni izbrane",
	BeepOn:         "Vklopi pisk",
	BeepOff:        "Izklopi pisk",
	CheckTCPPort:   "Preverjaj vrata TCP",
	UsePing:        "Uporabi ping",
	PingEvery:      "Ping vsakih",
	Timeout:        "Časovna omejitev",
	SlowAbove:      "Počasno nad",
	PortLabel:      "Vrata:",
	PingEveryMs:    "Ping vsakih, ms:",
	TimeoutMs:      "Časovna omejitev, ms:",
	SlowAboveMs:    "Počasno nad, ms:",
	SlowAboveMsOff: "Počasno nad, ms (0 - izklop):",

	RemoveHostsTitle: "Odstrani gostitelje",
	RemoveHost:       "Odstranim gostitelja?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Odstranim %d %s?", n, i18n.Plural("sl", n, "gostitelja", "gostitelja", "gostitelje", "gostiteljev"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... in še %d", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Izberite gostitelja",
	PingHistory: "Zgodovina pingov: ",
	Ms:          "ms",

	Name:          "Ime:",
	Host:          "Gostitelj:",
	TCPPort:       "Vrata TCP:",
	CheckInstead:  "Preverjaj namesto pinga",
	Port:          "Vrata:",
	SlowTooltip:   "Gostitelj je prikazan rumeno, ko njegov povprečni čas pinga preseže to vrednost. 0 - izklop",
	Notify:        "Obveščanje:",
	BeepDownBack:  "Pisk ob izpadu in vrnitvi",
	NotifyTooltip: "Izpad pomeni 3 neuspele pinge zapored. Naslov okna prikazuje, koliko takih gostiteljev je nedosegljivih.",
	Share:         "Deljenje:",
	ShareOn:       "Objavljaj čas pinga na u00.io",
	ShareTooltip:  "Časi pinga se pošiljajo na javno stran. Vidi jih lahko vsak, ki ima povezavo.",
	ShareOpen:     "Odpri",
	ShareCopy:     "Kopiraj povezavo",
	LinkCopied:    "Povezava kopirana",
	MenuShareOpen: "Odpri na u00.io",
	MenuShareCopy: "Kopiraj povezavo u00.io",
	ShareStart:    "Objavi na u00.io",
	ShareStop:     "Ustavi objavljanje",

	DowntimeTitle: "Izpadi",
	NoDowntime:    "V zadnjih 24 urah ni bilo izpadov",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Nedosegljiv enkrat v zadnjih 24 urah:"
		}
		return fmt.Sprintf("Nedosegljiv %d-krat v zadnjih 24 urah:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Izpadi v zadnjih 24 urah: %d", n) },
	Started:        "Začetek",
	Ended:          "Konec",
	Duration:       "Trajanje",
	StillDown:      "še vedno nedosegljiv",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Odpri konfiguracijo",
	Configs:            "Konfiguracije:",
	NewConfig:          "Nova...",
	SaveConfigAs:       "Shrani kot...",
	RemoveConfig:       "Odstrani",
	HostsCount:         "Število gostiteljev",
	Hosts:              "Gostitelji",
	NewConfigTitle:     "Nova konfiguracija",
	SaveConfigAsTitle:  "Shrani konfiguracijo kot",
	CopySuffix:         " kopija",
	RemoveConfigTitle:  "Odstrani konfiguracijo",
	RemoveConfigAsk:    "Odstranim izbrano konfiguracijo?",
	CannotRemoveOpened: "Trenutno odprte konfiguracije ni mogoče odstraniti.",

	ExtraColumns:   "Dodatni stolpci:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Od (zadnja sprememba)",
	AlwaysOnTop:    "Vedno na vrhu",
	Language:       "Jezik:",
	LanguageSystem: "Kot v sistemu",
	Theme:          "Tema:",
	ThemeDark:      "Temna",
	ThemeLight:     "Svetla",

	AboutTitle:   func(name string) string { return "O programu " + name },
	Version:      "Različica",
	Author:       "Avtor:",
	License:      "Licenca:",
	VisitWebsite: "Obišči spletno mesto",
	Close:        "Zapri",

	ExportTitle: "Izvoz zgodovine",
	CSVFiles:    "Datoteke CSV",
	SavedTo:     func(path string) string { return "Shranjeno v " + path },
	ConfigLoadFailed: func(err string) string {
		return "Seznama gostiteljev ni bilo mogoče prebrati, odprt je bil nov prazen seznam. Datoteka je ostala nespremenjena:\n\n" + err
	},

	Install:   "Namesti",
	Update:    "Posodobi",
	Uninstall: "Odstrani",
	InstallAsk: func(dir string) string {
		return "AltPing bo kopiran v " + dir + " in dodan v meni Start, na namizje in na seznam nameščenih aplikacij. Če je tam že nameščen, bo posodobljen. Nastavitve in seznami gostiteljev ostanejo nespremenjeni.\n\nAltPing se bo nato znova zagnal od tam."
	},
	InstallFailed: func(err string) string { return "Namestitev programa AltPing ni uspela:\n\n" + err },
	Installed:     "AltPing je nameščen",
	UninstallAsk: func(dir string) string {
		return "Odstranim AltPing iz tega računalnika?\n\nNastavitve in seznami gostiteljev ostanejo v " + dir
	},
	UninstallRunning: "AltPing se izvaja. Zaprite ga in poskusite znova.",
	Uninstalled:      "AltPing je bil odstranjen.",
}

var hr = Strings{
	Error: "Pogreška",

	DownCount: func(n int) string {
		return fmt.Sprintf("(%d %s)", n, i18n.Plural("hr", n, "nedostupan", "nedostupna", "nedostupnih"))
	},
	ToolConfigs:     "Konfiguracije (O)",
	ToolAddHost:     "Dodaj host (A)",
	ToolEditHost:    "Uredi host (E)",
	ToolRemoveHosts: "Ukloni odabrane hostove (Del)",
	ToolDetails:     "Detalji (D)",
	ToolStart:       "Pokreni pingiranje",
	ToolStop:        "Zaustavi pingiranje",
	ToolTray:        "Minimiziraj u područje obavijesti (T)",
	ToolShareOpen:   "Otvori odabrane hostove na u00.io (U)",
	ToolStyle:       "Stil (S)",
	KPIMedian:       "Medijan pinga",
	KPIUptime:       "Dostupnost 24 h",
	TrayShow:        "Prikaži",
	TrayQuit:        "Izlaz",
	Settings:        "Postavke",
	Help:            "Pomoć",
	About:           "O programu",
	ModeUDP:         "Način UDP",
	ModeICMP:        "Način ICMP",
	StatusHosts:     "Hostovi",
	StatusOK:        "OK",
	StatusSlow:      "Spori",
	StatusDown:      "Nedostupni",
	StatusStopped:   "Pingiranje zaustavljeno",

	Columns: ColumnStrings{
		Name: "Naziv", IP: "IP", Time: "Vrijeme ms", Trend: "Zadnjih 5 min", Loss: "Gubitak %", Min: "Min ms", Avg: "Prosj. ms",
		Max: "Maks ms", Jitter: "Jitter ms", Since: "Od", Details: "Stanje",
	},
	StateOK:         "OK",
	StateTimeout:    "ISTEKLO VRIJEME",
	StateNoResolve:  "NIJE MOGUĆE RAZRIJEŠITI NAZIV HOSTA",
	StatePortClosed: "PORT ZATVOREN",
	NoHosts:         "Još nema hostova",
	PressAToAdd:     "Pritisnite A za dodavanje hosta",
	HowItWorks:      "Kako funkcionira",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "d"},

	MenuEdit:       "Uredi... (E)",
	MenuRemove:     "Ukloni (Del)",
	MenuDowntime:   "Prekidi...",
	MenuExport:     "Izvezi povijest...",
	MenuChange:     "Promijeni odabrano",
	BeepOn:         "Uključi zvučni signal",
	BeepOff:        "Isključi zvučni signal",
	CheckTCPPort:   "Provjeravaj TCP port",
	UsePing:        "Koristi ping",
	PingEvery:      "Ping svakih",
	Timeout:        "Vremensko ograničenje",
	SlowAbove:      "Sporo iznad",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping svakih, ms:",
	TimeoutMs:      "Vremensko ograničenje, ms:",
	SlowAboveMs:    "Sporo iznad, ms:",
	SlowAboveMsOff: "Sporo iznad, ms (0 - isklj.):",

	RemoveHostsTitle: "Uklanjanje hostova",
	RemoveHost:       "Ukloniti host?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Ukloniti %d %s?", n, i18n.Plural("hr", n, "host", "hosta", "hostova"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... i još %d", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Odaberite host",
	PingHistory: "Povijest pinga: ",
	Ms:          "ms",

	Name:          "Naziv:",
	Host:          "Host:",
	TCPPort:       "TCP port:",
	CheckInstead:  "Provjeravaj umjesto pinga",
	Port:          "Port:",
	SlowTooltip:   "Host se prikazuje žuto kada mu prosječno vrijeme pinga prijeđe ovu vrijednost. 0 - isklj.",
	Notify:        "Obavijesti:",
	BeepDownBack:  "Zvučni signal pri prekidu i povratku",
	NotifyTooltip: "Prekid znači 3 neuspjela pinga zaredom. Naslov prozora prikazuje koliko je takvih hostova nedostupno.",
	Share:         "Dijeljenje:",
	ShareOn:       "Objavljuj vrijeme pinga na u00.io",
	ShareTooltip:  "Vremena pinga šalju se na javnu stranicu. Svatko tko ima vezu može ih vidjeti.",
	ShareOpen:     "Otvori",
	ShareCopy:     "Kopiraj vezu",
	LinkCopied:    "Veza je kopirana",
	MenuShareOpen: "Otvori na u00.io",
	MenuShareCopy: "Kopiraj vezu na u00.io",
	ShareStart:    "Objavi na u00.io",
	ShareStop:     "Zaustavi objavljivanje",

	DowntimeTitle: "Prekidi",
	NoDowntime:    "Nema prekida u posljednja 24 sata",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Nedostupan jednom u posljednja 24 sata:"
		}
		return fmt.Sprintf("Nedostupan %d %s u posljednja 24 sata:", n, i18n.Plural("hr", n, "put", "puta", "puta"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Prekidi u posljednja 24 sata: %d", n) },
	Started:        "Početak",
	Ended:          "Kraj",
	Duration:       "Trajanje",
	StillDown:      "još nedostupan",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Otvori konfiguraciju",
	Configs:            "Konfiguracije:",
	NewConfig:          "Nova...",
	SaveConfigAs:       "Spremi kao...",
	RemoveConfig:       "Ukloni",
	HostsCount:         "Broj hostova",
	Hosts:              "Hostovi",
	NewConfigTitle:     "Nova konfiguracija",
	SaveConfigAsTitle:  "Spremi konfiguraciju kao",
	CopySuffix:         " kopija",
	RemoveConfigTitle:  "Uklanjanje konfiguracije",
	RemoveConfigAsk:    "Ukloniti odabranu konfiguraciju?",
	CannotRemoveOpened: "Trenutno otvorena konfiguracija ne može se ukloniti.",

	ExtraColumns:   "Dodatni stupci:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Od (zadnja promjena)",
	AlwaysOnTop:    "Uvijek na vrhu",
	Language:       "Jezik:",
	LanguageSystem: "Kao u sustavu",
	Theme:          "Tema:",
	ThemeDark:      "Tamna",
	ThemeLight:     "Svijetla",

	AboutTitle:   func(name string) string { return "O programu " + name },
	Version:      "Verzija",
	Author:       "Autor:",
	License:      "Licenca:",
	VisitWebsite: "Posjeti web-mjesto",
	Close:        "Zatvori",

	ExportTitle: "Izvoz povijesti",
	CSVFiles:    "CSV datoteke",
	SavedTo:     func(path string) string { return "Spremljeno u " + path },
	ConfigLoadFailed: func(err string) string {
		return "Popis hostova nije bilo moguće pročitati, otvoren je novi prazan. Datoteka je ostala nepromijenjena:\n\n" + err
	},

	Install:   "Instaliraj",
	Update:    "Ažuriraj",
	Uninstall: "Deinstaliraj",
	InstallAsk: func(dir string) string {
		return "AltPing će se kopirati u " + dir + " i dodati u izbornik Start, na radnu površinu i na popis instaliranih aplikacija. Ako je ondje već instaliran, bit će ažuriran. Postavke i popisi hostova ostaju nepromijenjeni.\n\nAltPing će se zatim ponovno pokrenuti odande."
	},
	InstallFailed: func(err string) string { return "AltPing nije moguće instalirati:\n\n" + err },
	Installed:     "AltPing je instaliran",
	UninstallAsk: func(dir string) string {
		return "Ukloniti AltPing s ovog računala?\n\nPostavke i popisi hostova ostaju u " + dir
	},
	UninstallRunning: "AltPing je pokrenut. Zatvorite ga i pokušajte ponovno.",
	Uninstalled:      "AltPing je uklonjen.",
}

var hu = Strings{
	Error: "Hiba",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d nem elérhető)", n) },
	ToolConfigs:     "Konfigurációk (O)",
	ToolAddHost:     "Gazdagép hozzáadása (A)",
	ToolEditHost:    "Gazdagép szerkesztése (E)",
	ToolRemoveHosts: "Kijelölt gazdagépek eltávolítása (Del)",
	ToolDetails:     "Részletek (D)",
	ToolStart:       "Pingelés indítása",
	ToolStop:        "Pingelés leállítása",
	ToolTray:        "Kis méret az értesítési területre (T)",
	ToolShareOpen:   "Kijelölt gazdagépek megnyitása az u00.io-n (U)",
	ToolStyle:       "Stílus (S)",
	KPIMedian:       "Medián ping",
	KPIUptime:       "Rendelkezésre állás 24 ó",
	TrayShow:        "Megjelenítés",
	TrayQuit:        "Kilépés",
	Settings:        "Beállítások",
	Help:            "Súgó",
	About:           "Névjegy",
	ModeUDP:         "UDP mód",
	ModeICMP:        "ICMP mód",
	StatusHosts:     "Gazdagépek",
	StatusOK:        "OK",
	StatusSlow:      "Lassú",
	StatusDown:      "Nem elérhető",
	StatusStopped:   "Pingelés leállítva",

	Columns: ColumnStrings{
		Name: "Név", IP: "IP", Time: "Idő ms", Trend: "Utolsó 5 p", Loss: "Veszt. %", Min: "Min ms", Avg: "Átl. ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Mióta", Details: "Állapot",
	},
	StateOK:         "OK",
	StateTimeout:    "IDŐTÚLLÉPÉS",
	StateNoResolve:  "A GAZDAGÉPNÉV NEM OLDHATÓ FEL",
	StatePortClosed: "A PORT ZÁRVA",
	NoHosts:         "Még nincsenek gazdagépek",
	PressAToAdd:     "Gazdagép hozzáadásához nyomja meg az A billentyűt",
	HowItWorks:      "Hogyan működik",

	Units: UnitStrings{Second: "mp", Minute: "p", Hour: "ó", Day: "n"},

	MenuEdit:       "Szerkesztés... (E)",
	MenuRemove:     "Eltávolítás (Del)",
	MenuDowntime:   "Kiesések...",
	MenuExport:     "Előzmények exportálása...",
	MenuChange:     "Kijelöltek módosítása",
	BeepOn:         "Hangjelzés be",
	BeepOff:        "Hangjelzés ki",
	CheckTCPPort:   "TCP port ellenőrzése",
	UsePing:        "Ping használata",
	PingEvery:      "Pingelési időköz",
	Timeout:        "Időkorlát",
	SlowAbove:      "Lassúsági küszöb",
	PortLabel:      "Port:",
	PingEveryMs:    "Pingelési időköz, ms:",
	TimeoutMs:      "Időkorlát, ms:",
	SlowAboveMs:    "Lassúsági küszöb, ms:",
	SlowAboveMsOff: "Lassúsági küszöb, ms (0 - ki):",

	RemoveHostsTitle: "Gazdagépek eltávolítása",
	RemoveHost:       "Eltávolítja a gazdagépet?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("Eltávolít %d gazdagépet?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... és még %d", n) },

	Periods:     [3]string{"5 p", "1 ó", "24 ó"},
	SelectAHost: "Válasszon ki egy gazdagépet",
	PingHistory: "Pingelőzmények: ",
	Ms:          "ms",

	Name:          "Név:",
	Host:          "Gazdagép:",
	TCPPort:       "TCP port:",
	CheckInstead:  "Ellenőrzés ping helyett",
	Port:          "Port:",
	SlowTooltip:   "A gazdagép sárgán jelenik meg, ha az átlagos pingideje meghaladja ezt az értéket. 0 - ki",
	Notify:        "Értesítés:",
	BeepDownBack:  "Hangjelzés kieséskor és visszatéréskor",
	NotifyTooltip: "Kiesés: 3 egymást követő sikertelen ping. Az ablak címsora mutatja, hány ilyen gazdagép nem elérhető.",
	Share:         "Megosztás:",
	ShareOn:       "Pingidő közzététele az u00.io-n",
	ShareTooltip:  "A pingidők egy nyilvános oldalra kerülnek. A hivatkozás birtokában bárki megtekintheti őket.",
	ShareOpen:     "Megnyitás",
	ShareCopy:     "Hivatkozás másolása",
	LinkCopied:    "Hivatkozás másolva",
	MenuShareOpen: "Megnyitás az u00.io-n",
	MenuShareCopy: "u00.io-hivatkozás másolása",
	ShareStart:    "Közzététel az u00.io-n",
	ShareStop:     "Közzététel leállítása",

	DowntimeTitle: "Kiesések",
	NoDowntime:    "Nem volt kiesés az elmúlt 24 órában",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Egyszer volt nem elérhető az elmúlt 24 órában:"
		}
		return fmt.Sprintf("%d alkalommal volt nem elérhető az elmúlt 24 órában:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Kiesések az elmúlt 24 órában: %d", n) },
	Started:        "Kezdete",
	Ended:          "Vége",
	Duration:       "Időtartam",
	StillDown:      "még mindig nem elérhető",
	DateTimeLayout: "01.02. 15:04:05",

	ConfigsTitle:       "Konfiguráció megnyitása",
	Configs:            "Konfigurációk:",
	NewConfig:          "Új...",
	SaveConfigAs:       "Mentés másként...",
	RemoveConfig:       "Eltávolítás",
	HostsCount:         "Gazdagépek száma",
	Hosts:              "Gazdagépek",
	NewConfigTitle:     "Új konfiguráció",
	SaveConfigAsTitle:  "Konfiguráció mentése másként",
	CopySuffix:         " másolat",
	RemoveConfigTitle:  "Konfiguráció eltávolítása",
	RemoveConfigAsk:    "Eltávolítja a kijelölt konfigurációt?",
	CannotRemoveOpened: "Az éppen megnyitott konfiguráció nem távolítható el.",

	ExtraColumns:   "További oszlopok:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Mióta (utolsó változás)",
	AlwaysOnTop:    "Mindig felül",
	Language:       "Nyelv:",
	LanguageSystem: "A rendszer szerint",
	Theme:          "Téma:",
	ThemeDark:      "Sötét",
	ThemeLight:     "Világos",

	AboutTitle:   func(name string) string { return name + " névjegye" },
	Version:      "Verzió",
	Author:       "Szerző:",
	License:      "Licenc:",
	VisitWebsite: "Webhely megnyitása",
	Close:        "Bezárás",

	ExportTitle: "Előzmények exportálása",
	CSVFiles:    "CSV-fájlok",
	SavedTo:     func(path string) string { return "Mentve ide: " + path },
	ConfigLoadFailed: func(err string) string {
		return "A gazdagéplista nem olvasható be, ezért egy új, üres lista nyílt meg. A fájl változatlan maradt:\n\n" + err
	},

	Install:   "Telepítés",
	Update:    "Frissítés",
	Uninstall: "Eltávolítás",
	InstallAsk: func(dir string) string {
		return "Az AltPing a következő helyre lesz másolva: " + dir + ". Bekerül a Start menübe, az asztalra és a telepített alkalmazások listájába. Ha már telepítve van ott, frissül. A beállítások és a gazdagéplisták változatlanok maradnak.\n\nEzután az AltPing onnan indul újra."
	},
	InstallFailed: func(err string) string { return "Az AltPing nem telepíthető:\n\n" + err },
	Installed:     "Az AltPing telepítve van",
	UninstallAsk: func(dir string) string {
		return "Eltávolítja az AltPinget erről a számítógépről?\n\nA beállítások és a gazdagéplisták itt maradnak: " + dir
	},
	UninstallRunning: "Az AltPing fut. Zárja be, és próbálja újra.",
	Uninstalled:      "Az AltPing el lett távolítva.",
}

var ro = Strings{
	Error: "Eroare",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d indisponibile)", n) },
	ToolConfigs:     "Configurații (O)",
	ToolAddHost:     "Adăugare gazdă (A)",
	ToolEditHost:    "Editare gazdă (E)",
	ToolRemoveHosts: "Eliminare gazde selectate (Del)",
	ToolDetails:     "Detalii (D)",
	ToolStart:       "Pornire ping",
	ToolStop:        "Oprire ping",
	ToolTray:        "Minimizare în zona de notificare (T)",
	ToolShareOpen:   "Deschidere gazde selectate pe u00.io (U)",
	ToolStyle:       "Stil (S)",
	KPIMedian:       "Ping median",
	KPIUptime:       "Disponibilitate 24 h",
	TrayShow:        "Afișare",
	TrayQuit:        "Ieșire",
	Settings:        "Setări",
	Help:            "Ajutor",
	About:           "Despre",
	ModeUDP:         "Mod UDP",
	ModeICMP:        "Mod ICMP",
	StatusHosts:     "Gazde",
	StatusOK:        "OK",
	StatusSlow:      "Lente",
	StatusDown:      "Indisponibile",
	StatusStopped:   "Ping oprit",

	Columns: ColumnStrings{
		Name: "Nume", IP: "IP", Time: "Timp ms", Trend: "Ultimele 5 min", Loss: "Pierderi %", Min: "Min ms", Avg: "Med ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "De la", Details: "Detalii",
	},
	StateOK:         "OK",
	StateTimeout:    "TIMP EXPIRAT",
	StateNoResolve:  "NUMELE GAZDEI NU POATE FI REZOLVAT",
	StatePortClosed: "PORT ÎNCHIS",
	NoHosts:         "Încă nu există gazde",
	PressAToAdd:     "Apăsați A pentru a adăuga o gazdă",
	HowItWorks:      "Cum funcționează",

	Units: UnitStrings{Second: "s", Minute: "m", Hour: "h", Day: "z"},

	MenuEdit:       "Editare... (E)",
	MenuRemove:     "Eliminare (Del)",
	MenuDowntime:   "Întreruperi...",
	MenuExport:     "Exportare istoric...",
	MenuChange:     "Modificare selectate",
	BeepOn:         "Activare semnal sonor",
	BeepOff:        "Dezactivare semnal sonor",
	CheckTCPPort:   "Verificare port TCP",
	UsePing:        "Utilizare ping",
	PingEvery:      "Ping la fiecare",
	Timeout:        "Timp de expirare",
	SlowAbove:      "Lent peste",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping la fiecare, ms:",
	TimeoutMs:      "Timp de expirare, ms:",
	SlowAboveMs:    "Lent peste, ms:",
	SlowAboveMsOff: "Lent peste, ms (0 - dezactivat):",

	RemoveHostsTitle: "Eliminare gazde",
	RemoveHost:       "Eliminați gazda?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Eliminați %d %s?", n, i18n.Plural("ro", n, "gazdă", "gazde", "de gazde"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... și încă %d", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Selectați o gazdă",
	PingHistory: "Istoric ping: ",
	Ms:          "ms",

	Name:          "Nume:",
	Host:          "Gazdă:",
	TCPPort:       "Port TCP:",
	CheckInstead:  "Verificare în loc de ping",
	Port:          "Port:",
	SlowTooltip:   "Gazda este afișată cu galben când timpul mediu de ping depășește această valoare. 0 - dezactivat",
	Notify:        "Notificare:",
	BeepDownBack:  "Semnal sonor la căderea și revenirea gazdei",
	NotifyTooltip: "Indisponibilă înseamnă 3 ping-uri eșuate consecutiv. Titlul ferestrei arată câte astfel de gazde sunt indisponibile.",
	Share:         "Partajare:",
	ShareOn:       "Publicare timp de ping pe u00.io",
	ShareTooltip:  "Timpii de ping sunt trimiși către o pagină publică. Oricine are linkul îi poate vedea.",
	ShareOpen:     "Deschidere",
	ShareCopy:     "Copiere link",
	LinkCopied:    "Link copiat",
	MenuShareOpen: "Deschidere pe u00.io",
	MenuShareCopy: "Copiere link u00.io",
	ShareStart:    "Publicare pe u00.io",
	ShareStop:     "Oprire publicare",

	DowntimeTitle: "Întreruperi",
	NoDowntime:    "Nicio întrerupere în ultimele 24 de ore",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Indisponibilă o dată în ultimele 24 de ore:"
		}
		return fmt.Sprintf("Indisponibilă de %d %s în ultimele 24 de ore:", n, i18n.Plural("ro", n, "dată", "ori", "de ori"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Întreruperi în ultimele 24 de ore: %d", n) },
	Started:        "Început",
	Ended:          "Sfârșit",
	Duration:       "Durată",
	StillDown:      "încă indisponibilă",
	DateTimeLayout: "02.01 15:04:05",

	ConfigsTitle:       "Deschidere configurație",
	Configs:            "Configurații:",
	NewConfig:          "Nouă...",
	SaveConfigAs:       "Salvare ca...",
	RemoveConfig:       "Eliminare",
	HostsCount:         "Nr. gazde",
	Hosts:              "Gazde",
	NewConfigTitle:     "Configurație nouă",
	SaveConfigAsTitle:  "Salvare configurație ca",
	CopySuffix:         " copie",
	RemoveConfigTitle:  "Eliminare configurație",
	RemoveConfigAsk:    "Eliminați configurația selectată?",
	CannotRemoveOpened: "Configurația deschisă în prezent nu poate fi eliminată.",

	ExtraColumns:   "Coloane suplimentare:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "De la (ultima modificare)",
	AlwaysOnTop:    "Mereu deasupra",
	Language:       "Limbă:",
	LanguageSystem: "Ca în sistem",
	Theme:          "Temă:",
	ThemeDark:      "Întunecată",
	ThemeLight:     "Luminoasă",

	AboutTitle:   func(name string) string { return "Despre " + name },
	Version:      "Versiune",
	Author:       "Autor:",
	License:      "Licență:",
	VisitWebsite: "Vizitați site-ul web",
	Close:        "Închidere",

	ExportTitle: "Exportare istoric",
	CSVFiles:    "Fișiere CSV",
	SavedTo:     func(path string) string { return "Salvat în " + path },
	ConfigLoadFailed: func(err string) string {
		return "Lista de gazde nu a putut fi citită, așa că a fost deschisă una nouă, goală. Fișierul a rămas neschimbat:\n\n" + err
	},

	Install:   "Instalare",
	Update:    "Actualizare",
	Uninstall: "Dezinstalare",
	InstallAsk: func(dir string) string {
		return "AltPing va fi copiat în " + dir + " și adăugat în meniul Start, pe desktop și în lista de aplicații instalate. Dacă este deja instalat acolo, va fi actualizat. Setările și listele de gazde rămân neschimbate.\n\nApoi AltPing va reporni de acolo."
	},
	InstallFailed: func(err string) string { return "AltPing nu a putut fi instalat:\n\n" + err },
	Installed:     "AltPing este instalat",
	UninstallAsk: func(dir string) string {
		return "Eliminați AltPing de pe acest computer?\n\nSetările și listele de gazde rămân în " + dir
	},
	UninstallRunning: "AltPing rulează. Închideți-l și încercați din nou.",
	Uninstalled:      "AltPing a fost eliminat.",
}

var bg = Strings{
	Error: "Грешка",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d недостъпни)", n) },
	ToolConfigs:     "Конфигурации (O)",
	ToolAddHost:     "Добавяне на хост (A)",
	ToolEditHost:    "Редактиране на хост (E)",
	ToolRemoveHosts: "Премахване на избраните хостове (Del)",
	ToolDetails:     "Подробности (D)",
	ToolStart:       "Започване на пинг",
	ToolStop:        "Спиране на пинг",
	ToolTray:        "Минимизиране в областта за уведомяване (T)",
	ToolShareOpen:   "Отваряне на избраните хостове в u00.io (U)",
	ToolStyle:       "Стил (S)",
	KPIMedian:       "Медиана на пинга",
	KPIUptime:       "Достъпност 24 ч",
	TrayShow:        "Показване",
	TrayQuit:        "Изход",
	Settings:        "Настройки",
	Help:            "Помощ",
	About:           "Относно",
	ModeUDP:         "Режим UDP",
	ModeICMP:        "Режим ICMP",
	StatusHosts:     "Хостове",
	StatusOK:        "OK",
	StatusSlow:      "Бавни",
	StatusDown:      "Недостъпни",
	StatusStopped:   "Пингът е спрян",

	Columns: ColumnStrings{
		Name: "Име", IP: "IP", Time: "Време мс", Trend: "Последни 5 мин", Loss: "Загуби %", Min: "Мин мс", Avg: "Ср мс",
		Max: "Макс мс", Jitter: "Джитър мс", Since: "От", Details: "Подробности",
	},
	StateOK:         "OK",
	StateTimeout:    "ИЗТЕКЛО ВРЕМЕ",
	StateNoResolve:  "ИМЕТО НА ХОСТА НЕ Е ОТКРИТО",
	StatePortClosed: "ПОРТЪТ Е ЗАТВОРЕН",
	NoHosts:         "Все още няма хостове",
	PressAToAdd:     "Натиснете A, за да добавите хост",
	HowItWorks:      "Как работи",

	Units: UnitStrings{Second: "с", Minute: "м", Hour: "ч", Day: "д"},

	MenuEdit:       "Редактиране... (E)",
	MenuRemove:     "Премахване (Del)",
	MenuDowntime:   "Прекъсвания...",
	MenuExport:     "Експортиране на историята...",
	MenuChange:     "Промяна на избраните",
	BeepOn:         "Включване на звуковия сигнал",
	BeepOff:        "Изключване на звуковия сигнал",
	CheckTCPPort:   "Проверка на TCP порт",
	UsePing:        "Използване на пинг",
	PingEvery:      "Пинг на всеки",
	Timeout:        "Време за изчакване",
	SlowAbove:      "Бавно над",
	PortLabel:      "Порт:",
	PingEveryMs:    "Пинг на всеки, мс:",
	TimeoutMs:      "Време за изчакване, мс:",
	SlowAboveMs:    "Бавно над, мс:",
	SlowAboveMsOff: "Бавно над, мс (0 - изкл.):",

	RemoveHostsTitle: "Премахване на хостове",
	RemoveHost:       "Да се премахне ли хостът?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Да се премахнат ли %d %s?", n, i18n.Plural("bg", n, "хост", "хоста"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... и още %d", n) },

	Periods:     [3]string{"5 мин", "1 ч", "24 ч"},
	SelectAHost: "Изберете хост",
	PingHistory: "История на пинга: ",
	Ms:          "мс",

	Name:          "Име:",
	Host:          "Хост:",
	TCPPort:       "TCP порт:",
	CheckInstead:  "Проверка вместо пинг",
	Port:          "Порт:",
	SlowTooltip:   "Хостът се показва в жълто, когато средното му време за пинг е над тази стойност. 0 - изкл.",
	Notify:        "Известяване:",
	BeepDownBack:  "Звуков сигнал при отпадане и възстановяване",
	NotifyTooltip: "Недостъпен означава 3 неуспешни пинга подред. Заглавието на прозореца показва колко такива хоста са недостъпни.",
	Share:         "Споделяне:",
	ShareOn:       "Публикуване на времето за пинг в u00.io",
	ShareTooltip:  "Времената за пинг се изпращат към публична страница. Всеки, който има връзката, може да ги види.",
	ShareOpen:     "Отваряне",
	ShareCopy:     "Копиране на връзката",
	LinkCopied:    "Връзката е копирана",
	MenuShareOpen: "Отваряне в u00.io",
	MenuShareCopy: "Копиране на връзката към u00.io",
	ShareStart:    "Публикуване в u00.io",
	ShareStop:     "Спиране на публикуването",

	DowntimeTitle: "Прекъсвания",
	NoDowntime:    "Няма прекъсвания през последните 24 часа",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Недостъпен веднъж през последните 24 часа:"
		}
		return fmt.Sprintf("Недостъпен %d %s през последните 24 часа:", n, i18n.Plural("bg", n, "път", "пъти"))
	},
	GroupDownTimes: func(n int) string {
		return fmt.Sprintf("Прекъсвания през последните 24 часа: %d", n)
	},
	Started:        "Начало",
	Ended:          "Край",
	Duration:       "Продължителност",
	StillDown:      "все още недостъпен",
	DateTimeLayout: "02.01 15:04:05",

	ConfigsTitle:       "Отваряне на конфигурация",
	Configs:            "Конфигурации:",
	NewConfig:          "Нова...",
	SaveConfigAs:       "Запиши като...",
	RemoveConfig:       "Премахване",
	HostsCount:         "Брой хостове",
	Hosts:              "Хостове",
	NewConfigTitle:     "Нова конфигурация",
	SaveConfigAsTitle:  "Записване на конфигурацията като",
	CopySuffix:         " копие",
	RemoveConfigTitle:  "Премахване на конфигурация",
	RemoveConfigAsk:    "Да се премахне ли избраната конфигурация?",
	CannotRemoveOpened: "Текущо отворената конфигурация не може да бъде премахната.",

	ExtraColumns:   "Допълнителни колони:",
	ShowMin:        "Мин",
	ShowJitter:     "Джитър",
	ShowSince:      "От (последна промяна)",
	AlwaysOnTop:    "Винаги отгоре",
	Language:       "Език:",
	LanguageSystem: "Както в системата",
	Theme:          "Тема:",
	ThemeDark:      "Тъмна",
	ThemeLight:     "Светла",

	AboutTitle:   func(name string) string { return "Относно " + name },
	Version:      "Версия",
	Author:       "Автор:",
	License:      "Лиценз:",
	VisitWebsite: "Посетете уебсайта",
	Close:        "Затваряне",

	ExportTitle: "Експортиране на историята",
	CSVFiles:    "CSV файлове",
	SavedTo:     func(path string) string { return "Записано в " + path },
	ConfigLoadFailed: func(err string) string {
		return "Списъкът с хостове не можа да бъде прочетен, затова е отворен нов празен. Файлът е оставен непроменен:\n\n" + err
	},

	Install:   "Инсталиране",
	Update:    "Актуализиране",
	Uninstall: "Деинсталиране",
	InstallAsk: func(dir string) string {
		return "AltPing ще бъде копиран в " + dir + " и добавен в менюто „Старт“, на работния плот и в списъка с инсталирани приложения. Ако вече е инсталиран там, ще бъде актуализиран. Настройките и списъците с хостове остават непроменени.\n\nСлед това AltPing ще се рестартира оттам."
	},
	InstallFailed: func(err string) string { return "AltPing не можа да бъде инсталиран:\n\n" + err },
	Installed:     "AltPing е инсталиран",
	UninstallAsk: func(dir string) string {
		return "Да се премахне ли AltPing от този компютър?\n\nНастройките и списъците с хостове остават в " + dir
	},
	UninstallRunning: "AltPing работи. Затворете го и опитайте отново.",
	Uninstalled:      "AltPing е премахнат.",
}

var el = Strings{
	Error: "Σφάλμα",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d μη διαθέσιμοι)", n) },
	ToolConfigs:     "Διαμορφώσεις (O)",
	ToolAddHost:     "Προσθήκη κεντρικού υπολογιστή (A)",
	ToolEditHost:    "Επεξεργασία κεντρικού υπολογιστή (E)",
	ToolRemoveHosts: "Κατάργηση επιλεγμένων κεντρικών υπολογιστών (Del)",
	ToolDetails:     "Λεπτομέρειες (D)",
	ToolStart:       "Έναρξη ping",
	ToolStop:        "Διακοπή ping",
	ToolTray:        "Ελαχιστοποίηση στην περιοχή ειδοποιήσεων (T)",
	ToolShareOpen:   "Άνοιγμα επιλεγμένων στο u00.io (U)",
	ToolStyle:       "Στυλ (S)",
	KPIMedian:       "Διάμεσο ping",
	KPIUptime:       "Διαθεσιμότητα 24 ω",
	TrayShow:        "Εμφάνιση",
	TrayQuit:        "Έξοδος",
	Settings:        "Ρυθμίσεις",
	Help:            "Βοήθεια",
	About:           "Πληροφορίες",
	ModeUDP:         "Λειτουργία UDP",
	ModeICMP:        "Λειτουργία ICMP",
	StatusHosts:     "Κεντρικοί υπολογιστές",
	StatusOK:        "OK",
	StatusSlow:      "Αργοί",
	StatusDown:      "Μη διαθέσιμοι",
	StatusStopped:   "Το ping διακόπηκε",

	Columns: ColumnStrings{
		Name: "Όνομα", IP: "IP", Time: "Χρόνος ms", Trend: "Τελ. 5 λεπτά", Loss: "Απώλ. %", Min: "Ελάχ. ms", Avg: "Μέσ. ms",
		Max: "Μέγ. ms", Jitter: "Διακύμ. ms", Since: "Από", Details: "Λεπτομέρειες",
	},
	StateOK:         "OK",
	StateTimeout:    "ΛΗΞΗ ΧΡΟΝΙΚΟΥ ΟΡΙΟΥ",
	StateNoResolve:  "ΑΔΥΝΑΜΙΑ ΕΠΙΛΥΣΗΣ ΟΝΟΜΑΤΟΣ",
	StatePortClosed: "ΘΥΡΑ ΚΛΕΙΣΤΗ",
	NoHosts:         "Δεν υπάρχουν ακόμη κεντρικοί υπολογιστές",
	PressAToAdd:     "Πατήστε A για να προσθέσετε κεντρικό υπολογιστή",
	HowItWorks:      "Πώς λειτουργεί",

	Units: UnitStrings{Second: "δ", Minute: "λ", Hour: "ώ", Day: "ημ"},

	MenuEdit:       "Επεξεργασία... (E)",
	MenuRemove:     "Κατάργηση (Del)",
	MenuDowntime:   "Διακοπές λειτουργίας...",
	MenuExport:     "Εξαγωγή ιστορικού...",
	MenuChange:     "Αλλαγή επιλεγμένων",
	BeepOn:         "Ενεργοποίηση ήχου",
	BeepOff:        "Απενεργοποίηση ήχου",
	CheckTCPPort:   "Έλεγχος θύρας TCP",
	UsePing:        "Χρήση ping",
	PingEvery:      "Ping κάθε",
	Timeout:        "Χρονικό όριο",
	SlowAbove:      "Αργό πάνω από",
	PortLabel:      "Θύρα:",
	PingEveryMs:    "Ping κάθε, ms:",
	TimeoutMs:      "Χρονικό όριο, ms:",
	SlowAboveMs:    "Αργό πάνω από, ms:",
	SlowAboveMsOff: "Αργό πάνω από, ms (0 - ανενεργό):",

	RemoveHostsTitle: "Κατάργηση κεντρικών υπολογιστών",
	RemoveHost:       "Κατάργηση του κεντρικού υπολογιστή;",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Κατάργηση %d %s;", n, i18n.Plural("el", n, "κεντρικού υπολογιστή", "κεντρικών υπολογιστών"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... και %d ακόμη", n) },

	Periods:     [3]string{"5λ", "1ώ", "24ώ"},
	SelectAHost: "Επιλέξτε κεντρικό υπολογιστή",
	PingHistory: "Ιστορικό ping: ",
	Ms:          "ms",

	Name:          "Όνομα:",
	Host:          "Κεντρικός υπολογιστής:",
	TCPPort:       "Θύρα TCP:",
	CheckInstead:  "Έλεγχος αντί για ping",
	Port:          "Θύρα:",
	SlowTooltip:   "Ο κεντρικός υπολογιστής εμφανίζεται με κίτρινο όταν ο μέσος χρόνος ping υπερβαίνει αυτή την τιμή. 0 - ανενεργό",
	Notify:        "Ειδοποίηση:",
	BeepDownBack:  "Ήχος κατά τη διακοπή και την επαναφορά",
	NotifyTooltip: "Μη διαθέσιμος σημαίνει 3 αποτυχημένα ping στη σειρά. Ο τίτλος του παραθύρου δείχνει πόσοι τέτοιοι κεντρικοί υπολογιστές είναι μη διαθέσιμοι.",
	Share:         "Κοινή χρήση:",
	ShareOn:       "Δημοσίευση χρόνου ping στο u00.io",
	ShareTooltip:  "Οι χρόνοι ping αποστέλλονται σε μια δημόσια σελίδα. Όποιος έχει τη σύνδεση μπορεί να τους δει.",
	ShareOpen:     "Άνοιγμα",
	ShareCopy:     "Αντιγραφή σύνδεσης",
	LinkCopied:    "Η σύνδεση αντιγράφηκε",
	MenuShareOpen: "Άνοιγμα στο u00.io",
	MenuShareCopy: "Αντιγραφή σύνδεσης u00.io",
	ShareStart:    "Δημοσίευση στο u00.io",
	ShareStop:     "Διακοπή δημοσίευσης",

	DowntimeTitle: "Διακοπές λειτουργίας",
	NoDowntime:    "Καμία διακοπή τις τελευταίες 24 ώρες",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Μη διαθέσιμος μία φορά τις τελευταίες 24 ώρες:"
		}
		return fmt.Sprintf("Μη διαθέσιμος %d %s τις τελευταίες 24 ώρες:", n, i18n.Plural("el", n, "φορά", "φορές"))
	},
	GroupDownTimes: func(n int) string {
		return fmt.Sprintf("Διακοπές τις τελευταίες 24 ώρες: %d", n)
	},
	Started:        "Έναρξη",
	Ended:          "Λήξη",
	Duration:       "Διάρκεια",
	StillDown:      "ακόμη μη διαθέσιμος",
	DateTimeLayout: "02/01 15:04:05",

	ConfigsTitle:       "Άνοιγμα διαμόρφωσης",
	Configs:            "Διαμορφώσεις:",
	NewConfig:          "Νέα...",
	SaveConfigAs:       "Αποθήκευση ως...",
	RemoveConfig:       "Κατάργηση",
	HostsCount:         "Πλήθος υπολογιστών",
	Hosts:              "Κεντρικοί υπολογιστές",
	NewConfigTitle:     "Νέα διαμόρφωση",
	SaveConfigAsTitle:  "Αποθήκευση διαμόρφωσης ως",
	CopySuffix:         " αντίγραφο",
	RemoveConfigTitle:  "Κατάργηση διαμόρφωσης",
	RemoveConfigAsk:    "Κατάργηση της επιλεγμένης διαμόρφωσης;",
	CannotRemoveOpened: "Δεν είναι δυνατή η κατάργηση της διαμόρφωσης που είναι ανοιχτή.",

	ExtraColumns:   "Επιπλέον στήλες:",
	ShowMin:        "Ελάχ.",
	ShowJitter:     "Διακύμανση",
	ShowSince:      "Από (τελευταία αλλαγή)",
	AlwaysOnTop:    "Πάντα στο προσκήνιο",
	Language:       "Γλώσσα:",
	LanguageSystem: "Όπως στο σύστημα",
	Theme:          "Θέμα:",
	ThemeDark:      "Σκούρο",
	ThemeLight:     "Ανοιχτόχρωμο",

	AboutTitle:   func(name string) string { return "Πληροφορίες για το " + name },
	Version:      "Έκδοση",
	Author:       "Δημιουργός:",
	License:      "Άδεια χρήσης:",
	VisitWebsite: "Επίσκεψη στην τοποθεσία web",
	Close:        "Κλείσιμο",

	ExportTitle: "Εξαγωγή ιστορικού",
	CSVFiles:    "Αρχεία CSV",
	SavedTo:     func(path string) string { return "Αποθηκεύτηκε στο " + path },
	ConfigLoadFailed: func(err string) string {
		return "Δεν ήταν δυνατή η ανάγνωση της λίστας κεντρικών υπολογιστών, οπότε ανοίχτηκε μια νέα κενή. Το αρχείο παρέμεινε ως είχε:\n\n" + err
	},

	Install:   "Εγκατάσταση",
	Update:    "Ενημέρωση",
	Uninstall: "Κατάργηση εγκατάστασης",
	InstallAsk: func(dir string) string {
		return "Το AltPing θα αντιγραφεί στο " + dir + " και θα προστεθεί στο μενού Έναρξη, στην επιφάνεια εργασίας και στη λίστα εγκατεστημένων εφαρμογών. Αν είναι ήδη εγκατεστημένο εκεί, θα ενημερωθεί. Οι ρυθμίσεις και οι λίστες κεντρικών υπολογιστών παραμένουν ως έχουν.\n\nΣτη συνέχεια το AltPing θα επανεκκινηθεί από εκεί."
	},
	InstallFailed: func(err string) string {
		return "Δεν ήταν δυνατή η εγκατάσταση του AltPing:\n\n" + err
	},
	Installed: "Το AltPing εγκαταστάθηκε",
	UninstallAsk: func(dir string) string {
		return "Κατάργηση του AltPing από αυτόν τον υπολογιστή;\n\nΟι ρυθμίσεις και οι λίστες κεντρικών υπολογιστών παραμένουν στο " + dir
	},
	UninstallRunning: "Το AltPing εκτελείται. Κλείστε το και δοκιμάστε ξανά.",
	Uninstalled:      "Το AltPing καταργήθηκε.",
}

var lt = Strings{
	Error: "Klaida",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d nepasiekiami)", n) },
	ToolConfigs:     "Konfigūracijos (O)",
	ToolAddHost:     "Pridėti mazgą (A)",
	ToolEditHost:    "Redaguoti mazgą (E)",
	ToolRemoveHosts: "Šalinti pažymėtus mazgus (Del)",
	ToolDetails:     "Išsami informacija (D)",
	ToolStart:       "Pradėti ping",
	ToolStop:        "Stabdyti ping",
	ToolTray:        "Sumažinti į pranešimų sritį (T)",
	ToolShareOpen:   "Atidaryti pažymėtus mazgus u00.io (U)",
	ToolStyle:       "Stilius (S)",
	KPIMedian:       "Pingo mediana",
	KPIUptime:       "Pasiekiamumas 24 val.",
	TrayShow:        "Rodyti",
	TrayQuit:        "Išeiti",
	Settings:        "Parametrai",
	Help:            "Žinynas",
	About:           "Apie",
	ModeUDP:         "UDP režimas",
	ModeICMP:        "ICMP režimas",
	StatusHosts:     "Mazgai",
	StatusOK:        "Gerai",
	StatusSlow:      "Lėti",
	StatusDown:      "Nepasiekiami",
	StatusStopped:   "Ping sustabdytas",

	Columns: ColumnStrings{
		Name: "Pavadinimas", IP: "IP", Time: "Laikas ms", Trend: "Pask. 5 min", Loss: "Praradimai %", Min: "Min ms", Avg: "Vid ms",
		Max: "Maks ms", Jitter: "Svyr. ms", Since: "Nuo", Details: "Informacija",
	},
	StateOK:         "OK",
	StateTimeout:    "BAIGĖSI SKIRTASIS LAIKAS",
	StateNoResolve:  "NEPAVYKO RASTI MAZGO VARDO",
	StatePortClosed: "PRIEVADAS UŽDARYTAS",
	NoHosts:         "Mazgų dar nėra",
	PressAToAdd:     "Paspauskite A, kad pridėtumėte mazgą",
	HowItWorks:      "Kaip tai veikia",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "val", Day: "d"},

	MenuEdit:       "Redaguoti... (E)",
	MenuRemove:     "Šalinti (Del)",
	MenuDowntime:   "Prastovos...",
	MenuExport:     "Eksportuoti istoriją...",
	MenuChange:     "Keisti pažymėtus",
	BeepOn:         "Įjungti garsinį signalą",
	BeepOff:        "Išjungti garsinį signalą",
	CheckTCPPort:   "Tikrinti TCP prievadą",
	UsePing:        "Naudoti ping",
	PingEvery:      "Ping kas",
	Timeout:        "Skirtasis laikas",
	SlowAbove:      "Lėta, jei daugiau nei",
	PortLabel:      "Prievadas:",
	PingEveryMs:    "Ping kas, ms:",
	TimeoutMs:      "Skirtasis laikas, ms:",
	SlowAboveMs:    "Lėta, jei daugiau nei, ms:",
	SlowAboveMsOff: "Lėta, jei daugiau nei, ms (0 - išj.):",

	RemoveHostsTitle: "Mazgų šalinimas",
	RemoveHost:       "Šalinti mazgą?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Šalinti %d %s?", n, i18n.Plural("lt", n, "mazgą", "mazgus", "mazgų"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... ir dar %d", n) },

	Periods:     [3]string{"5 min", "1 val", "24 val"},
	SelectAHost: "Pasirinkite mazgą",
	PingHistory: "Ping istorija: ",
	Ms:          "ms",

	Name:          "Pavadinimas:",
	Host:          "Mazgas:",
	TCPPort:       "TCP prievadas:",
	CheckInstead:  "Tikrinti vietoj ping",
	Port:          "Prievadas:",
	SlowTooltip:   "Mazgas rodomas geltonai, kai jo vidutinis ping laikas viršija šią reikšmę. 0 - išjungta",
	Notify:        "Pranešti:",
	BeepDownBack:  "Garsinis signalas, kai tampa nepasiekiamas arba vėl pasiekiamas",
	NotifyTooltip: "Nepasiekiamas reiškia 3 nepavykusius ping iš eilės. Lango pavadinime rodoma, kiek tokių mazgų yra nepasiekiami.",
	Share:         "Bendrinti:",
	ShareOn:       "Skelbti ping laiką u00.io",
	ShareTooltip:  "Ping laikai siunčiami į viešą puslapį. Juos gali matyti visi, kas turi saitą.",
	ShareOpen:     "Atidaryti",
	ShareCopy:     "Kopijuoti saitą",
	LinkCopied:    "Saitas nukopijuotas",
	MenuShareOpen: "Atidaryti u00.io",
	MenuShareCopy: "Kopijuoti u00.io saitą",
	ShareStart:    "Skelbti u00.io",
	ShareStop:     "Nebeskelbti",

	DowntimeTitle: "Prastovos",
	NoDowntime:    "Per pastarąsias 24 valandas prastovų nebuvo",
	DownTimes: func(n int) string {
		return fmt.Sprintf("Nepasiekiamas %d %s per pastarąsias 24 valandas:", n, i18n.Plural("lt", n, "kartą", "kartus", "kartų"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Prastovos per pastarąsias 24 valandas: %d", n) },
	Started:        "Pradžia",
	Ended:          "Pabaiga",
	Duration:       "Trukmė",
	StillDown:      "vis dar nepasiekiamas",
	DateTimeLayout: "01-02 15:04:05",

	ConfigsTitle:       "Atidaryti konfigūraciją",
	Configs:            "Konfigūracijos:",
	NewConfig:          "Nauja...",
	SaveConfigAs:       "Įrašyti kaip...",
	RemoveConfig:       "Šalinti",
	HostsCount:         "Mazgų sk.",
	Hosts:              "Mazgai",
	NewConfigTitle:     "Nauja konfigūracija",
	SaveConfigAsTitle:  "Įrašyti konfigūraciją kaip",
	CopySuffix:         " kopija",
	RemoveConfigTitle:  "Konfigūracijos šalinimas",
	RemoveConfigAsk:    "Šalinti pažymėtą konfigūraciją?",
	CannotRemoveOpened: "Šiuo metu atidarytos konfigūracijos pašalinti negalima.",

	ExtraColumns:   "Papildomi stulpeliai:",
	ShowMin:        "Min",
	ShowJitter:     "Svyravimas",
	ShowSince:      "Nuo (paskutinis pokytis)",
	AlwaysOnTop:    "Visada viršuje",
	Language:       "Kalba:",
	LanguageSystem: "Kaip sistemoje",
	Theme:          "Tema:",
	ThemeDark:      "Tamsi",
	ThemeLight:     "Šviesi",

	AboutTitle:   func(name string) string { return "Apie " + name },
	Version:      "Versija",
	Author:       "Autorius:",
	License:      "Licencija:",
	VisitWebsite: "Apsilankyti svetainėje",
	Close:        "Uždaryti",

	ExportTitle: "Istorijos eksportas",
	CSVFiles:    "CSV failai",
	SavedTo:     func(path string) string { return "Įrašyta į " + path },
	ConfigLoadFailed: func(err string) string {
		return "Nepavyko perskaityti mazgų sąrašo, todėl atidarytas naujas tuščias. Failas paliktas nepakeistas:\n\n" + err
	},

	Install:   "Įdiegti",
	Update:    "Naujinti",
	Uninstall: "Pašalinti",
	InstallAsk: func(dir string) string {
		return "AltPing bus nukopijuotas į " + dir + " ir įtrauktas į meniu Pradžia, darbalaukį ir įdiegtų programų sąrašą. Jei jis ten jau įdiegtas, bus atnaujintas. Parametrai ir mazgų sąrašai liks nepakeisti.\n\nTada AltPing bus paleistas iš naujo iš ten."
	},
	InstallFailed: func(err string) string { return "Nepavyko įdiegti AltPing:\n\n" + err },
	Installed:     "AltPing įdiegtas",
	UninstallAsk: func(dir string) string {
		return "Pašalinti AltPing iš šio kompiuterio?\n\nParametrai ir mazgų sąrašai liks aplanke " + dir
	},
	UninstallRunning: "AltPing veikia. Uždarykite jį ir bandykite dar kartą.",
	Uninstalled:      "AltPing pašalintas.",
}

var lv = Strings{
	Error: "Kļūda",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d nepieejami)", n) },
	ToolConfigs:     "Konfigurācijas (O)",
	ToolAddHost:     "Pievienot resursdatoru (A)",
	ToolEditHost:    "Rediģēt resursdatoru (E)",
	ToolRemoveHosts: "Noņemt atlasītos resursdatorus (Del)",
	ToolDetails:     "Detaļas (D)",
	ToolStart:       "Sākt ping",
	ToolStop:        "Apturēt ping",
	ToolTray:        "Minimizēt paziņojumu apgabalā (T)",
	ToolShareOpen:   "Atvērt atlasītos resursdatorus u00.io (U)",
	ToolStyle:       "Stils (S)",
	KPIMedian:       "Pinga mediāna",
	KPIUptime:       "Pieejamība 24 h",
	TrayShow:        "Rādīt",
	TrayQuit:        "Iziet",
	Settings:        "Iestatījumi",
	Help:            "Palīdzība",
	About:           "Par",
	ModeUDP:         "UDP režīms",
	ModeICMP:        "ICMP režīms",
	StatusHosts:     "Resursdatori",
	StatusOK:        "Labi",
	StatusSlow:      "Lēni",
	StatusDown:      "Nepieejami",
	StatusStopped:   "Ping apturēts",

	Columns: ColumnStrings{
		Name: "Nosaukums", IP: "IP", Time: "Laiks ms", Trend: "Pēd. 5 min", Loss: "Zudumi %", Min: "Min ms", Avg: "Vid ms",
		Max: "Maks ms", Jitter: "Svārst. ms", Since: "Kopš", Details: "Detaļas",
	},
	StateOK:         "OK",
	StateTimeout:    "NOILDZE",
	StateNoResolve:  "NEVAR ATRISINĀT RESURSDATORA NOSAUKUMU",
	StatePortClosed: "PORTS AIZVĒRTS",
	NoHosts:         "Vēl nav neviena resursdatora",
	PressAToAdd:     "Nospiediet A, lai pievienotu resursdatoru",
	HowItWorks:      "Kā tas darbojas",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "d"},

	MenuEdit:       "Rediģēt... (E)",
	MenuRemove:     "Noņemt (Del)",
	MenuDowntime:   "Dīkstāves...",
	MenuExport:     "Eksportēt vēsturi...",
	MenuChange:     "Mainīt atlasītos",
	BeepOn:         "Ieslēgt skaņas signālu",
	BeepOff:        "Izslēgt skaņas signālu",
	CheckTCPPort:   "Pārbaudīt TCP portu",
	UsePing:        "Izmantot ping",
	PingEvery:      "Ping ik pēc",
	Timeout:        "Noildze",
	SlowAbove:      "Lēni virs",
	PortLabel:      "Ports:",
	PingEveryMs:    "Ping ik pēc, ms:",
	TimeoutMs:      "Noildze, ms:",
	SlowAboveMs:    "Lēni virs, ms:",
	SlowAboveMsOff: "Lēni virs, ms (0 - izsl.):",

	RemoveHostsTitle: "Resursdatoru noņemšana",
	RemoveHost:       "Vai noņemt resursdatoru?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Vai noņemt %d %s?", n, i18n.Plural("lv", n, "resursdatoru", "resursdatorus", "resursdatorus"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... un vēl %d", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Atlasiet resursdatoru",
	PingHistory: "Ping vēsture: ",
	Ms:          "ms",

	Name:          "Nosaukums:",
	Host:          "Resursdators:",
	TCPPort:       "TCP ports:",
	CheckInstead:  "Pārbaudīt ping vietā",
	Port:          "Ports:",
	SlowTooltip:   "Resursdators tiek rādīts dzeltenā krāsā, ja tā vidējais ping laiks pārsniedz šo vērtību. 0 - izslēgts",
	Notify:        "Paziņot:",
	BeepDownBack:  "Skaņas signāls, kad kļūst nepieejams vai atkal pieejams",
	NotifyTooltip: "Nepieejams nozīmē 3 neveiksmīgus ping pēc kārtas. Loga virsrakstā redzams, cik šādu resursdatoru ir nepieejami.",
	Share:         "Kopīgot:",
	ShareOn:       "Publicēt ping laiku u00.io",
	ShareTooltip:  "Ping laiki tiek nosūtīti uz publisku lapu. Ikviens, kam ir saite, var tos redzēt.",
	ShareOpen:     "Atvērt",
	ShareCopy:     "Kopēt saiti",
	LinkCopied:    "Saite nokopēta",
	MenuShareOpen: "Atvērt u00.io",
	MenuShareCopy: "Kopēt u00.io saiti",
	ShareStart:    "Publicēt u00.io",
	ShareStop:     "Pārtraukt publicēšanu",

	DowntimeTitle: "Dīkstāves",
	NoDowntime:    "Pēdējās 24 stundās dīkstāvju nav bijis",
	DownTimes: func(n int) string {
		return fmt.Sprintf("Nepieejams %d %s pēdējās 24 stundās:", n, i18n.Plural("lv", n, "reizi", "reizes", "reižu"))
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Dīkstāves pēdējās 24 stundās: %d", n) },
	Started:        "Sākums",
	Ended:          "Beigas",
	Duration:       "Ilgums",
	StillDown:      "joprojām nepieejams",
	DateTimeLayout: "02.01. 15:04:05",

	ConfigsTitle:       "Atvērt konfigurāciju",
	Configs:            "Konfigurācijas:",
	NewConfig:          "Jauna...",
	SaveConfigAs:       "Saglabāt kā...",
	RemoveConfig:       "Noņemt",
	HostsCount:         "Resursdatoru sk.",
	Hosts:              "Resursdatori",
	NewConfigTitle:     "Jauna konfigurācija",
	SaveConfigAsTitle:  "Saglabāt konfigurāciju kā",
	CopySuffix:         " kopija",
	RemoveConfigTitle:  "Konfigurācijas noņemšana",
	RemoveConfigAsk:    "Vai noņemt atlasīto konfigurāciju?",
	CannotRemoveOpened: "Pašlaik atvērto konfigurāciju nevar noņemt.",

	ExtraColumns:   "Papildu kolonnas:",
	ShowMin:        "Min",
	ShowJitter:     "Svārstības",
	ShowSince:      "Kopš (pēdējās izmaiņas)",
	AlwaysOnTop:    "Vienmēr virspusē",
	Language:       "Valoda:",
	LanguageSystem: "Kā sistēmā",
	Theme:          "Motīvs:",
	ThemeDark:      "Tumšs",
	ThemeLight:     "Gaišs",

	AboutTitle:   func(name string) string { return "Par " + name },
	Version:      "Versija",
	Author:       "Autors:",
	License:      "Licence:",
	VisitWebsite: "Apmeklēt tīmekļa vietni",
	Close:        "Aizvērt",

	ExportTitle: "Vēstures eksportēšana",
	CSVFiles:    "CSV faili",
	SavedTo:     func(path string) string { return "Saglabāts: " + path },
	ConfigLoadFailed: func(err string) string {
		return "Neizdevās nolasīt resursdatoru sarakstu, tāpēc tika atvērts jauns, tukšs saraksts. Fails netika mainīts:\n\n" + err
	},

	Install:   "Instalēt",
	Update:    "Atjaunināt",
	Uninstall: "Atinstalēt",
	InstallAsk: func(dir string) string {
		return "AltPing tiks iekopēts mapē " + dir + " un pievienots izvēlnei Sākt, darbvirsmai un instalēto programmu sarakstam. Ja tas tur jau ir instalēts, tas tiks atjaunināts. Iestatījumi un resursdatoru saraksti paliks nemainīti.\n\nPēc tam AltPing tiks restartēts no turienes."
	},
	InstallFailed: func(err string) string { return "AltPing neizdevās instalēt:\n\n" + err },
	Installed:     "AltPing ir instalēts",
	UninstallAsk: func(dir string) string {
		return "Vai noņemt AltPing no šī datora?\n\nIestatījumi un resursdatoru saraksti paliks mapē " + dir
	},
	UninstallRunning: "AltPing darbojas. Aizveriet to un mēģiniet vēlreiz.",
	Uninstalled:      "AltPing ir noņemts.",
}

var et = Strings{
	Error: "Viga",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d kättesaamatu)", n) },
	ToolConfigs:     "Konfiguratsioonid (O)",
	ToolAddHost:     "Lisa host (A)",
	ToolEditHost:    "Muuda hosti (E)",
	ToolRemoveHosts: "Eemalda valitud hostid (Del)",
	ToolDetails:     "Üksikasjad (D)",
	ToolStart:       "Alusta pingimist",
	ToolStop:        "Peata pingimine",
	ToolTray:        "Minimeeri teavitusalasse (T)",
	ToolShareOpen:   "Ava valitud hostid saidil u00.io (U)",
	ToolStyle:       "Stiil (S)",
	KPIMedian:       "Pingi mediaan",
	KPIUptime:       "Kättesaadavus 24 h",
	TrayShow:        "Kuva",
	TrayQuit:        "Välju",
	Settings:        "Sätted",
	Help:            "Spikker",
	About:           "Teave",
	ModeUDP:         "UDP-režiim",
	ModeICMP:        "ICMP-režiim",
	StatusHosts:     "Hostid",
	StatusOK:        "OK",
	StatusSlow:      "Aeglane",
	StatusDown:      "Kättesaamatu",
	StatusStopped:   "Pingimine peatatud",

	Columns: ColumnStrings{
		Name: "Nimi", IP: "IP", Time: "Aeg ms", Trend: "Viim. 5 min", Loss: "Kadu %", Min: "Min ms", Avg: "Kesk ms",
		Max: "Maks ms", Jitter: "Värin ms", Since: "Alates", Details: "Üksikasjad",
	},
	StateOK:         "OK",
	StateTimeout:    "AJALÕPP",
	StateNoResolve:  "HOSTINIME EI SAA LAHENDADA",
	StatePortClosed: "PORT SULETUD",
	NoHosts:         "Hoste veel pole",
	PressAToAdd:     "Hosti lisamiseks vajutage A",
	HowItWorks:      "Kuidas see töötab",

	Units: UnitStrings{Second: "s", Minute: "min", Hour: "h", Day: "p"},

	MenuEdit:       "Muuda... (E)",
	MenuRemove:     "Eemalda (Del)",
	MenuDowntime:   "Katkestused...",
	MenuExport:     "Ekspordi ajalugu...",
	MenuChange:     "Muuda valitud",
	BeepOn:         "Helisignaal sees",
	BeepOff:        "Helisignaal väljas",
	CheckTCPPort:   "Kontrolli TCP-porti",
	UsePing:        "Kasuta pingi",
	PingEvery:      "Pingi iga",
	Timeout:        "Ajalõpp",
	SlowAbove:      "Aeglane üle",
	PortLabel:      "Port:",
	PingEveryMs:    "Pingi iga, ms:",
	TimeoutMs:      "Ajalõpp, ms:",
	SlowAboveMs:    "Aeglane üle, ms:",
	SlowAboveMsOff: "Aeglane üle, ms (0 - väljas):",

	RemoveHostsTitle: "Hostide eemaldamine",
	RemoveHost:       "Kas eemaldada host?",
	RemoveHosts: func(n int) string {
		return fmt.Sprintf("Kas eemaldada %d %s?", n, i18n.Plural("et", n, "host", "hosti"))
	},
	AndMore: func(n int) string { return fmt.Sprintf("... ja veel %d", n) },

	Periods:     [3]string{"5 min", "1 h", "24 h"},
	SelectAHost: "Valige host",
	PingHistory: "Pingi ajalugu: ",
	Ms:          "ms",

	Name:          "Nimi:",
	Host:          "Host:",
	TCPPort:       "TCP-port:",
	CheckInstead:  "Kontrolli pingi asemel",
	Port:          "Port:",
	SlowTooltip:   "Host kuvatakse kollasena, kui selle keskmine pingiaeg on sellest suurem. 0 - väljas",
	Notify:        "Teavitamine:",
	BeepDownBack:  "Helisignaal katkestusel ja taastumisel",
	NotifyTooltip: "Kättesaamatu tähendab 3 järjestikust ebaõnnestunud pingi. Akna tiitliribal on näha, mitu sellist hosti on kättesaamatud.",
	Share:         "Jagamine:",
	ShareOn:       "Avalda pingiaeg saidil u00.io",
	ShareTooltip:  "Pingiajad saadetakse avalikule lehele. Kõik, kellel on link, näevad neid.",
	ShareOpen:     "Ava",
	ShareCopy:     "Kopeeri link",
	LinkCopied:    "Link kopeeriti",
	MenuShareOpen: "Ava saidil u00.io",
	MenuShareCopy: "Kopeeri u00.io link",
	ShareStart:    "Avalda saidil u00.io",
	ShareStop:     "Lõpeta avaldamine",

	DowntimeTitle: "Katkestused",
	NoDowntime:    "Viimase 24 tunni jooksul katkestusi polnud",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Viimase 24 tunni jooksul kättesaamatu üks kord:"
		}
		return fmt.Sprintf("Viimase 24 tunni jooksul kättesaamatu %d korda:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Katkestusi viimase 24 tunni jooksul: %d", n) },
	Started:        "Algus",
	Ended:          "Lõpp",
	Duration:       "Kestus",
	StillDown:      "endiselt kättesaamatu",
	DateTimeLayout: "02.01 15:04:05",

	ConfigsTitle:       "Ava konfiguratsioon",
	Configs:            "Konfiguratsioonid:",
	NewConfig:          "Uus...",
	SaveConfigAs:       "Salvesta nimega...",
	RemoveConfig:       "Eemalda",
	HostsCount:         "Hostide arv",
	Hosts:              "Hostid",
	NewConfigTitle:     "Uus konfiguratsioon",
	SaveConfigAsTitle:  "Konfiguratsiooni salvestamine nimega",
	CopySuffix:         " koopia",
	RemoveConfigTitle:  "Konfiguratsiooni eemaldamine",
	RemoveConfigAsk:    "Kas eemaldada valitud konfiguratsioon?",
	CannotRemoveOpened: "Praegu avatud konfiguratsiooni ei saa eemaldada.",

	ExtraColumns:   "Lisaveerud:",
	ShowMin:        "Min",
	ShowJitter:     "Värin",
	ShowSince:      "Alates (viimane muutus)",
	AlwaysOnTop:    "Alati pealmine",
	Language:       "Keel:",
	LanguageSystem: "Nagu süsteemis",
	Theme:          "Kujundus:",
	ThemeDark:      "Tume",
	ThemeLight:     "Hele",

	AboutTitle:   func(name string) string { return "Teave: " + name },
	Version:      "Versioon",
	Author:       "Autor:",
	License:      "Litsents:",
	VisitWebsite: "Külasta veebisaiti",
	Close:        "Sule",

	ExportTitle: "Ajaloo eksportimine",
	CSVFiles:    "CSV-failid",
	SavedTo:     func(path string) string { return "Salvestatud: " + path },
	ConfigLoadFailed: func(err string) string {
		return "Hostide loendit ei saanud lugeda, avati uus tühi loend. Faili ei muudetud:\n\n" + err
	},

	Install:   "Installi",
	Update:    "Värskenda",
	Uninstall: "Desinstalli",
	InstallAsk: func(dir string) string {
		return "AltPing kopeeritakse kausta " + dir + " ning lisatakse menüüsse Start, töölauale ja installitud rakenduste loendisse. Kui see on seal juba installitud, siis seda värskendatakse. Sätted ja hostide loendid jäävad muutmata.\n\nSeejärel käivitub AltPing sealt uuesti."
	},
	InstallFailed: func(err string) string { return "AltPingi ei saanud installida:\n\n" + err },
	Installed:     "AltPing on installitud",
	UninstallAsk: func(dir string) string {
		return "Kas eemaldada AltPing sellest arvutist?\n\nSätted ja hostide loendid jäävad kausta " + dir
	},
	UninstallRunning: "AltPing töötab. Sulgege see ja proovige uuesti.",
	Uninstalled:      "AltPing on eemaldatud.",
}

var tr = Strings{
	Error: "Hata",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d erişilemiyor)", n) },
	ToolConfigs:     "Yapılandırmalar (O)",
	ToolAddHost:     "Ana bilgisayar ekle (A)",
	ToolEditHost:    "Ana bilgisayarı düzenle (E)",
	ToolRemoveHosts: "Seçili ana bilgisayarları kaldır (Del)",
	ToolDetails:     "Ayrıntılar (D)",
	ToolStart:       "Ping'i başlat",
	ToolStop:        "Ping'i durdur",
	ToolTray:        "Bildirim alanına küçült (T)",
	ToolShareOpen:   "Seçili ana bilgisayarları u00.io'da aç (U)",
	ToolStyle:       "Stil (S)",
	KPIMedian:       "Medyan ping",
	KPIUptime:       "24 saat erişilebilirlik",
	TrayShow:        "Göster",
	TrayQuit:        "Çıkış",
	Settings:        "Ayarlar",
	Help:            "Yardım",
	About:           "Hakkında",
	ModeUDP:         "UDP modu",
	ModeICMP:        "ICMP modu",
	StatusHosts:     "Ana bilgisayarlar",
	StatusOK:        "Tamam",
	StatusSlow:      "Yavaş",
	StatusDown:      "Erişilemiyor",
	StatusStopped:   "Ping durduruldu",

	Columns: ColumnStrings{
		Name: "Ad", IP: "IP", Time: "Süre ms", Trend: "Son 5 dk", Loss: "Kayıp %", Min: "Min ms", Avg: "Ort ms",
		Max: "Maks ms", Jitter: "Jitter ms", Since: "Başlangıç", Details: "Ayrıntılar",
	},
	StateOK:         "OK",
	StateTimeout:    "ZAMAN AŞIMI",
	StateNoResolve:  "ANA BİLGİSAYAR ADI ÇÖZÜMLENEMİYOR",
	StatePortClosed: "BAĞLANTI NOKTASI KAPALI",
	NoHosts:         "Henüz ana bilgisayar yok",
	PressAToAdd:     "Ana bilgisayar eklemek için A tuşuna basın",
	HowItWorks:      "Nasıl çalışır",

	Units: UnitStrings{Second: "sn", Minute: "dk", Hour: "sa", Day: "g"},

	MenuEdit:       "Düzenle... (E)",
	MenuRemove:     "Kaldır (Del)",
	MenuDowntime:   "Kesintiler...",
	MenuExport:     "Geçmişi dışa aktar...",
	MenuChange:     "Seçilenleri değiştir",
	BeepOn:         "Sesli uyarı açık",
	BeepOff:        "Sesli uyarı kapalı",
	CheckTCPPort:   "TCP bağlantı noktasını denetle",
	UsePing:        "Ping kullan",
	PingEvery:      "Ping aralığı",
	Timeout:        "Zaman aşımı",
	SlowAbove:      "Yavaş eşiği",
	PortLabel:      "Bağlantı noktası:",
	PingEveryMs:    "Ping aralığı, ms:",
	TimeoutMs:      "Zaman aşımı, ms:",
	SlowAboveMs:    "Yavaş eşiği, ms:",
	SlowAboveMsOff: "Yavaş eşiği, ms (0 - kapalı):",

	RemoveHostsTitle: "Ana bilgisayarları kaldır",
	RemoveHost:       "Ana bilgisayar kaldırılsın mı?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("%d ana bilgisayar kaldırılsın mı?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... ve %d tane daha", n) },

	Periods:     [3]string{"5 dk", "1 sa", "24 sa"},
	SelectAHost: "Bir ana bilgisayar seçin",
	PingHistory: "Ping geçmişi: ",
	Ms:          "ms",

	Name:          "Ad:",
	Host:          "Ana bilgisayar:",
	TCPPort:       "TCP bağlantı noktası:",
	CheckInstead:  "Ping yerine denetle",
	Port:          "Bağlantı noktası:",
	SlowTooltip:   "Ortalama ping süresi bu değeri aştığında ana bilgisayar sarı gösterilir. 0 - kapalı",
	Notify:        "Bildir:",
	BeepDownBack:  "Erişim kesildiğinde veya geri geldiğinde sesli uyarı",
	NotifyTooltip: "Erişilemiyor, art arda 3 başarısız ping anlamına gelir. Pencere başlığı bu durumdaki ana bilgisayarların sayısını gösterir.",
	Share:         "Paylaş:",
	ShareOn:       "Ping süresini u00.io'da yayımla",
	ShareTooltip:  "Ping süreleri herkese açık bir sayfaya gönderilir. Bağlantıya sahip olan herkes bunları görebilir.",
	ShareOpen:     "Aç",
	ShareCopy:     "Bağlantıyı kopyala",
	LinkCopied:    "Bağlantı kopyalandı",
	MenuShareOpen: "u00.io'da aç",
	MenuShareCopy: "u00.io bağlantısını kopyala",
	ShareStart:    "u00.io'da yayımla",
	ShareStop:     "Yayımlamayı durdur",

	DowntimeTitle: "Kesintiler",
	NoDowntime:    "Son 24 saatte kesinti yok",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Son 24 saatte bir kez erişilemedi:"
		}
		return fmt.Sprintf("Son 24 saatte %d kez erişilemedi:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Son 24 saatteki kesintiler: %d", n) },
	Started:        "Başlangıç",
	Ended:          "Bitiş",
	Duration:       "Süre",
	StillDown:      "hâlâ erişilemiyor",
	DateTimeLayout: "02.01 15:04:05",

	ConfigsTitle:       "Yapılandırmayı aç",
	Configs:            "Yapılandırmalar:",
	NewConfig:          "Yeni...",
	SaveConfigAs:       "Farklı kaydet...",
	RemoveConfig:       "Kaldır",
	HostsCount:         "Ana bilgisayar sayısı",
	Hosts:              "Ana bilgisayarlar",
	NewConfigTitle:     "Yeni yapılandırma",
	SaveConfigAsTitle:  "Yapılandırmayı farklı kaydet",
	CopySuffix:         " kopya",
	RemoveConfigTitle:  "Yapılandırmayı kaldır",
	RemoveConfigAsk:    "Seçili yapılandırma kaldırılsın mı?",
	CannotRemoveOpened: "Şu anda açık olan yapılandırma kaldırılamaz.",

	ExtraColumns:   "Ek sütunlar:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Başlangıç (son değişiklik)",
	AlwaysOnTop:    "Her zaman üstte",
	Language:       "Dil:",
	LanguageSystem: "Sistemdeki gibi",
	Theme:          "Tema:",
	ThemeDark:      "Koyu",
	ThemeLight:     "Açık",

	AboutTitle:   func(name string) string { return name + " hakkında" },
	Version:      "Sürüm",
	Author:       "Yazar:",
	License:      "Lisans:",
	VisitWebsite: "Web sitesini ziyaret et",
	Close:        "Kapat",

	ExportTitle: "Geçmişi dışa aktar",
	CSVFiles:    "CSV dosyaları",
	SavedTo:     func(path string) string { return "Kaydedildi: " + path },
	ConfigLoadFailed: func(err string) string {
		return "Ana bilgisayar listesi okunamadı, yeni ve boş bir liste açıldı. Dosyaya dokunulmadı:\n\n" + err
	},

	Install:   "Yükle",
	Update:    "Güncelleştir",
	Uninstall: "Kaldır",
	InstallAsk: func(dir string) string {
		return "AltPing şu klasöre kopyalanacak: " + dir + ". Başlat menüsüne, masaüstüne ve yüklü uygulamalar listesine eklenecek. Zaten orada yüklüyse güncelleştirilir. Ayarlar ve ana bilgisayar listeleri olduğu gibi kalır.\n\nArdından AltPing oradan yeniden başlatılacak."
	},
	InstallFailed: func(err string) string { return "AltPing yüklenemedi:\n\n" + err },
	Installed:     "AltPing yüklendi",
	UninstallAsk: func(dir string) string {
		return "AltPing bu bilgisayardan kaldırılsın mı?\n\nAyarlar ve ana bilgisayar listeleri şu klasörde kalır: " + dir
	},
	UninstallRunning: "AltPing çalışıyor. Kapatıp yeniden deneyin.",
	Uninstalled:      "AltPing kaldırıldı.",
}

var ind = Strings{
	Error: "Kesalahan",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d mati)", n) },
	ToolConfigs:     "Konfigurasi (O)",
	ToolAddHost:     "Tambah host (A)",
	ToolEditHost:    "Edit host (E)",
	ToolRemoveHosts: "Hapus host yang dipilih (Del)",
	ToolDetails:     "Detail (D)",
	ToolStart:       "Mulai ping",
	ToolStop:        "Hentikan ping",
	ToolTray:        "Perkecil ke baki sistem (T)",
	ToolShareOpen:   "Buka host yang dipilih di u00.io (U)",
	ToolStyle:       "Gaya (S)",
	KPIMedian:       "Median ping",
	KPIUptime:       "Ketersediaan 24 jam",
	TrayShow:        "Tampilkan",
	TrayQuit:        "Keluar",
	Settings:        "Pengaturan",
	Help:            "Bantuan",
	About:           "Tentang",
	ModeUDP:         "Mode UDP",
	ModeICMP:        "Mode ICMP",
	StatusHosts:     "Host",
	StatusOK:        "OK",
	StatusSlow:      "Lambat",
	StatusDown:      "Mati",
	StatusStopped:   "Ping dihentikan",

	Columns: ColumnStrings{
		Name: "Nama", IP: "IP", Time: "Waktu ms", Trend: "5 mnt terakhir", Loss: "Hilang %", Min: "Min ms", Avg: "Rata ms",
		Max: "Maks ms", Jitter: "Jitter ms", Since: "Sejak", Details: "Detail",
	},
	StateOK:         "OK",
	StateTimeout:    "WAKTU HABIS",
	StateNoResolve:  "NAMA HOST TIDAK DAPAT DIRESOLUSI",
	StatePortClosed: "PORT TERTUTUP",
	NoHosts:         "Belum ada host",
	PressAToAdd:     "Tekan A untuk menambahkan host",
	HowItWorks:      "Cara kerjanya",

	Units: UnitStrings{Second: "d", Minute: "m", Hour: "j", Day: "h"},

	MenuEdit:       "Edit... (E)",
	MenuRemove:     "Hapus (Del)",
	MenuDowntime:   "Waktu mati...",
	MenuExport:     "Ekspor riwayat...",
	MenuChange:     "Ubah yang dipilih",
	BeepOn:         "Bunyi aktif",
	BeepOff:        "Bunyi nonaktif",
	CheckTCPPort:   "Periksa port TCP",
	UsePing:        "Gunakan ping",
	PingEvery:      "Ping setiap",
	Timeout:        "Batas waktu",
	SlowAbove:      "Lambat di atas",
	PortLabel:      "Port:",
	PingEveryMs:    "Ping setiap, ms:",
	TimeoutMs:      "Batas waktu, ms:",
	SlowAboveMs:    "Lambat di atas, ms:",
	SlowAboveMsOff: "Lambat di atas, ms (0 - nonaktif):",

	RemoveHostsTitle: "Hapus host",
	RemoveHost:       "Hapus host ini?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("Hapus %d host?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... dan %d lainnya", n) },

	Periods:     [3]string{"5m", "1j", "24j"},
	SelectAHost: "Pilih host",
	PingHistory: "Riwayat ping: ",
	Ms:          "ms",

	Name:          "Nama:",
	Host:          "Host:",
	TCPPort:       "Port TCP:",
	CheckInstead:  "Periksa alih-alih ping",
	Port:          "Port:",
	SlowTooltip:   "Tampilkan host dengan warna kuning jika rata-rata waktu ping-nya di atas nilai ini. 0 - nonaktif",
	Notify:        "Beri tahu:",
	BeepDownBack:  "Bunyi saat mati atau kembali aktif",
	NotifyTooltip: "Mati berarti 3 ping gagal berturut-turut. Judul jendela menunjukkan berapa banyak host seperti itu yang mati.",
	Share:         "Bagikan:",
	ShareOn:       "Publikasikan waktu ping di u00.io",
	ShareTooltip:  "Waktu ping dikirim ke halaman publik. Siapa pun yang memiliki tautan dapat melihatnya.",
	ShareOpen:     "Buka",
	ShareCopy:     "Salin tautan",
	LinkCopied:    "Tautan disalin",
	MenuShareOpen: "Buka di u00.io",
	MenuShareCopy: "Salin tautan u00.io",
	ShareStart:    "Publikasikan di u00.io",
	ShareStop:     "Hentikan publikasi",

	DowntimeTitle: "Waktu mati",
	NoDowntime:    "Tidak ada waktu mati dalam 24 jam terakhir",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Mati sekali dalam 24 jam terakhir:"
		}
		return fmt.Sprintf("Mati %d kali dalam 24 jam terakhir:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Gangguan dalam 24 jam terakhir: %d", n) },
	Started:        "Mulai",
	Ended:          "Selesai",
	Duration:       "Durasi",
	StillDown:      "masih mati",
	DateTimeLayout: "02/01 15:04:05",

	ConfigsTitle:       "Buka konfigurasi",
	Configs:            "Konfigurasi:",
	NewConfig:          "Baru...",
	SaveConfigAs:       "Simpan sebagai...",
	RemoveConfig:       "Hapus",
	HostsCount:         "Jumlah host",
	Hosts:              "Host",
	NewConfigTitle:     "Konfigurasi baru",
	SaveConfigAsTitle:  "Simpan konfigurasi sebagai",
	CopySuffix:         " salinan",
	RemoveConfigTitle:  "Hapus konfigurasi",
	RemoveConfigAsk:    "Hapus konfigurasi yang dipilih?",
	CannotRemoveOpened: "Konfigurasi yang sedang dibuka tidak dapat dihapus.",

	ExtraColumns:   "Kolom tambahan:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Sejak (perubahan terakhir)",
	AlwaysOnTop:    "Selalu di atas",
	Language:       "Bahasa:",
	LanguageSystem: "Sesuai sistem",
	Theme:          "Tema:",
	ThemeDark:      "Gelap",
	ThemeLight:     "Terang",

	AboutTitle:   func(name string) string { return "Tentang " + name },
	Version:      "Versi",
	Author:       "Penulis:",
	License:      "Lisensi:",
	VisitWebsite: "Kunjungi situs web",
	Close:        "Tutup",

	ExportTitle: "Ekspor riwayat",
	CSVFiles:    "File CSV",
	SavedTo:     func(path string) string { return "Disimpan ke " + path },
	ConfigLoadFailed: func(err string) string {
		return "Daftar host tidak dapat dibaca, daftar kosong baru telah dibuka. File dibiarkan apa adanya:\n\n" + err
	},

	Install:   "Instal",
	Update:    "Perbarui",
	Uninstall: "Copot instalan",
	InstallAsk: func(dir string) string {
		return "AltPing akan disalin ke " + dir + " dan ditambahkan ke menu Start, desktop, dan daftar aplikasi terinstal. Jika sudah terinstal di sana, AltPing akan diperbarui. Pengaturan dan daftar host tetap seperti semula.\n\nSetelah itu AltPing akan dimulai ulang dari sana."
	},
	InstallFailed: func(err string) string { return "AltPing tidak dapat diinstal:\n\n" + err },
	Installed:     "AltPing telah terinstal",
	UninstallAsk: func(dir string) string {
		return "Hapus AltPing dari komputer ini?\n\nPengaturan dan daftar host tetap ada di " + dir
	},
	UninstallRunning: "AltPing sedang berjalan. Tutup lalu coba lagi.",
	Uninstalled:      "AltPing telah dihapus.",
}

var vi = Strings{
	Error: "Lỗi",

	DownCount:       func(n int) string { return fmt.Sprintf("(%d mất kết nối)", n) },
	ToolConfigs:     "Cấu hình (O)",
	ToolAddHost:     "Thêm máy chủ (A)",
	ToolEditHost:    "Sửa máy chủ (E)",
	ToolRemoveHosts: "Xóa các máy chủ đã chọn (Del)",
	ToolDetails:     "Chi tiết (D)",
	ToolStart:       "Bắt đầu ping",
	ToolStop:        "Dừng ping",
	ToolTray:        "Thu nhỏ xuống khay hệ thống (T)",
	ToolShareOpen:   "Mở các máy chủ đã chọn trên u00.io (U)",
	ToolStyle:       "Kiểu giao diện (S)",
	KPIMedian:       "Ping trung vị",
	KPIUptime:       "Khả dụng 24 giờ",
	TrayShow:        "Hiện",
	TrayQuit:        "Thoát",
	Settings:        "Cài đặt",
	Help:            "Trợ giúp",
	About:           "Giới thiệu",
	ModeUDP:         "Chế độ UDP",
	ModeICMP:        "Chế độ ICMP",
	StatusHosts:     "Máy chủ",
	StatusOK:        "OK",
	StatusSlow:      "Chậm",
	StatusDown:      "Mất kết nối",
	StatusStopped:   "Đã dừng ping",

	Columns: ColumnStrings{
		Name: "Tên", IP: "IP", Time: "T.gian ms", Trend: "5 ph qua", Loss: "Mất %", Min: "Min ms", Avg: "TB ms",
		Max: "Max ms", Jitter: "Jitter ms", Since: "Từ lúc", Details: "Chi tiết",
	},
	StateOK:         "OK",
	StateTimeout:    "HẾT THỜI GIAN CHỜ",
	StateNoResolve:  "KHÔNG PHÂN GIẢI ĐƯỢC TÊN MÁY CHỦ",
	StatePortClosed: "CỔNG ĐÓNG",
	NoHosts:         "Chưa có máy chủ nào",
	PressAToAdd:     "Nhấn A để thêm máy chủ",
	HowItWorks:      "Cách hoạt động",

	Units: UnitStrings{Second: "g", Minute: "ph", Hour: "h", Day: "ng"},

	MenuEdit:       "Sửa... (E)",
	MenuRemove:     "Xóa (Del)",
	MenuDowntime:   "Thời gian mất kết nối...",
	MenuExport:     "Xuất lịch sử...",
	MenuChange:     "Thay đổi mục đã chọn",
	BeepOn:         "Bật âm báo",
	BeepOff:        "Tắt âm báo",
	CheckTCPPort:   "Kiểm tra cổng TCP",
	UsePing:        "Dùng ping",
	PingEvery:      "Ping mỗi",
	Timeout:        "Thời gian chờ",
	SlowAbove:      "Chậm khi trên",
	PortLabel:      "Cổng:",
	PingEveryMs:    "Ping mỗi, ms:",
	TimeoutMs:      "Thời gian chờ, ms:",
	SlowAboveMs:    "Chậm khi trên, ms:",
	SlowAboveMsOff: "Chậm khi trên, ms (0 - tắt):",

	RemoveHostsTitle: "Xóa máy chủ",
	RemoveHost:       "Xóa máy chủ này?",
	RemoveHosts:      func(n int) string { return fmt.Sprintf("Xóa %d máy chủ?", n) },
	AndMore:          func(n int) string { return fmt.Sprintf("... và %d mục khác", n) },

	Periods:     [3]string{"5 ph", "1 giờ", "24 giờ"},
	SelectAHost: "Chọn một máy chủ",
	PingHistory: "Lịch sử ping: ",
	Ms:          "ms",

	Name:          "Tên:",
	Host:          "Máy chủ:",
	TCPPort:       "Cổng TCP:",
	CheckInstead:  "Kiểm tra thay cho ping",
	Port:          "Cổng:",
	SlowTooltip:   "Hiển thị máy chủ màu vàng khi thời gian ping trung bình vượt quá giá trị này. 0 - tắt",
	Notify:        "Thông báo:",
	BeepDownBack:  "Phát âm báo khi mất hoặc có lại kết nối",
	NotifyTooltip: "Mất kết nối nghĩa là 3 lần ping thất bại liên tiếp. Tiêu đề cửa sổ cho biết có bao nhiêu máy chủ như vậy đang mất kết nối.",
	Share:         "Chia sẻ:",
	ShareOn:       "Công khai thời gian ping trên u00.io",
	ShareTooltip:  "Thời gian ping được gửi đến một trang công khai. Bất kỳ ai có liên kết đều có thể xem.",
	ShareOpen:     "Mở",
	ShareCopy:     "Sao chép liên kết",
	LinkCopied:    "Đã sao chép liên kết",
	MenuShareOpen: "Mở trên u00.io",
	MenuShareCopy: "Sao chép liên kết u00.io",
	ShareStart:    "Công khai trên u00.io",
	ShareStop:     "Ngừng công khai",

	DowntimeTitle: "Thời gian mất kết nối",
	NoDowntime:    "Không mất kết nối trong 24 giờ qua",
	DownTimes: func(n int) string {
		if n == 1 {
			return "Mất kết nối một lần trong 24 giờ qua:"
		}
		return fmt.Sprintf("Mất kết nối %d lần trong 24 giờ qua:", n)
	},
	GroupDownTimes: func(n int) string { return fmt.Sprintf("Số lần mất kết nối trong 24 giờ qua: %d", n) },
	Started:        "Bắt đầu",
	Ended:          "Kết thúc",
	Duration:       "Thời lượng",
	StillDown:      "vẫn mất kết nối",
	DateTimeLayout: "02/01 15:04:05",

	ConfigsTitle:       "Mở cấu hình",
	Configs:            "Cấu hình:",
	NewConfig:          "Mới...",
	SaveConfigAs:       "Lưu thành...",
	RemoveConfig:       "Xóa",
	HostsCount:         "Số máy chủ",
	Hosts:              "Máy chủ",
	NewConfigTitle:     "Cấu hình mới",
	SaveConfigAsTitle:  "Lưu cấu hình thành",
	CopySuffix:         " bản sao",
	RemoveConfigTitle:  "Xóa cấu hình",
	RemoveConfigAsk:    "Xóa cấu hình đã chọn?",
	CannotRemoveOpened: "Không thể xóa cấu hình đang mở.",

	ExtraColumns:   "Cột bổ sung:",
	ShowMin:        "Min",
	ShowJitter:     "Jitter",
	ShowSince:      "Từ lúc (thay đổi gần nhất)",
	AlwaysOnTop:    "Luôn ở trên cùng",
	Language:       "Ngôn ngữ:",
	LanguageSystem: "Theo hệ thống",
	Theme:          "Chủ đề:",
	ThemeDark:      "Tối",
	ThemeLight:     "Sáng",

	AboutTitle:   func(name string) string { return "Giới thiệu " + name },
	Version:      "Phiên bản",
	Author:       "Tác giả:",
	License:      "Giấy phép:",
	VisitWebsite: "Truy cập trang web",
	Close:        "Đóng",

	ExportTitle: "Xuất lịch sử",
	CSVFiles:    "Tệp CSV",
	SavedTo:     func(path string) string { return "Đã lưu vào " + path },
	ConfigLoadFailed: func(err string) string {
		return "Không đọc được danh sách máy chủ, đã mở một danh sách trống mới. Tệp được giữ nguyên:\n\n" + err
	},

	Install:   "Cài đặt",
	Update:    "Cập nhật",
	Uninstall: "Gỡ cài đặt",
	InstallAsk: func(dir string) string {
		return "AltPing sẽ được sao chép vào " + dir + " và thêm vào menu Start, màn hình nền và danh sách ứng dụng đã cài đặt. Nếu đã được cài đặt ở đó, AltPing sẽ được cập nhật. Cài đặt và danh sách máy chủ được giữ nguyên.\n\nSau đó AltPing sẽ khởi động lại từ vị trí đó."
	},
	InstallFailed: func(err string) string { return "Không thể cài đặt AltPing:\n\n" + err },
	Installed:     "Đã cài đặt AltPing",
	UninstallAsk: func(dir string) string {
		return "Gỡ AltPing khỏi máy tính này?\n\nCài đặt và danh sách máy chủ vẫn được giữ trong " + dir
	},
	UninstallRunning: "AltPing đang chạy. Hãy đóng ứng dụng và thử lại.",
	Uninstalled:      "Đã gỡ AltPing.",
}

var catalog = i18n.NewCatalog(en, map[string]Strings{
	"ru": ru, "pl": pl, "sr": sr, "de": de, "fr": fr, "es": es, "it": it, "pt": pt, "zh": zh, "ja": ja, "ko": ko,
	"nl": nl, "sv": sv, "da": da, "nb": nb, "fi": fi, "cs": cs, "sk": sk, "sl": sl, "hr": hr, "hu": hu,
	"ro": ro, "bg": bg, "el": el, "lt": lt, "lv": lv, "et": et, "tr": tr, "id": ind, "vi": vi,
})

// T returns the texts in the language of the application
func T() *Strings {
	return catalog.Get(ui.Language())
}

// SetLanguage switches the application to the language of the settings, "" - the system's
func SetLanguage(lang string) {
	if lang == "" {
		lang = ui.SystemLanguage()
	}
	ui.SetLanguage(lang)
}

// withoutHotkey cuts the key from a toolbar hint: "Add host (A)" is the
// title "Add host" of the dialog the button opens
func withoutHotkey(s string) string {
	if i := strings.LastIndex(s, " ("); i > 0 && strings.HasSuffix(s, ")") {
		return s[:i]
	}
	return s
}
