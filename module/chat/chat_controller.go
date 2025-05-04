package chat

import (
	"fmt"
	"net/http"

	"app/provider/openai"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/zutil"
)

func (h *Index) ANY(c *znet.Context) error {
	data, _ := c.GetDataRawBytes()
	if len(data) == 0 {
		data = zstring.String2Bytes(c.DefaultFormOrQuery("text", ""))
	}

	return h.chat(c, data)
}

func (h *Index) chat(c *znet.Context, data []byte) (err error) {
	if len(data) == 0 {
		return zerror.InvalidInput.Text("对话内容不能为空")
	}

	r := h.mu.RLock()
	defer h.mu.RUnlock(r)

	var (
		nodesKeys    []string
		lastErr      error
		nodes        = h.pool
		reserveNodes = h.reservePool
		ctx          = c.Request.Context()
		pools        = []string{"nodes", "reserveNodes"}
		nodeName     = "unknown"
		stream       = zjson.GetBytes(data, "stream").Bool()
		model        = zjson.GetBytes(data, "model").String()
	)

	defer func() {
		i, ok, _ := h.total.ProvideGet(nodeName, func() (*zutil.Int64, bool) {
			return zutil.NewInt64(0), true
		})
		if ok && i != nil {
			i.Add(1)
		}
	}()

	if model == "" {
		for _, k := range pools {
			if ctx.Err() != nil {
				continue
			}
			pool := nodes
			if k == "reserveNodes" {
				pool = reserveNodes
			}
			err = pool.Run(func(node openai.Openai) (normal bool, err error) {
				if ctx.Err() != nil {
					return true, ctx.Err()
				}

				defer func() {
					lastErr = err
				}()

				nodeName = node.Name()

				zlog.Tips("Node", nodeName, "--", nodes.Keys())
				normal, err = h.node(c, stream, node, data)
				return
			})
			if err == nil {
				break
			}
		}
	} else {
		for _, k := range pools {
			runModel := model
		nn:
			for {
				pool := nodes
				nodesKeys = h.modelMaps[runModel]
				if k == "reserveNodes" {
					pool = reserveNodes
					nodesKeys = h.reserveModelMaps[runModel]
				}

				if len(nodesKeys) == 0 {
					return zerror.InvalidInput.Text(fmt.Sprintf("%s model %s not supported", k, runModel))
				}

				err = pool.RunByKeys(nodesKeys, func(node openai.Openai) (normal bool, err error) {
					if ctx.Err() != nil {
						return true, ctx.Err()
					}

					defer func() {
						lastErr = err
					}()

					nodeName = node.Name()
					normal, err = h.node(c, stream, node, data)
					if err == nil {
						zlog.Tips("Node", nodeName, "--", nodesKeys)
					} else {
						zlog.Warn("Node", nodeName, err.Error(), "--", nodesKeys)
					}
					return
				}, conf.Balancer)
				if err == nil || len(conf.Fallback) == 0 {
					break nn
				}
				var ok bool
				runModel, ok = conf.Fallback[runModel]
				if !ok {
					break nn
				}
			}
			if err == nil {
				break
			}
		}
	}

	if lastErr != nil {
		return lastErr
	}

	return err
}

func (h *Index) node(c *znet.Context, stream bool, node openai.Openai, data []byte) (normal bool, err error) {
	if stream {
		sse := znet.NewSSE(c)
		done, err := node.Stream(c.Request.Context(), data, func(data string, raw []byte) {
			sse.SendByte("", raw)
		})
		if err != nil {
			return false, err
		}

		go func() {
			<-done
			sse.Stop()
		}()

		sse.Push()

		return true, err
	}

	resp, err := node.Generate(c.Request.Context(), data)
	if err != nil {
		return false, err
	}

	if !zjson.ValidBytes(data) {
		c.String(http.StatusOK, resp.Get("choices.0.message.content").String())
		return true, nil
	}

	c.Byte(http.StatusOK, resp.Bytes())
	c.SetContentType("application/json")

	return true, nil
}
