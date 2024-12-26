package peer2peer

import "net"

// Message represents any arbitary data that is being send over each transport between two nodes in the network
type Message struct {
	from    net.Addr
	Payload []byte
}
