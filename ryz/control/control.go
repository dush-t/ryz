package control

import (
	"log"

	"github.com/dush-t/ryz/ryz/core"

	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// SimpleControl implements the Control interface.
// I could not think of a better name
type SimpleControl struct {
	Client             core.P4RClient
	ArbitrationDone    bool
	IsMaster           bool
	DigestChannel      chan *p4V1.StreamMessageResponse_Digest
	ArbitrationChannel chan *p4V1.StreamMessageResponse_Arbitration
}

// StartMessageRouter will start a goroutine that takes incoming messages from the stream
// and then sends them to corresponding channels based on message types
func (sc *SimpleControl) StartMessageRouter() {
	IncomingMessageChannel := sc.Client.GetMessageChannels().IncomingMessageChannel
	go func() {
		for {
			in := <-IncomingMessageChannel
			update := in.GetUpdate()

			switch update.(type) {
			case *p4V1.StreamMessageResponse_Arbitration:
				sc.ArbitrationChannel <- update.(*p4V1.StreamMessageResponse_Arbitration)
			case *p4V1.StreamMessageResponse_Digest:
				sc.DigestChannel <- update.(*p4V1.StreamMessageResponse_Digest)
			default:
				log.Println("Message has unknown type")
			}
		}
	}()
}
