package main

import (
	"context"

	"app/internal"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zutil"
	"github.com/sohaha/zlsgo/zutil/daemon"
	"github.com/zlsgo/app_core/common"
	"github.com/zlsgo/app_core/service"
)

func init() {
	service.AppName = "HubAI"
	service.ConfFileName = "config.toml"

	zcli.Version = "0.1.0"
	zcli.Name = service.AppName
	zcli.EnableDetach = true
}

func main() {
	if conf, err := setup(); err != nil {
		if conf == nil || !conf.Base.Debug {
			zcli.Error("%s", err.Error())
		} else {
			zlog.Errorf("%+v\n", err)
		}
	}
}

func setup() (conf *service.Conf, err error) {
	_ = zutil.Loadenv()
	err = zutil.TryCatch(func() (err error) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		di := internal.InitDI(ctx)

		err = zcli.LaunchServiceRun(zcli.Name, "", func() {
			conf, err = internal.Init(di, true)
			common.Fatal(err)
			common.Fatal(internal.Start(di))
		}, &daemon.Config{Context: ctx})

		_, _ = di.Invoke(internal.Stop)
		return err
	})

	return
}
