package agent

import (
	"context"
	"fmt"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/models"
	netTools "github.com/f044fs3t5w3f/metrics/pkg/net"
	pb "github.com/f044fs3t5w3f/metrics/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func ReportBatchWithRPC(c pb.MetricsClient, batch MetricsBatch) {
	metrics := make([]*pb.Metric, len(batch))

	for i, metric := range batch {
		builder := pb.Metric_builder{}
		builder.Id = metric.ID

		switch metric.MType {
		case models.Counter:
			builder.Type = pb.Metric_COUNTER
			if metric.Delta != nil {
				builder.Delta = *metric.Delta
			}

		case models.Gauge:
			builder.Type = pb.Metric_GAUGE
			if metric.Value != nil {
				builder.Value = *metric.Value
			}
		}

		metrics[i] = builder.Build()
	}

	request := &pb.UpdateMetricsRequest{}
	request.SetMetrics(metrics)

	ctx := context.Background()
	ip, err := netTools.GetIP()
	if err == nil {
		md := metadata.New(map[string]string{"x-real-ip": ip})
		ctx = metadata.NewOutgoingContext(ctx, md)
	} else {
		logger.Log.Error("netTools.GetIP", zap.Error(err))
	}

	resp, err := c.UpdateMetrics(ctx, request)
	fmt.Println(resp)
	if err != nil {
		logger.Log.Error("grpc request failed", zap.Error(err))
	}
}
