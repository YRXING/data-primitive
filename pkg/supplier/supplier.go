package supplier

import (
	"context"
	"encoding/json"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
)

type supplier struct {
	client agent.AgentClient
	address string
	funcs map[string]interface{}
}

func NewSupplier() *supplier {
	s := &supplier{
		address: "10.10.102.1",
	}

	s.funcs = map[string]interface{}{
		"GetProducts": s.GetProducts,
		"Func1": s.Func1,
		"Func2": s.Func2,
	}
	return s
}

func (s *supplier) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error)  {
	switch p.Type {
	case agent.PacketType_INVOKE:
		res, err:= util.Call(s.funcs,p.GetInvoke().FuncName,p.GetInvoke().Args)
		//s.client = util.NewClient(p.SendAddress)
		if err != nil {
			// change []reflect.Value to []interface{}
			data := make([]interface{},0)
			for _, v := range res {
				data = append(data,v.Interface())
			}

			bytes,_ := json.Marshal(data)
			// make return packet
			pkt := &agent.Packet{
				Type: agent.PacketType_TRANSPORT,
				SourceAddress: s.address,
				Payload: &agent.Packet_Transport{
					Transport: &agent.Transport{
						Data: string(bytes),
					},
				},
			}
			return pkt,nil
		}
	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil,nil
}

func (s *supplier) GetProducts()  {
	
}

func (s *supplier) Func1()  {

}

func (s *supplier) Func2()  {

}
