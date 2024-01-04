package system

import (
	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
)

type Service struct {
	id       machine.Identifier
	hostSvc  service.Host
	boardSvc service.Motherboard
	cpuSvc   service.CPU
	memSvc   service.Memory
	netSvc   service.Network
	storage  service.Storage
}

func New(id machine.Identifier, hostSvc service.Host, boardSvc service.Motherboard, cpuSvc service.CPU, memSvc service.Memory, netSvc service.Network, storage service.Storage) *Service {
	return &Service{
		id:       id,
		hostSvc:  hostSvc,
		boardSvc: boardSvc,
		cpuSvc:   cpuSvc,
		memSvc:   memSvc,
		netSvc:   netSvc,
		storage:  storage,
	}
}
