package bank

import (
	"context"
	"encoding/json"
	. "github.com/YRXING/data-primitive/pkg/constants"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
	"github.com/opentracing/opentracing-go"
	"log"
)

type bank struct {
	address string
	name string
	funcs   map[string]interface{}
	parentSpan opentracing.Span
	receivedPacket *agent.Packet
}

var _ DigitalObject = &bank{}

func NewBank() *bank {
	b := &bank{
		address: "127.0.0.1:8082",
	}

	b.funcs = map[string]interface{}{
		"GetLoan": b.GetLoan,
	}
	return b
}

func (b *bank) Run() error{
	go agent.RunServer(BANK_SERVICE,b.address, b)

	return nil
}

func (b *bank) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error) {
	// get the server side root span
	b.parentSpan = opentracing.SpanFromContext(ctx)
	b.receivedPacket = p

	switch p.Type {
	case agent.PacketType_INVOKE:
		//util.ProcessInvokePacket(b,p)
		res, err := util.Call(b.GetFuncs(), p.GetInvoke().FuncName, p.GetInvoke().Args)
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
		pkt := util.GenerateDataPacket(b.GetAddress(), bytes)
		return pkt, nil
	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}

func (b *bank) GetLoan(bytes []byte) *Capital{
	var (
		f Form
		res *Capital
	)

	err := json.Unmarshal(bytes, &f)
	if err != nil {
		return nil
	}

	span := opentracing.StartSpan("GetLoan",opentracing.FollowsFrom(b.parentSpan.Context()))
	defer span.Finish()

	switch f.Type {
	case ACCOUNT_RECEIVABLE:
		res = &Capital{
			BankName: b.name,
			Num: f.Num,
		}
	default:
		res = nil
	}
	return res
}

func (b *bank) GetAddress() string  {
	return b.address
}

func (b *bank) GetFuncs() map[string]interface{} {
	return b.funcs
}