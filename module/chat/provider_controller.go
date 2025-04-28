package chat

import (
	"maps"
	"net/http"

	"app/provider/openai"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zhttp"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zpool"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zutil"
)

func (h *Index) GETProvider(c *znet.Context) any {
	r := h.mu.RLock()
	pool := h.pool
	h.mu.RUnlock(r)

	nodes := make(map[string]bool, pool.Len())
	pool.WalkNodes(func(node openai.Openai, available bool) (normal bool) {
		nodes[node.Name()] = available
		return true
	})
	return nodes
}

func (h *Index) POSTProvider(c *znet.Context) error {
	var body []byte

	body, _ = c.GetDataRawBytes()
	if len(body) == 0 {
		url := c.GetJSON("url").String()
		header := zhttp.Header{}
		if len(conf.keyArr) > 0 {
			header["Authorization"] = "bearer " + conf.keyArr[0]
		}
		resp, err := zhttp.Get(url, header, nil, c.Request.Context())
		if err != nil {
			return err
		}
		body = resp.Bytes()
	}

	parseErrs := make(map[string]string, 0)

	pool, modelMaps, _, errs := ParseNode(body, false)
	if len(errs) > 0 {
		parseErrs = errs
	}

	reservePool, reserveModelMaps, _, reserveErrs := ParseNode(body, true)
	if len(reserveErrs) > 0 {
		maps.Copy(parseErrs, reserveErrs)
	}

	if len(parseErrs) > 0 {
		c.ApiJSON(http.StatusBadRequest, "", parseErrs)
		return nil
	}

	h.mu.Lock()
	h.pool = pool
	h.modelMaps = modelMaps
	h.reservePool = reservePool
	h.reserveModelMaps = reserveModelMaps
	h.mu.Unlock()
	h.chatModels(c)

	_ = zfile.WriteFile(providerFile, body)

	return nil
}

func (h *Index) GETTotal(c *znet.Context) any {
	models := map[string]int64{}
	h.total.ForEach(func(k string, v *zutil.Int64) bool {
		models[k] = v.Load()
		return true
	})
	return models
}

func (h *Index) GETNodes(c *znet.Context) ztype.Map {
	r := h.mu.RLock()
	pools := []*zpool.Balancer[openai.Openai]{h.pool, h.reservePool}
	h.mu.RUnlock(r)

	failedNodes := []string{}
	nodes := []string{}
	for i := range pools {
		pools[i].WalkNodes(func(node openai.Openai, available bool) (normal bool) {
			if !available {
				failedNodes = append(failedNodes, node.Name())
			} else {
				nodes = append(nodes, node.Name())
			}

			return available
		})
	}

	return ztype.Map{
		"available": nodes,
		"failed":    failedNodes,
	}
}
