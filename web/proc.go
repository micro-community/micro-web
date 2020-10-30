package web

import (
	"fmt"
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
	s.render(w, r, indexTemplate, data)
}

func registryHandler(w http.ResponseWriter, r *http.Request) {

}

func callHandler(w http.ResponseWriter, r *http.Request) {

}
