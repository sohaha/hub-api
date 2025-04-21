package openai

import (
	"context"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/zlsgo/zllm/agent"
)

type LLM struct {
	name  string
	model string
	agent agent.LLMAgent
}

func New(name string, model string, agent agent.LLMAgent) *LLM {
	return &LLM{
		name:  name,
		model: model,
		agent: agent,
	}
}

func (l *LLM) Name() string {
	return l.name
}

func (l *LLM) Model() string {
	return l.model
}

func (l *LLM) Generate(ctx context.Context, json []byte) (*zjson.Res, error) {
	if zjson.ValidBytes(json) {
		json, _ = zjson.SetBytes(json, "model", l.model)
	}
	return l.agent.Generate(ctx, json)
}

func (l *LLM) Stream(ctx context.Context, json []byte, callback func(data string, raw []byte)) (<-chan *zjson.Res, error) {
	if zjson.ValidBytes(json) {
		json, _ = zjson.SetBytes(json, "model", l.model)
	}
	return l.agent.Stream(ctx, json, callback)
}

var _ Openai = (*LLM)(nil)

type Openai interface {
	Name() string
	Model() string
	Generate(ctx context.Context, json []byte) (*zjson.Res, error)
	Stream(ctx context.Context, json []byte, callback func(data string, raw []byte)) (<-chan *zjson.Res, error)
}
