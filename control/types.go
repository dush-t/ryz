package control

import (
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

type ControlTable interface {
	Table(string) TableControl
	Digest(string) DigestControl
	Counter(string) CounterControl
}

// Control represents a controller's control over a switch.
type Control interface {
	ControlTable
	PerformArbitration()
	IsMaster() bool
	SetMastershipStatus(bool)
	Run()
	InstallProgram(string, string) error
}

// CounterData represents a counter's value at a given index
type CounterData struct {
	ByteCount   int64
	PacketCount int64
	Index       int64
}

type TableEntry p4V1.TableEntry

type DirectCounterData struct {
	TableEntry  *TableEntry
	ByteCount   int64
	PacketCount int64
}
