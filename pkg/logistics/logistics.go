package logistics

import (
	"context"
	"github.com/YRXING/data-primitive/proto/agent"
	"github.com/opentracing/opentracing-go"
)

type logistics struct {
	address    string
	funcs      map[string]interface{}
	parentSpan opentracing.Span
}

func NewLogistics() *logistics {

	return &logistics{}
}

func Run() {

}

func (l *logistics) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error) {
	// get the server side root span
	l.parentSpan = opentracing.SpanFromContext(ctx)

	switch p.Type {
	case agent.PacketType_INVOKE:

	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil, nil
}
