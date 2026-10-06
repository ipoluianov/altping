#!/usr/bin/env bash
# Takes the screenshot of the main window on model data, without showing anything
# on the screen: the app runs in a virtual X server (Xvfb), under KWin when there
# is one (for the window frame), with the data of build/_screenshot/demo.go.
# Usage: build/_screenshot.sh [out.png]
# The default is screenshot.png in the bin/ directory of the build (as _build.sh names it).
# Needs Linux, Xvfb, python3 and ImageMagick; without them it says so and exits 0.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
SRC="$ROOT/build/_screenshot"

VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
OUT="${1:-bin/${BUILD_STAMP:-$(date +%Y-%m-%d-%H-%M)}-${VERSION}/screenshot.png}"

# demo.go sets the window to 1840x900; the "1h" button of the details is at
# this point of it, and the selected host is this many rows down
CLICK_PERIOD_X=1757
CLICK_PERIOD_Y=72
ROWS_DOWN=11

skip() {
  echo "Screenshot skipped: $*"
  exit 0
}
[ "$(uname -s)" = "Linux" ] || skip "needs Linux"
command -v Xvfb >/dev/null || skip "no Xvfb (apt install xvfb)"
command -v python3 >/dev/null || skip "no python3"
command -v xwininfo >/dev/null || skip "no xwininfo (apt install x11-utils)"
if command -v magick >/dev/null; then
  IM=(magick)
  IMPORT=(magick import)
elif command -v convert >/dev/null && command -v import >/dev/null; then
  IM=(convert)
  IMPORT=(import)
else
  skip "no ImageMagick (apt install imagemagick)"
fi

echo "Screenshot ${OUT}"
TMP="$(mktemp -d)"
PIDS=()
cleanup() {
  # Only the processes started here: by pid, the KWin session by its group
  for pid in "${PIDS[@]}"; do
    kill -- "$pid" 2>/dev/null || true
  done
  wait 2>/dev/null || true
  rm -rf "$TMP"
}
trap cleanup EXIT

# The hosts ping by the model: the network calls of host.go are replaced, and the
# generated history must go on without the breaks a start adds (host.go, history.go)
sed \
  -e 's|^\tc\.addGapToHistory()$|\t// screenshot: the generated history goes on without a break|' \
  -e 's|c\.checkIP()|c.demoCheckIP()|' \
  -e 's|c\.connectTCP(port, c\.configHost\.Timeout())|demoPing(c, port)|' \
  -e 's|c\.pingServer\.PingHost(c\.IP, 64, int(c\.configHost\.Timeout()\.Milliseconds()), c\.chanStop)|demoPing(c, "")|' \
  system/host.go > "$TMP/host.go"
sed '/The program was not stopped properly/,/^\t}$/d' system/history.go > "$TMP/history.go"
# Fail loudly when the sources changed and a replacement no longer applies
[ "$(grep -c 'demoPing(c, ' "$TMP/host.go")" = 2 ] &&
  grep -q 'demoCheckIP()' "$TMP/host.go" &&
  grep -q 'screenshot: the generated history' "$TMP/host.go" &&
  ! grep -q 'not stopped properly' "$TMP/history.go" ||
  { echo "Screenshot: system/host.go or system/history.go changed, update the replacements in $0"; exit 1; }

printf '{"Replace":{"%s":"%s","%s":"%s","%s":"%s"}}' \
  "$ROOT/system/host.go" "$TMP/host.go" \
  "$ROOT/system/history.go" "$TMP/history.go" \
  "$ROOT/system/demo.go" "$SRC/demo.go" > "$TMP/overlay.json"
CGO_ENABLED=0 GOOS= GOARCH= go build -overlay "$TMP/overlay.json" -o "$TMP/altping" .

# Xvfb picks a free display itself and writes its number
Xvfb -displayfd 3 -screen 0 1920x1080x24 -nolisten tcp 3>"$TMP/display" >"$TMP/xvfb.log" 2>&1 &
PIDS+=($!)
for _ in $(seq 50); do
  [ -s "$TMP/display" ] && break
  sleep 0.1
done
[ -s "$TMP/display" ] || { echo "Screenshot: Xvfb did not start"; cat "$TMP/xvfb.log"; exit 1; }
export DISPLAY=":$(cat "$TMP/display")"

# KWin draws the frame. It gets its own session bus and a copy of the look
# settings, so it does not touch the desktop session.
if command -v kwin_x11 >/dev/null && command -v dbus-run-session >/dev/null; then
  mkdir -p "$TMP/kwin-config"
  for f in kdeglobals kwinrc breezerc; do
    [ -f "$HOME/.config/$f" ] && cp "$HOME/.config/$f" "$TMP/kwin-config/"
  done
  XDG_CONFIG_HOME="$TMP/kwin-config" setsid dbus-run-session -- kwin_x11 >"$TMP/kwin.log" 2>&1 &
  PIDS+=("-$!")
  sleep 2
else
  echo "Screenshot: no kwin_x11, the window is taken without a frame"
fi

mkdir -p "$TMP/home"
HOME="$TMP/home" "$TMP/altping" >"$TMP/app.log" 2>&1 &
PIDS+=($!)

WIN=""
for _ in $(seq 100); do
  WIN="$(python3 -I "$SRC/x11.py" find "Alt Ping" 2>/dev/null)" && break
  sleep 0.2
done
[ -n "$WIN" ] || { echo "Screenshot: the window did not appear"; tail -20 "$TMP/app.log"; exit 1; }
sleep 3

for _ in $(seq "$ROWS_DOWN"); do
  python3 -I "$SRC/x11.py" key "$WIN" Down
  sleep 0.1
done
python3 -I "$SRC/x11.py" click "$WIN" "$CLICK_PERIOD_X" "$CLICK_PERIOD_Y"
sleep 3

# The frame is the top-level ancestor of the window (the window itself without a WM)
ROOTWIN="$(xwininfo -root | sed -n 's/.*Window id: \(0x[0-9a-f]*\).*/\1/p')"
FRAME="$WIN"
while :; do
  PARENT="$(xwininfo -id "$FRAME" -tree | sed -n 's/.*Parent window id: \(0x[0-9a-f]*\).*/\1/p')"
  [ "$PARENT" = "$ROOTWIN" ] && break
  FRAME="$PARENT"
done
# Taken from the screen: under compositing KWin draws the frame there, not into the frame window
geometry() { xwininfo -id "$FRAME" | sed -n "s/.*$1: *\([0-9-]*\).*/\1/p"; }
X0="$(geometry 'Absolute upper-left X')"
Y0="$(geometry 'Absolute upper-left Y')"
"${IMPORT[@]}" -window root -crop "$(geometry Width)x$(geometry Height)+${X0}+${Y0}" +repage "$TMP/window.png"

# A soft shadow around it, as a compositor would draw
mkdir -p "$(dirname "$OUT")"
"${IM[@]}" "$TMP/window.png" \( +clone -background black -shadow 45x16+0+8 \) +swap \
  -background none -layers merge +repage "$OUT"
echo "Screenshot done: ${OUT}"
