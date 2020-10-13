package control

import (
	"github.com/dush-t/ryz/core/entities"
)

// TableControl will contain all the information required
// to do stuff with a P4 table
type TableControl struct {
	table   *entities.Table
	control *SimpleControl
}

// InsertEntryRaw is used to insert an entry into a table using raw values of match fields and params
func (tc TableControl) InsertEntryRaw(action string, mf []entities.Match, params [][]byte) error {
	entityTypeAction := entities.EntityTypes.ACTION
	actions := *(tc.control.Client.GetEntities(entityTypeAction))
	actionID := actions[action].(*entities.Action).ID

	insertMessage := tc.table.InsertEntryMessage(actionID, mf, params)
	return tc.control.Client.WriteUpdate(insertMessage)
}

// InsertEntry is a less ugly version of InsertEntryRaw, since it supports pretty input (which can
// be serialized from incoming json, yay!). This method will call the table's Transformer function
// to convert the data into raw match fields and params, and will then call InsertEntryRaw with these
// values.
func (tc TableControl) InsertEntry(action string, data map[string]interface{}) error {
	mf, params := tc.table.Transformer(data)
	return tc.InsertEntryRaw(action, mf, params)
}

// RegisterTransformer will register a transformer function for the table managed by this
// TableControl instance
func (tc TableControl) RegisterTransformer(transformer entities.TableEntryTransformer) {
	tc.table.RegisterTransformer(transformer)
}
