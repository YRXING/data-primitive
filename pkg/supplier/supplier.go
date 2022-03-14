package supplier

import (
	"context"
	"encoding/json"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
	"log"
)

type supplier struct {
	address    string
	totalStock int
	totalFunds int
	funcs      map[string]interface{}
}

func NewSupplier() *supplier {
	s := &supplier{
		address: "127.0.0.1:8080",
	}

	s.funcs = map[string]interface{}{
		"GetProducts": s.GetProducts,
		"Func1":       s.Func1,
		"Func2":       s.Func2,
	}
	return s
}

func (s *supplier) Run() error {
	go agent.RunServer(s.address,s)

	return nil
}

func (s *supplier) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error) {
	switch p.Type {
	case agent.PacketType_INVOKE:
		res, err := util.Call(s.funcs, p.GetInvoke().FuncName, p.GetInvoke().Args)
		if err != nil {
			log.Println(err)
			return nil,err
		}
		//change []reflect.Value to []interface{}
		data := make([]interface{}, 0)
		for _, v := range res {
			data = append(data, v.Interface())
		}

		bytes, err := json.Marshal(data)
		if err != nil{
			log.Println(err)
			return nil,err
		}
		// make return packet
		pkt := util.GenerateDataPacket(s.address, bytes)
		return pkt, nil
	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}


func (s *supplier) GetProducts(bytes []byte) string {
	var (
		o Order
		res string
	)
	err := json.Unmarshal(bytes, &o)
	if err != nil {
		return ""
	}

	switch o.OrderType {
	case FINACINGWAREHOUSE_ORDER:
		res = "I have received this type of order"
	default:
		res = "unknown type order"
	}
	return res
}

func (s *supplier) Func1() {

}

func (s *supplier) Func2() {

}
