package chat

import (
	"app/provider/openai"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zpool"
)

func ParseNode(config []byte, fallback bool) (nodes *zpool.Balancer[openai.Openai], modelMaps map[string][]string, inlayErrors, loadErrors map[string]string) {
	loadErrors = make(map[string]string)
	inlayErrors = make(map[string]string)
	modelMaps = make(map[string][]string, 0)
	nodes = zpool.NewBalancer[openai.Openai]()

	if config == nil {
		config, _ = zfile.ReadFile(providerFile)
	}

	zjson.ParseBytes(config).ForEach(func(key *zjson.Res, value *zjson.Res) bool {
		name := key.String()
		data := value.Map()
		if fallback && data.Get("fallback").Bool() || !fallback && !data.Get("fallback").Bool() {
			err := openai.ParseMap(name, data, nodes, &modelMaps)
			if err != nil {
				loadErrors[name] = err.Error()
			}
		}

		return true
	})

	return nodes, modelMaps, inlayErrors, loadErrors
}
