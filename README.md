
# DragonflyDB HA Replication Deployment on Fly.io

This project creates a Golang application that manages a highly available [DragonflyDB](https://www.dragonflydb.io) deployment using replication on [Fly.io](https://fly.io).

This project leans heavily on the Fly.io's [posgres-flex](https://github.com/fly-apps/postgres-flex) project for how to achieve high availability on Fly.io.

## How it works

The application starts the following processes:
- HAProxy - used to dictate if traffic can be routed to it (only if the node is a master)
- [DragonflyDB](https://www.dragonflydb.io) 
- [Redis Sentinel](https://redis.io/docs/management/sentinel/) - handles failover and replica management 
- Admin API for health checks
- Monitoring process - updates Fly machine metadata to reflect the current master node's IP

## Configuration

Environment variables:

Variables prefixed with `FLY_` are set by Fly.io and are not configurable.

| Name | Description | Default |
| --- | --- | --- |
| FLY_APP_NAME | The name of the application | local |
| FLY_MACHINE_ID | The ID of the machine |  |
| FLY_NAMESERVER | The nameserver to use for DNS lookups | fdaa::3 |
| PRIMARY_REGION | The primary region to use for the cluster |  |
| DRAGONFLY_DIR | The directory to store DragonflyDB snapshots. Also supports s3 paths. | /data |
| AWS_ENDPOINT_URL_S3 | The endpoint URL for the S3 bucket. Only used if the directory is an S3 path |  |
| DRAGONFLY_MASTER_NAME | The name of the master node used by Redis Sentinel | mymaster |
| DRAGONFLY_QUORUM | The number of nodes required for a quorum. Minimum of 3 nodes are required for a quorum of 2| 2 |
| SENTINEL_DOWN_AFTER_MILLISECONDS | Milliseconds to wait before marking the master as down and attempting to failover | 60000 |
| OTEL_EXPORTER_OTLP_ENDPOINT | The endpoint for the OpenTelemetry Collector. Will not start if not set |  | 
| DEBUG | Enable debug logging | false |
