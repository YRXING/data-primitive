package util

import (
	"github.com/YRXING/data-primitive/proto/agent"
	"google.golang.org/grpc"
	"log"
)

func NewClient(address string) agent.AgentClient {
	conn,err := grpc.Dial(address)
	if err != nil {
		log.Fatalf("did not connect: %v",err)
	}

	c := agent.NewAgentClient(conn)

	return c
}