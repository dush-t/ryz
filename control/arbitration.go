package control

import (
	"log"

	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
	"google.golang.org/genproto/googleapis/rpc/code"
)

// PerformArbitration will send an Arbitration Request down the
// Stream Channel to participate in arbitration
func (sc *SimpleControl) PerformArbitration() {
	msgChannels := sc.Client.GetMessageChannels()
	outChan := msgChannels.OutgoingMessageChannel
	arbitrationData := sc.Client.GetArbitrationData()

	request := &p4V1.StreamMessageRequest{
		Update: &p4V1.StreamMessageRequest_Arbitration{&p4V1.MasterArbitrationUpdate{
			DeviceId:   arbitrationData.DeviceID,
			ElectionId: &(arbitrationData.ElectionID),
		}},
	}

	outChan <- request
}

// StartArbitrationUpdateListener will start a goroutine to listen on the
// ArbitrationChannel to check for any updates in arbitration
func (sc *SimpleControl) StartArbitrationUpdateListener() {
	go func() {
		update := <-sc.ArbitrationChannel
		if update.Arbitration.Status.Code != int32(code.Code_OK) {
			sc.SetMastershipStatus(false)
			log.Println("Arbitration was done. Control did not acquire mastership.")
		} else {
			sc.SetMastershipStatus(true)
			log.Println("Arbitration was done. Control acquired mastership")
		}
	}()
}
