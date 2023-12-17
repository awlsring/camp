package handler

import "github.com/awlsring/camp/internal/app/campd/ports/service"

type SystemHandler struct {
	cpuSvc  service.CPU
	hostSvc service.Host
	memSvc  service.Memory
	moboSvc service.Motherboard
	netSvc  service.Network
	strSvc  service.Storage
}

func New(cpuSvc service.CPU, hostSvc service.Host, memSvc service.Memory, moboSvc service.Motherboard, netSvc service.Network, strSvc service.Storage) *SystemHandler {
	return &SystemHandler{
		cpuSvc:  cpuSvc,
		hostSvc: hostSvc,
		memSvc:  memSvc,
		moboSvc: moboSvc,
		netSvc:  netSvc,
		strSvc:  strSvc,
	}
}
