#!/bin/bash
# usage: render2.sh [name-filter...]
cd "$(dirname "$0")" && python3 gen2.py "$@" | while read k; do
 google-chrome --headless=new --disable-gpu --no-sandbox --hide-scrollbars --force-device-scale-factor=2 --window-size=1600,960 --screenshot="$PWD/r2/$k.png" "file://$PWD/r2/$k.html" >/dev/null 2>&1; echo "$k"; done
