package logistics

import (
	"context"
	"encoding/json"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type logistics struct {
	address    string
	funcs      map[string]interface{}
	parentSpan opentracing.Span
	receivedPacket *agent.Packet
}

func NewLogistics() *logistics {
	l := &logistics{
		address: "127.0.0.1:8083",
	}

	return l
}

func Run() {

}

func (l *logistics) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error) {
	// get the server side root span
	l.parentSpan = opentracing.SpanFromContext(ctx)

	// store the received packet for subsequent use
	l.receivedPacket = p

	switch p.Type {
	case agent.PacketType_INVOKE:
		//util.ProcessInvokePacket(s,p)
		res, err := util.Call(l.GetFuncs(), p.GetInvoke().FuncName, p.GetInvoke().Args)
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
		pkt := util.GenerateDataPacket(l.GetAddress(), bytes)
		return pkt, nil
	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}

func (l *logistics) GetFuncs() map[string]interface{} {
	return l.funcs
}

func (l *logistics) GetAddress() string  {
	return l.address
}