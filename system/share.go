package system

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ipoluianov/altping/app"
	"github.com/ipoluianov/altping/config"
)

// sharer sends the ping results of the shared hosts to u00.io, one request
// at a time from its own goroutine, so a slow site never holds up pinging.
// u00.io takes at most 10 requests a second from an address: the requests
// are spaced by shareSendInterval, and a host that is pinged faster than its
// turn comes sends only its latest value.
type sharer struct {
	mtx     sync.Mutex
	pending map[string]shareItem // the value to send by API key
	order   []string             // the keys in the order their values came
	status  map[string]ShareStatus
	wake    chan struct{}
	client  *http.Client
}

// shareItem is what the page of a host shows: the value and its name
type shareItem struct {
	value string
	name  string
}

// ShareStatus is how the last value of a host went to u00.io
type ShareStatus struct {
	At  time.Time // when it was sent; zero - nothing was sent yet
	Err error     // why it failed, nil when it got there
}

// GetShareStatus returns how the last value of the API key went
func GetShareStatus(key string) ShareStatus {
	shareSender.mtx.Lock()
	defer shareSender.mtx.Unlock()
	return shareSender.status[key]
}

// shareUnit is the unit of the values, shown next to them
const shareUnit = "ms"

// shareSendInterval keeps the requests under the u00.io limit, with a margin
const shareSendInterval = time.Second / 8

var shareSender = newSharer()

func newSharer() *sharer {
	c := &sharer{
		pending: make(map[string]shareItem),
		status:  make(map[string]ShareStatus),
		wake:    make(chan struct{}, 1),
		client:  &http.Client{Timeout: 5 * time.Second},
	}
	go c.thSend()
	return c
}

// push queues the value of the key, replacing one not sent yet
func (c *sharer) push(key string, value shareItem) {
	c.mtx.Lock()
	if _, queued := c.pending[key]; !queued {
		c.order = append(c.order, key)
	}
	c.pending[key] = value
	c.mtx.Unlock()
	select {
	case c.wake <- struct{}{}:
	default:
	}
}

// next takes the oldest queued value, ok is false when there is none
func (c *sharer) next() (key string, value shareItem, ok bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if len(c.order) == 0 {
		return "", shareItem{}, false
	}
	key = c.order[0]
	c.order = c.order[1:]
	value = c.pending[key]
	delete(c.pending, key)
	return key, value, true
}

func (c *sharer) thSend() {
	for range c.wake {
		for {
			key, value, ok := c.next()
			if !ok {
				break
			}
			c.send(key, value)
			time.Sleep(shareSendInterval)
		}
	}
}

func (c *sharer) send(key string, item shareItem) {
	err := c.get(shareSetURL(key, item))
	if err != nil {
		fmt.Println("Share error:", err)
	}
	c.mtx.Lock()
	c.status[key] = ShareStatus{At: time.Now(), Err: err}
	c.mtx.Unlock()
}

func (c *sharer) get(url string) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

// shareSetURL makes the request that sets the value of the page:
// https://u00.io/set/API_KEY?/=VALUE&/_name=NAME&/_uom=UNIT&/provider=AltPing
// The parameter names are paths on the page: "/" is the value, "/_name" its name,
// "/provider" the program that sends it.
func shareSetURL(key string, item shareItem) string {
	u := config.ShareSite + "/set/" + key + "?/=" + url.QueryEscape(item.value)
	if item.name != "" {
		u += "&/_name=" + url.QueryEscape(item.name)
	}
	return u + "&/_uom=" + url.QueryEscape(shareUnit) + "&/provider=" + url.QueryEscape(app.DisplayName)
}

// shareValue is what a ping result looks like on the page: the time in ms,
// or the error text, which the chart shows as a break
func shareValue(pingTime time.Duration, err error) string {
	if err != nil {
		return err.Error()
	}
	return strconv.FormatFloat(float64(pingTime)/float64(time.Millisecond), 'f', 1, 64)
}
