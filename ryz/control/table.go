package control

import (
	"github.com/dush-t/ryz/ryz/core"
)

// TableControl will contain all the information required
// to do stuff with a P4 table
type TableControl struct {
	table   *core.Table
	control *SimpleControl
}

func (tc *TableControl) InsertEntry(action string, mf []core.Match, params [][]byte) error {
	return tc.table.InsertEntry(tc.control.Client, action, mf, params)
}
