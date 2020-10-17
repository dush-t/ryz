package control

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
