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
	switchControl, err := control.NewControl(defaultAddr, defaultDeviceID, electionID)
	if err != nil {
		log.Fatal("Error initializing control over device", err)
	}
	switchControl.Run()
	switchControl.InstallProgram(binPath, p4InfoPath)
	log.Println("Installed p4 program on target")

	// Registering transformers for tables
	switchControl.Table("ingress.ipv4_lpm").RegisterTransformer(ipv4LpmTransform)
	switchControl.Table("ingress.forward").RegisterTransformer(forwardTransformer)
	switchControl.Table("egress.send_frame").RegisterTransformer(sendFrameTransformer)

	// Populating tables with entries
	// Yikes
	ipv4Data1 := map[string]interface{}{"ip": string("10.0.0.10"), "port": uint32(1)}
	ipv4Data2 := map[string]interface{}{"ip": string("10.0.1.10"), "port": uint32(2)}
	switchControl.Table("ingress.ipv4_lpm").InsertEntry("ingress.set_nhop", ipv4Data1)
	switchControl.Table("ingress.ipv4_lpm").InsertEntry("ingress.set_nhop", ipv4Data2)

	sendFrameData1 := map[string]interface{}{"port": uint32(1), "mac": string("00:aa:bb:00:00:00")}
	sendFrameData2 := map[string]interface{}{"port": uint32(2), "mac": string("00:aa:bb:00:00:01")}
	switchControl.Table("egress.send_frame").InsertEntry("egress.rewrite_mac", sendFrameData1)
	switchControl.Table("egress.send_frame").InsertEntry("egress.rewrite_mac", sendFrameData2)

	forwardData1 := map[string]interface{}{"ip": string("10.0.0.10"), "mac": string("00:04:00:00:00:00")}
	forwardData2 := map[string]interface{}{"ip": string("10.0.1.10"), "mac": string("00:04:00:00:00:01")}
	switchControl.Table("ingress.forward").InsertEntry("ingress.set_dmac", forwardData1)
	switchControl.Table("ingress.forward").InsertEntry("ingress.set_dmac", forwardData2)

	stopCh := signals.RegisterSignalHandlers()

	log.Println("Press Ctrl+C to quit")
	<-stopCh
	log.Println("Stopped")
}
