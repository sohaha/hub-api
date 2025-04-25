package chat

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zpool"
	"github.com/zlsgo/app_core/service"
)

type Conf struct {
	Key          string                 `z:"key"`
	keyArr       []string               `z:"-"`
	Balancer     zpool.BalancerStrategy `z:"algorithm"`
	TestInterval int64                  `z:"test_interval"`
}

func (Conf) ConfKey() string {
	return "chat"
}

func (Conf) DisableWrite() bool {
	return false
}

var conf = &Conf{
	Key:          "sk-sb123",
	Balancer:     zpool.StrategyRandom,
	TestInterval: 60000,
}

const providerFile = "./provider.json"

func New() (p *Module) {
	service.DefaultConf = append(service.DefaultConf, conf)

	return &Module{
		ModuleLifeCycle: service.ModuleLifeCycle{
			OnLoad: func(di zdi.Invoker) (any, error) {
				if !zfile.FileExist(providerFile) {
					zfile.WriteFile(providerFile, []byte(`{ 
    "Openai": {
        "base": "https://api.openai.com/v1",
        "models": {"4o-mini": "gpt-4o-mini"},
        "key": "sk-1,sk-2",
        "cooldown": 6000,
        "weight": 1,
        "max": 10
    }
}`))
				}

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
