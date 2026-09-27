package system

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type PingServer struct {
	mtx             sync.Mutex
	srv             *icmp.PacketConn
	srv6            *icmp.PacketConn // nil when the system has no IPv6
	source          uint16
	nextSequenceNum uint16
	mode            string

	activeRequestsByUniqueId map[string]*PingRequest // key = hex-encoded unique key
}

type PingRequest struct {
	Source   uint16
	Sequence uint16

	Addr      string
	DataSize  int
	TimeoutMs int
	SentTime  time.Time
	RecvTime  time.Time

	ResultPeer     net.Addr
	ResultResponse *icmp.Message
	ResultErr      error
}

func init() {
}

func NewPingServer() *PingServer {
	var c PingServer
	c.activeRequestsByUniqueId = make(map[string]*PingRequest)
	c.source = uint16(os.Getpid() & 0xFFFF)
	c.nextSequenceNum = 1
	return &c
}

func (c *PingServer) Mode() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.mode
}

func (c *PingServer) getNextSequenceNum() uint16 {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	seq := c.nextSequenceNum
	c.nextSequenceNum++
	if c.nextSequenceNum == 0xFFFF {
		c.nextSequenceNum = 1
	}
	return seq
}

func (c *PingServer) Start() {
	var err error

	useUdpSocket := false
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		weDontHaveRoot := os.Geteuid() != 0
		if weDontHaveRoot {
			useUdpSocket = true
		}
	}

	mode := ""
	network4, network6 := "ip4:icmp", "ip6:ipv6-icmp"
	if useUdpSocket {
		network4, network6 = "udp4", "udp6"
		mode = "udp"
	} else {
		mode = "icmp"
	}
	c.srv, err = icmp.ListenPacket(network4, "0.0.0.0")
	if err != nil {
		fmt.Println("Err", err)
		c.setMode("")
		return
	}
	// IPv6 is optional: without it only IPv6 hosts fail
	c.srv6, err = icmp.ListenPacket(network6, "::")
	if err != nil {
		fmt.Println("IPv6 ping is not available:", err)
		c.srv6 = nil
	}
	c.setMode(mode)
	go c.thReceive(c.srv, ipv4.ICMPTypeEchoReply.Protocol(), ipv4.ICMPTypeEchoReply)
	if c.srv6 != nil {
		go c.thReceive(c.srv6, ipv6.ICMPTypeEchoReply.Protocol(), ipv6.ICMPTypeEchoReply)
	}
}

// setMode is locked: the hosts read the mode while the server is started and stopped
func (c *PingServer) setMode(mode string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.mode = mode
}

func (c *PingServer) Stop() error {
	c.setMode("")
	if c.srv6 != nil {
		c.srv6.Close()
	}
	if c.srv != nil {
		return c.srv.Close()
	}
	return nil
}

// thReceive reads the replies until srv is closed. It gets its own socket,
// as after Stop and Start c.srv is already the new one.
func (c *PingServer) thReceive(srv *icmp.PacketConn, protocol int, replyType icmp.Type) {
	var err error
	rb := make([]byte, 1500)

	for {
		var n int
		var peer net.Addr
		n, peer, err = srv.ReadFrom(rb)
		if err != nil {
			break
		}
		var rm *icmp.Message
		rm, err = icmp.ParseMessage(protocol, rb[:n])
		if err != nil {
			continue
		}

		if rm.Type != replyType {
			continue
		}

		var echo *icmp.Echo
		echo, ok := rm.Body.(*icmp.Echo)
		if !ok {
			err = errors.New("error")
			continue
		}

		if echo == nil {
			err = errors.New("error")
			continue
		}

		// get first 8 bytes of data as unique key
		if len(echo.Data) < 8 {
			continue
		}
		uniqueKey := echo.Data[0:8]
		key := hex.EncodeToString(uniqueKey)

		c.mtx.Lock()
		req, ok := c.activeRequestsByUniqueId[key]
		if ok {
			req.ResultResponse = rm
			req.RecvTime = time.Now()
			req.ResultPeer = peer
		}
		c.mtx.Unlock()
	}
}

func (c *PingServer) PingHost(addr string, frameSize int, timeoutMs int, chanStop chan struct{}) (result time.Duration, peer net.Addr, err error) {
	if frameSize < 8 || frameSize > 1400 {
		err = errors.New("wrong data frame length")
		return
	}

	dataFrame := make([]byte, frameSize)

	uniqueKey := make([]byte, 8)
	_, err = rand.Read(uniqueKey)
	if err != nil {
		return
	}

	copy(dataFrame[0:8], uniqueKey)

	if len(addr) < 1 {
		err = errors.New("wrong address")
		return
	}

	if timeoutMs < 1 {
		err = errors.New("wrong timeout")
		return
	}

	if timeoutMs > 10000 {
		err = errors.New("wrong timeout")
		return
	}

	// The host passes the address it has resolved; a name is resolved here, IPv4 first
	ipAddr := net.ParseIP(addr)
	if ipAddr == nil {
		var IPs []net.IP
		IPs, err = net.LookupIP(addr)
		if err != nil {
			return
		}
		ipAddr = pickIP(IPs)
	}
	if ipAddr == nil {
		err = errors.New("cannot lookup address")
		return
	}
	isIPv6 := ipAddr.To4() == nil
	srv := c.srv
	var echoType icmp.Type = ipv4.ICMPTypeEcho
	if isIPv6 {
		srv = c.srv6
		echoType = ipv6.ICMPTypeEchoRequest
		if srv == nil {
			err = errors.New("IPv6 is not available")
			return
		}
	}

	source := c.source
	sequenceNum := c.getNextSequenceNum()

	wm := icmp.Message{
		Type: echoType, Code: 0,
		Body: &icmp.Echo{
			ID:   int(source),
			Seq:  int(sequenceNum),
			Data: dataFrame,
		},
	}
	var wb []byte
	wb, err = wm.Marshal(nil)
	if err != nil {
		return
	}

	var destAddr net.Addr

	switch c.Mode() {
	case "udp":
		destAddr = &net.UDPAddr{IP: ipAddr}
	case "icmp":
		destAddr = &net.IPAddr{IP: ipAddr}
	default:
		err = errors.New("ping server not started")
		return
	}

	startTime := time.Now()

	var req *PingRequest
	req = &PingRequest{
		Source:    source,
		Sequence:  sequenceNum,
		Addr:      addr,
		DataSize:  len(dataFrame),
		TimeoutMs: timeoutMs,
		SentTime:  startTime,
	}
	c.mtx.Lock()
	uniqueKeyStr := hex.EncodeToString(uniqueKey)
	c.activeRequestsByUniqueId[uniqueKeyStr] = req
	c.mtx.Unlock()

	defer func() {
		c.mtx.Lock()
		delete(c.activeRequestsByUniqueId, uniqueKeyStr)
		c.mtx.Unlock()
	}()

	if _, err = srv.WriteTo(wb, destAddr); err != nil {
		return
	}

	var response *icmp.Message

	for {
		c.mtx.Lock()
		if req.ResultResponse != nil {
			response = req.ResultResponse
			c.mtx.Unlock()
			break
		}
		c.mtx.Unlock()
		if time.Since(startTime) > time.Duration(timeoutMs)*time.Millisecond {
			err = errors.New("timeout")
			return
		}
		time.Sleep(10 * time.Millisecond)

		select {
		case <-chanStop:
			err = errors.New("stopped")
			return
		default:
		}
	}

	if response == nil {
		err = errors.New("no response")
		return
	}

	result = req.RecvTime.Sub(req.SentTime)
	peer = req.ResultPeer
	return
}

// pickIP returns the first IPv4 address, or the first IPv6 one if there is no IPv4
func pickIP(ips []net.IP) net.IP {
	for _, ip := range ips {
		if ip4 := ip.To4(); ip4 != nil {
			return ip4
		}
	}
	if len(ips) > 0 {
		return ips[0]
	}
	return nil
}
