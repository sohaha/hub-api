package chat

import (
	"reflect"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/app_core/service"
)

type Module struct {
	log *zlog.Logger

	service.ModuleLifeCycle
}

var (
	_                = reflect.TypeOf(&Module{})
	_ service.Module = &Module{}
)

func (p *Module) Name() string {
	return "Chat"
}
