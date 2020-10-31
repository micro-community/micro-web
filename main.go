package main

import (
	"github.com/micro-community/micro-webui/auth"
	"github.com/micro-community/micro-webui/server/httpweb"
	"github.com/micro-community/micro-webui/web"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/service/server/mock"
	"github.com/urfave/cli/v2"
)

//Meta Fields of micro web
var (
	// Default server name
	Name = "web"
	// Default address to bind to
	Address = ":80"
	// The namespace to serve
	// Example:
	// Namespace + /[Service]/foo/bar
	// Host: Namespace.Service Endpoint: /foo/bar
	Namespace = "micro"
	Type      = "web"
	Resolver  = "path"
	Handler   = "meta"
)

func main() {

	//register new cmds
	cmd.Register(web.Commands()...)

	cmd.DefaultCmd.Init(cmd.Before(func(ctx *cli.Context) error {
		return nil
	}))

	srv := service.New(service.Name(Name))
	//replace default server
	server.DefaultServer = mock.NewServer(server.WrapHandler(auth.NewAuthHandlerWrapper()))

	//init api server
	api := httpweb.NewServer(Address)

	srv.Init(service.AfterStop(func() error {
		// Stop HttpWeb after srv stop
		if err := api.Stop(); err != nil {
			logger.Fatal(err)
		}
		return nil
	}))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

}
