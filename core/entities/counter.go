package entities

import (
	p4ConfigV1 "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// Counter stores all the information we need about a counter
type Counter struct {
	ID   uint32
	Size int64
}

// ReadValueAtIndexRequest returns a request entity that can be used to
// read a counter value at a certain index
func (c *Counter) ReadValueAtIndexRequest(index int64) *p4V1.Entity {
	entry := &p4V1.CounterEntry{
		CounterId: c.ID,
		Index:     &p4V1.Index{Index: index},
	}
	entity := &p4V1.Entity{
		Entity: &p4V1.Entity_CounterEntry{entry},
	}
	return entity
}

// ReadWildcardRequest returns a request entity that can be used to read
// counter value at all indices
func (c *Counter) ReadWildcardRequest() *p4V1.Entity {
	entry := &p4V1.CounterEntry{
		CounterId: c.ID,
	}
	entity := &p4V1.Entity{
		Entity: &p4V1.Entity_CounterEntry{entry},
	}
	return entity
}

// Type returns the type of the Entity represented by DigestEntry,
// so that it can implement the Entity interface
func (c *Counter) Type() EntityType {
	return EntityTypes.COUNTER
}

// GetID will return the ID of the DigestEntry Entity
func (c *Counter) GetID() uint32 {
	return c.ID
}

// GetCounter will return a Counter entity from a p4config counter entity
func GetCounter(counter *p4ConfigV1.Counter) Counter {
	return Counter{
		ID:   counter.Preamble.Id,
		Size: counter.Size,
	}
}
