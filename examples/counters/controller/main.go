package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dush-t/ryz/control"
	"github.com/dush-t/ryz/signals"
	p4V1 "github.com/p4lang/p4runtime/go/p4/v1"
)

const (
	defaultAddr     = "127.0.0.1:50051"
	defaultDeviceID = 0
)

func printCounterData(c *control.CounterData) {
	log.Printf("Received %d packets on port %d till now\n", c.PacketCount, c.Index)
}

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

	// Registering transformer for table
	switchControl.Table("ingress.ipv4_lpm").RegisterTransformer(ipv4LpmTransform)

	// Populating tables with entries
	ipv4Data1 := map[string]interface{}{
		"ip":   string("10.0.0.10"),
		"port": uint32(1),
		"mac":  string("00:04:00:00:00:00"),
	}

	ipv4Data2 := map[string]interface{}{
		"ip":   string("10.0.1.10"),
		"port": uint32(2),
		"mac":  string("00:04:00:00:00:01"),
	}

	switchControl.Table("ingress.ipv4_lpm").InsertEntry("ingress.ipv4_forward", ipv4Data1)
	switchControl.Table("ingress.ipv4_lpm").InsertEntry("ingress.ipv4_forward", ipv4Data2)

	// Read the counters every 10 seconds
	counter := switchControl.Counter("ingress.port_counter")
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println(" ")
			log.Println("Reading counter value from switch")
			result, err := counter.ReadValues()
			if err != nil {
				log.Println("Error occured while reading counter value", err)
				break
			}
			for _, cdata := range result {
				printCounterData(cdata)
			}
		}
	}()

	stopCh := signals.RegisterSignalHandlers()

	log.Println("Press Ctrl+C to quit")
	<-stopCh
	log.Println("Stopped")
}
