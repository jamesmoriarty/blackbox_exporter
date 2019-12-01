package prober

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/blackbox_exporter/config"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func ProbeGRPCHealthV1(ctx context.Context, target string, module config.Module, registry *prometheus.Registry, logger log.Logger) bool {
	healthConfig := module.GRPCHealthV1

	parts := strings.Split(target, "/")
  host, service := parts[0], parts[1]

	level.Info(logger).Log("msg", "dailing", "host", host)

	var options []grpc.DialOption
	if healthConfig.Plaintext {
		options = append(options, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(host, options...)

	if err != nil {
		level.Error(logger).Log("msg", "failed to connect", "err", err)
		return false
	}

	defer conn.Close()

	health := healthpb.NewHealthClient(conn)

	req := &healthpb.HealthCheckRequest{Service: service}
	res, err := health.Check(ctx, req)
	if err != nil {
		level.Error(logger).Log("msg", "error calling health.Check", "err", err)
		return false
	}

	return res.Status == healthpb.HealthCheckResponse_SERVING
}
