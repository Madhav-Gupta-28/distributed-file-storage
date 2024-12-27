package peer2peer

import "net"

// RPC represents any arbitary data that is being send over each transport between two nodes in the network
type RPC struct {
	from    net.Addr
	Payload []byte
}
