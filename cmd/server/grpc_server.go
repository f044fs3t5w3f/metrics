package main

import (
	"context"
	"net"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/service"
	pb "github.com/f044fs3t5w3f/metrics/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedMetricsServer
	Service *service.Service
	Subnet  *net.IPNet
}

func (s *GRPCServer) checkIP(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Internal, "missing metadata")
	}
	ipSlice := md.Get("x-real-ip")
	if len(ipSlice) != 1 {
		return status.Errorf(codes.Internal, "missing metadata")
	}
	ip := net.ParseIP(ipSlice[0])
	if !s.Subnet.Contains(ip) {
		return status.Errorf(codes.PermissionDenied, "Forbidden")
	}
	return nil
}

func (s *GRPCServer) UpdateMetrics(ctx context.Context, req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	err := s.checkIP(ctx)
	if err != nil {
		return nil, err
	}
	metrics := make([]*models.Metrics, len(req.GetMetrics()))
	for i, metric := range req.GetMetrics() {
		metrics[i] = &models.Metrics{
			ID: metric.GetId(),
		}
		switch metric.GetType() {
		case pb.Metric_COUNTER:
			metrics[i].MType = models.Counter
			delta := metric.GetDelta()
			metrics[i].Delta = &delta
		case pb.Metric_GAUGE:
			metrics[i].MType = models.Gauge
			value := metric.GetValue()
			metrics[i].Value = &value
		default:
			return nil, status.Errorf(codes.InvalidArgument, `unsupported type for metric with id %s `, metric.GetId())
		}
	}
	err = s.Service.UpdateMetrics(ctx, metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
	return &pb.UpdateMetricsResponse{}, nil
}
