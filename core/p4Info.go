package core

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	p4ConfigV1 "github.com/p4lang/p4runtime/go/p4/config/v1"
)

const invalidID = 0

func getP4Info(p4InfoPath string) (*p4ConfigV1.P4Info, error) {
	bytes, err := ioutil.ReadFile(p4InfoPath)
	if err != nil {
		return nil, err
	}

	p4Info := &p4ConfigV1.P4Info{}
	if err = proto.UnmarshalText(string(bytes), p4Info); err != nil {
		return nil, err
	}
	
	return p4Info, nil
}

func getTableID(p4Info *p4ConfigV1.P4Info, name string) uint32 {
	if p4Info == nil {
		return 0
	}
	for _, table := range p4Info.Tables {
		if table.Preamble.Name == name {
			return table.Preamble.Id
		}
	}

	return 0
}

func getActionID(p4Info *p4ConfigV1.P4Info, name string) uint32 {
	if p4Info == nil {
		return 0
	}
	for _, action := range p4Info.Actions {
		if action.Preamble.Name == name {
			return action.Preamble.Id
		}
	}

	return 0
}
