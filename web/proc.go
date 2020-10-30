package web

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"sort"
	"strings"

	"github.com/micro-community/micro-webui/resolver"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	"golang.org/x/net/publicsuffix"
)

type webService struct {
	Name string
	Link string
	Icon string // TODO: lookup icon
}

// ServeHTTP serves the web dashboard and proxies where appropriate
func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Host) == 0 {
		r.URL.Host = r.Host
	}

	if len(r.URL.Scheme) == 0 {
		r.URL.Scheme = "http"
	}

	// no host means dashboard
	host := r.URL.Hostname()
	if len(host) == 0 {
		h, _, err := net.SplitHostPort(r.Host)
		if err != nil && strings.Contains(err.Error(), "missing port in address") {
			host = r.Host
		} else if err == nil {
			host = h
		}
	}

	// check again
	if len(host) == 0 {
		s.Router.ServeHTTP(w, r)
		return
	}

	// check based on host set
	if len(Host) > 0 && Host == host {
		s.Router.ServeHTTP(w, r)
		return
	}

	// an ip instead of hostname means dashboard
	ip := net.ParseIP(host)
	if ip != nil {
		s.Router.ServeHTTP(w, r)
		return
	}

	// namespace matching host means dashboard
	parts := strings.Split(host, ".")
	reverse(parts)
	namespace := strings.Join(parts, ".")

	// replace mu since we know its ours
	if strings.HasPrefix(namespace, "mu.micro") {
		namespace = strings.Replace(namespace, "mu.micro", "go.micro", 1)
	}

	// web dashboard if namespace matches
	if namespace == Namespace+"."+Type {
		s.Router.ServeHTTP(w, r)
		return
	}

	// if a host has no subdomain serve dashboard
	v, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil || v == host {
		s.Router.ServeHTTP(w, r)
		return
	}

	// check if its a web request
	if _, _, isWeb := s.resolver.Info(r); isWeb {
		s.Router.ServeHTTP(w, r)
		return
	}
	// otherwise serve the proxy
	s.prx.ServeHTTP(w, r)
}

// proxy is a http reverse proxy
func (s *srv) proxy() *proxy {
	director := func(r *http.Request) {
		kill := func() {
			r.URL.Host = ""
			r.URL.Path = ""
			r.URL.Scheme = ""
			r.Host = ""
			r.RequestURI = ""
		}

		// check to see if the endpoint was encoded in the request context
		// by the auth wrapper
		var endpoint *resolver.Endpoint
		if val, ok := (r.Context().Value(resolver.Endpoint{})).(*resolver.Endpoint); ok {
			endpoint = val
		}

		// TODO: better error handling
		var err error
		if endpoint == nil {
			if endpoint, err = s.resolver.Resolve(r); err != nil {
				fmt.Printf("Failed to resolve url: %v: %v\n", r.URL, err)
				kill()
				return
			}
		}

		r.Header.Set(BasePathHeader, "/"+endpoint.Name)
		r.URL.Host = endpoint.Host
		r.URL.Path = endpoint.Path
		r.URL.Scheme = "http"
		r.Host = r.URL.Host
	}

	return &proxy{
		Router:   &httputil.ReverseProxy{Director: director},
		Director: director,
	}
}

func (s *srv) indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		return
	}

	services, err := s.registry.ListServices(registry.ListContext(r.Context()))
	if err != nil {
		logger.Errorf("Error listing services: %v", err)
	}

	// if the resolver is subdomain, we will need the domain
	domain, _ := publicsuffix.EffectiveTLDPlusOne(r.URL.Hostname())

	var webServices []webService
	for _, srv := range services {
		// not a web app
		comps := strings.Split(srv.Name, ".web.")
		if len(comps) == 1 {
			continue
		}
		name := comps[1]

		link := fmt.Sprintf("/%v/", name)
		if Resolver == "subdomain" && len(domain) > 0 {
			link = fmt.Sprintf("https://%v.%v", name, domain)
		}

		// in the case of 3 letter things e.g m3o convert to M3O
		if len(name) <= 3 && strings.ContainsAny(name, "012345789") {
			name = strings.ToUpper(name)
		}

		webServices = append(webServices, webService{Name: name, Link: link})
	}

	sort.Slice(webServices, func(i, j int) bool { return webServices[i].Name < webServices[j].Name })

	type templateData struct {
		HasWebServices bool
		WebServices    []webService
	}

	data := templateData{len(webServices) > 0, webServices}
	s.render(w, r, indexTemplate, data)
}

func (s *srv) registryHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *srv) callHandler(w http.ResponseWriter, r *http.Request) {

}
