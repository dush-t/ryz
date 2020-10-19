package control

import (
	"errors"

	"github.com/dush-t/ryz/core"

	"github.com/dush-t/ryz/core/entities"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

func getDirectCounterData(entity *p4V1.Entity) DirectCounterData {
	dcEntry := entity.GetDirectCounterEntry()
	tableEntry := TableEntry(*dcEntry.TableEntry)
	return DirectCounterData{
		TableEntry:  &tableEntry,
		ByteCount:   dcEntry.Data.ByteCount,
		PacketCount: dcEntry.Data.PacketCount,
	}
}

func getMultipleDCValuesSync(c core.P4RClient, req []*p4V1.Entity) ([]*DirectCounterData, error) {
	res, err := c.ReadEntitiesSync(req)
	if err != nil {
		return nil, err
	}

	result := make([]*DirectCounterData, 0)
	for _, item := range res {
		if item == nil {
			continue
		}
		dcData := getDirectCounterData(item)
		result = append(result, &dcData)
	}
	return result, nil
}

func streamMultipleDCValues(c core.P4RClient, req []*p4V1.Entity) (chan *DirectCounterData, error) {
	dcCounterEntityCh, err := c.ReadEntities(req)
	if err != nil {
		return nil, err
	}

	dcDataChannel := make(chan *DirectCounterData, 100)
	go func() {
		defer close(dcDataChannel)
		for e := range dcCounterEntityCh {
			dcCounterData := getDirectCounterData(e)
			dcDataChannel <- &dcCounterData
		}
	}()

	return dcDataChannel, nil
}

// ReadDirectCounterValueOnEntry will read the value of the DirectCounter against an entry in the table
func (tc TableControl) ReadDirectCounterValueOnEntry(matches []entities.Match) (*DirectCounterData, error) {
	entity := tc.table.DirectCounterForTableEntryMessage(matches)
	entityList := []*p4V1.Entity{entity}

	res, err := tc.control.Client.ReadEntitiesSync(entityList)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("No counter entries found")
	}
	// Why is the value in res[1] instead of res[0]? I have no clue
	result := getDirectCounterData(res[1])
	return &result, nil
}

// ReadDirectCounterValuesSync will return a list of DirectCounter values against all the entries of the table.
func (tc TableControl) ReadDirectCounterValuesSync() ([]*DirectCounterData, error) {
	entity := tc.table.AllDirectCountersForTableMessage()
	entityList := []*p4V1.Entity{entity}

	return getMultipleDCValuesSync(tc.control.Client, entityList)
}

// StreamDirectCounterValues does what ReadDirectCounterValuesSync does, except that instead of returning a list
// of values, it returns a channel on which these values can be asynchronously sent.
func (tc TableControl) StreamDirectCounterValues() (chan *DirectCounterData, error) {
	entity := tc.table.AllDirectCountersForTableMessage()
	entityList := []*p4V1.Entity{entity}

	return streamMultipleDCValues(tc.control.Client, entityList)
}
