package entities

// Entity represents a P4 Entity
type Entity interface {
	Type() EntityType
	GetID() uint32
}
