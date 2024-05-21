package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/sequenceplanner/rclgo/pkg/rclgo"
	"github.com/sequenceplanner/rclgo/pkg/rclgo/ros2/std_msgs"
)

const (
	rosTopic   = "chatter"
	natsTopic  = "chatter"
	natsServer = "nats://localhost:4222"
)

func main() {
	// Set up context for managing shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle system interrupt signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	// Initialize ROS 2 node
	node, err := rclgo.NewNode(ctx, "go_ros2_nats_node", "")
	if err != nil {
		log.Fatalf("Failed to create ROS 2 node: %v", err)
	}
	defer node.Close()

	// Initialize NATS connection
	nc, err := nats.Connect(natsServer)
	if err != nil {
		log.Fatalf("Failed to connect to NATS server: %v", err)
	}
	defer nc.Close()

	// ROS 2 subscription callback
	callback := func(msg *std_msgs.String) {
		// Rebroadcast the message over NATS
		err := nc.Publish(natsTopic, []byte(msg.Data))
		if err != nil {
			log.Printf("Failed to publish message to NATS: %v", err)
		}
		log.Printf("Received ROS message: %s, rebroadcast to NATS", msg.Data)
	}

	// Create ROS 2 subscriber
	_, err = node.NewSubscription(ctx, std_msgs.StringTypeSupport, rosTopic, callback)
	if err != nil {
		log.Fatalf("Failed to create ROS 2 subscription: %v", err)
	}

	// Keep the service running until interrupted
	<-ctx.Done()
	fmt.Println("Shutting down...")
}
