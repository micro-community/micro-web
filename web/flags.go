package web

import (
	"github.com/micro/micro/v3/service/config"
)

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
