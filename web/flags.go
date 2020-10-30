package web

import (
	"strings"

	"github.com/micro/micro/v3/plugin"
	"github.com/micro/micro/v3/service"
	"github.com/urfave/cli/v2"
)

func resolveContext(ctx *cli.Context) {

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("resolver")) > 0 {
		Resolver = ctx.String("resolver")
	}
	if len(ctx.String("type")) > 0 {
		Type = ctx.String("type")
	}
	if len(ctx.String("namespace")) > 0 {
		// remove the service type from the namespace to allow for
		// backwards compatability
		Namespace = strings.TrimSuffix(ctx.String("namespace"), "."+Type)
	}
	// Init plugins
	for _, p := range plugin.Plugins() {
		p.Init(ctx)
	}
}

//Commands for `micro web`
func Commands(options ...service.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "web",
		Usage: "Run the web dashboard",
		Action: func(c *cli.Context) error {
			Run(c, options...)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the web UI address e.g 0.0.0.0:8082",
				EnvVars: []string{"MICRO_WEB_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "namespace",
				Usage:   "Set the namespace used by the Web proxy e.g. m3o.org",
				EnvVars: []string{"MICRO_WEB_NAMESPACE"},
			},
			&cli.StringFlag{
				Name:    "resolver",
				Usage:   "Set the resolver to route to services e.g path, domain",
				EnvVars: []string{"MICRO_WEB_RESOLVER"},
			},
			&cli.StringFlag{
				Name:    "auth_login_url",
				EnvVars: []string{"MICRO_AUTH_LOGIN_URL"},
				Usage:   "The relative URL where a user can login",
			},
		},
	}

	for _, p := range plugin.Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []*cli.Command{command}
}
