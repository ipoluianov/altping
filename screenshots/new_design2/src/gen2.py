#!/usr/bin/env python3
"""Second round of AltPing mockups: new layouts (tiles, lanes, sidebar) and
styles aimed at different audiences. Reuses the helpers of gen.py."""
import math, os, random, sys

OUT = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, OUT)
import gen
from gen import HOSTS, SEL, series, spark_path, icon, LOGO, BASE_CSS, chart, table, kpis, COLS_CLASSIC, COLS_CARDS
from themes import v

COLS_COMPACT = [("Name", "l name", "name"), ("IP", "l ip", "ip"), ("Ping", "r", "time"),
                ("Last 5 min", "l spk", "spark"), ("Loss", "r", "loss"), ("Avg", "r", "avg"),
                ("Since", "r", "since"), ("Status", "l det", "det")]

EXTRA_CSS = """
.tiles{flex:1;display:grid;gap:var(--gap);grid-auto-rows:1fr;min-width:0}
.tile{background:var(--panel);border:1px solid var(--pborder);border-radius:var(--pradius);padding:12px 14px;display:flex;flex-direction:column;min-width:0;position:relative;overflow:hidden}
.tile.sel{border-color:var(--accent);box-shadow:0 0 0 1px var(--accent)}
.tile .t1{display:flex;align-items:center;font-size:14px;font-weight:600;white-space:nowrap}
.tile .t1 .nm{overflow:hidden;text-overflow:ellipsis}
.tile .tip{font-size:11px;color:var(--muted);font-family:var(--numfont);margin:1px 0 0 18px}
.tile .tv{display:flex;align-items:baseline;gap:4px;margin-top:auto}
.tile .tv b{font-size:var(--tvsize,30px);font-weight:800;font-family:var(--numfont);letter-spacing:-.02em;line-height:1.1}
.tile .tv small{font-size:12px;color:var(--muted)}
.tile .tv .pill{margin-left:auto;align-self:center}
.tile.slow .tv b{color:var(--slow)}.tile.down .tv b{color:var(--down)}
.tile .tsp{margin-top:4px}
.tile .tsp svg{display:block;width:100%}
.tile .tf{display:flex;justify-content:space-between;font-size:11px;color:var(--muted);margin-top:6px;font-family:var(--numfont)}
.sig{margin-left:auto;display:inline-flex;gap:2px;align-items:flex-end;height:14px}
.sig i{width:3px;background:var(--muted);opacity:.3;border-radius:1px}
.sig i.on{background:var(--ok);opacity:1}
.tile.slow .sig i.on{background:var(--slow)}
.face{font-family:'DejaVu Sans',sans-serif;font-size:15px;margin-right:8px}
.lanes{display:flex;flex-direction:column}
.lane{flex:1;display:grid;grid-template-columns:230px 1fr 170px;align-items:center;border-bottom:1px solid var(--rowsep);padding:0 14px;min-height:0}
.lane.sel{background:var(--sel);box-shadow:inset 3px 0 0 var(--accent)}
.lane .ln{font-size:14px;font-weight:600;white-space:nowrap}
.lane .li{font-size:11px;color:var(--muted);font-family:var(--numfont)}
.lane>svg{display:block;width:100%}
.lane .ln svg{vertical-align:middle}
.lane .lv{text-align:right;font-family:var(--numfont);color:var(--ok);line-height:1}
.lane .lv b{font-size:34px;font-weight:700}
.lane .lv small{font-size:12px;margin-left:4px;opacity:.7}
.lane .lv div{font-size:10px;opacity:.7;margin-top:3px;letter-spacing:.06em}
.lane.slow .lv{color:var(--slow)}.lane.down .lv{color:var(--down)}
.side{width:196px;background:var(--side);border-right:1px solid var(--border);padding:10px 8px;font-size:13px;display:flex;flex-direction:column;gap:2px}
.side h4{font-size:11px;color:var(--muted);font-weight:600;padding:10px 8px 4px}
.side div{display:flex;align-items:center;gap:8px;padding:6px 8px;border-radius:6px}
.side div.on{background:var(--sideon);color:var(--sideontext);font-weight:600}
.side div span{margin-left:auto;color:var(--muted);font-size:12px}
.side div.on span{color:inherit;opacity:.8}
.side div i{width:8px;height:8px;border-radius:50%;display:block}
.body{flex:1;display:flex;min-height:0}
.body .col{flex:1;display:flex;flex-direction:column;min-width:0}
.tb.mac .wc{margin:0;order:-1}
.tb.mac .wc i:nth-child(1){background:#ff5f57}.tb.mac .wc i:nth-child(2){background:#febc2e}.tb.mac .wc i:nth-child(3){background:#28c840}
"""


def spark2(t, idx, host, h=40):
    w = 300
    state, base = host[9], host[11]
    col = {"ok": "var(--accent)", "slow": "var(--slow)", "down": "var(--down)"}[state]
    if t.get("spark_state_color"):
        col = {"ok": "var(--ok)", "slow": "var(--slow)", "down": "var(--down)"}[state]
    out = [f'<svg viewBox="0 0 {w} {h}" height="{h}" preserveAspectRatio="none">']
    if state == "down":
        out.append(f'<line x1="0" x2="{w}" y1="{h-4}" y2="{h-4}" stroke="{col}" stroke-width="2" vector-effect="non-scaling-stroke" stroke-dasharray="6 5"/></svg>')
        return "".join(out)
    vals = series(idx * 7 + 3, base, spikes=2 + idx % 3)
    gid = f"s{idx}"
    out.append(f'<defs><linearGradient id="{gid}" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="{col}" '
               f'stop-opacity="{t.get("spark_fill", .35)}"/><stop offset="1" stop-color="{col}" stop-opacity="0"/></linearGradient></defs>')
    kind = t.get("spark", "area")
    if kind == "bars":
        lo, hi = min(vals), sorted(vals)[int(len(vals) * .95)]
        bw = w / len(vals)
        for i, x in enumerate(vals):
            bh = max(2, (h - 2) * min(1, (x - lo) / (hi - lo + 1e-9)) * .9 + 3)
            out.append(f'<rect x="{i*bw:.1f}" y="{h-bh:.1f}" width="{bw-1.5:.1f}" height="{bh:.1f}" fill="{col}"/>')
    else:
        if kind in ("area", "glow"):
            out.append(f'<path d="{spark_path(vals, w, h, close=True)}" fill="url(#{gid})"/>')
        cls = ' class="glow"' if kind == "glow" else ""
        out.append(f'<path{cls} d="{spark_path(vals, w, h)}" fill="none" stroke="{col}" stroke-width="{t.get("spark_sw", 1.5)}" '
                   f'vector-effect="non-scaling-stroke" stroke-linejoin="round"/>')
    out.append("</svg>")
    return "".join(out)


def sig(host):
    if host[9] == "down":
        n = 0
    else:
        b = host[11]
        n = 4 if b < 20 else 3 if b < 60 else 2 if b < 150 else 1
    return '<span class="sig">' + "".join(
        f'<i class="{"on" if i < n else ""}" style="height:{5+i*3}px"></i>' for i in range(4)) + "</span>"


FACES = {"ok": "(◕‿◕)", "slow": "(◔_◔)", "down": "(×_×)"}


def tiles(t):
    n = t.get("tile_cols", 3)
    h = t.get("spark_h", 40)
    out = [f'<div class="tiles" style="grid-template-columns:repeat({n},1fr)">']
    for i, hst in enumerate(HOSTS):
        st = hst[9]
        lead = f'<span class="face">{FACES[st]}</span>' if t.get("faces") else f'<span class="dot {st}"></span>'
        bell = icon("bell", 13, "bell") if hst[10] else ""
        val = (f'<b>{hst[2]}</b><small>ms</small>' if st != "down" else f'<b>{t.get("down_text", "TIMEOUT")}</b>')
        pill = gen.status_cell(t, hst) if t.get("tile_pill", True) else ""
        loss = hst[3]
        out.append(f'<div class="tile {st} {"sel" if i == SEL else ""}">'
                   f'<div class="t1">{lead}<span class="nm">{hst[0]}</span>&nbsp;{bell}{sig(hst)}</div>'
                   f'<div class="tip">{hst[1]}</div>'
                   f'<div class="tv">{val}{pill}</div>'
                   f'<div class="tsp">{spark2(t, i, hst, h)}</div>'
                   f'<div class="tf"><span>loss {loss}%</span><span>jitter {hst[7]}</span><span>{hst[8]}</span></div></div>')
    out.append("</div>")
    return "".join(out)


def ecg(idx, host, w=600, h=44):
    st = host[9]
    col = {"ok": "var(--ok)", "slow": "var(--slow)", "down": "var(--down)"}[st]
    mid = h * 0.62
    if st == "down":
        d = f"M0,{mid} L{w},{mid}"
    else:
        r = random.Random(idx * 5 + 2)
        period = 34 + math.log10(host[11] + 1) * 22
        x = r.random() * period * .5
        d = f"M0,{mid}"
        while x < w:
            a = (h * .55) * (0.8 + r.random() * .4)
            if r.random() < .08:
                a *= .45
            p = period * (0.92 + r.random() * .16)
            d += (f" L{x:.1f},{mid} q3,-5 6,0 l4,0 l2,3 l3,{-a:.1f} l3,{a+6:.1f} l2,-9 l5,0 q5,-9 10,0")
            x += p
        d += f" L{w},{mid}"
    return (f'<svg viewBox="0 0 {w} {h}" height="{h}" preserveAspectRatio="none"><defs>'
            f'<linearGradient id="e{idx}" x1="0" y1="0" x2="1" y2="0"><stop offset="0" stop-color="{col}" stop-opacity=".15"/>'
            f'<stop offset=".85" stop-color="{col}" stop-opacity="1"/><stop offset="1" stop-color="#fff"/></linearGradient></defs>'
            f'<path d="{d}" fill="none" stroke="{col if st == "down" else f"url(#e{idx})"}" stroke-width="1.8" vector-effect="non-scaling-stroke" '
            f'stroke-linejoin="round" style="filter:drop-shadow(0 0 3px {col})"/></svg>')


def lanes(t):
    out = ['<div class="left lanes">']
    for i, h in enumerate(HOSTS):
        st = h[9]
        bell = icon("bell", 13, "bell") if h[10] else ""
        if st == "down":
            val = "<b>- - -</b><div>NO RESPONSE · 14m 39s</div>"
        else:
            val = f"<b>{h[2]}</b><small>ms</small><div>LOSS {h[3]}% · JIT {h[7]}</div>"
        out.append(f'<div class="lane {st} {"sel" if i == SEL else ""}"><div><div class="ln">{h[0]} {bell}</div>'
                   f'<div class="li">{h[1]}</div></div>{ecg(i, h)}<div class="lv">{val}</div></div>')
    out.append("</div>")
    return "".join(out)


def sidebar():
    def it(name, n, on=False, dot=None):
        d = f'<i style="background:{dot}"></i>' if dot else ""
        return f'<div class="{"on" if on else ""}">{d}{name}<span>{n}</span></div>'
    return ('<div class="side"><h4>Host lists</h4>' + it("Office", 12, True) + it("Servers", 8) + it("Home lab", 5)
            + it("Clients", 23) + "<h4>Smart filters</h4>" + it("Down", 1, dot="var(--down)")
            + it("Slow", 1, dot="var(--slow)") + it("With alerts", 3, dot="var(--accent)")
            + it("Lossy", 2, dot="var(--muted)") + "</div>")


def page(t):
    layout = t.get("layout", "classic")
    cols = {"classic": COLS_CLASSIC, "cards": COLS_CARDS, "compact": COLS_COMPACT}[
        t.get("cols", "cards" if layout == "cards" else "classic")]
    content_h = 880 - 38 - 54 - 38 - 2
    if layout == "cards":
        content_h -= 96 + 14
    ch = content_h - 2 * 14 - 46 - 70 - 10 - t.get("chart_dh", 0)
    if layout == "tiles":
        left = tiles(t)
    elif layout == "lanes":
        left = lanes(t)
    else:
        left = f'<div class="left">{table(t, cols)}</div>'
    right = "" if t.get("noright") else f"""<div class="right">
   <div class="rhead"><div><div class="t1">Backup, Singapore</div><div class="t2">198.51.100.7 · ICMP · every 1 s</div></div>
    <div class="seg"><span>5m</span><span class="on">1h</span><span>24h</span></div></div>
   <div class="stats"><div class="stat"><div class="sl">Now</div><div class="sv">207<small>ms</small></div></div>
    <div class="stat"><div class="sl">Avg</div><div class="sv">182<small>ms</small></div></div>
    <div class="stat"><div class="sl">Loss</div><div class="sv" style="color:var(--slow)">1.0<small>%</small></div></div>
    <div class="stat"><div class="sl">Jitter</div><div class="sv">6.9<small>ms</small></div></div></div>
   <div style="padding:0 10px">{chart(t, 410, ch)}</div></div>"""
    tools = "".join(f'<div class="tbtn">{icon(n)}</div>' for n in ("add", "edit", "del"))
    rtools = (f'<div class="tbtn on">{icon("chart")}</div><div class="tbtn dis">{icon("share")}</div>'
              f'<div class="tbtn dis">{icon("play")}</div><div class="tbtn">{icon("stop")}</div><div class="sepv"></div>'
              f'<div class="tbtn">{icon("open")}</div><div class="tbtn">{icon("save")}</div>')
    title = t.get("title_html", '<span class="ttl"><b>Alt Ping</b> · Office</span>')
    mac = t.get("mac")
    tb = (f'<div class="tb {"mac" if mac else ""}">{"" if mac else LOGO}{title}<div class="wc"><i></i><i></i><i></i></div></div>')
    toolbar = f'<div class="tools">{tools}<div class="sp"></div><div class="search">{icon("search",16)}Filter hosts…</div>{rtools}</div>'
    main = f'<div class="main">{left}{right}</div>'
    if layout == "sidebar":
        body = f'<div class="body">{sidebar()}<div class="col">{toolbar}{main}</div></div>'
    else:
        body = toolbar + (kpis(t) if layout == "cards" else "") + main
    return f"""<!doctype html><html><head><meta charset="utf-8"><style>
:root{{{t['vars']}}}
{BASE_CSS}{EXTRA_CSS}
{t.get('css','')}
</style></head><body>
<div class="win">{tb}{body}
 <div class="status"><span class="c"><i style="background:var(--ok)"></i><b>10</b> up</span>
  <span class="c"><i style="background:var(--slow)"></i><b>1</b> slow</span>
  <span class="c"><i style="background:var(--down)"></i><b>1</b> down</span>
  <span>12 hosts · UDP mode</span><span class="sp"></span><a>Settings</a><a>Help</a><a>About</a></div>
 {t.get('overlay','')}
</div></body></html>"""


T = {}
MONO = "'Noto Sans Mono', monospace"
BEV = "border:2px solid;border-color:#fff #808080 #808080 #fff"
SUNK = "border:2px solid;border-color:#808080 #fff #fff #808080"

# 01 Gamer HUD — tiles with cut corners, acid yellow + cyan, signal bars
T["01_gamer_hud"] = dict(
    layout="tiles", spark="glow", status="pill", chart="glow", spark_fill=.3, spark_h=38, tile_pill=False,
    vars=v(backdrop="linear-gradient(135deg,#0a0a0f 0%,#16111f 60%,#1b0f2e 100%)",
           bg="#0c0c12", titlebar="#0c0c12", toolbar="#0c0c12", statusbar="#0c0c12",
           panel="#14141d", pborder="#26263a", thead="#14141d", thtext="#8a8aa8",
           text="#f2f2ff", muted="#7f7f9c", border="#26263a", rowsep="#1c1c2a", sel="rgba(247,255,0,.08)",
           accent="#00f0ff", accent2="#f7ff00", icon="#c9c9e8", btnon="#f7ff00", btnontext="#0c0c12", input="#101018",
           ok="#39ff88", slow="#ffb300", down="#ff2a6d", downrow="#55556e",
           okbg="rgba(57,255,136,.12)", slowbg="rgba(255,179,0,.14)", downbg="rgba(255,42,109,.16)",
           dotglow="0 0 8px 1px", statbg="#101018", grid="#1f1f30", chart1="#00f0ff", chart2="#f7ff00",
           link="#f7ff00", radius="4px", pradius="0px", btnr="0px",
           shadow="0 40px 100px rgba(0,0,0,.8), 0 0 0 1px #f7ff00, 0 0 60px rgba(247,255,0,.12)"),
    css="""
.tile,.right{clip-path:polygon(0 0,calc(100% - 16px) 0,100% 16px,100% 100%,16px 100%,0 calc(100% - 16px))}
.tile{border-left:3px solid var(--accent)}
.tile.slow{border-left-color:var(--slow)}.tile.down{border-left-color:var(--down);background:linear-gradient(135deg,rgba(255,42,109,.16),#14141d 60%)}
.tile.sel{background:linear-gradient(135deg,rgba(247,255,0,.12),#14141d 55%);border-color:#f7ff00;box-shadow:none}
.tile .t1{text-transform:uppercase;letter-spacing:.06em;font-size:13px}
.tile .tv b{font-style:italic;font-size:34px}
.tb .ttl b{letter-spacing:.2em;text-transform:uppercase;color:#f7ff00}
.rhead .t1{text-transform:uppercase;letter-spacing:.06em}
.tbtn.on,.seg span.on{clip-path:polygon(0 0,calc(100% - 6px) 0,100% 6px,100% 100%,6px 100%,0 calc(100% - 6px))}
""")

# 02 Heartbeat monitor — every host is an ECG lane with a big vital number
T["02_heartbeat_monitor"] = dict(
    layout="lanes", spark="glow", status="pill", chart="glow", chart_sw=1.5,
    vars=v(backdrop="radial-gradient(900px 600px at 50% 0%,#0b2a2a 0%,transparent 70%),#020607",
           bg="#04090a", titlebar="#04090a", toolbar="#04090a", statusbar="#04090a",
           panel="#060d0e", pborder="#12302f", thead="#060d0e", thtext="#4d8b88",
           text="#d8fff9", muted="#4d8b88", border="#12302f", rowsep="#0e2423", sel="rgba(0,255,200,.07)",
           accent="#00ffc8", accent2="#00b3ff", icon="#8fd9d0", btnon="#00ffc8", btnontext="#04090a", input="#071112",
           ok="#19ff9c", slow="#ffd400", down="#ff3355", downrow="#3a5a58",
           okbg="rgba(25,255,156,.12)", slowbg="rgba(255,212,0,.14)", downbg="rgba(255,51,85,.16)",
           statbg="#071112", grid="#0f2a29", chart1="#19ff9c", chart2="#00b3ff", link="#00ffc8",
           numfont=MONO, radius="10px", pradius="8px",
           shadow="0 40px 100px rgba(0,0,0,.8), 0 0 0 1px #12302f, 0 0 80px rgba(0,255,200,.10)"),
    css="""
.lanes{background-image:linear-gradient(rgba(0,255,200,.05) 1px,transparent 1px),linear-gradient(90deg,rgba(0,255,200,.05) 1px,transparent 1px);background-size:20px 20px}
.lane.down{background:rgba(255,51,85,.08)}
.lane .lv b{text-shadow:0 0 12px currentColor}
""")

# 03 NOC wall — big solid status tiles, readable from across the room
T["03_noc_wall"] = dict(
    layout="tiles", noright=True, tile_cols=4, spark="area", status="pill", chart="area", spark_h=64,
    spark_state_color=True, spark_fill=.45, tile_pill=False, down_text="DOWN",
    vars=v(backdrop="#05070a", bg="#0a0d12", titlebar="#0a0d12", toolbar="#0a0d12", statusbar="#0a0d12",
           panel="#0f2a1c", pborder="#17422b", thead="#111", thtext="#888",
           text="#f4fff8", muted="rgba(255,255,255,.6)", border="#1b222c", rowsep="#1b222c", sel="transparent",
           accent="#ffffff", accent2="#4ade80", icon="#c7d0db", btnon="#ffffff", btnontext="#0a0d12", input="#10141b",
           ok="#4ade80", slow="#fbbf24", down="#ffffff", downrow="#fff",
           okbg="transparent", slowbg="transparent", downbg="transparent",
           statbg="#111", grid="#222", chart1="#4ade80", chart2="#4ade80", link="#4ade80",
           radius="10px", pradius="10px", gap="12px", tvsize="56px",
           shadow="0 40px 100px rgba(0,0,0,.8), 0 0 0 1px rgba(255,255,255,.08)"),
    css="""
.tile{padding:16px 20px}
.tile .t1{font-size:19px;font-weight:700}
.tile .tip{font-size:13px}.tile .tf{font-size:13px}
.tile.slow{background:#3b2c05;border-color:#7a5b0a}
.tile.down{background:linear-gradient(135deg,#e11d48,#b91c1c);border-color:#ff6b81;box-shadow:0 0 40px rgba(225,29,72,.45)}
.tile.down .dot{background:#fff;box-shadow:none}
.tile.sel{box-shadow:0 0 0 2px #fff}
.sig{height:18px}.sig i{width:4px}
""")

# 04 Desktop native light — sidebar with host lists, calm system look
T["04_desktop_native"] = dict(
    layout="sidebar", cols="compact", mac=True, spark="area", status="pill", chart="area", spark_w=120,
    spark_fill=.2, chart_fill=.25, chart_sw=1.5,
    vars=v(backdrop="linear-gradient(160deg,#7fb2e5 0%,#b9a3e3 50%,#f0b3a8 100%)",
           bg="#ffffff", titlebar="#ecebee", toolbar="#ffffff", statusbar="#f6f5f7",
           panel="#ffffff", pborder="#e4e3e8", thead="#fafafb", thtext="#86868b",
           text="#1d1d1f", muted="#86868b", border="#e4e3e8", rowsep="#f0eff3", sel="#e3effd",
           accent="#0a7aff", accent2="#5ac8fa", icon="#515154", btnon="#e9e9ed", btnontext="#0a7aff", input="#f1f0f4",
           ok="#28a745", slow="#f29d00", down="#ff3b30", downrow="#b7b7bd",
           okbg="#e9f8ee", slowbg="#fff5df", downbg="#ffebea",
           statbg="#f7f7f9", grid="#efeff3", chart1="#0a7aff", chart2="#5ac8fa", link="#0a7aff",
           thcase="none", thls="0", thsize="12px", thw="500", wc="#ccc", side="rgba(236,235,240,.92)",
           sideon="#0a7aff", sideontext="#fff", rowh="36px",
           shadow="0 30px 80px rgba(30,30,80,.35), 0 0 0 1px rgba(0,0,0,.12)"),
    css=".tools{height:50px}")

# 05 Retro 95 — gray bevels, navy title bar, perfmon-style chart
T["05_retro_95"] = dict(
    spark="line", status="stxt", chart="line", chart_fill=0, spark_sw=1.2, chart_sw=1.2,
    vars=v(backdrop="#008080", bg="#c0c0c0", titlebar="linear-gradient(90deg,#000080,#1084d0)", toolbar="#c0c0c0",
           statusbar="#c0c0c0", panel="#ffffff", pborder="#808080", thead="#c0c0c0", thtext="#000",
           text="#000000", muted="#404040", border="#808080", rowsep="#e4e4e4", sel="#000080",
           accent="#000080", accent2="#1084d0", icon="#000", btnon="#000080", btnontext="#fff", input="#fff",
           ok="#008000", slow="#b36b00", down="#d00000", downrow="#808080",
           okbg="transparent", slowbg="transparent", downbg="transparent",
           statbg="#c0c0c0", grid="#005000", chart1="#00ff00", chart2="#00ff00", link="#000080",
           radius="0px", pradius="0px", btnr="0px", thcase="none", thls="0", thsize="12px", thw="400",
           rowh="30px", wc="#c0c0c0", shadow="2px 2px 0 #000", winborder="#dfdfdf"),
    css=f"""
.win{{border:2px solid;border-color:#fff #404040 #404040 #fff;padding:3px}}
.tb{{height:28px;color:#fff;font-weight:700}}.tb .ttl{{text-align:left;left:40px}}.tb .ttl b{{color:#fff}}
.tb .wc{{gap:2px}}.tb .wc i{{border-radius:0;width:20px;height:18px;{BEV}}}
.tools{{border-bottom:1px solid #808080;box-shadow:0 1px 0 #fff}}
.tbtn{{{BEV};width:36px;height:36px;margin-right:2px}}
.tbtn.on{{{SUNK};background:#dcdcdc;color:#000}}
.left,.right,.search{{{SUNK}}}
.right{{background:#c0c0c0}}
th{{{BEV};height:26px}}
td{{font-size:12px}}
tbody tr.sel td,tbody tr.sel td.ip,tbody tr.sel .stxt,tbody tr.sel .lossv{{color:#fff}}
tbody tr.sel td:first-child{{box-shadow:none;outline:1px dotted #ff0;outline-offset:-2px}}
tbody tr.sel svg path{{stroke:#fff}}
.dot{{border-radius:0}}
.stat{{{SUNK};background:#c0c0c0}}.stat .sl{{text-transform:none;letter-spacing:0;color:#000}}
.seg{{border:none;background:none;gap:2px}}.seg span{{{BEV};color:#000;border-radius:0}}.seg span.on{{{SUNK};background:#dcdcdc;color:#000}}
.chart{{background:#000;{SUNK}}}.chart .ax{{fill:#00c000}}.chart .avgline{{stroke:#ff0}}.chart .avglbl{{fill:#ff0}}
.chart .lastpt{{stroke:#000}}
.status{{border-top:none}}.status .c,.status>span:not(.sp){{{SUNK};padding:2px 10px;color:#000;border-width:1px}}
.status a{{text-decoration:underline}}
""",
    title_html='<span class="ttl"><b>Alt Ping - Office</b></span>')

# 06 Blueprint — engineering drawing: cyan paper, white hairlines, mono labels
T["06_blueprint"] = dict(
    spark="line", status="stxt", chart="line", chart_fill=.12, spark_sw=1.2, chart_sw=1.1,
    vars=v(font=MONO, numfont=MONO, backdrop="#0a2a57",
           bg="#0f3f85", titlebar="transparent", toolbar="transparent", statusbar="transparent",
           panel="rgba(8,40,90,.35)", pborder="rgba(255,255,255,.55)", thead="transparent", thtext="#cfe4ff",
           text="#ffffff", muted="#9cc3f5", border="rgba(255,255,255,.4)", rowsep="rgba(255,255,255,.14)",
           sel="rgba(255,255,255,.14)", accent="#ffffff", accent2="#9cd3ff", icon="#fff",
           btnon="#fff", btnontext="#0f3f85", input="transparent",
           ok="#8dffc0", slow="#ffe27a", down="#ff8f8f", downrow="#7fa8dc",
           okbg="transparent", slowbg="transparent", downbg="transparent",
           statbg="transparent", grid="rgba(255,255,255,.14)", chart1="#ffffff", chart2="#9cd3ff", link="#fff",
           radius="2px", pradius="0px", btnr="0px", thw="500",
           shadow="0 30px 80px rgba(0,0,0,.5), 0 0 0 6px #0f3f85, 0 0 0 7px rgba(255,255,255,.7)", winborder="rgba(255,255,255,.7)"),
    css="""
.win{background-image:linear-gradient(rgba(255,255,255,.07) 1px,transparent 1px),linear-gradient(90deg,rgba(255,255,255,.07) 1px,transparent 1px),linear-gradient(rgba(255,255,255,.12) 1px,transparent 1px),linear-gradient(90deg,rgba(255,255,255,.12) 1px,transparent 1px);background-size:20px 20px,20px 20px,100px 100px,100px 100px}
table{font-size:12px}td,th{padding:0 8px}
th{border-bottom:2px solid rgba(255,255,255,.7)}
.stat,.search,.seg{border-style:dashed}
.dot{border-radius:0;transform:rotate(45deg)}
.tblock{position:absolute;left:28px;bottom:56px;border:1px solid rgba(255,255,255,.7);font-size:10px;display:grid;grid-template-columns:auto auto auto;color:#cfe4ff;letter-spacing:.08em;background:#0f3f85}
.tblock div{padding:5px 10px;border:1px solid rgba(255,255,255,.35)}.tblock b{color:#fff;font-size:13px}
""",
    overlay='<div class="tblock"><div>PROJECT<br><b>ALTPING</b></div><div>SHEET<br><b>OFFICE</b></div><div>SCALE<br><b>1 s / PING</b></div>'
            '<div>HOSTS<br><b>12</b></div><div>REV<br><b>B</b></div><div>DATE<br><b>2026-10-06</b></div></div>',
    title_html='<span class="ttl"><b>ALT PING</b> / NETWORK LATENCY DIAGRAM / OFFICE</span>')

# 07 Trading terminal — black, amber, dense, green/red numbers
T["07_trading_terminal"] = dict(
    spark="bars", status="stxt", chart="bars", heat=True, chart_fill=.9,
    vars=v(font=MONO, numfont=MONO, backdrop="#000",
           bg="#000000", titlebar="#1a1206", toolbar="#000", statusbar="#1a1206",
           panel="#000", pborder="#3a2a0c", thead="#ff9e1b", thtext="#000",
           text="#ffb84d", muted="#a36a14", border="#3a2a0c", rowsep="#1c1406", sel="#2b1d05",
           accent="#ff9e1b", accent2="#ffd27a", icon="#ff9e1b", btnon="#ff9e1b", btnontext="#000", input="#0a0703",
           ok="#00e676", slow="#ffea00", down="#ff1744", warn="#ff6d00", downrow="#6b4a12",
           okbg="transparent", slowbg="transparent", downbg="transparent",
           statbg="#0a0703", grid="#221806", chart1="#ff9e1b", chart2="#ff9e1b", link="#ffd27a",
           radius="0px", pradius="0px", btnr="0px", rowh="32px", gap="8px", thw="700",
           shadow="0 0 0 1px #3a2a0c, 0 40px 100px rgba(0,0,0,.9)", winborder="#3a2a0c"),
    css="""
table{font-size:12.5px}td,th{padding:0 8px}th{height:26px}
.stxt.ok{color:var(--ok)}
tr.sel td{color:#fff}
.tick{position:absolute;left:0;right:0;bottom:38px;height:26px;background:#0a0703;border-top:1px solid #3a2a0c;display:flex;align-items:center;gap:26px;padding:0 14px;font-size:12px;white-space:nowrap;overflow:hidden}
.tick b{color:#fff;font-weight:500}.up{color:#00e676}.dn{color:#ff1744}
.main{padding-bottom:34px}
""",
    chart_dh=26,
    overlay='<div class="tick"><span><b>ROUTER</b> 0.4 <span class="up">▼0.1</span></span><span><b>NAS</b> 0.8 <span class="dn">▲0.1</span></span>'
            '<span><b>PRINTER</b> <span class="dn">TIMEOUT 14m</span></span><span><b>CLOUDFLARE</b> 8.6 <span class="up">▼1.4</span></span>'
            '<span><b>GOOGLE</b> 15 <span class="up">▼0.3</span></span><span><b>GITHUB</b> 38 <span class="dn">▲2.0</span></span>'
            '<span><b>VPN</b> 26 <span class="up">▼4.1</span></span><span><b>WEBSHOP</b> 61 <span class="dn">▲3.0 SLOW</span></span>'
            '<span><b>MAIL</b> 30 <span class="up">▼0.2</span></span><span><b>BACKUP-SG</b> 207 <span class="dn">▲25</span></span></div>',
    title_html='<span class="ttl"><b>ALTPING</b> &lt;GO&gt; · OFFICE · 23:31:07</span>')

# 08 Arctic — cool calm palette, dashboard cards
T["08_arctic"] = dict(
    layout="cards", spark="area", status="pill", chart="area", spark_w=120, spark_fill=.3, chart_fill=.35, chart_sw=1.5,
    vars=v(backdrop="linear-gradient(160deg,#88c0d0 0%,#5e81ac 60%,#4c566a 100%)",
           bg="#2e3440", titlebar="#2e3440", toolbar="#2e3440", statusbar="#2b303b",
           panel="#3b4252", pborder="#434c5e", thead="#3b4252", thtext="#a7b1c5",
           text="#eceff4", muted="#a7b1c5", border="#434c5e", rowsep="#434c5e", sel="#4c566a",
           accent="#88c0d0", accent2="#81a1c1", icon="#d8dee9", btnon="#88c0d0", btnontext="#2e3440", input="#353b49",
           ok="#a3be8c", slow="#ebcb8b", down="#bf616a", warn="#d08770", downrow="#7b869c",
           okbg="rgba(163,190,140,.16)", slowbg="rgba(235,203,139,.16)", downbg="rgba(191,97,106,.25)",
           statbg="#353b49", grid="#434c5e", chart1="#88c0d0", chart2="#b48ead", link="#88c0d0",
           rowh="36px", shadow="0 30px 80px rgba(20,30,50,.5)"),
    css=".kpi .kv{font-size:26px}")

# 09 Warm retro — earthy dark palette loved by terminal dwellers
T["09_warm_retro"] = dict(
    spark="bars", status="pill", chart="area", heat=True, chart_fill=.4, chart_sw=1.5,
    vars=v(font=MONO, numfont=MONO, backdrop="linear-gradient(135deg,#d79921 0%,#cc241d 100%)",
           bg="#282828", titlebar="#1d2021", toolbar="#282828", statusbar="#1d2021",
           panel="#32302f", pborder="#504945", thead="#3c3836", thtext="#bdae93",
           text="#ebdbb2", muted="#a89984", border="#504945", rowsep="#3c3836", sel="#504945",
           accent="#83a598", accent2="#fabd2f", icon="#d5c4a1", btnon="#fabd2f", btnontext="#282828", input="#1d2021",
           ok="#b8bb26", slow="#fabd2f", down="#fb4934", warn="#fe8019", downrow="#7c6f64",
           okbg="rgba(184,187,38,.14)", slowbg="rgba(250,189,47,.14)", downbg="rgba(251,73,52,.18)",
           statbg="#282828", grid="#3c3836", chart1="#fabd2f", chart2="#fe8019", link="#83a598",
           radius="8px", pradius="6px", btnr="4px", shadow="0 30px 80px rgba(60,20,0,.5)"),
    css="table{font-size:12px}td,th{padding:0 8px}.pill{border-radius:4px}")

# 10 Paper — ink on warm paper, serif headings, colour only for trouble
T["10_paper_ink"] = dict(
    spark="line", status="stxt", chart="line", chart_fill=.06, spark_sw=1.1, chart_sw=1,
    vars=v(backdrop="#d9d3c7", bg="#f6f2ea", titlebar="#f6f2ea", toolbar="#f6f2ea", statusbar="#f6f2ea",
           panel="#fbf8f2", pborder="#1a1a1a", thead="#fbf8f2", thtext="#1a1a1a",
           text="#1a1a1a", muted="#6f6a60", border="#1a1a1a", rowsep="#ddd6c8", sel="#ece5d6",
           accent="#1a1a1a", accent2="#1a1a1a", icon="#1a1a1a", btnon="#1a1a1a", btnontext="#f6f2ea", input="#fbf8f2",
           ok="#1a1a1a", slow="#a35b00", down="#c8102e", downrow="#9a9488",
           okbg="transparent", slowbg="transparent", downbg="transparent",
           statbg="#fbf8f2", grid="#e2dccf", chart1="#1a1a1a", chart2="#1a1a1a", link="#1a1a1a",
           radius="3px", pradius="0px", btnr="0px", thw="700", wc="#b9b2a3",
           shadow="0 30px 70px rgba(60,50,30,.35), 0 0 0 1px #1a1a1a"),
    css="""
.tb .ttl b,.rhead .t1,.stat .sv{font-family:'DejaVu Serif',serif}
.tb .ttl b{font-size:16px;font-style:italic}
.rhead .t1{font-size:19px}
th{border-bottom:2px solid #1a1a1a;border-top:1px solid #1a1a1a}
.left,.right{border-width:1px 0;background:transparent}
.right{border:none;border-left:1px solid #1a1a1a}
.dot{box-shadow:none;background:transparent;border:1.5px solid #1a1a1a}
.dot.slow{background:var(--slow);border-color:var(--slow)}.dot.down{background:var(--down);border-color:var(--down)}
tr.down td.name .nm{color:var(--down);font-weight:700}
.stxt.ok{color:#6f6a60;font-weight:400}
.stat{border:none;border-top:2px solid #1a1a1a;background:transparent;padding-left:0}
.chart .lastpt{fill:#c8102e;filter:none;stroke:#fbf8f2}
.status a{text-decoration:underline}
""",
    title_html='<span class="ttl"><b>Alt Ping</b> — Office network, evening edition</span>')

# 11 High contrast — large type, thick borders, status readable without colour
T["11_high_contrast"] = dict(
    cols="compact", chart_dh=24, spark="line", status="pill", chart="line", chart_fill=0, spark_sw=2, chart_sw=1.6, spark_w=130,
    vars=v(backdrop="#000", bg="#000", titlebar="#000", toolbar="#000", statusbar="#000",
           panel="#000", pborder="#fff", thead="#000", thtext="#ffff00",
           text="#ffffff", muted="#d0d0d0", border="#fff", rowsep="#666", sel="#003b8e",
           accent="#ffff00", accent2="#ffff00", icon="#fff", btnon="#ffff00", btnontext="#000", input="#000",
           ok="#00ff66", slow="#ffd000", down="#ff5c5c", downrow="#d0d0d0",
           okbg="#00ff66", slowbg="#ffd000", downbg="#ff5c5c",
           statbg="#000", grid="#444", chart1="#ffff00", chart2="#ffff00", link="#ffff00",
           radius="0px", pradius="0px", btnr="0px", rowh="46px", thsize="14px", thw="700",
           shadow="0 0 0 3px #fff", winborder="#fff"),
    css="""
table{font-size:16px}td.name{font-weight:700}
.left,.right,.search,.stat,.seg{border-width:2px}
.pill{color:#000!important;font-size:13px;height:26px;border-radius:0}
.pill.ok i{border-radius:50%}.pill.slow i{border-radius:0;transform:rotate(45deg)}
.pill.down i{border-radius:0;width:10px;height:3px}
.tbtn{border:2px solid #fff;margin-right:6px}
.status{font-size:15px}.status a{text-decoration:underline}
.chart .ax{font-size:13px;fill:#fff}.stat .sv{font-size:22px}.rhead .t1{font-size:18px}
tbody tr.sel td:first-child{box-shadow:inset 6px 0 0 #ffff00}
""")

# 12 Bubblegum — friendly rounded tiles with little faces for status
T["12_bubblegum"] = dict(
    layout="tiles", faces=True, spark="area", status="pill", chart="area", spark_h=34, spark_fill=.45,
    chart_fill=.45, chart_sw=2, spark_sw=2, down_text="zzz…", spark_state_color=True,
    vars=v(backdrop="linear-gradient(135deg,#ffd6e8 0%,#e3d5ff 50%,#c9f0ff 100%)",
           bg="#fff7fb", titlebar="#fff7fb", toolbar="#fff7fb", statusbar="#ffffff",
           panel="#ffffff", pborder="#ffe0ee", thead="#fff", thtext="#b58aa5",
           text="#4a3350", muted="#a98ba8", border="#ffe0ee", rowsep="#ffeef6", sel="#fff0f7",
           accent="#ff6fae", accent2="#a58bff", icon="#8d6b96", btnon="linear-gradient(135deg,#ff8cc0,#a58bff)",
           btnontext="#fff", input="#fff",
           ok="#2fcf9a", slow="#ffa63d", down="#ff5d8f", downrow="#c9b3c9",
           okbg="#e4faf2", slowbg="#fff1df", downbg="#ffe6ee",
           statbg="#fff7fb", grid="#fbe9f2", chart1="#ff6fae", chart2="#a58bff", link="#d9488f",
           radius="26px", pradius="22px", btnr="14px", wc="#f3cfe0",
           shadow="0 30px 80px rgba(200,100,170,.35), 0 0 0 1px rgba(255,255,255,.8)"),
    css="""
.tile{box-shadow:0 6px 18px rgba(255,120,180,.12);border-width:2px}
.tile.ok .face{color:var(--ok)}.tile.slow .face{color:var(--slow)}.tile.down .face{color:var(--down)}
.tile.down{background:#fff5f8;border-color:#ffc2d6}
.tile.sel{border-color:#ff6fae;box-shadow:0 0 0 3px rgba(255,111,174,.25)}
.pill{border-radius:12px}.right{box-shadow:0 6px 18px rgba(255,120,180,.12);border-width:2px}
.stat{border-radius:16px}
""")

if __name__ == "__main__":
    only = sys.argv[1:]
    os.makedirs(os.path.join(OUT, "r2"), exist_ok=True)
    for key, t in T.items():
        if only and not any(o in key for o in only):
            continue
        with open(os.path.join(OUT, "r2", f"{key}.html"), "w") as f:
            f.write(page(t))
        print(key)
