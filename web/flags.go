package web

import (
	"strings"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/urfave/cli/v2"
)

//ResolveContext for web
func ResolveContext(ctx *cli.Context) error {

	if len(ctx.String("web_name")) > 0 {
		Name = ctx.String("web_name")
	}
	if len(ctx.String("web_address")) > 0 {
		Address = ctx.String("web_address")
	}
	if len(ctx.String("resolver")) > 0 {
		Resolver = ctx.String("resolver")
	}
	if len(ctx.String("type")) > 0 {
		Type = ctx.String("type")
	}
	if len(ctx.String("namespace")) > 0 {
		// remove the service type from the namespace to allow for backwards compatability
		Namespace = strings.TrimSuffix(ctx.String("namespace"), "."+Type)
	}

	return nil
}

//Flags for `micro web flag`
func Flags(options ...service.Option) []cli.Flag {

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "web_address",
			Usage:   "Set the web UI address e.g 0.0.0.0:8082",
			EnvVars: []string{"MICRO_WEB_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "web_name",
			Usage:   "Set the server_name of web",
			EnvVars: []string{"MICRO_WEB_NAME"},
		},
		&cli.StringFlag{
			Name:    "web_resolver",
			Usage:   "Set the resolver to route to services e.g path, domain",
			EnvVars: []string{"MICRO_WEB_RESOLVER"},
		},
		&cli.StringFlag{
			Name:    "auth_login_url",
			EnvVars: []string{"MICRO_AUTH_LOGIN_URL"},
			Usage:   "The relative URL where a user can login",
		},
	}

	return flags
}

//ParseEnv from env
func ParseEnv() {

	resolverV, err := config.Get("resolver")
	resolverValue := resolverV.String("")

	if err == nil && resolverValue != "" {
		Resolver = resolverValue
	}

	typeV, err := config.Get("type")
	typeValue := typeV.String("")

	if err == nil && typeValue != "" {
		Type = typeValue
	}

}
