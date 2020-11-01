// Package web is a web dashboard
package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/micro-community/micro-webui/handler/meta"
	"github.com/micro-community/micro-webui/resolver"
	"github.com/micro-community/micro-webui/resolver/path"
	"github.com/micro-community/micro-webui/router"
	regRouter "github.com/micro-community/micro-webui/router/registry"
	"github.com/micro-community/micro-webui/server"
	"github.com/micro-community/micro-webui/server/httpweb"

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

type srvWeb struct {
	svc      *service.Service
	api      server.Server
	rr       resolver.Resolver
	rt       router.Router
	registry registry.Registry
	logged   bool
}

func New(address string, service *service.Service) *srvWeb {

	rr := path.NewResolver(resolver.WithServicePrefix(Namespace), resolver.WithHandler(Handler))
	rt := regRouter.NewRouter(router.WithResolver(rr), router.WithRegistry(registry.DefaultRegistry))

	return &srvWeb{
		api: httpweb.NewServer(address),
		rr:  rr,
		rt:  rt,
		svc: service,
	}

}

//Run run micro web
func (s *srvWeb) Run() error {

	logger.Init(logger.WithFields(map[string]interface{}{"service": "web"}))

	//	ResolveContext(ctx)

	var h http.Handler
	r := mux.NewRouter()
	h = r

	logger.Infof("Registering API & Web Handler at %s", APIPath)

	//rt := regRouter.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}

		response := fmt.Sprintf(`{"version": "%s"}`, s.svc.Version())
		w.Write([]byte(response))
	})

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	r.HandleFunc("/client", s.CallHandler)
	r.HandleFunc("/services", s.RegistryHandler)
	r.HandleFunc("/service/{name}", s.RegistryHandler)
	//r.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(p)

	r.PathPrefix(APIPath).Handler(meta.NewMetaHandler(s.svc.Client(), s.rt, Namespace))

	// register all the http handler plugins
	for _, p := range plugin.Plugins() {
		if v := p.Handler(); v != nil {
			h = v(h)
		}
	}
	// append the auth wrapper
	//h = auth.Wrapper(rr, Namespace)(h)

	// register the handler
	s.api.Handle("/", h)

	// Start API
	return s.api.Start()
}

// s.HandleFunc("/favicon.ico", faviconHandler)
// s.HandleFunc("/client", s.callHandler)
// s.HandleFunc("/services", s.registryHandler)
// s.HandleFunc("/service/{name}", s.registryHandler)
// s.HandleFunc("/rpc", handler.RPC)
// s.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(p)
// s.HandleFunc("/", s.indexHandler)

func (s *srvWeb) Stop() error {
	return s.api.Stop()
}
