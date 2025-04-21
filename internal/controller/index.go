package controller

import (
	"github.com/zlsgo/app_core/service"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type Index struct {
	service.App
}

func (h *Index) Init(r *znet.Engine) error {
	// 开放静态资源目录
	r.Static("/static/", zfile.RealPath("./static"))

	return nil
}

func (h *Index) GET(c *znet.Context) (ztype.Map, error) {
	return ztype.Map{"hello": "world"}, nil
}
