# AltPing

A desktop monitor for host availability. Add hosts to a list and AltPing
pings them all at once. For each host it shows the response time, packet
loss and a trend line, and it keeps a day of history. It runs on Linux,
Windows and macOS as a single binary, and on Linux and macOS it works
without root rights.

Website and documentation: https://altbins.pro/altping/

![AltPing main window: a list of hosts with response times, a 5-minute trend, loss/jitter statistics and the ping history chart of the selected host](https://github.com/ipoluianov/altping/releases/latest/download/screenshot.png)

## Features

- **Checks.** ICMP ping over IPv4 and IPv6, or a TCP connect to a port.
  The interval, timeout and "slow" threshold are set per host.
- **Statistics.** Loss, min/avg/max and jitter over the last 5 minutes,
  plus how long a host has been in its current state. Click a column
  header to sort by it.
- **History.** The last 24 hours are saved to disk, so they are still
  there after a restart. The chart shows 5 minutes, 1 hour or 24 hours
  for up to 4 hosts at once. The outage list shows when each host was
  down, and the history can be exported to CSV.
- **Alerts.** A host can be set to beep when it goes down or comes back.
  The taskbar button then flashes, and the window title shows how many
  hosts are down.
- **Public page.** One click sends a host's ping times to a public page
  with a live chart on [u00.io](https://u00.io).
- **Host lists.** Keep several lists, for example "Office" and "Servers",
  and change settings for many hosts at once.
- **Interface.** System tray, always on top, keyboard shortcuts, light
  and dark themes, 12 languages.
- **One copy at a time.** Starting AltPing again brings the running window
  to the front, even from the tray.

## Installation

Download a build for your system from the
[releases](https://github.com/ipoluianov/altping/releases):

- **Windows:** `altping.exe`.
- **macOS** (Apple Silicon): the `.dmg` image.
- **Linux** (x86-64 and ARM64): a `.deb` or `.rpm` package, or a user-level
  install (no root) that picks the build for your machine:

  ```sh
  curl -fsSL https://github.com/ipoluianov/altping/releases/latest/download/linux-install.sh | bash
  ```

Settings, host lists and history are stored in `~/.altbins/.altping/`.
Removing the application leaves them in place.

### Windows

The exe is not signed, so on the first start SmartScreen may say "Windows
protected your PC". Click **More info**, then **Run anyway**.

A downloaded copy offers to install itself into `%USERPROFILE%\.altbins`
with the **Install** link in the status bar. It adds AltPing to the Start
menu, the desktop and "Installed apps". When a newer version is started
from a download, the link reads **Update**. Remove AltPing from "Installed
apps" or with the **Uninstall** link.

### macOS

Open the `.dmg` and drag AltPing to Applications. The app is not notarized,
so macOS blocks the first start. Open **System Settings → Privacy &
Security**, find the message about AltPing and click **Open Anyway**.
Or remove the quarantine flag in Terminal:

```sh
xattr -dr com.apple.quarantine /Applications/AltPing.app
```

### Linux

- **Ping without root.** AltPing uses the unprivileged ICMP sockets of the
  kernel. They are allowed by default on systemd 243 and later (Ubuntu
  20.04+, Debian 11+, Fedora, Arch). If ping fails for every host and the
  status bar shows no mode, allow them for all users:

  ```sh
  echo 'net.ipv4.ping_group_range = 0 2147483647' | sudo tee /etc/sysctl.d/99-ping.conf
  sudo sysctl --system
  ```

  TCP checks work in any case.
- **System tray.** AltPing shows a StatusNotifierItem icon. KDE, XFCE and
  Cinnamon show it as is; GNOME needs the
  [AppIndicator](https://extensions.gnome.org/extension/615/appindicator-support/)
  extension. Without a tray, "hide to tray" minimizes the window.
- **Display.** AltPing needs X11, and runs under Wayland through XWayland,
  which GNOME and KDE have. glibc-based distributions only (not Alpine).
- **Uninstall.** Packages are removed by the package manager. A user-level
  install adds **Uninstall** to the menu entry (right-click) and leaves
  `~/.altbins/.altping-uninstall.sh`.

## Building

You need Go (see the version in [go.mod](go.mod)). No cgo is required.

```sh
go build .
```

`build/all.sh` (or `build/all.bat`) builds every platform into `bin/`,
together with the deb/rpm packages, the dmg and the installer script.

## License

[MIT](LICENSE) © Ivan Poluianov
