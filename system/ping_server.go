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
)

type PingServer struct {
	mtx             sync.Mutex
	srv             *icmp.PacketConn
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
	if useUdpSocket {
		c.srv, err = icmp.ListenPacket("udp4", "0.0.0.0")
		mode = "udp"
	} else {
		c.srv, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		mode = "icmp"
	}
	if err != nil {
		fmt.Println("Err", err)
		c.mode = ""
		return
	}
	c.mode = mode
	go c.thReceive()
}

func (c *PingServer) Stop() error {
	if c.srv != nil {
		return c.srv.Close()
	}
	return nil
}

func (c *PingServer) thReceive() {
	var err error
	rb := make([]byte, 1500)

	for {
		var n int
		var peer net.Addr
		n, peer, err = c.srv.ReadFrom(rb)
		if err != nil {
			break
		}
		var rm *icmp.Message
		rm, err = icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
		if err != nil {
			continue
		}

		if rm.Type != ipv4.ICMPTypeEchoReply {
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

	c.mode = ""
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

	var IPs []net.IP
	IPs, err = net.LookupIP(addr)
	if err != nil {
		return
	}

	var ipAddr net.IP

	for _, ip := range IPs {
		ipv4 := ip.To4()
		if ipv4 != nil {
			ipAddr = ipv4
		}
	}

	if len(ipAddr) == 0 {
		err = errors.New("cannot lookup address")
		return
	}

	source := c.source
	sequenceNum := c.getNextSequenceNum()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
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

	switch c.mode {
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

	if _, err = c.srv.WriteTo(wb, destAddr); err != nil {
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
