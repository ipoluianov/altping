package forms

import (
	"fmt"

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
	TrayShow        string
	TrayQuit        string
	Settings        string
	Help            string
	About           string
	ModeUDP         string
	ModeICMP        string

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
}

// ColumnStrings are the names of the table columns
type ColumnStrings struct {
	Name    string
	IP      string
	Time    string
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
	TrayShow:        "Show",
	TrayQuit:        "Quit",
	Settings:        "Settings",
	Help:            "Help",
	About:           "About",
	ModeUDP:         "UDP Mode",
	ModeICMP:        "ICMP Mode",

	Columns: ColumnStrings{
		Name: "Name", IP: "IP", Time: "Time ms", Loss: "Loss %", Min: "Min ms", Avg: "Avg ms",
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
	TrayShow:        "Показать",
	TrayQuit:        "Выход",
	Settings:        "Настройки",
	Help:            "Справка",
	About:           "О программе",
	ModeUDP:         "Режим UDP",
	ModeICMP:        "Режим ICMP",

	Columns: ColumnStrings{
		Name: "Имя", IP: "IP", Time: "Время мс", Loss: "Потери %", Min: "Мин мс", Avg: "Сред мс",
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
	TrayShow:        "Pokaż",
	TrayQuit:        "Zakończ",
	Settings:        "Ustawienia",
	Help:            "Pomoc",
	About:           "O programie",
	ModeUDP:         "Tryb UDP",
	ModeICMP:        "Tryb ICMP",

	Columns: ColumnStrings{
		Name: "Nazwa", IP: "IP", Time: "Czas ms", Loss: "Straty %", Min: "Min ms", Avg: "Śr. ms",
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
	TrayShow:        "Прикажи",
	TrayQuit:        "Изађи",
	Settings:        "Подешавања",
	Help:            "Помоћ",
	About:           "О програму",
	ModeUDP:         "UDP режим",
	ModeICMP:        "ICMP режим",

	Columns: ColumnStrings{
		Name: "Име", IP: "IP", Time: "Време мс", Loss: "Губитак %", Min: "Мин мс", Avg: "Прос мс",
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
	TrayShow:        "Anzeigen",
	TrayQuit:        "Beenden",
	Settings:        "Einstellungen",
	Help:            "Hilfe",
	About:           "Über",
	ModeUDP:         "UDP-Modus",
	ModeICMP:        "ICMP-Modus",

	Columns: ColumnStrings{
		Name: "Name", IP: "IP", Time: "Zeit ms", Loss: "Verlust %", Min: "Min ms", Avg: "Mittel ms",
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
	TrayShow:        "Afficher",
	TrayQuit:        "Quitter",
	Settings:        "Paramètres",
	Help:            "Aide",
	About:           "À propos",
	ModeUDP:         "Mode UDP",
	ModeICMP:        "Mode ICMP",

	Columns: ColumnStrings{
		Name: "Nom", IP: "IP", Time: "Temps ms", Loss: "Perte %", Min: "Min ms", Avg: "Moy ms",
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
	TrayShow:        "Mostrar",
	TrayQuit:        "Salir",
	Settings:        "Ajustes",
	Help:            "Ayuda",
	About:           "Acerca de",
	ModeUDP:         "Modo UDP",
	ModeICMP:        "Modo ICMP",

	Columns: ColumnStrings{
		Name: "Nombre", IP: "IP", Time: "Tiempo ms", Loss: "Pérdida %", Min: "Mín ms", Avg: "Media ms",
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
	TrayShow:        "Mostra",
	TrayQuit:        "Esci",
	Settings:        "Impostazioni",
	Help:            "Aiuto",
	About:           "Informazioni",
	ModeUDP:         "Modalità UDP",
	ModeICMP:        "Modalità ICMP",

	Columns: ColumnStrings{
		Name: "Nome", IP: "IP", Time: "Tempo ms", Loss: "Perdita %", Min: "Min ms", Avg: "Media ms",
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
	TrayShow:        "Mostrar",
	TrayQuit:        "Sair",
	Settings:        "Configurações",
	Help:            "Ajuda",
	About:           "Sobre",
	ModeUDP:         "Modo UDP",
	ModeICMP:        "Modo ICMP",

	Columns: ColumnStrings{
		Name: "Nome", IP: "IP", Time: "Tempo ms", Loss: "Perda %", Min: "Mín ms", Avg: "Média ms",
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
	TrayShow:        "显示",
	TrayQuit:        "退出",
	Settings:        "设置",
	Help:            "帮助",
	About:           "关于",
	ModeUDP:         "UDP 模式",
	ModeICMP:        "ICMP 模式",

	Columns: ColumnStrings{
		Name: "名称", IP: "IP", Time: "时间 ms", Loss: "丢包 %", Min: "最小 ms", Avg: "平均 ms",
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
	TrayShow:        "表示",
	TrayQuit:        "終了",
	Settings:        "設定",
	Help:            "ヘルプ",
	About:           "バージョン情報",
	ModeUDP:         "UDP モード",
	ModeICMP:        "ICMP モード",

	Columns: ColumnStrings{
		Name: "名前", IP: "IP", Time: "時間 ms", Loss: "損失 %", Min: "最小 ms", Avg: "平均 ms",
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
	TrayShow:        "표시",
	TrayQuit:        "종료",
	Settings:        "설정",
	Help:            "도움말",
	About:           "정보",
	ModeUDP:         "UDP 모드",
	ModeICMP:        "ICMP 모드",

	Columns: ColumnStrings{
		Name: "이름", IP: "IP", Time: "시간 ms", Loss: "손실 %", Min: "최소 ms", Avg: "평균 ms",
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
}

var catalog = i18n.NewCatalog(en, map[string]Strings{
	"ru": ru, "pl": pl, "sr": sr, "de": de, "fr": fr, "es": es, "it": it, "pt": pt, "zh": zh, "ja": ja, "ko": ko,
})

// T returns the texts in the language of the application
func T() *Strings {
	return catalog.Get(ui.Language())
}

// The library translates its own buttons (OK, Cancel...) to Russian and Chinese only
func init() {
	for lang, s := range map[string]ui.UIStrings{
		"pl": {OK: "OK", Cancel: "Anuluj", Yes: "Tak", No: "Nie"},
		"sr": {OK: "У реду", Cancel: "Откажи", Yes: "Да", No: "Не"},
		"de": {OK: "OK", Cancel: "Abbrechen", Yes: "Ja", No: "Nein"},
		"fr": {OK: "OK", Cancel: "Annuler", Yes: "Oui", No: "Non"},
		"es": {OK: "Aceptar", Cancel: "Cancelar", Yes: "Sí", No: "No"},
		"it": {OK: "OK", Cancel: "Annulla", Yes: "Sì", No: "No"},
		"pt": {OK: "OK", Cancel: "Cancelar", Yes: "Sim", No: "Não"},
		"ja": {OK: "OK", Cancel: "キャンセル", Yes: "はい", No: "いいえ"},
		"ko": {OK: "확인", Cancel: "취소", Yes: "예", No: "아니요"},
	} {
		ui.RegisterUIStrings(lang, s)
	}
}

// SetLanguage switches the application to the language of the settings, "" - the system's
func SetLanguage(lang string) {
	if lang == "" {
		lang = ui.SystemLanguage()
	}
	ui.SetLanguage(lang)
}
