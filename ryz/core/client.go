package core

import (
	"google.golang.org/grpc"

	p4ConfigV1 "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// Client contains all the information required to handle a client
type Client struct {
	p4V1.P4RuntimeClient
	deviceID      uint64
	electionID    p4V1.Uint128
	p4Info        *p4ConfigV1.P4Info
	MessageStream chan *p4V1.StreamMessageResponse
	IsMaster      bool
}

// Init will create a new gRPC connection and initialize the client
func (c *Client) Init(addr string, p4Info *p4ConfigV1.P4Info, deviceID uint64, electionID p4V1.Uint128) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	p4RtC := p4V1.NewP4RuntimeClient(conn)

	streamMsgs := make(chan *p4V1.StreamMessageResponse, 20)

	c.P4RuntimeClient = p4RtC
	c.deviceID = deviceID
	c.electionID = electionID
	c.p4Info = p4Info
	c.MessageStream = streamMsgs
	c.IsMaster = false

	return nil
}

func (c *Client) Run()

// NewClient will create a new P4 Runtime Client
func NewClient(addr, p4InfoPath string, deviceID uint64, electionID p4V1.Uint128) (*Client, error) {
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
