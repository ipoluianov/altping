#!/usr/bin/env python3
"""Generates HTML mockups of AltPing in several visual styles."""
import math, random, os, sys

OUT = os.path.dirname(os.path.abspath(__file__))

# name, ip, time, loss, min, avg, max, jitter, since, state, bell, base
HOSTS = [
    ("Router", "192.168.1.1", "0.4", "0", "0.4", "0.4", "1.0", "0.1", "14h 11m", "ok", False, 0.4),
    ("Core switch", "192.168.1.2", "0.5", "0", "0.5", "0.6", "1.2", "0.1", "14h 11m", "ok", False, 0.5),
    ("NAS", "192.168.1.20", "0.8", "0", "0.7", "0.9", "1.8", "0.2", "14h 11m", "ok", False, 0.8),
    ("Printer, 2nd floor", "192.168.1.42", "—", "100", "—", "—", "—", "—", "14m 39s", "down", True, 0),
    ("Wi-Fi AP, meeting room", "192.168.1.61", "5.1", "1.0", "3.2", "4.8", "12", "1.6", "3m 6s", "ok", False, 4.8),
    ("Cloudflare DNS", "1.1.1.1", "8.6", "0", "8.4", "10", "27", "0.8", "2h 21m", "ok", False, 10),
    ("google.com", "142.250.185.78", "15", "0", "14", "15", "31", "1.2", "2h 49m", "ok", False, 15),
    ("github.com", "140.82.121.4", "38", "0", "36", "38", "82", "2.0", "46m 33s", "ok", False, 38),
    ("VPN gateway", "203.0.113.10", "26", "0", "24", "30", "58", "2.6", "7m 15s", "ok", True, 30),
    ("Web shop", "203.0.113.25", "61", "0", "58", "63", "116", "5.3", "21m 47s", "slow", True, 63),
    ("Mail", "203.0.113.40", "30", "0", "29", "30", "35", "1.3", "1h 7m", "ok", False, 30),
    ("Backup, Singapore", "198.51.100.7", "207", "1.0", "171", "182", "223", "6.9", "43s", "ok", False, 182),
]
SEL = 11


def series(seed, base, n=60, spikes=2, jit=0.12):
    r = random.Random(seed)
    v = []
    for i in range(n):
        x = base * (1 + r.gauss(0, jit))
        if r.random() < spikes / n:
            x *= 1.6 + r.random() * 1.4
        v.append(max(x, base * 0.5))
    return v


def spark_path(vals, w, h, pad=2, close=False):
    lo, hi = min(vals), max(vals)
    if hi - lo < 1e-9:
        hi = lo + 1
    pts = []
    for i, v in enumerate(vals):
        x = i * (w / (len(vals) - 1))
        y = pad + (h - 2 * pad) * (1 - (v - lo) / (hi - lo))
        pts.append((x, y))
    d = "M" + " L".join(f"{x:.1f},{y:.1f}" for x, y in pts)
    if close:
        d += f" L{w},{h} L0,{h} Z"
    return d


BLOCKS = "▁▂▃▄▅▆▇█"


def spark_blocks(vals):
    vals = vals[::3]
    sv = sorted(vals); lo, hi = sv[0], sv[int(len(sv)*.9)]
    return "".join(BLOCKS[min(7, int((v - lo) / (hi - lo + 1e-9) * 8))] for v in vals)


def spark(t, idx, host):
    w, h = t.get("spark_w", 150), 24
    state, base = host[9], host[11]
    if state == "down":
        if t["spark"] == "blocks":
            return f'<span class="blk down">{"▁" * 20}</span>'
        return (f'<svg width="{w}" height="{h}"><rect x="0" y="{h/2-2}" width="{w}" height="4" rx="2" '
                f'fill="var(--down)" opacity=".85"/></svg>')
    vals = series(idx * 7 + 3, base, spikes=2 + idx % 3)
    loss = [i for i in ([9, 14] if host[3] not in ("0",) else [])]
    if t["spark"] == "blocks":
        return f'<span class="blk">{spark_blocks(vals)}</span>'
    col = "var(--slow)" if state == "slow" else "var(--accent)"
    line = spark_path(vals, w, h)
    area = spark_path(vals, w, h, close=True)
    gid = f"g{idx}"
    out = [f'<svg width="{w}" height="{h}" style="overflow:visible">']
    out.append(f'<defs><linearGradient id="{gid}" x1="0" y1="0" x2="0" y2="1">'
               f'<stop offset="0" stop-color="{col}" stop-opacity="{t.get("spark_fill", .35)}"/>'
               f'<stop offset="1" stop-color="{col}" stop-opacity="0"/></linearGradient></defs>')
    if t["spark"] in ("area", "glow"):
        out.append(f'<path d="{area}" fill="url(#{gid})"/>')
    flt = ' class="glow"' if t["spark"] == "glow" else ""
    out.append(f'<path{flt} d="{line}" fill="none" stroke="{col}" stroke-width="{t.get("spark_sw", 1.3)}" '
               f'stroke-linejoin="round" stroke-linecap="round"/>')
    for li in loss:
        x = li * w / 59
        out.append(f'<rect x="{x:.1f}" y="{h-3}" width="2" height="5" rx="1" fill="var(--down)"/>')
    out.append("</svg>")
    return "".join(out)


def status_cell(t, host):
    st = host[9]
    label = {"ok": "OK", "slow": "SLOW", "down": "TIMEOUT"}[st]
    s = t["status"]
    if s == "pill":
        return f'<span class="pill {st}"><i></i>{label}</span>'
    if s == "led":
        return f'<span class="led {st}"></span>'
    return f'<span class="stxt {st}">{label}</span>'


def heat_color(ms):
    try:
        v = float(ms)
    except ValueError:
        return "var(--down)"
    if v < 5:
        return "var(--ok)"
    if v < 40:
        return "var(--text)"
    if v < 100:
        return "var(--slow)"
    return "var(--warn, var(--slow))"


def uptime_strip(idx, host, cells=48):
    r = random.Random(idx * 13 + 1)
    out = ['<span class="strip">']
    for i in range(cells):
        c = "ok"
        x = r.random()
        if host[9] == "down" and i >= cells - 1:
            c = "down"
        elif x < 0.025:
            c = "down"
        elif x < 0.08 and host[11] > 20:
            c = "slow"
        if idx == 4 and i in (30, 31):
            c = "slow"
        out.append(f'<i class="{c}"></i>')
    out.append("</span>")
    return "".join(out)


def history():
    r = random.Random(42)
    v = []
    for i in range(360):
        base = 205 if i < 60 else (210 if i < 200 else 185)
        x = base + r.gauss(0, 7)
        if r.random() < 0.05:
            x += r.random() * 80
        if r.random() < 0.006:
            x += 120 + r.random() * 160
        v.append(x)
    return v


def chart(t, w, h):
    vals = history()
    lo, hi = 150, 500
    pl, pr, pt, pb = 40, 8, 26, 36
    cw, ch = w - pl - pr, h - pt - pb
    pts = [(pl + i * cw / (len(vals) - 1), pt + ch * (1 - (min(v, hi) - lo) / (hi - lo))) for i, v in enumerate(vals)]
    line = "M" + " L".join(f"{x:.1f},{y:.1f}" for x, y in pts)
    area = line + f" L{pl+cw},{pt+ch} L{pl},{pt+ch} Z"
    out = [f'<svg width="{w}" height="{h}" class="chart">']
    out.append('<defs><linearGradient id="cg" x1="0" y1="0" x2="0" y2="1">'
               f'<stop offset="0" stop-color="var(--chart1)" stop-opacity="{t.get("chart_fill", .45)}"/>'
               '<stop offset="1" stop-color="var(--chart2)" stop-opacity="0.02"/></linearGradient>'
               '<linearGradient id="cl" x1="0" y1="0" x2="1" y2="0">'
               '<stop offset="0" stop-color="var(--chart1)"/><stop offset="1" stop-color="var(--chart2)"/></linearGradient>'
               '</defs>')
    if t["chart"] == "synth":
        # perspective grid floor
        for i in range(0, 13):
            y = pt + ch * (i / 12) ** 1.6
            out.append(f'<line x1="{pl}" x2="{pl+cw}" y1="{y:.1f}" y2="{y:.1f}" class="grid"/>')
        for i in range(-8, 17):
            x0 = pl + cw / 8 * i
            out.append(f'<line x1="{pl+cw/2+(x0-pl-cw/2)*0.25:.1f}" y1="{pt}" x2="{x0:.1f}" y2="{pt+ch}" class="grid"/>')
    else:
        for i in range(8):
            y = pt + ch * i / 7
            out.append(f'<line x1="{pl}" x2="{pl+cw}" y1="{y:.1f}" y2="{y:.1f}" class="grid"/>')
    for i in range(8):
        y = pt + ch * i / 7
        val = hi - (hi - lo) * i / 7
        out.append(f'<text x="{pl-8}" y="{y+4:.1f}" class="ax" text-anchor="end">{val:.0f}</text>')
    labels = ["22:40", "22:50", "23:00", "23:10", "23:20", "23:30"]
    for i, l in enumerate(labels):
        x = pl + cw * (i + 0.5) / 6
        out.append(f'<text x="{x:.1f}" y="{pt+ch+18}" class="ax" text-anchor="middle">{l}</text>')
    # loss markers
    for i in (61, 62, 200, 288):
        x = pts[i][0]
        out.append(f'<line x1="{x:.1f}" x2="{x:.1f}" y1="{pt}" y2="{pt+ch}" class="lossmark"/>')
    if t["chart"] == "bars":
        bw = cw / len(vals)
        for (x, y) in pts:
            out.append(f'<rect x="{x:.1f}" y="{y:.1f}" width="{max(bw-0.3,0.8):.2f}" height="{pt+ch-y:.1f}" fill="url(#cg)"/>')
        out.append(f'<path d="{line}" fill="none" stroke="var(--chart1)" stroke-width="1" opacity=".9"/>')
    else:
        out.append(f'<path d="{area}" fill="url(#cg)"/>')
        cls = ' class="cglow"' if t["chart"] in ("glow", "synth") else ""
        out.append(f'<path{cls} d="{line}" fill="none" stroke="url(#cl)" stroke-width="{t.get("chart_sw", 1.2)}" stroke-linejoin="round"/>')
    # avg line + label
    ay = pt + ch * (1 - (190 - lo) / (hi - lo))
    out.append(f'<line x1="{pl}" x2="{pl+cw}" y1="{ay:.1f}" y2="{ay:.1f}" class="avgline"/>')
    out.append(f'<text x="{pl+6}" y="{ay+16:.1f}" class="avglbl">avg 190 ms</text>')
    # last point
    lx, ly = pts[-1]
    out.append(f'<circle cx="{lx:.1f}" cy="{ly:.1f}" r="4" class="lastpt"/>')
    out.append("</svg>")
    return "".join(out)


ICONS = {
    "add": '<circle cx="12" cy="12" r="9"/><path d="M12 8v8M8 12h8"/>',
    "edit": '<path d="M4 20h4L19 9l-4-4L4 16z"/><path d="M14 6l4 4"/>',
    "del": '<path d="M5 7h14M10 7V4h4v3M7 7l1 13h8l1-13"/>',
    "chart": '<path d="M4 4v16h16"/><path d="M7 15l4-5 3 3 5-7"/>',
    "share": '<path d="M14 4h6v6M20 4l-9 9M18 14v6H4V6h6"/>',
    "play": '<path d="M7 4l13 8-13 8z"/>',
    "stop": '<rect x="6" y="6" width="12" height="12" rx="2"/>',
    "open": '<path d="M3 7h7l2 2h9v10H3z"/>',
    "save": '<path d="M12 4v11M7 10l5 5 5-5M4 20h16"/>',
    "bell": '<path d="M6 16V11a6 6 0 0 1 12 0v5l2 2H4z"/><path d="M10 20h4"/>',
    "search": '<circle cx="11" cy="11" r="6"/><path d="M20 20l-4.5-4.5"/>',
}


def icon(n, size=20, cls=""):
    return (f'<svg class="ic {cls}" width="{size}" height="{size}" viewBox="0 0 24 24" fill="none" '
            f'stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">{ICONS[n]}</svg>')


LOGO = ('<svg width="20" height="20" viewBox="0 0 24 24"><defs><linearGradient id="lg" x1="0" y1="0" x2="1" y2="1">'
        '<stop offset="0" stop-color="var(--accent)"/><stop offset="1" stop-color="var(--accent2)"/></linearGradient></defs>'
        '<rect width="24" height="24" rx="6" fill="url(#lg)"/>'
        '<path d="M4 13h4l2-5 3 9 2-6 1 2h4" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>')


def table(t, cols):
    head = "".join(f'<th class="{c[1]}">{c[0]}</th>' for c in cols)
    rows = []
    for i, h in enumerate(HOSTS):
        cls = ["sel"] if i == SEL else []
        cls.append(h[9])
        cells = []
        for c in cols:
            k = c[2]
            if k == "name":
                bell = icon("bell", 13, "bell") if h[10] else '<span class="nobell"></span>'
                dot = f'<span class="dot {h[9]}"></span>' if t["status"] != "led" else status_cell(t, h)
                v = f'{dot}{bell}<span class="nm">{h[0]}</span>'
            elif k == "ip":
                v = h[1]
            elif k == "time":
                if t.get("heat"):
                    v = f'<b style="color:{heat_color(h[2])}">{h[2]}</b>'
                else:
                    v = h[2]
            elif k == "spark":
                v = spark(t, i, h)
            elif k == "loss":
                v = h[3] if h[3] == "0" else f'<span class="lossv">{h[3]}</span>'
            elif k == "min":
                v = h[4]
            elif k == "avg":
                v = h[5]
            elif k == "max":
                v = h[6]
            elif k == "jit":
                v = h[7]
            elif k == "since":
                v = h[8]
            elif k == "det":
                v = status_cell(t, h) if t["status"] != "led" else f'<span class="stxt {h[9]}">{ {"ok":"OK","slow":"SLOW","down":"TIMEOUT"}[h[9]] }</span>'
            elif k == "strip":
                v = uptime_strip(i, h)
            cells.append(f'<td class="{c[1]}">{v}</td>')
        rows.append(f'<tr class="{" ".join(cls)}">{"".join(cells)}</tr>')
    return f'<table><thead><tr>{head}</tr></thead><tbody>{"".join(rows)}</tbody></table>'


COLS_CLASSIC = [("Name", "l name", "name"), ("IP", "l ip", "ip"), ("Time ms", "r", "time"),
                ("Last 5 min", "l spk", "spark"), ("Loss %", "r", "loss"), ("Min", "r", "min"),
                ("Avg", "r", "avg"), ("Max", "r", "max"), ("Jitter", "r", "jit"),
                ("Since", "r", "since"), ("Status", "l det", "det")]
COLS_CARDS = [("Host", "l name", "name"), ("Address", "l ip", "ip"), ("Ping", "r", "time"),
              ("Last 5 min", "l spk", "spark"), ("Loss", "r", "loss"), ("Avg", "r", "avg"),
              ("Jitter", "r", "jit"), ("Last 24 h", "l strp", "strip"), ("Status", "l det", "det")]


def kpis(t):
    items = [("Hosts", "12", "monitored", ""), ("Up", "10", "83%", "ok"), ("Slow", "1", "Web shop", "slow"),
             ("Down", "1", "Printer, 2nd floor", "down"), ("Median ping", "22 ms", "▼ 3 ms vs 1h", ""),
             ("Uptime 24h", "99.2%", "2 outages", "ok")]
    return '<div class="kpis">' + "".join(
        f'<div class="kpi {c}"><div class="kl">{a}</div><div class="kv">{b}</div><div class="ks">{s}</div></div>'
        for a, b, s, c in items) + "</div>"


BASE_CSS = """
*{box-sizing:border-box;margin:0;padding:0}
html,body{width:1600px;height:960px;overflow:hidden}
body{font-family:var(--font);color:var(--text);background:var(--backdrop);display:flex;align-items:center;justify-content:center;-webkit-font-smoothing:antialiased}
.win{width:1520px;height:880px;background:var(--bg);border-radius:var(--radius);box-shadow:var(--shadow);display:flex;flex-direction:column;overflow:hidden;position:relative;border:1px solid var(--winborder)}
.tb{height:38px;display:flex;align-items:center;padding:0 14px;gap:10px;background:var(--titlebar);font-size:13px;color:var(--muted);position:relative}
.tb .ttl{position:absolute;left:0;right:0;text-align:center;pointer-events:none}
.tb .ttl b{color:var(--text);font-weight:600}
.tb .wc{margin-left:auto;display:flex;gap:8px}
.tb .wc i{width:12px;height:12px;border-radius:50%;background:var(--wc);display:block}
.tools{height:54px;display:flex;align-items:center;padding:0 14px;gap:4px;background:var(--toolbar);border-bottom:1px solid var(--border)}
.tbtn{width:38px;height:38px;border-radius:var(--btnr);display:flex;align-items:center;justify-content:center;color:var(--icon)}
.tbtn.on{background:var(--btnon);color:var(--btnontext)}
.tbtn.dis{opacity:.35}
.sp{flex:1}
.search{height:32px;width:260px;border-radius:var(--btnr);background:var(--input);border:1px solid var(--border);display:flex;align-items:center;gap:8px;padding:0 10px;color:var(--muted);font-size:13px;margin-right:12px}
.sepv{width:1px;height:24px;background:var(--border);margin:0 8px}
.main{flex:1;display:flex;gap:var(--gap);padding:var(--gap);min-height:0}
.left{flex:1;background:var(--panel);border-radius:var(--pradius);overflow:hidden;border:1px solid var(--pborder);min-width:0}
.right{width:430px;background:var(--panel);border-radius:var(--pradius);border:1px solid var(--pborder);display:flex;flex-direction:column;overflow:hidden}
table{width:100%;border-collapse:collapse;font-size:13px}
th{height:34px;font-weight:var(--thw);color:var(--thtext);background:var(--thead);font-size:var(--thsize);text-transform:var(--thcase);letter-spacing:var(--thls);padding:0 10px;border-bottom:1px solid var(--border);white-space:nowrap}
td{height:var(--rowh);padding:0 10px;border-bottom:1px solid var(--rowsep);white-space:nowrap;font-variant-numeric:tabular-nums}
.l{text-align:left}.r{text-align:right}
td.r,td.ip{font-family:var(--numfont)}
td.ip{color:var(--muted)}
tbody tr.sel{background:var(--sel)}
tbody tr.sel td:first-child{box-shadow:inset 3px 0 0 var(--accent)}
td.name{font-weight:var(--namew)}
td.name>*{vertical-align:middle}
.dot{display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:10px;background:var(--ok);box-shadow:var(--dotglow) var(--ok)}
.dot.slow{background:var(--slow);box-shadow:var(--dotglow) var(--slow)}
.dot.down{background:var(--down);box-shadow:var(--dotglow) var(--down)}
.bell{color:var(--muted);margin-right:6px}
.nobell{display:inline-block;width:19px}
tr.down td{color:var(--downrow)}
tr.down td.name .nm{color:var(--text)}
.lossv{color:var(--slow)}
tr.down .lossv{color:var(--down);font-weight:600}
.pill{display:inline-flex;align-items:center;gap:6px;height:22px;padding:0 9px;border-radius:11px;font-size:11px;font-weight:700;letter-spacing:.04em}
.pill i{width:6px;height:6px;border-radius:50%;background:currentColor}
.pill.ok{color:var(--ok);background:var(--okbg)}
.pill.slow{color:var(--slow);background:var(--slowbg)}
.pill.down{color:var(--down);background:var(--downbg)}
.stxt{font-size:12px;font-weight:600}
.stxt.ok{color:var(--muted)}.stxt.slow{color:var(--slow)}.stxt.down{color:var(--down)}
.glow{filter:drop-shadow(0 0 3px var(--accent))}
.cglow{filter:drop-shadow(0 0 4px var(--chart1))}
.rhead{display:flex;align-items:center;padding:12px 14px 4px;gap:8px}
.rhead .t1{font-size:15px;font-weight:700}
.rhead .t2{font-size:12px;color:var(--muted)}
.seg{margin-left:auto;display:flex;background:var(--input);border-radius:var(--btnr);padding:3px;border:1px solid var(--border)}
.seg span{font-size:12px;padding:4px 10px;border-radius:calc(var(--btnr) - 2px);color:var(--muted)}
.seg span.on{background:var(--btnon);color:var(--btnontext);font-weight:600}
.stats{display:flex;gap:8px;padding:10px 14px}
.stat{flex:1;background:var(--statbg);border-radius:calc(var(--pradius) - 2px);padding:8px 10px;border:1px solid var(--pborder)}
.stat .sl{font-size:11px;color:var(--muted);text-transform:uppercase;letter-spacing:.06em}
.stat .sv{font-size:18px;font-weight:700;font-family:var(--numfont);margin-top:2px}
.stat .sv small{font-size:11px;color:var(--muted);font-weight:500;margin-left:2px}
.chart .grid{stroke:var(--grid);stroke-width:1}
.chart .ax{fill:var(--muted);font-size:11px;font-family:var(--numfont)}
.chart .lossmark{stroke:var(--down);stroke-width:1.5;stroke-dasharray:2 3;opacity:.55}
.chart .avgline{stroke:var(--text);stroke-width:1;stroke-dasharray:4 4;opacity:.35}
.chart .avglbl{fill:var(--muted);font-size:11px}
.chart .lastpt{fill:var(--chart2);stroke:var(--panel);stroke-width:2;filter:drop-shadow(0 0 6px var(--chart2))}
.status{height:38px;display:flex;align-items:center;padding:0 16px;gap:18px;font-size:13px;background:var(--statusbar);border-top:1px solid var(--border);color:var(--muted)}
.status b{color:var(--text);font-weight:600}
.status .c{display:flex;align-items:center;gap:6px}
.status .c i{width:8px;height:8px;border-radius:50%;display:block}
.status a{color:var(--link);text-decoration:none}
.kpis{display:flex;gap:var(--gap);padding:var(--gap) var(--gap) 0}
.kpi{flex:1;background:var(--panel);border:1px solid var(--pborder);border-radius:var(--pradius);padding:12px 16px;position:relative;overflow:hidden}
.kpi .kl{font-size:11px;color:var(--muted);text-transform:uppercase;letter-spacing:.08em}
.kpi .kv{font-size:28px;font-weight:800;font-family:var(--numfont);margin-top:2px;letter-spacing:-.02em}
.kpi .ks{font-size:12px;color:var(--muted)}
.kpi.ok .kv{color:var(--ok)}.kpi.slow .kv{color:var(--slow)}.kpi.down .kv{color:var(--down)}
.kpi::after{content:"";position:absolute;left:0;top:0;bottom:0;width:3px;background:var(--kbar,transparent)}
.kpi.ok{--kbar:var(--ok)}.kpi.slow{--kbar:var(--slow)}.kpi.down{--kbar:var(--down)}
.strip{display:inline-flex;gap:1px;vertical-align:middle}
.strip i{width:3px;height:16px;border-radius:1px;background:var(--ok);opacity:.85}
.strip i.slow{background:var(--slow)}.strip i.down{background:var(--down)}
.blk{font-family:var(--numfont);color:var(--accent);letter-spacing:-1px;font-size:14px}
.blk.down{color:var(--down)}
tr.slow .blk{color:var(--slow)}
.led{display:inline-block;width:10px;height:10px;border-radius:2px;margin-right:10px;background:var(--ok);box-shadow:0 0 8px var(--ok)}
.led.slow{background:var(--slow);box-shadow:0 0 8px var(--slow)}
.led.down{background:var(--down);box-shadow:0 0 10px var(--down);animation:none}
.badge{position:absolute;font-family:var(--font)}
"""


def page(t):
    layout = t.get("layout", "classic")
    cols = COLS_CARDS if layout == "cards" else COLS_CLASSIC
    content_h = 880 - 38 - 54 - 38 - 2
    if layout == "cards":
        content_h -= 96 + 14
    ch = content_h - 2 * 14 - 46 - 70 - 10
    tools = "".join(f'<div class="tbtn">{icon(n)}</div>' for n in ("add", "edit", "del"))
    rtools = (f'<div class="tbtn on">{icon("chart")}</div><div class="tbtn dis">{icon("share")}</div>'
              f'<div class="tbtn dis">{icon("play")}</div><div class="tbtn">{icon("stop")}</div><div class="sepv"></div>'
              f'<div class="tbtn">{icon("open")}</div><div class="tbtn">{icon("save")}</div>')
    title = t.get("title_html", '<span class="ttl"><b>Alt Ping</b> · Office</span>')
    html = f"""<!doctype html><html><head><meta charset="utf-8"><style>
:root{{{t['vars']}}}
{BASE_CSS}
{t.get('css','')}
</style></head><body>
{t.get('pre','')}
<div class="win">
 <div class="tb">{LOGO}{title}<div class="wc"><i></i><i></i><i></i></div></div>
 <div class="tools">{tools}<div class="sp"></div><div class="search">{icon("search",16)}Filter hosts…</div>{rtools}</div>
 {kpis(t) if layout == "cards" else ""}
 <div class="main">
  <div class="left">{table(t, cols)}</div>
  <div class="right">
   <div class="rhead"><div><div class="t1">Backup, Singapore</div><div class="t2">198.51.100.7 · ICMP · every 1 s</div></div>
    <div class="seg"><span>5m</span><span class="on">1h</span><span>24h</span></div></div>
   <div class="stats"><div class="stat"><div class="sl">Now</div><div class="sv">207<small>ms</small></div></div>
    <div class="stat"><div class="sl">Avg</div><div class="sv">182<small>ms</small></div></div>
    <div class="stat"><div class="sl">Loss</div><div class="sv" style="color:var(--slow)">1.0<small>%</small></div></div>
    <div class="stat"><div class="sl">Jitter</div><div class="sv">6.9<small>ms</small></div></div></div>
   <div style="padding:0 10px">{chart(t, 410, ch)}</div>
  </div>
 </div>
 <div class="status"><span class="c"><i style="background:var(--ok)"></i><b>10</b> up</span>
  <span class="c"><i style="background:var(--slow)"></i><b>1</b> slow</span>
  <span class="c"><i style="background:var(--down)"></i><b>1</b> down</span>
  <span>12 hosts · UDP mode</span><span class="sp"></span><a>Settings</a><a>Help</a><a>About</a></div>
 {t.get('overlay','')}
</div>
{t.get('post','')}
</body></html>"""
    return html


if __name__ == "__main__":
    sys.path.insert(0, OUT)
    from themes import THEMES
    for key, t in THEMES.items():
        with open(os.path.join(OUT, f"{key}.html"), "w") as f:
            f.write(page(t))
        print(key)
