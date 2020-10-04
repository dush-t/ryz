package core

import (
	p4ConfigV1 "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// EntityClient defines any client that can interact with P4 switch entities such
// as tables, actions, counters, etc
type EntityClient interface {
	GetEntities(EntityType) *map[string]Entity
}

// P4RClient represents a p4Runtime client. Most methods are just getters since Go's
// interface implementation does not allow non-function members
type P4RClient interface {
	EntityClient
	// To initialize the client
	Init(addr string, p4Info *p4ConfigV1.P4Info, deviceID uint64, electionID p4V1.Uint128) error

	// Run will do whatever is needed to ensure that the client is active
	// once it is initialized.
	Run()

	// GetMessageChannels will return the message channels used by the client
	GetMessageChannels() MessageChannels

	// GetArbitrationData will return the data required to perform arbitration
	// for the client
	GetArbitrationData() ArbitrationData

	// GetStreamChannel will return the StreamChannel instance associated with the client
	GetStreamChannel() p4V1.P4Runtime_StreamChannelClient

	// P4Info will return the P4Info struct associated to the client
	P4Info() *p4ConfigV1.P4Info

	// IsMaster returns true if the client is master
	IsMaster() bool

	// SetMastershipStatus sets the mastership status of the client
	SetMastershipStatus(bool)

	// WriteUpdate is used to update an entity on the switch. Refer to the P4Runtime spec to know more.
	WriteUpdate(update *p4V1.Update) error
}
