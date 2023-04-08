package queue

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/micromdm/micromdm/pkg/httputil"
)

type ListDeviceCommandOption struct {
	UDID string `json:"udid"`
}

type CommandDTO struct {
	UUID           string    `json:"uuid"`
	Payload        []byte    `json:"payload"`
	CreatedAt      time.Time `json:"created_at"`
	LastSentAt     time.Time `json:"last_sent_at"`
	Acknowledged   time.Time `json:"acknowledged"`
	TimesSent      int       `json:"times_sent"`
	LastStatus     string    `json:"last_status"`
	FailureMessage []byte    `json:"failure_message"`
}

type DeviceCommandDTO struct {
	DeviceUDID string       `json:"device_udid"`
	Commands   []CommandDTO `json:"commands"`
	Completed  []CommandDTO `json:"completed"`
	Failed     []CommandDTO `json:"failed"`
	NotNow     []CommandDTO `json:"not_now"`
}

type getDeviceCommandRequest struct{ Opts ListDeviceCommandOption }
type getDeviceCommandResponse struct {
	DeviceCommand DeviceCommandDTO `json:"device_command"`
	Err           error            `json:"err,omitempty"`
}

func (r getDeviceCommandResponse) Failed() error { return r.Err }

func decodeListDeviceCommandRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var opts ListDeviceCommandOption
	err := httputil.DecodeJSONRequest(r, &opts)
	req := getDeviceCommandRequest{
		Opts: opts,
	}
	return req, err
}

func decodeListDeviceCommandResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var resp getDeviceCommandResponse
	err := httputil.DecodeJSONResponse(r, &resp)
	return resp, err
}

func MakeListDeviceCommandEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getDeviceCommandRequest)
		dto, err := svc.ListDeviceCommands(ctx, req.Opts)
		return getDeviceCommandResponse{
			DeviceCommand: dto,
			Err:           err,
		}, nil
	}
}

func (e Endpoints) ListDeviceCommand(ctx context.Context, opts ListDeviceCommandOption) (DeviceCommandDTO, error) {
	request := getDeviceCommandRequest{opts}
	response, err := e.ListDeviceCommandEndpoint(ctx, request.Opts)
	if err != nil {
		return DeviceCommandDTO{}, err
	}
	return response.(getDeviceCommandResponse).DeviceCommand, response.(getDeviceCommandResponse).Err
}
