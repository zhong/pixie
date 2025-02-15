package dhcp

import (
	"net"

	"github.com/p3lim/pixie/pkg/log"
)

// https://datatracker.ietf.org/doc/html/rfc2131#section-2
/*
	0                   1                   2                   3
	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|     op (1)    |   htype (1)   |   hlen (1)    |   hops (1)    |
	+---------------+---------------+---------------+---------------+
	|                            xid (4)                            |
	+-------------------------------+-------------------------------+
	|           secs (2)            |           flags (2)           |
	+-------------------------------+-------------------------------+
	|                          ciaddr  (4)                          |
	+---------------------------------------------------------------+
	|                          yiaddr  (4)                          |
	+---------------------------------------------------------------+
	|                          siaddr  (4)                          |
	+---------------------------------------------------------------+
	|                          giaddr  (4)                          |
	+---------------------------------------------------------------+
	|                                                               |
	|                          chaddr  (16)                         |
	|                                                               |
	|                                                               |
	+---------------------------------------------------------------+
	|                                                               |
	|                          sname   (64)                         |
	+---------------------------------------------------------------+
	|                                                               |
	|                          file    (128)                        |
	+---------------------------------------------------------------+
	|                                                               |
	|                          options (variable)                   |
	+---------------------------------------------------------------+
*/

type Message struct {
	options map[Option][]byte
	raw     []byte
}

func ParseMessage(msg []byte) (*Message, error) {
	m := &Message{
		raw:     msg,
		options: make(map[Option][]byte),
	}

	if err := m.parseOptions(msg[240:]); err != nil {
		return nil, err
	}

	return m, nil
}

// GetOP returns the "op code" from the DHCP Message.
// The value can be 1 (for BOOTREQUEST) or 2 (for BOOTREPLY).
func (m Message) GetOP() byte {
	return m.raw[0]
}

// GetHTYPE returns the hardware address type from the DHCP Message.
// The most common value is 1, for ethernet.
func (m Message) GetHTYPE() byte {
	return m.raw[1]
}

// GetHLEN returns the length of the hardware address from the DHCP Message.
// This value is typically 16, for ethernet MAC.
func (m Message) GetHLEN() byte {
	return m.raw[2]
}

// GetHOPS returns the "ops" field from the DHCP Message.
// This is set by the client, typically used by relay agents.
func (m Message) GetHOPS() byte {
	return m.raw[3]
}

// GetXID returns the transaction ID from the DHCP Message.
// This is a random number chosen by the client, used by both client and server to associate
// messages and responses between them.
func (m Message) GetXID() []byte {
	return m.raw[4:8]
}

// GetSECS returns the seconds elapsed since the client began the address/renewal process from the
// DHCP Message.
func (m Message) GetSECS() []byte {
	return m.raw[8:10]
}

// GetFLAGS returns additional flags from the DHCP Message.
// See Flags for more info.
func (m Message) GetFLAGS() Flags {
	return Flags(m.raw[10:12])
}

// GetCIADDR returns the client's IP address from the DHCP Message.
// This is only filled if the client is in BOUND, RENEW, or REBINDING state and can respond to ARP
// requests.
func (m Message) GetCIADDR() net.IP {
	return net.IP(m.raw[12:16])
}

// GetYIADDR returns the server's IP address from the DHCP Message.
func (m Message) GetYIADDR() net.IP {
	return net.IP(m.raw[16:20])
}

// GetSIADDR returns the next-server's IP address from the DHCP Message.
func (m Message) GetSIADDR() net.IP {
	return net.IP(m.raw[20:24])
}

// GetGIADDR returns the relay agent's IP address from the DHCP Message.
func (m Message) GetGIADDR() net.IP {
	return net.IP(m.raw[24:28])
}

// GetCHADDR returns the client's hardware address from the DHCP Message.
// This is typically the ethernet MAC.
func (m Message) GetCHADDR() net.HardwareAddr {
	return net.HardwareAddr(m.raw[28 : 28+m.GetHLEN()])
}

// GetSNAME returns the server's hostname from the DHCP Message.
func (m Message) GetSNAME() string {
	return string(m.raw[44:108])
}

// GetFILE returns the boot file name from the DHCP Message.
func (m Message) GetFILE() string {
	return string(m.raw[108:236])
}

// GetMagicCookie returns the "magic cookie" prefixed to the options field of the DHCP Message.
// This field always contains the (decimal) values 99, 130, 83, and 99.
func (m Message) GetMagicCookie() []byte {
	return m.raw[236:240]
}

func (m Message) DebugLog() {
	log.Debug("------------")
	log.Debugf("op: %v", m.GetOP())
	log.Debugf("htype: %v", m.GetHTYPE())
	log.Debugf("hlen: %v", m.GetHLEN())
	log.Debugf("hops: %v", m.GetHOPS())
	log.Debugf("xid: %v", m.GetXID())
	log.Debugf("secs: %v", m.GetSECS())
	log.Debugf("flag (broadcast): %v", m.GetFLAGS().Broadcast())
	log.Debugf("ciaddr: %v", m.GetCIADDR())
	log.Debugf("yiaddr: %v", m.GetYIADDR())
	log.Debugf("siaddr: %v", m.GetSIADDR())
	log.Debugf("giaddr: %v", m.GetGIADDR())
	log.Debugf("chaddr: %v", m.GetCHADDR())
	log.Debugf("sname: %v", m.GetSNAME())
	log.Debugf("file: %v", m.GetFILE())
	log.Debugf("cookie: %v", m.GetMagicCookie())

	for opt, value := range m.options {
		log.Debugf("option %d: %v", opt, value)
	}

	log.Debug("------------")
}
