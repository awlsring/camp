package handler

import (
	"github.com/awlsring/camp/apps/local/internal/core/ports/service"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

var _ camplocal.Handler = &Handler{}

type Handler struct {
	mSvc service.Machine
}

func NewHandler(mSvc service.Machine) *Handler {
	return &Handler{
		mSvc: mSvc,
	}
}
