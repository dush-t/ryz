package entities

import (
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// DirectCounterForTableEntryMessage will return a request entity that can be sent to
// the device to get the value of the DirectCounter against a given table entry.
func (t *Table) DirectCounterForTableEntryMessage(matches []Match) *p4V1.Entity {
	tableEntry := &p4V1.TableEntry{
		TableId: t.ID,
	}
	for idx, m := range matches {
		tableEntry.Match = append(tableEntry.Match, m.get(uint32(idx+1)))
	}

	dcEntry := &p4V1.DirectCounterEntry{
		TableEntry: tableEntry,
	}
	entity := &p4V1.Entity{
		Entity: &p4V1.Entity_DirectCounterEntry{dcEntry},
	}
	return entity
}

// AllDirectCountersForTableMessage will return a request entity that can be sent to
// the device to get values of DirectCounters against all entries of a table.
func (t *Table) AllDirectCountersForTableMessage() *p4V1.Entity {
	tableEntry := &p4V1.TableEntry{
		TableId: t.ID,
	}
	dcEntry := &p4V1.DirectCounterEntry{
		TableEntry: tableEntry,
	}
	entity := &p4V1.Entity{
		Entity: &p4V1.Entity_DirectCounterEntry{dcEntry},
	}
	return entity
}

// AllDirectCountersMessage will return a request entity that can be sent to the device
// to get values of all DirectCounters against all entries in all the tables.
func AllDirectCountersMessage() *p4V1.Entity {
	tableEntry := &p4V1.TableEntry{
		TableId: 0,
	}
	dcEntry := &p4V1.DirectCounterEntry{
		TableEntry: tableEntry,
	}
	entity := &p4V1.Entity{
		Entity: &p4V1.Entity_DirectCounterEntry{dcEntry},
	}
	return entity
}
