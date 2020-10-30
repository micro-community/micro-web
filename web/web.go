// Package web is a web dashboard
package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"

	"github.com/micro-community/micro-webui/handler"
	"github.com/micro-community/micro-webui/namespace"
	"github.com/micro-community/micro-webui/resolver"
	"github.com/micro-community/micro-webui/resolver/path"
	"github.com/micro-community/micro-webui/router"

	regRouter "github.com/micro-community/micro-webui/router/registry"

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

type srv struct {
	*mux.Router
	// registry we use
	registry registry.Registry
	// the resolver
	resolver *resolver.Resolver
	// the namespace resolver
	nsResolver *namespace.Resolver
	// the proxy server
	prx *proxy

	logged bool
}


//Run run micro web
func Run(ctx *cli.Context, srvOpts ...service.Option) {

	logger.Init(logger.WithFields(map[string]interface{}{"service": "web"}))

	resolveContext(ctx)

	// create the router
	//	var h http.Handler
	r := mux.NewRouter()
	//	h = r

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

	r.HandleFunc("/client", s.callHandler)
	r.HandleFunc("/services", s.registryHandler)
	r.HandleFunc("/service/{name}", s.registryHandler)
	//r.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(p)
	r.PathPrefix(APIPath).Handler(handler.Meta(srv.Client(), rt, Namespace))

		// register all the http handler plugins
	for _, p := range plugin.Plugins() {
		if v := p.Handler(); v != nil {
			h = v(h)
		}
	}

	h = auth.Wrapper(rr, Namespace)(h)


}

func (s *srv) render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"format": format,
		"Title":  strings.Title,
		"First": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.Title(string(s[0]))
		},
	}).Parse(layoutTemplate)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
	t, err = t.Parse(tmpl)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	// If the user is logged in, render Account instead of Login
	loginTitle := "Login"
	user := ""

	if c, err := r.Cookie(TokenCookieName); err == nil && c != nil {
		token := strings.TrimPrefix(c.Value, TokenCookieName+"=")
		//	if acc, err := s.auth.Inspect(token); err == nil {
		if len(token) > 0 && s.logged {
			loginTitle = "Account"
			//user = acc.ID
		}
	}

	if err := t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"LoginTitle": loginTitle,
		"LoginURL":   loginURL,
		"StatsURL":   statsURL,
		"Results":    data,
		"User":       user,
	}); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
	}
}
