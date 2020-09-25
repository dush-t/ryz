package core

import (
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// P4RClient represents a p4Runtime client
type P4RClient interface {
	// To initialize the client
	Init()

	// Used to request information from the stream channel
	RequestFromStream(p4V1.StreamMessageRequest)
}
