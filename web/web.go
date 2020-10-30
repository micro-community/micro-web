// Package web is a web dashboard
package web

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/micro/micro/v3/plugin"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
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
	// Base path sent to web service.
	// This is stripped from the request path
	// Allows the web service to define absolute paths
	ProxyPath             = "/{service:[a-zA-Z0-9]+}"
	BasePathHeader        = "X-Micro-Web-Base-Path"
	statsURL              string
	loginURL              string
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"

	// Host name the web dashboard is served on
	Host, _ = os.Hostname()
)

type srv struct {
	*mux.Router
	// the proxy server
	prx *proxy
	// auth service
	auth auth.Auth
}

type reg struct {
	registry.Registry

	sync.RWMutex
	lastPull time.Time
	services []*registry.Service
}

//Run run micro web
func Run(ctx *cli.Context, srvOpts ...service.Option) {

	logger.Init(logger.WithFields(map[string]interface{}{"service": "web"}))

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

	// create the router
	var h http.Handler
	r := mux.NewRouter()
	h = r

}

// ServeHTTP serves the web dashboard and proxies where appropriate
func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func format(v *registry.Value) string {
	if v == nil || len(v.Values) == 0 {
		return "{}"
	}
	var f []string
	for _, k := range v.Values {
		f = append(f, formatEndpoint(k, 0))
	}
	return fmt.Sprintf("{\n%s}", strings.Join(f, ""))
}

func (s *srv) indexHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *srv) registryHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *srv) callHandler(w http.ResponseWriter, r *http.Request) {

}
