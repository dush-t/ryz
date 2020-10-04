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

// Table will return a TableControl struct
func (sc *SimpleControl) Table(tableName string) TableControl {
	tables := *(sc.Client.GetEntities(core.EntityTypes.TABLE))
	table := tables[tableName].(*core.Table)

	return TableControl{
		table:   table,
		control: sc,
	}
}

// SetMastershipStatus will call a method of the same name on P4RClient.
// We need to keep track of mastership to reason about which control can be
// used for what.
func (sc *SimpleControl) SetMastershipStatus(status bool) {
	sc.Client.SetMastershipStatus(status)
}

// IsMaster will return true if the control has mastership
func (sc *SimpleControl) IsMaster() bool {
	return sc.Client.IsMaster()
}
