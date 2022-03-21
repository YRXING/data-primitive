package distributor

import (
	"context"
	"encoding/json"
	"github.com/YRXING/data-primitive/pkg/trace"
	"github.com/opentracing/opentracing-go"
	"log"

	. "github.com/YRXING/data-primitive/pkg/constants"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
)

type distributor struct {
	address string
	funcs      map[string]interface{}
	parentSpan opentracing.Span
}

var _ DigitalObject = &distributor{}

func NewDistributor() *distributor {
	return &distributor{
		address: "127.0.0.1:8081",
	}
}

func (d *distributor) Run() error {
	// run gRPC server
	go agent.RunServer(DISTRIBUTOR_SERVICE, d.address, d)
	// get supplier information
	tracer,closer := trace.NewTracer(DISTRIBUTOR_SERVICE)
	defer closer.Close()

	conn := util.NewConn(tracer,"127.0.0.1:8080",context.Background())
	defer conn.Close()

	c := agent.NewAgentClient(conn)
	o := &Order{
		OrderType:  NORMAL,
		OrderPrice: 10,
		OrderCount: 10,
	}

	bytes, _ := json.Marshal(o)
	p := util.GenerateInvokePacket(d.address, "GetProducts", bytes)
	res, err := c.Interact(context.Background(), p)
	if err != nil {
		return err
	}
	log.Println("distributor get the result: ", res)
	return nil
}

func (d *distributor) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error) {
	// get the server side root span context
	d.parentSpan = opentracing.SpanFromContext(ctx)

	switch p.Type {
	case agent.PacketType_INVOKE:

	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}

func (d *distributor) GetAddress() string  {
	return d.address
}

func (d *distributor) GetFuncs() map[string]interface{}  {
	return d.funcs
}