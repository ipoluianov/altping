package system

import (
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type PingServer struct {
	srv          *icmp.PacketConn
	useUdpSocket bool

	activeRequests map[string]*PingRequest
}

type PingRequest struct {
	Addr         string
	DataSize     int
	TimeoutMs    int
	UseUdpSocket bool
}

var server *PingServer

func init() {
	server, _ = NewPingServer(false)
	go server.thReceive()
}

func GetServer() *PingServer {
	return server
}

func NewPingServer(useUdpSocket bool) (*PingServer, error) {
	var c PingServer
	var err error
	c.activeRequests = make(map[string]*PingRequest)
	c.useUdpSocket = useUdpSocket
	if useUdpSocket {
		c.srv, err = icmp.ListenPacket("udp4", "0.0.0.0")
	} else {
		c.srv, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (ps *PingServer) Close() error {
	if ps.srv != nil {
		return ps.srv.Close()
	}
	return nil
}

func (c *PingServer) thReceive() {
	var err error
	for {
		rb := make([]byte, 1500)

		err = c.srv.SetDeadline(time.Now().Add(100 * time.Millisecond))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var n int
		var peer net.Addr
		n, peer, err = c.srv.ReadFrom(rb)
		_ = peer
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		var rm *icmp.Message
		rm, err = icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if rm.Type == ipv4.ICMPTypeDestinationUnreachable {
			err = errors.New("destination unreachable")
			fmt.Println("Error:", err)
			continue
		}

		if rm.Type != ipv4.ICMPTypeEchoReply {
			err = errors.New("error")
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Received ICMP Echo Reply from", peer)
	}

}

func (c *PingServer) PingHost(addr string, dataFrame []byte, timeoutMs int, source uint16, sequenceNum uint16, useUdpSocket bool) (result int, peer net.Addr, err error) {
	if len(dataFrame) < 1 || len(dataFrame) > 1400 {
		err = errors.New("wrong data frame length")
		return
	}

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
	if useUdpSocket {
		destAddr = &net.UDPAddr{IP: ipAddr}
	} else {
		destAddr = &net.IPAddr{IP: ipAddr}
	}

	var req *PingRequest
	req = &PingRequest{
		Addr:         addr,
		DataSize:     len(dataFrame),
		TimeoutMs:    timeoutMs,
		UseUdpSocket: useUdpSocket,
	}
	c.activeRequests[addr] = req

	if _, err = c.srv.WriteTo(wb, destAddr); err != nil {
		return
	}

	return
}
