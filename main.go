package main

import (
	"github.com/micro-community/micro-webui/auth"
	_ "github.com/micro-community/micro-webui/profile"
	"github.com/micro-community/micro-webui/web"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/service/server/mock"
)

func main() {

	srv := service.New(service.Name(web.Name))
	//replace default server
	server.DefaultServer = mock.NewServer(server.WrapHandler(auth.NewAuthHandlerWrapper()))

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
