package control

type Control interface {
	PerformArbitration() (bool, error)
}
