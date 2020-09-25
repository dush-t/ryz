package core

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	p4ConfigV1 "github.com/p4lang/p4runtime/go/p4/config/v1"
)

func getDeviceConfig(binPath string) ([]byte, error) {
	return ioutil.ReadFile(binPath)
}

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
