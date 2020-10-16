package control

import (
	"github.com/dush-t/ryz/core/entities"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// DigestControl contains all the information needed to handle
// DigestEntries in the controller
type DigestControl struct {
	control *SimpleControl
	digest  *entities.DigestEntry
}

func (dc *DigestControl) getDigestEntryConfig(maxListSize int32, maxTimeoutNs, ackTimeoutNs int64) *p4V1.DigestEntry {
	return &p4V1.DigestEntry{
		DigestId: dc.digest.ID,
		Config: &p4V1.DigestEntry_Config{
			MaxTimeoutNs: maxTimeoutNs,
			MaxListSize:  maxListSize,
			AckTimeoutNs: ackTimeoutNs,
		},
	}
}

// Insert will insert a DigestEntry in the switch
func (dc DigestControl) Insert(maxListSize int32, maxTimeoutNs, ackTimeoutNs int64) error {
	entry := dc.getDigestEntryConfig(maxListSize, maxTimeoutNs, ackTimeoutNs)
	update := dc.digest.ConfigureMessage(entry)
	return dc.control.Client.WriteUpdate(update)
}

// Modify will modify an existing DigestEntry on the switch
func (dc DigestControl) Modify(maxListSize int32, maxTimeoutNs, ackTimeoutNs int64) error {
	entry := dc.getDigestEntryConfig(maxListSize, maxTimeoutNs, ackTimeoutNs)
	update := dc.digest.ModifyMessage(entry)
	return dc.control.Client.WriteUpdate(update)
}

// Delete will delete the selected DigestEntry from the switch
// i.e. no digest messages corresponding to that DigestEntry will
// be received.
func (dc DigestControl) Delete() error {
	update := dc.digest.DeleteMessage()
	return dc.control.Client.WriteUpdate(update)

}

// Acknowledge will acknowledge that the control plane did receive
// a DigestList of a given ID
func (dc DigestControl) Acknowledge(digestList *p4V1.DigestList) {
	message := dc.digest.AcknowledgeMessage(digestList)
	reqChannel := dc.control.Client.GetMessageChannels().OutgoingMessageChannel
	reqChannel <- message
}
