package handler

import (
	"net/http"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/router"
)

type metaHandler struct {
	c  client.Client
	r  router.Router
	ns string
}

func (m *metaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logger.Info("i'm a meta handler")

	// service, err := m.r.Route(r)
	// if err != nil {
	// 	er := errors.InternalServerError(m.ns, err.Error())
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(500)
	// 	w.Write([]byte(er.Error()))
	// 	return
	// }

}

// Meta is a http.Handler that routes based on endpoint metadata
func Meta(s *service.Service, r router.Router, ns string) http.Handler {
	return &metaHandler{
		c:  s.Client(),
		r:  r,
		ns: ns,
	}
}
