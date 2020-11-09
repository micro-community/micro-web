package main

import (
	_ "github.com/micro-community/micro-webui/profile"

	"github.com/micro-community/micro-webui/web"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {

	//replace default cmd
	cmd.DefaultCmd = cmd.New(cmd.Flags(web.Flags()...))
	//path parameter from ctx
	cmd.DefaultCmd.Init(cmd.Before(web.ResolveContext))

	srv := service.New(service.Name(web.Name), service.Version("latest"))

	web.ParseEnv()

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
