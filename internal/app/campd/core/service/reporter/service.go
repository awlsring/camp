package reporter

import (
	"github.com/awlsring/camp/internal/app/campd/ports/gateway"
	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/machine"
)

type Service struct {
	id        machine.Identifier
	reporting gateway.Reporting
	hostSvc   service.Host
	boardSvc  service.Motherboard
	cpuSvc    service.CPU
	memSvc    service.Memory
	netSvc    service.Network
	storage   service.Storage
}
