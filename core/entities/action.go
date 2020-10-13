package entities

import (
	p4ConfigV1 "github.com/p4lang/p4runtime/go/p4/config/v1"
)

type Action struct {
	Name string
	ID   uint32
}

func (a *Action) GetID() uint32 {
	return a.ID
}

func (a *Action) Type() EntityType {
	return EntityTypes.ACTION
}

func GetAction(ac *p4ConfigV1.Action) Action {
	return Action{
		Name: ac.Preamble.Name,
		ID:   ac.Preamble.Id,
	}
}
