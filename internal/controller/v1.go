package controller

import (
	"reflect"

	"github.com/sohaha/zlsgo/znet"
	"github.com/zlsgo/app_core/service"
)

type V1 struct {
	Path string
	service.App
}

var _ = reflect.TypeOf(V1{})

func (h *V1) Init(r *znet.Engine) error {
	r.Any("/*", func(c *znet.Context) error {
		var ok bool
		url, err := r.GenerateURL(c.Request.Method, c.Request.URL.Path, map[string]string{})
		if err == nil {
			ok = c.Engine.FindHandle(c, c.Request, url, true)
		}
		if !ok {
			r.NotFoundHandler(c)
			return nil
		}
		return nil
	})
	return nil
}
