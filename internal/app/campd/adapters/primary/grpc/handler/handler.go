package handler

import (
	"github.com/awlsring/camp/internal/app/campd/ports/service"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
)

type Handler struct {
	cpuSvc  service.CPU
	hostSvc service.Host
	memSvc  service.Memory
	moboSvc service.Motherboard
	netSvc  service.Network
	strSvc  service.Storage
	campd.UnimplementedCampdServer
}

func New(cpuSvc service.CPU, hostSvc service.Host, memSvc service.Memory, moboSvc service.Motherboard, netSvc service.Network, strSvc service.Storage) campd.CampdServer {
	return &Handler{
		cpuSvc:  cpuSvc,
		hostSvc: hostSvc,
		memSvc:  memSvc,
		moboSvc: moboSvc,
		netSvc:  netSvc,
		strSvc:  strSvc,
	}
}
