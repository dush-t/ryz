package main

import (
	"flag"
	"log"

	"github.com/dush-t/ryz/control"
	"github.com/dush-t/ryz/signals"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

const (
	defaultAddr     = "127.0.0.1:50051"
	defaultDeviceID = 0
)

func main() {
	var binPath string
	flag.StringVar(&binPath, "bin", "", "Path to P4 binary")
	var p4InfoPath string
	flag.StringVar(&p4InfoPath, "p4Info", "", "Path to p4Info")

	flag.Parse()

	if binPath == "" || p4InfoPath == "" {
		log.Fatal("Missing flags: bin or p4Info")
	}

	electionID := p4V1.Uint128{High: 0, Low: 1}
	control, err := control.NewControl(defaultAddr, p4InfoPath, defaultDeviceID, electionID)
	if err != nil {
		log.Fatal("Error initializing control over device", err)
	}
	control.Run()
	stopCh := signals.RegisterSignalHandlers()

	log.Println("Press Ctrl+C to quit")
	<-stopCh
	log.Println("Stopped")
}
