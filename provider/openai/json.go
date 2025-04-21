package openai

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zpool"
	"github.com/sohaha/zlsgo/zreflect"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zllm/agent"
	"github.com/zlsgo/zllm/message"
)

type jsonLLM struct {
	agent     agent.LLMAgent
	realModel string
}

func (d *jsonLLM) Generate(ctx context.Context, data []byte) (resp *zjson.Res, err error) {
	return d.agent.Generate(ctx, data)
}

func (d *jsonLLM) Stream(ctx context.Context, data []byte, callback func(string, []byte)) (done <-chan *zjson.Res, err error) {
	return d.agent.Stream(ctx, data, callback)
}

func (d *jsonLLM) ParseResponse(resp *zjson.Res) (*agent.Response, error) {
	return d.agent.ParseResponse(resp)
}

func (d *jsonLLM) PrepareRequest(messages *message.Messages, options ...func(ztype.Map) ztype.Map) (body []byte, err error) {
	return d.agent.PrepareRequest(messages, options...)
}

var _ agent.LLMAgent = &jsonLLM{}

func ParseMap(name string, data ztype.Map, nodes *zpool.Balancer[Openai], modelsMap *map[string][]string) (err error) {
	base := data.Get("base").String()
	keys := data.Get("key").String()
	apiurl := data.Get("apiurl")
	cooldown := data.Get("cooldown").Int()
	if cooldown == 0 {
		cooldown = 6000
	}
	weight := data.Get("weight").Int()
	if weight == 0 {
		weight = 1
	}
	maxConns := data.Get("max").Int()
	stream := data.Get("stream")

	modelData := data.Get("models")
	models := modelData.Map()
	if len(models) == 0 {
		return errors.New("models is empty")
	}
	if zreflect.ValueOf(modelData.Value()).Kind() == reflect.Slice {
		models = make(ztype.Map)
		for _, v := range modelData.SliceString() {
			models[v] = v
		}
	}
	for model := range models {
		nodeName := fmt.Sprintf("%s:%s", strings.ReplaceAll(name, ":", "\\:"), model)
		modelName := model
		realModel := models.Get(model).String()
		if realModel != "" {
			modelName = realModel
		}
		if _, ok := (*modelsMap)[model]; !ok {
			(*modelsMap)[model] = []string{}
		}
		(*modelsMap)[model] = append((*modelsMap)[model], nodeName)

		nodes.Add(nodeName, New(
			nodeName,
			modelName,
			agent.NewOpenAIProvider(func(oa *agent.OpenAIOptions) {
				oa.Model = modelName
				oa.APIKey = keys
				if base != "" {
					oa.BaseURL = base
				}
				if apiurl.Exists() {
					oa.APIURL = apiurl.String()
				}
				if stream.Exists() {
					oa.Stream = stream.Bool()
				}
			}),
		), func(opts *zpool.BalancerNodeOptions) {
			opts.Weight = uint64(weight)
			opts.Cooldown = int64(cooldown)
			opts.MaxConns = int64(maxConns)
		})
	}
	return
}
