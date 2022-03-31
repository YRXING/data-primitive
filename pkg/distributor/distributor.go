package distributor

import (
	"context"
	"encoding/json"
	"github.com/YRXING/data-primitive/pkg/trace"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"

	. "github.com/YRXING/data-primitive/pkg/constants"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
)

type distributor struct {
	address    string
	name       string
	funcs      map[string]interface{}
	parentSpan opentracing.Span
}

var _ DigitalObject = &distributor{}

func NewDistributor() *distributor {
	d := &distributor{
		address: "127.0.0.1:8081",
		name:    "distributorA",
	}
	d.funcs = map[string]interface{}{
		"GetPaymentPromise": d.GetPaymentPromise,
	}
	return d
}

func generateOrder(name string,ch chan <-*Order, wg *sync.WaitGroup) {
	var o *Order

	t := []string{NORMAL,ACCOUNT_RECEIVABLE}
	for i := 0; i < len(t); i++ {
		o = &Order{
			OrderType: t[i],
			OrderPrice:      10,
			OrderCount:      10,
			DistributorName: name,
		}
		ch <- o
		time.Sleep(5*time.Second)
	}

	close(ch)
	wg.Done()
}

func (d *distributor) Run() error {
	// run gRPC server
	go agent.RunServer(DISTRIBUTOR_SERVICE, d.address, d)
	tracer, closer := trace.NewTracer(DISTRIBUTOR_SERVICE)
	defer closer.Close()
	time.Sleep(3*time.Second)
	log.Println("start simulation process...")
	// get supplier information from register center
	log.Infof("finding supplier....")
	conn := util.NewConn(tracer, "127.0.0.1:8080", context.Background())
	defer conn.Close()

	c := agent.NewAgentClient(conn)
	log.Infof("supplier find: supplierA, establish connection successfully")

	ch := make(chan *Order)
	var wg sync.WaitGroup
	wg.Add(2)
	// order producer
	go generateOrder(d.name,ch,&wg)

	// order consumer
	go func(ch <-chan *Order) {
		for o := range ch {
			log.Infof("perceived new order:%+v",o)
			bytes, err := json.Marshal(o)
			if err != nil {
				log.Errorf("serialization failed, %v",err)
			}

			log.Printf("sending data to supplierA: %s",bytes)
			p := util.GenerateInvokePacket(d.address, "GetProducts", bytes)
			res, err := c.Interact(context.Background(), p)
			if err != nil {
				log.Error("get result failed: ",err)
				return
			}
			log.Println("distributor get the products: ", res)

			var resProducts Products
			json.Unmarshal(res.GetTransport().Data,&resProducts)

			if resProducts.OrderState == SUCCESS {
				switch o.OrderType {
				case NORMAL:
				case ACCOUNT_RECEIVABLE:
					// pay for products
					log.Info("get bank information from order")
					conn = util.NewConn(tracer, "127.0.0.1:8082", context.Background())
					c = agent.NewAgentClient(conn)
					log.Infof("bank find: bankA, establish connection successfully")
					capital := &Capital{
						BankName: "bankA",
						Num:      o.OrderPrice * o.OrderCount,
					}
					log.Println("prepare capital for products...")
					bytes, _ = json.Marshal(capital)
					log.Printf("sending capital to bankA: %s",bytes)
					p = util.GenerateInvokePacket(d.address, "PayForProducts", bytes)
					res, err = c.Interact(opentracing.ContextWithSpan(context.Background(), d.parentSpan), p)
					log.Println("the payment result:",res)
				}
			}
			time.Sleep(5*time.Second)
		}
		wg.Done()
	}(ch)

	wg.Wait()
	log.Println("simulation process finished.")
	return nil
}

func (d *distributor) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error) {
	// get the server side root span context
	d.parentSpan = opentracing.SpanFromContext(ctx)

	switch p.Type {
	case agent.PacketType_INVOKE:
		res, err := util.Call(d.GetFuncs(), p.GetInvoke().FuncName, p.GetInvoke().Args)
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
		pkt := util.GenerateDataPacket(d.GetAddress(), bytes)
		return pkt, nil
	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}

func (d *distributor) GetPaymentPromise(bytes []byte) *PaymentPromise {
	var (
		p PaymentPromise
	)
	err := json.Unmarshal(bytes, &p)
	if err != nil {
		return nil
	}
	log.Infof("get a payment promise:%+v",p)
	span := opentracing.StartSpan("GetPaymentPromise", opentracing.FollowsFrom(d.parentSpan.Context()))
	defer span.Finish()

	// verify the order
	if p.DistributorName == d.name {
		log.Infof("I promise to pay for products, signatured!")
		p.Signatured = true
	}
	return &p
}

func (d *distributor) GetAddress() string {
	return d.address
}

func (d *distributor) GetFuncs() map[string]interface{} {
	return d.funcs
}
