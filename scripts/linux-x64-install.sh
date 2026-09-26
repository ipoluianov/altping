#!/usr/bin/env bash
# Installer for __APP__ __TAG__ (linux/amd64)
# Usage: curl -fsSL https://github.com/__REPO__/releases/latest/download/linux-x64-install.sh | bash
set -euo pipefail

APP="__APP__"
TAG="__TAG__"
REPO="__REPO__"
URL="https://github.com/${REPO}/releases/download/${TAG}/${APP}-${TAG}-linux-amd64.tar.gz"

INSTALL_DIR="${HOME}/.altbins/${APP}"
APPS_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/applications"
DESKTOP_FILE="${APPS_DIR}/${APP}.desktop"

die() { echo "Error: $*" >&2; exit 1; }

[ "$(uname -s)" = "Linux" ] || die "this installer is for Linux only"
case "$(uname -m)" in
  x86_64|amd64) ;;
  *) die "unsupported architecture: $(uname -m) (expected x86_64)" ;;
esac

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "Downloading ${APP} ${TAG}..."
if command -v curl >/dev/null 2>&1; then
  curl -fL --progress-bar -o "$TMP/app.tar.gz" "$URL"
elif command -v wget >/dev/null 2>&1; then
  wget -q --show-progress -O "$TMP/app.tar.gz" "$URL"
else
  die "curl or wget is required"
fi

tar -xzf "$TMP/app.tar.gz" -C "$TMP"
[ -f "$TMP/$APP" ] || die "binary '$APP' not found in archive"

mkdir -p "$INSTALL_DIR" "$APPS_DIR"

# cp + mv: atomic replace, works even if the app is currently running
cp "$TMP/$APP" "$INSTALL_DIR/.$APP.new"
chmod 755 "$INSTALL_DIR/.$APP.new"
mv -f "$INSTALL_DIR/.$APP.new" "$INSTALL_DIR/$APP"

# Icon: use one from the archive if present, otherwise a theme icon
ICON="network-wired"
for ext in png svg; do
  if [ -f "$TMP/$APP.$ext" ]; then
    cp "$TMP/$APP.$ext" "$INSTALL_DIR/icon.$ext"
    ICON="$INSTALL_DIR/icon.$ext"
    break
  fi
done

# --- uninstall script ---
UNINSTALL="$INSTALL_DIR/uninstall.sh"
{
  echo '#!/usr/bin/env bash'
  printf 'APP=%q\n' "$APP"
  printf 'INSTALL_DIR=%q\n' "$INSTALL_DIR"
  printf 'APPS_DIR=%q\n' "$APPS_DIR"
  printf 'DESKTOP_FILE=%q\n' "$DESKTOP_FILE"
  cat <<'EOF'
set -u
MSG="Uninstall ${APP}?"

if [ -t 0 ]; then
  read -r -p "$MSG [y/N] " ans
  [[ "$ans" =~ ^[Yy] ]] || exit 0
elif command -v kdialog >/dev/null 2>&1; then
  kdialog --title "$APP" --yesno "$MSG" || exit 0
elif command -v zenity >/dev/null 2>&1; then
  zenity --question --title="$APP" --text="$MSG" || exit 0
fi

pkill -x "$APP" 2>/dev/null || true
rm -f "$DESKTOP_FILE"
rm -rf "$INSTALL_DIR"
rmdir "$(dirname "$INSTALL_DIR")" 2>/dev/null || true

command -v update-desktop-database >/dev/null 2>&1 && update-desktop-database "$APPS_DIR" 2>/dev/null
{ kbuildsycoca6 || kbuildsycoca5; } >/dev/null 2>&1 || true

if [ -t 1 ]; then
  echo "${APP} has been removed."
else
  command -v notify-send >/dev/null 2>&1 && notify-send "$APP" "${APP} has been removed."
fi
exit 0
EOF
} > "$UNINSTALL"
chmod 755 "$UNINSTALL"

# --- .desktop file (GNOME / KDE) ---
# Exec values are quoted per the Desktop Entry spec
esc() { printf '%s' "$1" | sed -e 's/[\\"`$]/\\&/g'; }

cat > "$DESKTOP_FILE" <<EOF
[Desktop Entry]
Type=Application
Name=${APP}
Comment=${APP} ${TAG}
Exec="$(esc "$INSTALL_DIR/$APP")"
Icon=${ICON}
Terminal=false
Categories=Network;Utility;
StartupNotify=false
Actions=Uninstall;

[Desktop Action Uninstall]
Name=Uninstall
Name[ru]=Удалить
Exec="$(esc "$UNINSTALL")"
Icon=edit-delete
EOF
chmod 644 "$DESKTOP_FILE"

command -v update-desktop-database >/dev/null 2>&1 && update-desktop-database "$APPS_DIR" 2>/dev/null || true
{ kbuildsycoca6 || kbuildsycoca5; } >/dev/null 2>&1 || true

echo
echo "Installed ${APP} ${TAG} to ${INSTALL_DIR}"
echo "Menu entry: ${DESKTOP_FILE}"
echo "Uninstall:  ${UNINSTALL}  (or right-click the app in the menu -> Uninstall)"
