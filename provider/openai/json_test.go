package openai

import (
	"context"
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zpool"
	"github.com/sohaha/zlsgo/ztype"
)

func Test_loadNode(t *testing.T) {
	tt := zlsgo.NewTest(t)

	modelsMap := make(map[string][]string, 0)
	nodes := zpool.NewBalancer[Openai]()
	errs := ParseMap("Unlimited", ztype.Map{
		"base": "https://fuck-u-altman-and-openai.deno.dev/v1",
		"models": ztype.Map{
			"deekseek-v3":  "gpt-4.1-mini",
			"gpt-4.1-nano": "",
		},
		"weight": 10,
		"max":    3,
	}, nodes, &modelsMap)
	tt.Log(errs)
	tt.Log(nodes.Len())
	tt.Log(nodes.Keys())

	data := []byte(`{"model":"deekseek-v3","messages":[{"role":"user","content":"树上有 9 只鸟，猎人开枪打死一只，树上还剩下多少只鸟？"}],"stream":true}`)
	nodekeys := modelsMap[zjson.ParseBytes(data).Get("model").String()]
	tt.Log("nodekeys", nodekeys)
	err := nodes.RunByKeys(nodekeys, func(node Openai) (normal bool, err error) {
		resp, err := node.Generate(context.Background(), data)
		tt.Log(resp.String())
		tt.NoError(err, true)
		return true, nil
	})
	tt.NoError(err, true)

	data = []byte(`{"model":"gpt-4.1-nano","messages":[{"role":"user","content":"树上有 9 只鸟，猎人开枪打死一只，树上还剩下多少只鸟？"}],"stream":true}`)
	nodekeys = modelsMap[zjson.ParseBytes(data).Get("model").String()]
	tt.Log("nodekeys", nodekeys)
	err = nodes.RunByKeys(nodekeys, func(node Openai) (normal bool, err error) {
		resp, err := node.Stream(context.Background(), data, func(data string, raw []byte) {
			tt.Log("--", data)
		})
		tt.NoError(err, true)
		tt.Log(<-resp)
		return true, nil
	})
	tt.NoError(err, true)
}

func Test_Prve(t *testing.T) {
	tt := zlsgo.NewTest(t)

	tt.Log(zarray.Slice[string]("", ","))
	tt.Log(zarray.Slice[string](",", ","))
	tt.Log(zarray.Slice[string]("   ", ","))
}
