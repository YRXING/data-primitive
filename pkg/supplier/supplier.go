package supplier

import (
	"context"
	"encoding/json"
	. "github.com/YRXING/data-primitive/pkg/constants"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
	"github.com/opentracing/opentracing-go"
	"log"
)

type supplier struct {
	address        string
	name           string
	totalStock     int
	totalFunds     int
	funcs          map[string]interface{}
	parentSpan     opentracing.Span
	receivedPacket *agent.Packet
}

var _ DigitalObject = &supplier{}

func NewSupplier() *supplier {
	s := &supplier{
		address: "127.0.0.1:8080",
		name:    "supplierA",
	}

	s.funcs = map[string]interface{}{
		"GetProducts": s.GetProducts,
	}
	return s
}

func (s *supplier) Run() error {
	go agent.RunServer(SUPPLIER_SERVICE, s.address, s)

	return nil
}

func (s *supplier) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error) {
	// get the server side root span
	s.parentSpan = opentracing.SpanFromContext(ctx)
	// store the received packet for subsequent use
	s.receivedPacket = p

	switch p.Type {
	case agent.PacketType_INVOKE:
		//util.ProcessInvokePacket(s,p)
		res, err := util.Call(s.GetFuncs(), p.GetInvoke().FuncName, p.GetInvoke().Args)
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
		pkt := util.GenerateDataPacket(s.GetAddress(), bytes)
		return pkt, nil
	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}

func (s *supplier) GetProducts(bytes []byte) *Products {
	var (
		o   Order
		res *Products
	)
	err := json.Unmarshal(bytes, &o)
	if err != nil {
		return ErrorProducts(s.name, "wrong data format")
	}

	span := opentracing.StartSpan("GetProducts", opentracing.FollowsFrom(s.parentSpan.Context()))
	defer span.Finish()

	switch o.OrderType {
	case NORMAL:
		res = SuccessProducts(s.name)
	case FINACING_WAREHOUSE:

	case ACCOUNT_RECEIVABLE:
		conn := util.NewConn(opentracing.GlobalTracer(),
			"127.0.0.1:8082",
			context.Background())
		defer conn.Close()
		c := agent.NewAgentClient(conn)
		// generate data
		f := &Form{
			Type:            ACCOUNT_RECEIVABLE,
			SupplierName:    s.name,
			DistributorName: o.DistributorName,
			Num:             10000,
		}
		bytes, _ := json.Marshal(f)
		p := util.GenerateInvokePacket(s.address, "GetLoan", bytes)
		resP, err := c.Interact(opentracing.ContextWithSpan(context.Background(), span), p)
		if err != nil || resP.GetTransport().Data == nil {
			res = ErrorProducts(s.name, "Insufficient funds!")
		}
		res = SuccessProducts(s.name)

	case ADVANCE:

	default:
		res = ErrorProducts("unknown", "unknown order type!")
	}
	return res
}

func (s *supplier) GetAddress() string {
	return s.address
}

func (s *supplier) GetFuncs() map[string]interface{} {
	return s.funcs
}
