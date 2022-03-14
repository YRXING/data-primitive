package agent

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"reflect"
)

func RunServer(addr string,server AgentServer) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterAgentServer(s, server)
	log.Printf("Starting %v gRPC server, listener on %v", reflect.TypeOf(server),addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
