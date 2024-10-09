package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/streamersonglist/dragonfly-flex/internal/dragonfly"
	"github.com/streamersonglist/dragonfly-flex/internal/fly"
	"github.com/streamersonglist/dragonfly-flex/internal/supervisor"
)

func main() {
	log.SetFlags(0)

	ctx := context.Background()

	svisor := supervisor.New("dragonfly", 5*time.Second)
	dataDir := os.Getenv("DRAGONFLY_DIR")
	if dataDir == "" {
		dataDir = "/data"
	}

	s3endpoint := os.Getenv("AWS_ENDPOINT_URL_S3")
	if strings.HasPrefix(dataDir, "s3://") && s3endpoint == "" {
		log.Fatalf("AWS_ENDPOINT_URL_S3 must be set when using an S3 path")
	}

	additionalOpts := ""
	if s3endpoint != "" {
		additionalOpts = fmt.Sprintf(" --s3_endpoint %s", s3endpoint)
	}

	node, err := dragonfly.NewNode()
	if err != nil {
		log.Fatalf("failed to create node: %s", err.Error())
	}

	ip, err := updateSentinelConfig(ctx, node)
	if err != nil {
		log.Fatalf("failed to update master IP: %s", err.Error())
	}

	if ip != nil && *ip != node.PrivateIP {
		additionalOpts = fmt.Sprintf("%s --replicaof [%s]:%d", additionalOpts, *ip, dragonfly.Port)
	}

	svisor.AddProcess("dragonfly", fmt.Sprintf("dragonfly --logtostderr --bind :: --port %d --dir %s%s", dragonfly.Port, dataDir, additionalOpts))
	startCollector(svisor)

	proxyEnv := map[string]string{
		"FLY_APP_NAME":      os.Getenv("FLY_APP_NAME"),
		"PRIMARY_REGION":    os.Getenv("PRIMARY_REGION"),
		"DF_LISTEN_ADDRESS": node.PrivateIP,
	}

	svisor.AddProcess("proxy", "/usr/sbin/haproxy -W -db -f /fly/haproxy.cfg", supervisor.WithEnv(proxyEnv), supervisor.WithRestart(0, 1*time.Second))
	svisor.AddProcess("admin", "/usr/local/bin/admin_server", supervisor.WithRestart(0, 5*time.Second))
	svisor.AddProcess("sentinel", "/usr/local/bin/redis-sentinel /fly/sentinel.conf", supervisor.WithRestart(0, 5*time.Second))
	svisor.AddProcess("monitor", "/usr/local/bin/monitor", supervisor.WithRestart(0, 5*time.Second))

	svisor.StopOnSignal(syscall.SIGINT, syscall.SIGTERM)

	if err := svisor.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func updateSentinelConfig(ctx context.Context, node *dragonfly.Node) (*string, error) {
	flyclient := fly.NewClient()
	machines, err := flyclient.ListMachines(ctx, node.AppName, nil, nil, nil)
	if err != nil {
		log.Printf("failed to list machines: %s", err.Error())
		return nil, err
	}

	var masterIP string = ""

	log.Printf("found %d machines", len(machines))

	for _, machine := range machines {
		value, exists := machine.Config.Metadata[dragonfly.MetadataKey]
		if exists {
			log.Printf("found master machine: %+v", machine)
			masterIP = value
			break
		}
	}

	// we just meed to make sure that only one machine is created at a time or else
	// we there's the potential of multiple masters being created at the same time
	if masterIP == "" {
		log.Printf("no master machine found in metadata, setting this machine as master")
		masterIP = node.PrivateIP
		err := flyclient.UpdateMachineMetadata(ctx, node.AppName, node.MachineID, dragonfly.MetadataKey, node.PrivateIP)
		if err != nil {
			log.Printf("failed to update master metadata: %s", err.Error())
			return nil, err
		}
	}

	downAfterMilliseconds := os.Getenv("SENTINEL_DOWN_AFTER_MILLISECONDS")
	if downAfterMilliseconds == "" {
		downAfterMilliseconds = "60000"
	}

	content, err := os.ReadFile("/fly/sentinel.conf")
	if err != nil {
		return nil, err
	}

	updated := bytes.ReplaceAll(content, []byte("$MASTER_IP"), []byte(masterIP))
	updated = bytes.ReplaceAll(updated, []byte("$SENTINEL_DOWN_AFTER_MILLISECONDS"), []byte(downAfterMilliseconds))

	err = os.WriteFile("/fly/sentinel.conf", updated, 0644)
	if err != nil {
		return nil, err
	}

	return &masterIP, nil
}

func startCollector(svisor *supervisor.Supervisor) {
	collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if collectorEndpoint == "" {
		return
	}

	content, err := os.ReadFile("/fly/collector.yaml")
	if err != nil {
		log.Fatalf("failed to read collector config: %s", err.Error())
	}

	updated := bytes.ReplaceAll(content, []byte("$OTEL_EXPORTER_OTLP_ENDPOINT"), []byte(collectorEndpoint))

	err = os.WriteFile("/fly/collector.yaml", updated, 0644)
	if err != nil {
		log.Fatalf("failed to write collector config: %s", err.Error())
	}

	svisor.AddProcess("otel-collector", "/otelcol-contrib --config /fly/collector.yaml", supervisor.WithRestart(0, 10*time.Second))
}
