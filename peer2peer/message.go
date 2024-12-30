package peer2peer

// RPC represents a remote procedure call
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}

const IncomingMessage = 0x1
const IncomingStream = 0x2
