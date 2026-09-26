#!/bin/bash
# Renders forms/icons/*.svg to 32x32 PNGs using headless Chrome
# (ImageMagick's built-in SVG renderer drops strokes).
set -e
cd "$(dirname "$0")/../forms/icons"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
for f in *.svg; do
    n=${f%.svg}
    printf '<html><body style="margin:0;background:transparent"><img src="file://%s" width="32" height="32" style="display:block"></body></html>' "$PWD/$f" > "$TMP/$n.html"
    google-chrome --headless=new --disable-gpu --hide-scrollbars --default-background-color=00000000 \
        --force-device-scale-factor=1 --window-size=32,32 --screenshot="$PWD/$n.png" "file://$TMP/$n.html" >/dev/null 2>&1
done
