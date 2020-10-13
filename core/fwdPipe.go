package core

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/dush-t/ryz/core/entities"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

func getDeviceConfig(binPath string) ([]byte, error) {
	return ioutil.ReadFile(binPath)
}

// SetFwdPipe will install the p4 compiled binary on the target device
func (c *Client) SetFwdPipe(binPath string, p4infoPath string) error {
	deviceConfig, err := getDeviceConfig(binPath)
	if err != nil {
		return fmt.Errorf("error when reading binary device config: %v", err)
	}
	p4Info, err := getP4Info(p4infoPath)
	if err != nil {
		return fmt.Errorf("error when reading P4Info text file: %v", err)
	}
	config := &p4V1.ForwardingPipelineConfig{
		P4Info:         p4Info,
		P4DeviceConfig: deviceConfig,
	}
	req := &p4V1.SetForwardingPipelineConfigRequest{
		DeviceId:   c.deviceID,
		ElectionId: &c.electionID,
		Action:     p4V1.SetForwardingPipelineConfigRequest_VERIFY_AND_COMMIT,
		Config:     config,
	}
	_, err = c.SetForwardingPipelineConfig(context.Background(), req)

	// Setup the entities on the client
	Tables := make(map[string]entities.Entity)
	for _, table := range p4Info.Tables {
		t := entities.GetTable(table)
		Tables[table.Preamble.Name] = entities.Entity(&t)
	}

	Actions := make(map[string]entities.Entity)
	for _, action := range p4Info.Actions {
		a := entities.GetAction(action)
		Actions[action.Preamble.Name] = entities.Entity(&a)
	}

	Entities := make(map[entities.EntityType]*(map[string]entities.Entity))
	Entities[entities.EntityTypes.TABLE] = &Tables
	Entities[entities.EntityTypes.ACTION] = &Actions
	c.Entities = Entities

	if err == nil {
		c.p4Info = p4Info
	}
	return err
}
