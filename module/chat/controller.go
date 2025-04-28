package chat

import (
	"context"
	"reflect"
	"strings"
	"time"

	"app/provider/openai"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zpool"
	"github.com/sohaha/zlsgo/zsync"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zutil"
)

type Index struct {
	pool             *zpool.Balancer[openai.Openai]
	reservePool      *zpool.Balancer[openai.Openai]
	modelMaps        map[string][]string
	reserveModelMaps map[string][]string
	mu               *zsync.RBMutex
	Path             string
	total            *zarray.Maper[string, *zutil.Int64]
}

var _ = reflect.TypeOf(&Index{})

func authMiddleware() func(c *znet.Context) error {
	return func(c *znet.Context) error {
		if len(conf.Key) != 0 {
			token := c.GetHeader("Authorization")
			if token == "" {
				token = c.DefaultFormOrQuery("token", "")
			} else {
				b := strings.SplitN(token, "Bearer ", 2)
				if len(b) == 2 {
					token = b[1]
				}
			}
			if !zarray.Contains(conf.keyArr, token) {
				return zerror.Unauthorized.Text("Invalid API key")
			}
		}
		c.Next()
		return nil
	}
}

func (h *Index) Init(r *znet.Engine) error {
	var inlayErrors map[string]string
	var reservenlayErrors map[string]string

	r.Use(authMiddleware())

	// 兼容 OPENAI 接口
	_ = r.POSTAndName("/completions", h.ANY, "/v1/chat/completions")
	_ = r.GETAndName("/models", h.chatModels, "/v1/models")

	h.mu = zsync.NewRBMutex()
	h.pool, h.modelMaps, inlayErrors, _ = ParseNode(nil, false)
	h.reservePool, h.reserveModelMaps, reservenlayErrors, _ = ParseNode(nil, true)
	h.total = zarray.NewHashMap[string, *zutil.Int64]()

	for _, v := range []map[string]string{inlayErrors, reservenlayErrors} {
		if len(v) > 0 {
			for name, err := range v {
				zlog.Error(name, err)
			}
		}
	}

	go func() {
		for {
			testInterval := min(conf.TestInterval, 5000)
			time.Sleep(time.Duration(testInterval) * time.Millisecond)
			h.testNodes()
		}
	}()

	return nil
}

func (h *Index) testNodes() {
	r := h.mu.RLock()
	pools := []*zpool.Balancer[openai.Openai]{h.pool, h.reservePool}
	h.mu.RUnlock(r)

	for i := range pools {
		pools[i].WalkNodes(func(node openai.Openai, available bool) (normal bool) {
			if available {
				return true
			}

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			resp, err := node.Generate(ctx, []byte("写一个10个字的冷笑话"))
			if err != nil {
				zlog.Error("test node: ", node.Name(), " err: ", err)
				return false
			}
			zlog.Info("test node: ", node.Name(), " resp: ", resp)
			return true
		})
	}
}

func (h *Index) chatModels(c *znet.Context) {
	r := h.mu.RLock()
	defer h.mu.RUnlock(r)

	pool := h.pool
	minLen := pool.Len()

	availableNodes := make(map[string]struct{}, minLen)
	pool.WalkNodes(func(node openai.Openai, available bool) (normal bool) {
		if !available {
			return false
		}

		availableNodes[node.Name()] = struct{}{}
		return true
	})

	models := make([]string, 0, minLen)
	for _, v := range []map[string][]string{h.modelMaps, h.reserveModelMaps} {
		for model, nodes := range v {
			for _, node := range nodes {
				if _, ok := availableNodes[node]; ok {
					models = append(models, model)
					break
				}
			}
		}
	}

	now := ztime.Time().Unix()
	resp, _ := zjson.Set(`{"object":"list"}`, "data", zarray.Map(zarray.Unique(models), func(_ int, v string) ztype.Map {
		return ztype.Map{
			"id":       v,
			"object":   "model",
			"owned_by": "system",
			"created":  now,
		}
	}, 10))
	c.String(200, resp)
	c.SetContentType(znet.ContentTypeJSON)
}
