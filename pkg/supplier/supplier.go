package supplier

import (
	"context"
	"encoding/json"
	. "github.com/YRXING/data-primitive/pkg/constants"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
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

		// we have only one value
		bytes, err := json.Marshal(data[0])
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
	log.Info("get an order ",o," start producing....")

	span := opentracing.StartSpan("GetProducts", opentracing.FollowsFrom(s.parentSpan.Context()))
	defer span.Finish()

	switch o.OrderType {
	case NORMAL:
		log.Info("products ready,start transportation...")
		res = SuccessProducts(s.name)
	case FINACING_WAREHOUSE:

	case ACCOUNT_RECEIVABLE:
		log.Info("insufficient funds,looking for a bank to make a loan...")

		log.Info("finding bank...")
		conn := util.NewConn(opentracing.GlobalTracer(),
			"127.0.0.1:8082",
			context.Background())
		defer conn.Close()
		c := agent.NewAgentClient(conn)
		log.Infof("bank find: bankA, establish connection successfully")
		// generate data
		f := &Form{
			Type:            ACCOUNT_RECEIVABLE,
			SupplierName:    s.name,
			DistributorName: o.DistributorName,
			Num:             10000,
		}
		log.Infof("generate form: %+v",f)
		bytes, err := json.Marshal(f)
		if err != nil {
			log.Errorf("serialization failed, %v",err)
			return ErrorProducts(s.name,err.Error())
		}
		log.Printf("start sending data: %s",bytes)
		p := util.GenerateInvokePacket(s.address, "GetLoan", bytes)
		resP, err := c.Interact(opentracing.ContextWithSpan(context.Background(), span), p)
		var capital Capital
		json.Unmarshal(resP.GetTransport().Data,&capital)
		if err != nil || capital.Num == 0 {
			log.Errorf("get loan failed!")
			res = ErrorProducts(s.name, "Insufficient funds!")
			return res
		}

		log.Info("products ready,start transportation...")
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
