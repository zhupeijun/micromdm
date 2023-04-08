package queue

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type Endpoints struct {
	ListDeviceCommandEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service, outer endpoint.Middleware, others ...endpoint.Middleware) Endpoints {
	return Endpoints{
		ListDeviceCommandEndpoint: endpoint.Chain(outer, others...)(MakeListDeviceCommandEndpoint(s)),
	}
}

func RegisterHTTPHandlers(r *mux.Router, e Endpoints, options ...httptransport.ServerOption) {
	// POST     /v1/device_command		get a list of devices managed by the server

	r.Methods("POST").Path("/v1/device_command").Handler(httptransport.NewServer(
		e.ListDeviceCommandEndpoint,
		decodeListDeviceCommandRequest,
		httputil.EncodeJSONResponse,
		options...,
	))
}
