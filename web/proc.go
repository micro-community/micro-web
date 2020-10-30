package web

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	"golang.org/x/net/publicsuffix"
)

type webService struct {
	Name string
	Link string
	Icon string // TODO: lookup icon
}

var (
	regs registry.Registry
)

// ServeHTTP serves the web dashboard and proxies where appropriate
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	return
	// otherwise serve the proxy
	//s.prx.ServeHTTP(w, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		return
	}

	services, err := regs.ListServices(regs.ListContext(r.Context()))
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
	render(w, r, indexTemplate, data)
}

func registryHandler(w http.ResponseWriter, r *http.Request) {

}

func callHandler(w http.ResponseWriter, r *http.Request) {

}
