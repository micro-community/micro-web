// Package web is a web dashboard
package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"

	"github.com/micro-community/micro-webui/handler/meta"
	"github.com/micro-community/micro-webui/resolver"
	"github.com/micro-community/micro-webui/resolver/path"
	"github.com/micro-community/micro-webui/router"

	regRouter "github.com/micro-community/micro-webui/router/registry"
	httpserver "github.com/micro-community/micro-webui/server/http"

	"github.com/micro/micro/v3/plugin"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
)

const (
	// BearerScheme used for Authorization header
	BearerScheme = "Bearer "
	// TokenCookieName is the name of the cookie which stores the auth token
	TokenCookieName = "micro-token"
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
	// Base path sent to web service.
	// This is stripped from the request path
	// Allows the web service to define absolute paths
	APIPath               = "/{service:[a-zA-Z0-9]+}"
	BasePathHeader        = "X-Micro-Web-Base-Path"
	statsURL              string
	loginURL              string
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"

	// Host name the web dashboard is served on
	Host, _ = os.Hostname()
)

//Run run micro web
func Run(ctx *cli.Context, srvOpts ...service.Option) {

	logger.Init(logger.WithFields(map[string]interface{}{"service": "web"}))

	resolveContext(ctx)

	var h http.Handler
	r := mux.NewRouter()
	h = r

	logger.Infof("Registering API Web Handler at %s", APIPath)

	//rt := regRouter.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}
		response := fmt.Sprintf(`{"version": "%s"}`, ctx.App.Version)
		w.Write([]byte(response))
	})

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	rr := path.NewResolver(resolver.WithServicePrefix(Namespace), resolver.WithHandler(Handler))
	rt := regRouter.NewRouter(router.WithResolver(rr), router.WithRegistry(registry.DefaultRegistry))
	// initialize service
	srv := service.New(service.Name(Name))

	// r.HandleFunc("/client", s.callHandler)
	// r.HandleFunc("/services", s.registryHandler)
	// r.HandleFunc("/service/{name}", s.registryHandler)
	//r.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(p)
	r.PathPrefix(APIPath).Handler(meta.NewMetaHandler(srv.Client(), rt, Namespace))

	// register all the http handler plugins
	for _, p := range plugin.Plugins() {
		if v := p.Handler(); v != nil {
			h = v(h)
		}
	}

	// append the auth wrapper
	//h = auth.Wrapper(rr, Namespace)(h)

	// create a new api server with wrappers
	api := httpserver.NewServer(Address)

	// register the handler
	api.Handle("/", h)

	// Start API
	if err := api.Start(); err != nil {
		logger.Fatal(err)
	}

	// Run server
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	// Stop API
	if err := api.Stop(); err != nil {
		logger.Fatal(err)
	}

}
