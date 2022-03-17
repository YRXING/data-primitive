package util

import (
	"context"
	"github.com/YRXING/data-primitive/pkg/trace"
	"github.com/YRXING/data-primitive/proto/agent"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
	"log"
	"time"
)

func NewConn(tracer opentracing.Tracer,address string) *grpc.ClientConn {
	var (
		conn *grpc.ClientConn
		err  error
	)

	// create a connect until it successes or timeout
	wait.Poll(3*time.Second, 5*time.Minute, func() (bool, error) {
		conn, err = grpc.DialContext(
			context.Background(),address,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(trace.ClientInterceptor(tracer)),
			)
		if err != nil {
			log.Println("can not connect: ", address)
			return false, err
		}
		return true, nil
	})

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
