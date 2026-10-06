package system

// Model data for the README screenshot: built in by build/_screenshot.sh with
// go build -overlay, never part of the product (the _ directory keeps it out
// of ./...). On start it writes a model config, settings and 24 hours of
// history into ConfigDirectory (run it with HOME set to a scratch directory),
// and the hosts "ping" by the same model instead of the network.

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"math"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ipoluianov/altping/config"
)

type demoProfile struct {
	id, name, hostname string
	port               int
	ip                 string
	base, jitter       float64 // ms
	loss               float64 // probability of a lost ping
	slowMs             int
	notify             bool
	// outages: [from, to) before now; to == 0 - still down
	outages [][2]time.Duration
}

var demoProfiles = []demoProfile{
	{id: "d1a0000000000001", name: "Router", hostname: "192.168.1.1", ip: "192.168.1.1", base: 0.35, jitter: 0.12},
	{id: "d1a0000000000002", name: "Core switch", hostname: "192.168.1.2", ip: "192.168.1.2", base: 0.5, jitter: 0.15},
	{id: "d1a0000000000003", name: "NAS", hostname: "nas.office.lan", ip: "192.168.1.20", base: 0.7, jitter: 0.3},
	{id: "d1a0000000000004", name: "Printer, 2nd floor", hostname: "192.168.1.45", ip: "192.168.1.45", base: 1.6, jitter: 0.8, notify: true,
		outages: [][2]time.Duration{{14*time.Minute + 32*time.Second, 0}, {7 * time.Hour, 6*time.Hour + 41*time.Minute}}},
	{id: "d1a0000000000005", name: "Wi-Fi AP, meeting room", hostname: "192.168.1.61", ip: "192.168.1.61", base: 3.2, jitter: 2.6, loss: 0.007},
	{id: "d1a0000000000006", name: "Cloudflare DNS", hostname: "1.1.1.1", ip: "1.1.1.1", base: 8.4, jitter: 1.1},
	{id: "d1a0000000000007", name: "google.com", hostname: "google.com", ip: "142.250.185.78", base: 13, jitter: 1.6},
	{id: "d1a0000000000008", name: "github.com", hostname: "github.com", ip: "140.82.121.4", base: 36, jitter: 3},
	{id: "d1a0000000000009", name: "VPN gateway", hostname: "vpn.example.com", ip: "203.0.113.10", base: 24, jitter: 4, loss: 0.002, notify: true,
		outages: [][2]time.Duration{{5*time.Hour + 12*time.Minute, 5*time.Hour + 6*time.Minute}}},
	{id: "d1a000000000000a", name: "Web shop", hostname: "shop.example.com", port: 443, ip: "203.0.113.25", base: 58, jitter: 9, slowMs: 50, notify: true},
	{id: "d1a000000000000b", name: "Mail", hostname: "mail.example.com", ip: "203.0.113.40", base: 29, jitter: 2.5},
	{id: "d1a000000000000c", name: "Backup, Singapore", hostname: "backup-sg.example.com", ip: "198.51.100.7", base: 171, jitter: 12, loss: 0.01},
}

var demoStart = time.Now()

func demoProfileOf(id string) *demoProfile {
	for i := range demoProfiles {
		if demoProfiles[i].id == id {
			return &demoProfiles[i]
		}
	}
	return nil
}

// demoRand returns a stable pseudo-random number in [0, 1) for the host, second and salt
func demoRand(id string, t time.Time, salt int) float64 {
	h := fnv.New64a()
	h.Write([]byte(id))
	h.Write([]byte(strconv.FormatInt(t.Unix(), 10)))
	h.Write([]byte{byte(salt)})
	// FNV barely mixes the last bytes into the high bits: finish with splitmix64
	x := h.Sum64()
	x ^= x >> 30
	x *= 0xbf58476d1ce4e5b9
	x ^= x >> 27
	x *= 0x94d049bb133111eb
	x ^= x >> 31
	return float64(x>>11) / float64(1<<53)
}

// demoSample is the model ping of the host at the time
func demoSample(p *demoProfile, t time.Time) (time.Duration, bool) {
	ago := demoStart.Sub(t)
	for _, o := range p.outages {
		if ago < o[0] && (o[1] == 0 || ago >= o[1]) {
			return 0, false
		}
	}
	// Congestion on the way out: a few minutes now and then with more jitter and losses
	jitter, loss := p.jitter, p.loss
	if p.base > 5 && demoRand(p.id, demoSlot(t, 3*time.Minute), 5) < 0.07 {
		jitter *= 3.5
		loss = loss*4 + 0.004
	}
	if demoRand(p.id, t, 1) < loss {
		return 0, false
	}
	// The route changes now and then: the level steps up for a while and back
	base := p.base
	if r := demoRand(p.id, demoSlot(t, 25*time.Minute), 6); r < 0.3 {
		base *= 1.04 + 0.5*r
	}
	// The time is never below the level, so the noise only adds, with a long tail; rare spikes
	ms := base + -math.Log(1-demoRand(p.id, t, 2))*jitter*0.6
	if demoRand(p.id, t, 3) < 0.004 {
		ms += base * (0.4 + 1.2*demoRand(p.id, t, 4))
	}
	return time.Duration(ms * float64(time.Millisecond)), true
}

// demoSlot is the start of the slot of the length the time is in: the same
// random numbers for the whole slot
func demoSlot(t time.Time, length time.Duration) time.Time {
	return t.Truncate(length)
}

func demoPing(c *Host, port string) (time.Duration, net.Addr, error) {
	p := demoProfileOf(c.ID)
	if p == nil {
		return 0, nil, errors.New("timeout")
	}
	now := time.Now()
	d, ok := demoSample(p, now)
	if !ok {
		select {
		case <-time.After(c.configHost.Timeout()):
		case <-c.chanStop:
		}
		return 0, nil, errors.New("timeout")
	}
	select {
	case <-time.After(d):
	case <-c.chanStop:
	}
	ip := net.ParseIP(p.ip)
	if port != "" {
		return d, &net.TCPAddr{IP: ip}, nil
	}
	return d, &net.IPAddr{IP: ip}, nil
}

func (c *Host) demoCheckIP() bool {
	if p := demoProfileOf(c.ID); p != nil {
		c.mtx.Lock()
		c.IP = p.ip
		c.mtx.Unlock()
	}
	return true
}

func demoWriteJSON(name string, v any) {
	bs, _ := json.MarshalIndent(v, "", "  ")
	os.WriteFile(filepath.Join(config.ConfigDirectory(), name), bs, 0644)
}

func init() {
	dir := config.ConfigDirectory()
	os.MkdirAll(filepath.Join(dir, "history"), 0755)

	cfg := config.Config{ID: "d1a00000000000ff", Name: "Office"}
	for _, p := range demoProfiles {
		cfg.Hosts = append(cfg.Hosts, &config.ConfigHost{
			ID: p.id, DisplayName: p.name, Hostname: p.hostname, Port: p.port,
			SlowMs: p.slowMs, Notify: p.notify,
		})
	}
	demoWriteJSON(cfg.ID+".ws", cfg)
	os.WriteFile(filepath.Join(dir, "last_config_id.txt"), []byte(cfg.ID), 0644)
	if _, err := os.Stat(filepath.Join(dir, "settings.json")); err != nil {
		demoWriteJSON("settings.json", config.Settings{ShowMin: true, ShowJitter: true, ShowSince: true, Language: "en"})
	}
	if _, err := os.Stat(filepath.Join(dir, "window.json")); err != nil {
		demoWriteJSON("window.json", config.WindowState{X: 40, Y: 60, Width: 1840, Height: 900, DetailsVisible: true, DetailsWidth: 520})
	}

	// History up to now, one sample a second, as the hosts would have written it
	from := demoStart.Add(-historyDepth + time.Minute).Truncate(time.Second)
	for _, p := range demoProfiles {
		p := p
		buf := make([]byte, 0, int(historyDepth/time.Second)*historyRecordSize)
		// The program was restarted in the morning: a break in the data
		restart := demoStart.Add(-14*time.Hour - 12*time.Minute).Truncate(time.Second)
		for t := from; t.Before(demoStart); t = t.Add(time.Second) {
			if t.Equal(restart) {
				buf = append(buf, encodeHistorySample(HistorySample{DT: t, Gap: true})...)
				t = t.Add(40 * time.Second)
			}
			d, ok := demoSample(&p, t)
			buf = append(buf, encodeHistorySample(HistorySample{DT: t, PingTime: d, OK: ok})...)
		}
		os.WriteFile(historyFilePath(p.id), buf, 0644)
	}
}
