package control

import (
	"errors"
	"log"

	"github.com/dush-t/ryz/core"
	"github.com/dush-t/ryz/core/entities"

	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

// SimpleControl implements the Control interface.
// I could not think of a better name
type SimpleControl struct {
	Client             core.P4RClient
	DigestChannel      chan *p4V1.StreamMessageResponse_Digest
	ArbitrationChannel chan *p4V1.StreamMessageResponse_Arbitration
	setupNotifChannel  chan bool
}

// StartMessageRouter will start a goroutine that takes incoming messages from the stream
// and then sends them to corresponding channels based on message types
func (sc *SimpleControl) StartMessageRouter() {
	IncomingMessageChannel := sc.Client.GetMessageChannels().IncomingMessageChannel
	go func() {
		for {
			in := <-IncomingMessageChannel
			update := in.GetUpdate()

			switch update.(type) {
			case *p4V1.StreamMessageResponse_Arbitration:
				sc.ArbitrationChannel <- update.(*p4V1.StreamMessageResponse_Arbitration)
			case *p4V1.StreamMessageResponse_Digest:
				sc.DigestChannel <- update.(*p4V1.StreamMessageResponse_Digest)
			default:
				log.Println("Message has unknown type")
			}
		}
	}()
}

// SetMastershipStatus will call a method of the same name on P4RClient.
// We need to keep track of mastership to reason about which control can be
// used for what.
func (sc *SimpleControl) SetMastershipStatus(status bool) {
	sc.setupNotifChannel <- status
	sc.Client.SetMastershipStatus(status)
}

// IsMaster will return true if the control has mastership
func (sc *SimpleControl) IsMaster() bool {
	return sc.Client.IsMaster()
}

// Run will do all the work required to actually get the control
// instance up and running
func (sc *SimpleControl) Run() {
	// Start running the client i.e. start the Stream Channel
	// on the client.
	sc.Client.Run()

	// Start the goroutine that will take messages from the
	// streamchannel and route it to appropriate goroutines for
	// handling.
	sc.StartMessageRouter()

	// Start the goroutine that listens to arbitration updates
	// and handles those updates.
	sc.StartArbitrationUpdateListener()

	// Perform arbitration
	sc.PerformArbitration()

	// Wait till arbitration is complete
	<-sc.setupNotifChannel
}

// InstallProgram will install a p4 compiled binary on a given target
func (sc *SimpleControl) InstallProgram(binPath, p4InfoPath string) error {
	if !sc.IsMaster() {
		return errors.New("Control does not have mastership, cannot install program on device")
	}
	return sc.Client.SetFwdPipe(binPath, p4InfoPath)
}

// NewControl will create a new Control instance
func NewControl(addr string, deviceID uint64, electionID p4V1.Uint128) (Control, error) {
	client, err := core.NewClient(addr, deviceID, electionID)
	if err != nil {
		return nil, err
	}
	digestChan := make(chan *p4V1.StreamMessageResponse_Digest, 10)
	arbitrationChan := make(chan *p4V1.StreamMessageResponse_Arbitration)
	setupNotifChan := make(chan bool)

	control := SimpleControl{
		Client:             client,
		DigestChannel:      digestChan,
		ArbitrationChannel: arbitrationChan,
		setupNotifChannel:  setupNotifChan,
	}

	return &control, nil
}

// Table will return a TableControl struct
func (sc *SimpleControl) Table(tableName string) TableControl {
	tables := *(sc.Client.GetEntities(entities.EntityTypes.TABLE))
	table := tables[tableName].(*entities.Table)

	return TableControl{
		table:   table,
		control: sc,
	}
}

// Digest will return a DigestControl struct
func (sc *SimpleControl) Digest(digestName string) DigestControl {
	digests := *(sc.Client.GetEntities(entities.EntityTypes.DIGEST))
	digest := digests[digestName].(*entities.DigestEntry)

	return DigestControl{
		digest:  digest,
		control: sc,
	}
}

// Counter will return a CounterControl struct
func (sc *SimpleControl) Counter(counterName string) CounterControl {
	counters := *(sc.Client.GetEntities(entities.EntityTypes.COUNTER))
	counter := counters[counterName].(*entities.Counter)

	return CounterControl{
		counter: counter,
		control: sc,
	}
}
