package main

import (
	"github.com/micro-community/micro-webui/auth"
	"github.com/micro-community/micro-webui/web"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/service/server/mock"
	"github.com/urfave/cli/v2"
)

func main() {

	cmd.Register(web.Commands()...)

	cmd.DefaultCmd.Init(cmd.Before(func(ctx *cli.Context) error {
		return prepareSomething()
	}))

	server.DefaultServer = mock.NewServer(server.WrapHandler(auth.NewAuthHandlerWrapper()))

	if err := cmd.DefaultCmd.Run(); err != nil {
		logger.Fatal(err)
	}
}

func prepareSomething() error {

	return nil
}
