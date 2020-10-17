package control

import (
	"errors"

	"github.com/dush-t/ryz/core/entities"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// CounterData represents a counter's value at a given index
type CounterData struct {
	ByteCount   int64
	PacketCount int64
	Index       int64
}

// CounterControl will contain all the information required
// to do stuff with a P4 counter
type CounterControl struct {
	control *SimpleControl
	counter *entities.Counter
}

func getCounterData(entity *p4V1.Entity) CounterData {
	counterEntry := entity.GetCounterEntry()
	return CounterData{
		ByteCount:   counterEntry.Data.ByteCount,
		PacketCount: counterEntry.Data.PacketCount,
		Index:       counterEntry.Index.Index,
	}
}

// ReadValueAtIndex will return the value of the target counter at the target
// index. Duh.
func (cc *CounterControl) ReadValueAtIndex(index int64) (*CounterData, error) {
	entity := cc.counter.ReadValueAtIndexRequest(index)
	entityList := []*p4V1.Entity{entity}

	res, err := cc.control.Client.ReadEntitiesSync(entityList)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("No counter entries found at given index")
	}

	result := getCounterData(res[0])
	return &result, nil
}

// ReadValues will return a list of counter values at all indices
func (cc *CounterControl) ReadValues() ([]*CounterData, error) {
	entity := cc.counter.ReadWildcardRequest()
	entityList := []*p4V1.Entity{entity}

	res, err := cc.control.Client.ReadEntitiesSync(entityList)

	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("Target counter does not have any entries")
	}

	result := make([]*CounterData, cc.counter.Size)
	for _, item := range res {
		counterData := getCounterData(item)
		result = append(result, &counterData)
	}

	return result, nil
}

// StreamValues is like ReadValues, except it will return a channel
// on which all the counter values will be sent. Use this when you
// want to get all values asynchronously
func (cc *CounterControl) StreamValues() (chan *CounterData, error) {
	entity := cc.counter.ReadWildcardRequest()
	entityList := []*p4V1.Entity{entity}

	counterEntityCh, err := cc.control.Client.ReadEntities(entityList)
	if err != nil {
		return nil, err
	}

	cdataChannel := make(chan *CounterData, cc.counter.Size)
	go func() {
		defer close(cdataChannel)
		for e := range counterEntityCh {
			counterData := getCounterData(e)
			cdataChannel <- &counterData
		}
	}()

	return cdataChannel, nil
}
