package core

import (
	"context"
	"io"
	"log"
)

// StartMessageChannels starts two goroutines - one listens to the stream
// channel and sends any messages received to the PullStreamMessages
// channel, and the other listens for messages on the PushStreamMessages
// channel and sends them to the stream. So in short, to send any stream
// messages the developer can send them to the PushStreamMessages channel
// and to receive any stream messages the developer can listen to the
// PullStreamMessages channel
func (c *Client) StartMessageChannels() {
	stream, err := c.StreamChannel(context.Background())
	if err != nil {
		log.Fatal("Unable to get StreamChannel for client")
	}

	// Start a goroutine that gets messages from the stream and sends those
	// messages to the IncomingMessageChannel channel
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Println("Failed to get message from stream:", err)
			}
			if err != nil {
				log.Println("Error receiving message from stream:", err)
			}

			c.IncomingMessageChannel <- in
		}
	}()

	// Start a goroutine that listens on OutgoingMessageChannel, and sends any message
	// received on that channel to the gRPC stream channel
	go func() {
		for {
			sendMess := <-c.OutgoingMessageChannel
			if err := stream.Send(sendMess); err != nil {
				log.Println("Unable to send message to stream")
			}
		}
	}()
}
