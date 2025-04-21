package internal

import (
	"app/module/chat"

	"github.com/zlsgo/app_core/service"
)

func RegModule() []service.Module {
	return []service.Module{
		chat.New(),
	}
}
