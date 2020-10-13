package control

// Control represents a controller's control over a switch.
type Control interface {
	PerformArbitration()
	IsMaster() bool
	SetMastershipStatus(bool)
	Run()
}
