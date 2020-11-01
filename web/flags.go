package web

import (
	"strings"

	"github.com/micro/micro/v3/plugin"
	"github.com/micro/micro/v3/service"
	"github.com/urfave/cli/v2"
)

//ResolveContext for web
func ResolveContext(ctx *cli.Context) {

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
}

//Flags for `micro web flag`
func Flags(options ...service.Option) []cli.Flag {

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Set the web UI address e.g 0.0.0.0:8082",
			EnvVars: []string{"MICRO_WEB_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "namespace",
			Usage:   "Set the namespace used by Web e.g. arch.wiki",
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
	}

	//apppend flags from plugins
	for _, p := range plugin.Plugins() {
		if flags := p.Flags(); len(flags) > 0 {
			flags = append(flags, flags...)
		}
	}

	return flags
}
