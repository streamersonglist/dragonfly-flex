package main

import (
	"context"
	"log"

	"github.com/streamersonglist/dragonfly-flex/internal/dragonfly"
	"github.com/streamersonglist/dragonfly-flex/internal/fly"
)

func main() {
	ctx := context.Background()

	log.SetFlags(0)

	node, err := dragonfly.NewNode()
	if err != nil {
		log.Fatalf("failed to create node: %s", err.Error())
	}

	node.Connect()
	defer node.Close()

	flyclient := fly.NewClient()

	monitor(ctx, node, flyclient)
}
