package core

import (
	"context"
	"log"

	"github.com/dush-t/ryz/core/entities"
	"google.golang.org/grpc"

	p4ConfigV1 "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// Client contains all the information required to handle a client
type Client struct {
	p4V1.P4RuntimeClient
	deviceID               uint64
	isMaster               bool
	electionID             p4V1.Uint128
	p4Info                 *p4ConfigV1.P4Info
	IncomingMessageChannel chan *p4V1.StreamMessageResponse
	OutgoingMessageChannel chan *p4V1.StreamMessageRequest
	streamChannel          p4V1.P4Runtime_StreamChannelClient
	Entities               map[entities.EntityType]*(map[string]entities.Entity)
}

// Init will create a new gRPC connection and initialize the client
func (c *Client) Init(addr string, p4Info *p4ConfigV1.P4Info, deviceID uint64, electionID p4V1.Uint128) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	p4RtC := p4V1.NewP4RuntimeClient(conn)
	resp, err := p4RtC.Capabilities(context.Background(), &p4V1.CapabilitiesRequest{})
	if err != nil {
		log.Fatal("Error in capabilities RPC", err)
	}
	log.Println("P4Runtime server version is", resp.P4RuntimeApiVersion)

	streamMsgs := make(chan *p4V1.StreamMessageResponse, 20)
	pushMsgs := make(chan *p4V1.StreamMessageRequest, 20)

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

	c.P4RuntimeClient = p4RtC
	c.deviceID = deviceID
	c.electionID = electionID
	c.p4Info = p4Info
	c.IncomingMessageChannel = streamMsgs
	c.OutgoingMessageChannel = pushMsgs
	c.Entities = Entities

	stream, streamInitErr := c.StreamChannel(context.Background())
	if streamInitErr != nil {
		return streamInitErr
	}

	c.streamChannel = stream

	return nil
}

// Run will do whatever is needed to ensure that the client is active
// once it is initialized.
func (c *Client) Run() {
	c.StartMessageChannels()
}

// WriteUpdate is used to update an entity on the
// switch. Refer to the P4Runtime spec to know more.
func (c *Client) WriteUpdate(update *p4V1.Update) error {
	req := &p4V1.WriteRequest{
		DeviceId:   c.deviceID,
		ElectionId: &c.electionID,
		Updates:    []*p4V1.Update{update},
	}

	_, err := c.Write(context.Background(), req)
	return err
}

// NewClient will create a new P4 Runtime Client
func NewClient(addr, p4InfoPath string, deviceID uint64, electionID p4V1.Uint128) (P4RClient, error) {
	p4Info, err := getP4Info(p4InfoPath)
	if err != nil {
		return nil, err
	}

	var client *Client
	initErr := client.Init(addr, p4Info, deviceID, electionID)
	if initErr != nil {
		return nil, initErr
	}

	return client, nil
}

/*
	Getters and Setters beyond this point
*/

// GetMessageChannels will return the message channels used by the client
func (c *Client) GetMessageChannels() MessageChannels {
	return MessageChannels{
		IncomingMessageChannel: c.IncomingMessageChannel,
		OutgoingMessageChannel: c.OutgoingMessageChannel,
	}
}

// GetArbitrationData will return the data required to perform arbitration
// for the client
func (c *Client) GetArbitrationData() ArbitrationData {
	return ArbitrationData{
		DeviceID:   c.deviceID,
		ElectionID: c.electionID,
	}
}

// GetStreamChannel will return the StreamChannel instance associated with the client
func (c *Client) GetStreamChannel() p4V1.P4Runtime_StreamChannelClient {
	return c.streamChannel
}

// P4Info will return the P4Info struct associated to the client
func (c *Client) P4Info() *p4ConfigV1.P4Info {
	return c.p4Info
}

// IsMaster returns true if the client is master
func (c *Client) IsMaster() bool {
	return c.isMaster
}

// SetMastershipStatus sets the mastership status of the client
func (c *Client) SetMastershipStatus(status bool) {
	c.isMaster = status
}

// GetEntities will return the Entities that the client has
func (c *Client) GetEntities(EntityType entities.EntityType) *map[string]entities.Entity {
	return c.Entities[EntityType]
}
