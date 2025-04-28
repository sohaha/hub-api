package chat

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zpool"
	"github.com/zlsgo/app_core/service"
)

type Conf struct {
	Key          string                 `z:"key"`
	keyArr       []string               `z:"-"`
	Balancer     zpool.BalancerStrategy `z:"algorithm"`
	TestInterval int64                  `z:"test_interval"`
	Fallback     map[string]string      `z:"fallback"`
}

func (Conf) ConfKey() string {
	return "chat"
}

func (Conf) DisableWrite() bool {
	return false
}

var conf = &Conf{
	Key:          "nmtx",
	Balancer:     zpool.StrategyRandom,
	TestInterval: 60000,
	Fallback: map[string]string{
		"claude-3.7-sonnet": "o4-mini",
		"o4-mini":           "gpt-4.1",
		"o3":                "gpt-4.1",
		"gpt-4.1":           "deepseek-v3",
	},
}

const providerFile = "./provider.json"

func New() (p *Module) {
	service.DefaultConf = append(service.DefaultConf, conf)

	return &Module{
		ModuleLifeCycle: service.ModuleLifeCycle{
			OnLoad: func(di zdi.Invoker) (any, error) {
				return nil, nil
			},
			OnStart: func(di zdi.Invoker) error {
				return nil
			},
			OnDone: func(di zdi.Invoker) error {
				conf.keyArr = zarray.Map(strings.Split(conf.Key, ","), func(i int, v string) string {
					return strings.TrimSpace(v)
				})
				return nil
			},
			OnStop: func(di zdi.Invoker) error {
				return nil
			},
			Service: &service.ModuleService{
				Controllers: []service.Controller{&Index{Path: "/chat"}},
				Tasks:       []service.Task{},
			},
		},
	}
}
