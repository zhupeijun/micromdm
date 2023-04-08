package queue

import (
	"context"
	"github.com/pkg/errors"
)

type Service interface {
	ListDeviceCommands(ctx context.Context, opt ListDeviceCommandOption) (DeviceCommandDTO, error)
}

type DeviceCommandService struct {
	store *Store
}

func (svc *DeviceCommandService) ListDeviceCommands(_ context.Context, opt ListDeviceCommandOption,
) (DeviceCommandDTO, error) {
	deviceCommand, err := svc.store.DeviceCommand(opt.UDID)
	if err != nil {
		if isNotFound(err) {
			return DeviceCommandDTO{}, nil
		}
		return DeviceCommandDTO{}, errors.Wrapf(err, "get device command from queue, udid: %s", opt.UDID)
	}

	dto := DeviceCommandDTO{
		DeviceUDID: deviceCommand.DeviceUDID,
		Commands:   svc.convertCommandsToDTOs(deviceCommand.Commands),
		Completed:  svc.convertCommandsToDTOs(deviceCommand.Completed),
		Failed:     svc.convertCommandsToDTOs(deviceCommand.Failed),
		NotNow:     svc.convertCommandsToDTOs(deviceCommand.NotNow),
	}
	return dto, err
}

func (svc *DeviceCommandService) convertCommandToDTO(command Command) CommandDTO {
	return CommandDTO{
		UUID:           command.UUID,
		Payload:        command.Payload,
		CreatedAt:      command.CreatedAt,
		LastSentAt:     command.LastSentAt,
		Acknowledged:   command.Acknowledged,
		TimesSent:      command.TimesSent,
		LastStatus:     command.LastStatus,
		FailureMessage: command.FailureMessage,
	}
}

func (svc *DeviceCommandService) convertCommandsToDTOs(commands []Command) []CommandDTO {
	var dto []CommandDTO
	for _, command := range commands {
		dto = append(dto, svc.convertCommandToDTO(command))
	}
	return dto
}

func New(store *Store) *DeviceCommandService {
	return &DeviceCommandService{store: store}
}
