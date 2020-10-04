package control

// Control represents a controller's control over a switch.
type Control interface {
	PerformArbitration() (bool, error)
	IsMaster() bool
	SetMastershipStatus(bool)
}
