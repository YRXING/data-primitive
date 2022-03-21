package constants

import (
	"context"
	"github.com/YRXING/data-primitive/proto/agent"
)

type DigitalObject interface {
	Run() error
	Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error)
	GetAddress() string
	GetFuncs() map[string]interface{}
}
