package entities

import (
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

type DigestEntry struct {
	ID uint32
}

// ConfigureMessage will return an update message that, when sent to the switch, will insert the
// DigestEntry configuration accordingly.
func (d *DigestEntry) ConfigureMessage(entry *p4V1.DigestEntry) *p4V1.Update {
	update := &p4V1.Update{
		Type: p4V1.Update_INSERT,
		Entity: &p4V1.Entity{
			Entity: &p4V1.Entity_DigestEntry{entry},
		},
	}

	return update
}

// ModifyMessage will return an update message that, when sent to the switch, will update the
// DigestEntry configuration accordingly
func (d *DigestEntry) ModifyMessage(entry *p4V1.DigestEntry) *p4V1.Update {
	update := &p4V1.Update{
		Type: p4V1.Update_MODIFY,
		Entity: &p4V1.Entity{
			Entity: &p4V1.Entity_DigestEntry{entry},
		},
	}

	return update
}

// DeleteMessage will return an update message that, when sent to the switch, will delete the
// corresponding DigestEntry item
func (d *DigestEntry) DeleteMessage() *p4V1.Update {
	entry := &p4V1.DigestEntry{
		DigestId: d.ID,
	}
	update := &p4V1.Update{
		Type: p4V1.Update_DELETE,
		Entity: &p4V1.Entity{
			Entity: &p4V1.Entity_DigestEntry{entry},
		},
	}

	return update
}

// AcknowledgeMessage will return a StreamMessageRequest message that, when sent to the switch, will
// acknowledge the reception of the digestList passed as the argument
func (d *DigestEntry) AcknowledgeMessage(digestList *p4V1.DigestList) *p4V1.StreamMessageRequest {
	return &p4V1.StreamMessageRequest{
		Update: &p4V1.StreamMessageRequest_DigestAck{&p4V1.DigestListAck{
			DigestId: d.ID,
			ListId:   digestList.ListId,
		}},
	}
}

// Type returns the type of the Entity represented by DigestEntry,
// so that it can implement the Entity interface
func (d *DigestEntry) Type() EntityType {
	return EntityTypes.DIGEST
}

// GetID will return the ID of the DigestEntry Entity
func (d *DigestEntry) GetID() uint32 {
	return d.ID
}
