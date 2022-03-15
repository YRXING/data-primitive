package distributor

import (
	"context"
	"encoding/json"
	"log"

	. "github.com/YRXING/data-primitive/pkg/constants"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
)

type distributor struct {
	address string
}

func NewDistributor() *distributor {
	return &distributor{
		address: "127.0.0.1:8081",
	}
}

func (d *distributor) Run() error {
	// run gRPC server
	go agent.RunServer(d.address,d)
	// get supplier information
	conn := util.NewConn("127.0.0.1:8080")
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
	switch p.Type {
	case agent.PacketType_INVOKE:

	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}


