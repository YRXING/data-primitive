package bank

import (
	"context"
	"encoding/json"
	. "github.com/YRXING/data-primitive/pkg/constants"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type bank struct {
	address        string
	name           string
	funcs          map[string]interface{}
	parentSpan     opentracing.Span
	receivedPacket *agent.Packet
}

var _ DigitalObject = &bank{}

func NewBank() *bank {
	b := &bank{
		address: "127.0.0.1:8082",
		name:    "bankA",
	}

	b.funcs = map[string]interface{}{
		"GetLoan":        b.GetLoan,
		"PayForProducts": b.PayForProducts,
	}
	return b
}

func (b *bank) Run() error {
	go agent.RunServer(BANK_SERVICE, b.address, b)

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

		// we have only one value
		bytes, err := json.Marshal(data[0])
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

func (b *bank) GetLoan(bytes []byte) *Capital {
	var (
		f   Form
		res *Capital
	)

	err := json.Unmarshal(bytes, &f)
	if err != nil {
		return nil
	}

	log.Info("get a form ",f," start processing....")

	span := opentracing.StartSpan("GetLoan", opentracing.FollowsFrom(b.parentSpan.Context()))
	defer span.Finish()

	switch f.Type {
	case ACCOUNT_RECEIVABLE:
		log.Infof("get distributor information from form...")
		conn := util.NewConn(opentracing.GlobalTracer(),
			"127.0.0.1:8081",
			context.Background())
		defer conn.Close()
		c := agent.NewAgentClient(conn)
		log.Infof("distributor find: distributorA, establish connection successfully")

		// generate data
		paymentPromise := &PaymentPromise{
			DistributorName: f.DistributorName,
			SupplierName:    f.SupplierName,
			Signatured:      false,
		}
		log.Infof("generate payment promise:%+v",paymentPromise)

		bytes, _ := json.Marshal(paymentPromise)
		log.Printf("start sending data %s....",bytes)
		p := util.GenerateInvokePacket(b.address, "GetPaymentPromise", bytes)
		resP, err := c.Interact(opentracing.ContextWithSpan(context.Background(), span), p)
		if err != nil || resP.GetTransport().Data == nil {
			log.Infof("the loan is not approved.")
			res = nil
		}

		// verify whether the payment promise is signatured
		log.Infof("verify whether the payment promise is signatured...")
		var resPaymentPromise PaymentPromise
		json.Unmarshal(resP.GetTransport().Data, &resPaymentPromise)
		log.Infof("get the payment promise from distributorA %s",resP.GetTransport().Data)
		if resPaymentPromise.Signatured == true {
			log.Infof("the loan is approved.")
			res = &Capital{
				BankName: b.name,
				Num:      f.Num,
			}
		}else {
			log.Infof("the loan is not approved")
			res = &Capital{
				BankName: b.name,
				Num: 0,
			}
		}
	default:
		res = &Capital{
			BankName: b.name,
			Num: 0,
		}
	}
	return res
}

func (b *bank) PayForProducts(bytes []byte) bool {
	var (
		c Capital
	)

	err := json.Unmarshal(bytes, &c)
	if err != nil {
		log.Println("capital received failed!")
		return false
	}
	log.Println("I have received capital: ", c)
	return true
}

func (b *bank) GetAddress() string {
	return b.address
}

func (b *bank) GetFuncs() map[string]interface{} {
	return b.funcs
}
