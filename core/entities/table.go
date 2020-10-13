package entities

import (
	p4ConfigV1 "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// Match represents a match field on a table Entity on
// a P4 device
type Match interface {
	get(ID uint32) *p4V1.FieldMatch
}

// ExactMatch implements the Match interface to provide
// Exact Match functionality
type ExactMatch struct {
	Value []byte
}

// LpmMatch implements the Match interface to provide
// Longest Prefix Match functionality
type LpmMatch struct {
	Value []byte
	PLen  int32
}

func (m *ExactMatch) get(ID uint32) *p4V1.FieldMatch {
	exact := &p4V1.FieldMatch_Exact{
		Value: m.Value,
	}
	mf := &p4V1.FieldMatch{
		FieldId:        ID,
		FieldMatchType: &p4V1.FieldMatch_Exact_{exact},
	}
	return mf
}

func (m *LpmMatch) get(ID uint32) *p4V1.FieldMatch {
	lpm := &p4V1.FieldMatch_LPM{
		Value:     m.Value,
		PrefixLen: m.PLen,
	}

	firstByteMasked := int(m.PLen / 8)
	if firstByteMasked != len(m.Value) {
		i := firstByteMasked
		r := m.PLen % 8
		m.Value[i] = m.Value[i] & (0xff << (8 - r))
		for i = i + 1; i < len(m.Value); i++ {
			m.Value[i] = 0
		}
	}

	mf := &p4V1.FieldMatch{
		FieldId:        ID,
		FieldMatchType: &p4V1.FieldMatch_Lpm{lpm},
	}

	return mf
}

// TableEntryTransformer can be registered with a table to convert json data
// into data compatible with the p4runtime protobuf spec.
type TableEntryTransformer func(map[string]interface{}) ([]Match, [][]byte)

// Table represents a table Entity
type Table struct {
	ID          uint32
	Name        string
	Transformer TableEntryTransformer
}

// InsertEntryMessage will insert a table entry in the P4 device
func (t *Table) InsertEntryMessage(actionID uint32, mfs []Match, params [][]byte) *p4V1.Update {
	directAction := &p4V1.Action{
		ActionId: actionID,
	}

	for idx, p := range params {
		param := &p4V1.Action_Param{
			ParamId: uint32(idx + 1),
			Value:   p,
		}
		directAction.Params = append(directAction.Params, param)
	}

	tableAction := &p4V1.TableAction{
		Type: &p4V1.TableAction_Action{directAction},
	}

	entry := &p4V1.TableEntry{
		TableId:         t.ID,
		Action:          tableAction,
		IsDefaultAction: (mfs == nil),
	}

	for idx, mf := range mfs {
		entry.Match = append(entry.Match, mf.get(uint32(idx+1)))
	}

	var updateType p4V1.Update_Type
	if mfs == nil {
		updateType = p4V1.Update_MODIFY
	} else {
		updateType = p4V1.Update_INSERT
	}
	update := &p4V1.Update{
		Type: updateType,
		Entity: &p4V1.Entity{
			Entity: &p4V1.Entity_TableEntry{entry},
		},
	}

	return update
}

// Type returns the type of the Entity represented by Table,
// so that it can implement the Entity interface
func (t *Table) Type() EntityType {
	return EntityTypes.TABLE
}

// GetID will return the ID of the table Entity
func (t *Table) GetID() uint32 {
	return t.ID
}

// RegisterTransformer will set the provided transformer to be used for
// a table entity.
func (t *Table) RegisterTransformer(transformer TableEntryTransformer) {
	t.Transformer = transformer
}

// GetTable will return a new table from a derived from a P4Info table
func GetTable(t *p4ConfigV1.Table) Table {
	return Table{
		Name: t.Preamble.Name,
		ID:   t.Preamble.Id,
	}
}
