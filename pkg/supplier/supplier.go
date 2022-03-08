package supplier

import (
	"context"
	"github.com/YRXING/data-primitive/pkg/util"
	"github.com/YRXING/data-primitive/proto/agent"
)

type supplier struct {
	client agent.AgentClient

}

func (s *supplier) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error)  {
	switch p.Type {
	case agent.PacketType_INVOKE:
		res := p.GetInvoke().FuncName
		s.client = util.NewClient(p.GetInvoke())
		
	}
}

func (s *supplier) GetProducts()  {
	
}
