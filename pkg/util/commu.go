package util

import (
	"github.com/YRXING/data-primitive/proto/agent"
	"google.golang.org/grpc"
	"log"
)

func NewConn(address string) *grpc.ClientConn {
	conn, err := grpc.Dial(address,grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func GenerateInvokePacket(sa, funcName string, args []byte) *agent.Packet {
	return &agent.Packet{
		Type:          agent.PacketType_INVOKE,
		SourceAddress: sa,
		Payload: &agent.Packet_Invoke{
			Invoke: &agent.Invoke{
				FuncName: funcName,
				Args:     args,
			},
		},
	}
}

func GenerateDataPacket(sa string, data []byte) *agent.Packet {
	return &agent.Packet{
		Type:          agent.PacketType_TRANSPORT,
		SourceAddress: sa,
		Payload: &agent.Packet_Transport{
			Transport: &agent.Transport{
				Data: data,
			},
		},
	}
}
