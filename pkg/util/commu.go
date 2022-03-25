package util

import (
	"context"
	"encoding/json"
	. "github.com/YRXING/data-primitive/pkg/constants"
	"github.com/YRXING/data-primitive/pkg/trace"
	"github.com/YRXING/data-primitive/proto/agent"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
	"log"
	"time"
)

func NewConn(tracer opentracing.Tracer, address string, ctx context.Context) *grpc.ClientConn {
	var (
		conn *grpc.ClientConn
		err  error
	)

	// create a connect until it successes or timeout
	wait.Poll(3*time.Second, 5*time.Minute, func() (bool, error) {
		conn, err = grpc.DialContext(
			ctx,
			address,
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

func ProcessInvokePacket(do DigitalObject, p *agent.Packet) (*agent.Packet, error) {
	res, err := Call(do.GetFuncs(), p.GetInvoke().FuncName, p.GetInvoke().Args)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//change []reflect.Value to []interface{}
	data := make([]interface{}, 0)
	for _, v := range res {
		data = append(data, v.Interface())
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// make return packet
	pkt := GenerateDataPacket(do.GetAddress(), bytes)
	return pkt, nil
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
