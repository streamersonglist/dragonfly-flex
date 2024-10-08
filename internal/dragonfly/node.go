package dragonfly

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streamersonglist/dragonfly-flex/internal/privnet"
)

type Node struct {
	AppName       string
	Port          int
	PrivateIP     string
	PrimaryRegion string
	MachineID     string
	MasterName    string
	Quorum        string
	Client        *redis.Client
	Debug         bool
}

var (
	MetadataKey     = "dragonfly_master"
	MasterRole      = "master"
	SlaveRole       = "slave"
	SentinelAddress = "127.0.0.1:26379"
	Port            = 6380
)

func NewNode() (*Node, error) {
	node := &Node{
		AppName:    "local",
		Port:       Port,
		MachineID:  os.Getenv("FLY_MACHINE_ID"),
		MasterName: "mymaster",
		Quorum:     "2",
	}

	if appName := os.Getenv("FLY_APP_NAME"); appName != "" {
		node.AppName = appName
	}

	if masterName := os.Getenv("DRAGONFLY_MASTER_NAME"); masterName != "" {
		node.MasterName = masterName
	}

	if quorum := os.Getenv("DRAGONFLY_QUORUM"); quorum != "" {
		node.Quorum = quorum
	}

	ipv6, err := privnet.PrivateIPv6()
	if err != nil {
		return nil, fmt.Errorf("failed to get private IPv6: %s", err.Error())
	}
	node.PrivateIP = ipv6.String()

	node.Debug = os.Getenv("DEBUG") == "true"

	node.PrimaryRegion = os.Getenv("PRIMARY_REGION")
	if node.PrimaryRegion == "" {
		return nil, fmt.Errorf("PRIMARY_REGION environment variable must be set")
	}

	log.Printf("node: %+v\n", node)

	return node, nil
}

func (n *Node) Connect() {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("127.0.0.1:%d", n.Port),
		Password:        "",
		MaxRetries:      3,
		MinRetryBackoff: 1 * time.Second,
	})
	n.Client = client
}

func (n *Node) Close() {
	if n.Client != nil {
		n.Client.Close()
	}
}

func (n *Node) CheckRole(ctx context.Context) (*string, error) {
	cmd := redis.NewSliceCmd(ctx, "ROLE")

	if n.Client == nil {
		n.Connect()
	}

	err := n.Client.Process(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("error processing ROLE command: %w", err)
	}

	resp, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("error running ROLE command: %w", err)
	}

	if n.Debug {
		fmt.Printf("ROLE command response: %+v\n", resp)
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf("no response from ROLE command")
	}

	if resp[0] == "master" {
		return &MasterRole, nil
	}

	if resp[0] == "slave" {
		return &SlaveRole, nil
	}

	return nil, nil
}

func (n *Node) GetSentinelMaster(ctx context.Context) (*string, error) {
	client := redis.NewSentinelClient(&redis.Options{
		Addr: SentinelAddress,
	})

	defer client.Close()

	cmd := client.GetMasterAddrByName(ctx, n.MasterName)

	err := client.Process(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("error processing GetMasterAddrByName command: %w", err)
	}

	result, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("error getting result from GetMasterAddrByName command: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no result from GetMasterAddrByName command")
	}

	return &result[0], nil
}

func (n *Node) SetSentinelMaster(ctx context.Context, ip string) error {
	client := redis.NewSentinelClient(&redis.Options{
		Addr: SentinelAddress,
	})

	defer client.Close()

	cmd := client.Monitor(ctx, n.MasterName, ip, "26379", n.Quorum)

	err := client.Process(ctx, cmd)
	if err != nil {
		return fmt.Errorf("error processing monitor command: %w", err)
	}

	return nil
}

func (n *Node) SetReplicaOf(ctx context.Context, ip string) error {
	cmd := redis.NewCmd(ctx, "REPLICAOF", ip, n.Port)

	if n.Client == nil {
		n.Connect()
	}

	err := n.Client.Process(ctx, cmd)
	if err != nil {
		return fmt.Errorf("error processing REPLICAOF command: %w", err)
	}

	return nil
}

func (n *Node) SetMaster(ctx context.Context) error {
	cmd := redis.NewCmd(ctx, "REPLICAOF", "NO ONE")

	if n.Client == nil {
		n.Connect()
	}

	err := n.Client.Process(ctx, cmd)
	if err != nil {
		return fmt.Errorf("error processing REPLICAOF NO ONE command: %w", err)
	}

	return nil
}
