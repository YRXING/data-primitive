package agent

import (
	"context"
	"github.com/YRXING/data-primitive/pkg/trace"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"reflect"
)

func RunServer(serviceName string, addr string, server AgentServer) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tracer, closer := trace.NewTracer(serviceName)
	defer closer.Close()

	s := grpc.NewServer(grpc.UnaryInterceptor(trace.ServerInterceptor(tracer)))
	RegisterAgentServer(s, server)
	grpc_health_v1.RegisterHealthServer(s, &HealthImpl{})
	log.Printf("Starting %v gRPC server, listener on %v", reflect.TypeOf(server), addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// HealthImpl dedicate to grpc service check, it realized the HealthServer interface
type HealthImpl struct{}

func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.Println("health checking")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}