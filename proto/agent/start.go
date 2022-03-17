package agent

import (
	"github.com/YRXING/data-primitive/pkg/trace"
	"google.golang.org/grpc"
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
	log.Printf("Starting %v gRPC server, listener on %v", reflect.TypeOf(server), addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
