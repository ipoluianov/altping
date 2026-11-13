# Changelog

## Unreleased

Changes since v0.0.15.

### Added

- Windows: a downloaded `altping.exe` installs itself into
  `%USERPROFILE%\.altbins` with the **Install** link in the status bar. It
  adds shortcuts to the Start menu and the desktop and an entry to
  "Installed apps".
- Windows: the link reads **Update** when an older version is installed,
  and **Uninstall** when this version or a newer one is installed.
- Linux: `.deb` and `.rpm` packages.
- Linux: `linux-install.sh` (replaces `linux-x64-install.sh`) installs the
  x86-64 or ARM64 build, whichever fits the machine.
- macOS: a `.dmg` with AltPing.app, signed and notarized, so it opens
  without the Gatekeeper warning.
- Only one copy runs at a time. Starting AltPing again brings the running
  window to the front, even from the tray.
- Styles: 18 looks for the whole window, picked with the palette button in
  the toolbar (or **S**) from a list of small previews. Besides the
  classic dark and light themes there are dark, light, monospaced and high
  contrast styles.
- Two more views of the hosts, chosen by the style: **tiles** with a big
  ping time, a signal bar and the last 5 minutes (Gamer HUD, NOC Wall,
  Bubblegum), and **heartbeat lanes**, where each host is the trace of a
  heart monitor drawn from its last 30 seconds of pings (Heartbeat).
- Command Center and Arctic show cards with the totals above the hosts:
  how many are up, slow and down, the median ping and the availability
  over the last 24 hours.
- The chart panel shows the time, average, loss and jitter of the
  selected host as cards above the chart.

### Changed

- New icon: white arrows on a blue tile, readable at 16×16.
- Windows: the release file is `altping.exe`, without a platform suffix.

### Removed

- Windows ARM64 and macOS Intel builds.

## 0.0.15 and earlier

The first public builds: ICMP and TCP checks, statistics, 24-hour
history with charts and CSV export, alerts, public pages on u00.io, host
lists, tray, themes and translations.
