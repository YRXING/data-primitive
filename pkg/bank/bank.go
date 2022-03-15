package bank

import (
	"context"
	"github.com/YRXING/data-primitive/proto/agent"
)

type bank struct {
	address string
	funcs      map[string]interface{}
}

func NewBank() *bank {
	b := &bank{
		address: "127.0.0.1:8082",
	}

	b.funcs = map[string]interface{}{
		"GetLoan": b.GetLoan,
	}
	return b
}

func (b *bank) Run()  {
	go agent.RunServer(b.address,b)
}

func (b *bank) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error) {
	switch p.Type {
	case agent.PacketType_INVOKE:


	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}

func (b *bank) GetLoan()  {

}