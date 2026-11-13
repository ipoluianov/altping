"""Visual styles for the AltPing mockups."""

BASE = dict(
    font="'Noto Sans', sans-serif", numfont="'Noto Sans', sans-serif",
    radius="12px", shadow="0 30px 80px rgba(0,0,0,.55), 0 0 0 1px rgba(255,255,255,.04)",
    winborder="transparent", btnr="8px", gap="14px", pradius="10px",
    thw="600", thsize="11px", thcase="uppercase", thls=".07em", rowh="38px", namew="500",
    dotglow="0 0 0 0", wc="rgba(255,255,255,.18)",
)


def v(**kw):
    d = dict(BASE)
    d.update(kw)
    return "".join(f"--{k}:{val};" for k, val in d.items())


THEMES = {}

# 1. Neon Midnight — deep navy, cyan→violet gradients, glowing traces
THEMES["01_neon_midnight"] = dict(
    spark="glow", status="pill", chart="glow", heat=True, spark_fill=.4, chart_sw=1.4,
    vars=v(backdrop="radial-gradient(1200px 700px at 15% 10%, #3b1d7a 0%, transparent 60%),"
                    "radial-gradient(1000px 700px at 90% 90%, #0e4d6e 0%, transparent 60%), #070a16",
           bg="#0b1022", titlebar="#0b1022", toolbar="#0b1022", statusbar="#0b1022",
           panel="#10172e", pborder="rgba(120,140,255,.10)", thead="#10172e", thtext="#6f7ba8",
           text="#e6e9ff", muted="#7d88b5", border="rgba(120,140,255,.10)", rowsep="rgba(120,140,255,.06)",
           sel="linear-gradient(90deg, rgba(34,211,238,.16), rgba(167,139,250,.10))",
           accent="#22d3ee", accent2="#a78bfa", icon="#a9b3e0", btnon="linear-gradient(135deg,#22d3ee,#a78bfa)",
           btnontext="#0b1022", input="#0e1429",
           ok="#34f5a4", slow="#fbbf24", down="#ff4d7d", warn="#ff8a4c", downrow="#5d6690",
           okbg="rgba(52,245,164,.12)", slowbg="rgba(251,191,36,.14)", downbg="rgba(255,77,125,.16)",
           dotglow="0 0 8px 1px", statbg="#0d1328", grid="rgba(120,140,255,.08)",
           chart1="#22d3ee", chart2="#a78bfa", link="#7dd3fc",
           shadow="0 40px 120px rgba(0,0,0,.7), 0 0 0 1px rgba(140,160,255,.12), 0 0 80px rgba(80,60,200,.25)"),
)

# 2. Aurora Glass — frosted translucent window over an aurora wallpaper
THEMES["02_aurora_glass"] = dict(
    spark="area", status="pill", chart="area", heat=False, spark_fill=.35, chart_fill=.5,
    vars=v(backdrop="radial-gradient(900px 600px at 10% 20%, #ff5fa2 0%, transparent 55%),"
                    "radial-gradient(900px 700px at 85% 15%, #5b8cff 0%, transparent 55%),"
                    "radial-gradient(1000px 700px at 60% 100%, #18e2b0 0%, transparent 60%), #1a1036",
           bg="rgba(16,14,40,.55)", titlebar="rgba(255,255,255,.04)", toolbar="rgba(255,255,255,.03)",
           statusbar="rgba(255,255,255,.04)", panel="rgba(255,255,255,.06)", pborder="rgba(255,255,255,.12)",
           thead="rgba(255,255,255,.04)", thtext="rgba(255,255,255,.6)",
           text="#ffffff", muted="rgba(255,255,255,.62)", border="rgba(255,255,255,.12)", rowsep="rgba(255,255,255,.06)",
           sel="rgba(255,255,255,.12)", accent="#9fd8ff", accent2="#ff9ad5", icon="rgba(255,255,255,.85)",
           btnon="rgba(255,255,255,.9)", btnontext="#2a1f5c", input="rgba(255,255,255,.08)",
           ok="#5effc4", slow="#ffd166", down="#ff6b8b", downrow="rgba(255,255,255,.4)",
           okbg="rgba(94,255,196,.16)", slowbg="rgba(255,209,102,.18)", downbg="rgba(255,107,139,.2)",
           statbg="rgba(255,255,255,.06)", grid="rgba(255,255,255,.08)", chart1="#9fd8ff", chart2="#ff9ad5",
           link="#ffffff", radius="18px", pradius="14px", btnr="10px",
           winborder="rgba(255,255,255,.22)",
           shadow="0 40px 120px rgba(10,0,40,.55), inset 0 1px 0 rgba(255,255,255,.25)"),
    css=".win{backdrop-filter:blur(40px) saturate(160%)}",
)

# 3. Phosphor Terminal — green CRT, monospace, block sparklines, scanlines
THEMES["03_phosphor_terminal"] = dict(
    spark="blocks", status="led", chart="bars", heat=False, chart_fill=.55,
    vars=v(font="'Noto Sans Mono', monospace", numfont="'Noto Sans Mono', monospace",
           backdrop="radial-gradient(circle at 50% 40%, #0a2010 0%, #020603 70%)",
           bg="#041008", titlebar="#06160b", toolbar="#041008", statusbar="#06160b",
           panel="#05130a", pborder="rgba(80,255,140,.18)", thead="#06180c", thtext="#3fbf6a",
           text="#8dffb3", muted="#3f9a5c", border="rgba(80,255,140,.18)", rowsep="rgba(80,255,140,.06)",
           sel="rgba(80,255,140,.12)", accent="#4dff88", accent2="#b6ff4d", icon="#6cff9c",
           btnon="#4dff88", btnontext="#041008", input="#06160b",
           ok="#4dff88", slow="#ffcc33", down="#ff3b3b", downrow="#2f6e44",
           okbg="transparent", slowbg="transparent", downbg="transparent",
           statbg="#06180c", grid="rgba(80,255,140,.08)", chart1="#4dff88", chart2="#b6ff4d",
           link="#b6ff4d", radius="6px", pradius="2px", btnr="2px", thw="700", namew="400",
           winborder="rgba(80,255,140,.35)",
           shadow="0 0 0 1px rgba(80,255,140,.15), 0 0 90px rgba(40,255,120,.18), 0 40px 100px rgba(0,0,0,.8)"),
    css="""
body{text-shadow:0 0 6px rgba(77,255,136,.45)}
.win::after{content:"";position:absolute;inset:0;pointer-events:none;
 background:repeating-linear-gradient(0deg, rgba(0,0,0,.18) 0 1px, transparent 1px 3px);mix-blend-mode:multiply}
.win::before{content:"";position:absolute;inset:0;pointer-events:none;z-index:2;
 background:radial-gradient(ellipse at center, transparent 60%, rgba(0,0,0,.45) 100%)}
th{border-bottom:1px dashed rgba(80,255,140,.35)}
.tb .ttl b::before{content:"> "}
.tb .ttl b::after{content:"_";animation:none;color:var(--accent)}
.stat .sv{font-size:16px}
table{font-size:12px}td,th{padding:0 7px}
""",
    title_html='<span class="ttl"><b>altping --list office</b></span>',
)

# 4. Clean Light — airy, soft shadows, status chips (Linear/Vercel-like)
THEMES["04_clean_light"] = dict(
    spark="area", status="pill", chart="area", heat=False, spark_fill=.22, chart_fill=.28, chart_sw=1.6, spark_sw=1.5,
    vars=v(backdrop="linear-gradient(135deg, #e0e7ff 0%, #fae8ff 50%, #dcfce7 100%)",
           bg="#f7f8fb", titlebar="#f7f8fb", toolbar="#f7f8fb", statusbar="#ffffff",
           panel="#ffffff", pborder="#e8eaf2", thead="#ffffff", thtext="#8a90a6",
           text="#171a2b", muted="#7b8199", border="#e8eaf2", rowsep="#f1f2f7",
           sel="#eef2ff", accent="#4f46e5", accent2="#06b6d4", icon="#4b5170",
           btnon="#4f46e5", btnontext="#ffffff", input="#ffffff",
           ok="#10b981", slow="#f59e0b", down="#ef4444", downrow="#b6bacb",
           okbg="#ecfdf5", slowbg="#fffbeb", downbg="#fef2f2",
           statbg="#f7f8fb", grid="#eef0f6", chart1="#4f46e5", chart2="#06b6d4", link="#4f46e5",
           wc="#d9dce8", radius="14px", pradius="12px",
           shadow="0 30px 80px rgba(60,60,120,.25), 0 0 0 1px rgba(0,0,0,.05)"),
    css=".left,.right,.kpi{box-shadow:0 1px 2px rgba(20,20,60,.04), 0 4px 16px rgba(20,20,60,.04)}",
)

# 5. Pastel Mocha — soft developer-favourite pastel palette on warm dark
THEMES["05_pastel_mocha"] = dict(
    spark="area", status="pill", chart="area", heat=True, spark_fill=.3, chart_fill=.4, chart_sw=1.6, spark_sw=1.5,
    vars=v(backdrop="linear-gradient(135deg, #f5c2e7 0%, #cba6f7 40%, #89b4fa 100%)",
           bg="#1e1e2e", titlebar="#181825", toolbar="#1e1e2e", statusbar="#181825",
           panel="#24273a", pborder="#313244", thead="#24273a", thtext="#9399b2",
           text="#cdd6f4", muted="#9399b2", border="#313244", rowsep="#2a2b3d",
           sel="#363a4f", accent="#cba6f7", accent2="#f5c2e7", icon="#bac2de",
           btnon="#cba6f7", btnontext="#1e1e2e", input="#181825",
           ok="#a6e3a1", slow="#f9e2af", down="#f38ba8", warn="#fab387", downrow="#6c7086",
           okbg="rgba(166,227,161,.14)", slowbg="rgba(249,226,175,.16)", downbg="rgba(243,139,168,.18)",
           statbg="#1e1e2e", grid="#313244", chart1="#89b4fa", chart2="#cba6f7", link="#89b4fa",
           radius="16px", pradius="14px", btnr="10px", thcase="none", thls="0", thsize="12px",
           shadow="0 30px 80px rgba(40,20,80,.45)"),
)

# 6. Command Center — dashboard layout: KPI cards + 24h uptime strips
THEMES["06_command_center"] = dict(
    layout="cards", spark="area", status="pill", chart="area", heat=True, spark_w=120, spark_fill=.3, chart_fill=.35,
    vars=v(backdrop="radial-gradient(1100px 600px at 80% 0%, #1f3b2a 0%, transparent 60%), #0c0f0d",
           bg="#101312", titlebar="#101312", toolbar="#101312", statusbar="#0d100f",
           panel="#161a19", pborder="#232927", thead="#161a19", thtext="#7c8783",
           text="#e8efec", muted="#7c8783", border="#232927", rowsep="#1d2220",
           sel="rgba(163,230,53,.08)", accent="#a3e635", accent2="#22d3ee", icon="#b9c4c0",
           btnon="#a3e635", btnontext="#101312", input="#131716",
           ok="#4ade80", slow="#facc15", down="#f87171", warn="#fb923c", downrow="#59625f",
           okbg="rgba(74,222,128,.12)", slowbg="rgba(250,204,21,.14)", downbg="rgba(248,113,113,.16)",
           dotglow="0 0 8px 1px", statbg="#121615", grid="#1f2523", chart1="#a3e635", chart2="#22d3ee",
           link="#a3e635", rowh="36px", radius="12px", pradius="10px",
           shadow="0 40px 100px rgba(0,0,0,.7), 0 0 0 1px rgba(255,255,255,.05)"),
    css=".kpi .kv{font-size:26px}",
)

# 7. Synthwave — sunset gradients, magenta/orange, grid horizon chart
THEMES["07_synthwave"] = dict(
    spark="glow", status="pill", chart="synth", heat=False, spark_fill=.45, chart_fill=.55, chart_sw=1.6, spark_sw=1.5,
    vars=v(backdrop="linear-gradient(180deg, #2b0b3f 0%, #5a1460 45%, #ff6a3d 100%)",
           bg="#160a26", titlebar="linear-gradient(90deg,#2a0f45,#3d0f4a)", toolbar="#160a26",
           statusbar="linear-gradient(90deg,#2a0f45,#3d0f4a)",
           panel="#1d0d33", pborder="rgba(255,90,200,.18)", thead="#1d0d33", thtext="#c48ad8",
           text="#ffe9fb", muted="#b783c9", border="rgba(255,90,200,.18)", rowsep="rgba(255,90,200,.07)",
           sel="linear-gradient(90deg, rgba(255,60,172,.22), rgba(255,140,60,.08))",
           accent="#ff3cac", accent2="#ffb340", icon="#f3b8ff",
           btnon="linear-gradient(135deg,#ff3cac,#ffb340)", btnontext="#1a0a2a", input="#1a0b2e",
           ok="#2bffd5", slow="#ffd23f", down="#ff2e63", downrow="#7a5a8c",
           okbg="rgba(43,255,213,.12)", slowbg="rgba(255,210,63,.14)", downbg="rgba(255,46,99,.18)",
           dotglow="0 0 10px 1px", statbg="#190b2c", grid="rgba(255,60,172,.22)",
           chart1="#ff3cac", chart2="#ffb340", link="#ffb340",
           shadow="0 40px 120px rgba(40,0,60,.7), 0 0 0 1px rgba(255,90,200,.25), 0 0 100px rgba(255,60,172,.25)"),
    css=".stat .sv{background:linear-gradient(90deg,#ff3cac,#ffb340);-webkit-background-clip:text;color:transparent}"
        ".rhead .t1{background:linear-gradient(90deg,#ff3cac,#ffb340);-webkit-background-clip:text;color:transparent}",
)
