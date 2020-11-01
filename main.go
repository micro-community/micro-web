package main

import (
	"github.com/micro-community/micro-webui/auth"
	"github.com/micro-community/micro-webui/web"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/service/server/mock"
	"github.com/urfave/cli/v2"
)

func main() {

	before := func(ctx *cli.Context) error {
		_ = web.ResolveContext(ctx)
		return nil
	}

	parseFlags := cmd.New(cmd.SetupOnly(), cmd.Flags(web.Flags()...), cmd.Before(before))

	//this is a workaround
	parseFlags.App().Flags = append(parseFlags.App().Flags, web.Flags()...)
	parseFlags.Run()

	srv := service.New(service.Name(web.Name))

	//replace default server
	server.DefaultServer = mock.NewServer(server.WrapHandler(auth.NewAuthHandlerWrapper()))

	//init api server
	mweb := web.New(web.Address, srv)

	srv.Init(service.AfterStop(func() error {
		// Stop HttpWeb after srv
		if err := mweb.Stop(); err != nil {
			logger.Fatal(err)
		}
		return nil
	}))

	if err := mweb.Run(); err != nil {
		logger.Fatal(err)
	}

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

}
