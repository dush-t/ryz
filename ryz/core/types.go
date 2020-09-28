package core

import (
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// MessageChannels contains the incoming and outgoing message
// channels of a client
type MessageChannels struct {
	IncomingMessageChannel chan *p4V1.StreamMessageResponse
	OutgoingMessageChannel chan *p4V1.StreamMessageRequest
}

// ArbitrationData contains information required to perform
// stream arbitration
type ArbitrationData struct {
	DeviceID   uint64
	ElectionID p4V1.Uint128
}
