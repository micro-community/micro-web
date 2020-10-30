package web

import (
	"errors"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/micro-community/micro-webui/resolver"
	"github.com/micro/micro/v3/service/router"
)

var re = regexp.MustCompile("^[a-zA-Z0-9]+([a-zA-Z0-9-]*[a-zA-Z0-9]*)?$")

type Resolver struct {
	// Options
	Options resolver.Options
	// selector to choose from a pool of nodes
	// Selector selector.Selector
	// router to lookup routes
	Router router.Router
}

func (r *Resolver) String() string {
	return "web/resolver"
}

// Resolve replaces the values of Host, Path, Scheme to calla backend service
// It accounts for subdomains for service names based on namespace
func (r *Resolver) Resolve(req *http.Request, opts ...resolver.ResolveOption) (*resolver.Endpoint, error) {
	// parse the options
	options := resolver.NewResolveOptions(opts...)

	parts := strings.Split(req.URL.Path, "/")
	if len(parts) < 2 {
		return nil, errors.New("unknown service")
	}

	if !re.MatchString(parts[1]) {
		return nil, resolver.ErrInvalidPath
	}

	name := parts[1]
	if len(r.Options.ServicePrefix) > 0 {
		name = r.Options.ServicePrefix + "." + name
	}

	// lookup the routes for the service
	query := []router.LookupOption{
		router.LookupNetwork(options.Domain),
	}

	routes, err := r.Router.Lookup(name, query...)
	if err == router.ErrRouteNotFound {
		return nil, resolver.ErrNotFound
	} else if err != nil {
		return nil, err
	} else if len(routes) == 0 {
		return nil, resolver.ErrNotFound
	}

	// select a random route to use
	// todo: update to use selector once go-micro has updated the interface
	// route, err := r.Selector.Select(routes...)
	// if err != nil {
	// 	return nil, err
	// }
	route := routes[rand.Intn(len(routes))]

	// we're done
	return &resolver.Endpoint{
		Name:   name,
		Method: req.Method,
		Host:   route.Address,
		Path:   "/" + strings.Join(parts[2:], "/"),
		Domain: options.Domain,
	}, nil
}

// NewResolver creates a new micro resolver
func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &Resolver{Options: resolver.NewOptions(opts...)}
}

// Info checks whether this is a web request.
// It returns host, namespace and whether its internal
func (r *Resolver) Info(req *http.Request) (string, string, bool) {
	// set to host
	host := req.URL.Hostname()

	// set as req.Host if blank
	if len(host) == 0 {
		host = req.Host
	}

	// split out ip
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	// determine the namespace of the request
	namespace := r.Namespace(req)

	// overide host if the namespace is go.micro.web, since
	// this will also catch localhost & 127.0.0.1, resulting
	// in a more consistent dev experience
	if host == "localhost" || host == "127.0.0.1" {
		host = "dev.m3o.org"
	}

	// if the type is path, always resolve using the path
	if r.Type == "path" {
		return host, namespace, true
	}

	// if the namespace is not the default (go.micro.web),
	// we always resolve using path
	if namespace != defaultNamespace {
		return host, namespace, true
	}

	// check for micro subdomains, we want to do subdomain routing
	// on these if the subdomoain routing has been specified
	if r.Type == "subdomain" && host != "dev.m3o.org" && strings.HasSuffix(host, ".m3o.org") {
		return host, namespace, false
	}

	// Check for services info path, also handled by micro web but
	// not a top level path. TODO: Find a better way of detecting and
	// handling the non-proxied paths.
	if strings.HasPrefix(req.URL.Path, "/service/") {
		return host, namespace, true
	}

	// Check if the request is a top level path
	isWeb := strings.Count(req.URL.Path, "/") == 1
	return host, namespace, isWeb
}
