package meta

import (
	"net/http"

	"github.com/micro-community/micro-webui/handler"
	"github.com/micro-community/micro-webui/handler/api"
	"github.com/micro-community/micro-webui/handler/web"
	"github.com/micro-community/micro-webui/router"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

type metaHandler struct {
	c  client.Client
	r  router.Router
	ns string
}

func (m *metaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logger.Info("i'm a meta handler")

	service, err := m.r.Route(r)
	if err != nil {
		er := errors.InternalServerError(m.ns, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	switch service.Endpoint.Handler {
	// web socket handler
	case web.Handler:
		web.WithService(service, handler.WithClient(m.c)).ServeHTTP(w, r)
	// api handler
	case api.Handler:
		api.WithService(service, handler.WithClient(m.c)).ServeHTTP(w, r)
	default:
		web.WithService(service, handler.WithClient(m.c)).ServeHTTP(w, r)
	}

}

// NewMetaHandler is a http.Handler that routes based on endpoint metadata
func NewMetaHandler(cli client.Client, r router.Router, ns string) http.Handler {
	return &metaHandler{
		c:  cli,
		r:  r,
		ns: ns,
	}
}
