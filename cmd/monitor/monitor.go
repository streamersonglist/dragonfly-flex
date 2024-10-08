package main

import (
	"context"
	"log"
	"time"

	"github.com/streamersonglist/dragonfly-flex/internal/dragonfly"
	"github.com/streamersonglist/dragonfly-flex/internal/fly"
)

var (
	monitorInterval = 30 * time.Second
)

func monitor(ctx context.Context, node *dragonfly.Node, flyClient *fly.Client) {
	ticker := time.NewTicker(monitorInterval)
	defer ticker.Stop()

	for range ticker.C {
		monitorNode(ctx, node, flyClient)
	}
}

// monitorNode asks sentinel for the current master and updates the machine metadata if it doesn't match.
// The metadata gets cloned to a new machine when the app scales up so its best to always keep this up to date
func monitorNode(ctx context.Context, node *dragonfly.Node, flyClient *fly.Client) {
	ip, err := node.GetSentinelMaster(ctx)
	if err != nil {
		log.Printf("failed to check node role: %s", err.Error())
		return
	}

	// if *ip != node.PrivateIP {
	// 	role, err := node.CheckRole(ctx)
	// 	if err != nil {
	// 		log.Printf("failed to check node role: %s", err.Error())
	// 		return
	// 	}

	// 	if role != nil && *role == "master" {
	// 		log.Printf("node role should be a replica, setting replica of %s", *ip)
	// 		err := node.SetReplicaOf(ctx, *ip)
	// 		if err != nil {
	// 			log.Printf("failed to set replica of: %s", err.Error())
	// 			return
	// 		}
	// 	}
	// }

	meta, err := flyClient.GetMachineMetadata(ctx, node.AppName, node.MachineID)
	if err != nil {
		log.Printf("failed to get machine metadata: %s", err.Error())
		return
	}

	if meta[dragonfly.MetadataKey] != *ip {
		log.Printf("machine metadata [%s] does not match node master [%s], updating...", meta[dragonfly.MetadataKey], *ip)
		err := flyClient.UpdateMachineMetadata(ctx, node.AppName, node.MachineID, dragonfly.MetadataKey, *ip)
		if err != nil {
			log.Printf("failed to update machine metadata: %s", err.Error())
			return
		}
	}
}
