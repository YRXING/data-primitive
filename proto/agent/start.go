package agent

import (
	"github.com/YRXING/data-primitive/pkg/trace"
	"github.com/YRXING/data-primitive/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
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
	grpc_health_v1.RegisterHealthServer(s, &util.HealthImpl{})
	log.Printf("Starting %v gRPC server, listener on %v", reflect.TypeOf(server), addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
