package handler

import (
	"github.com/awlsring/camp/internal/app/local/ports/service"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
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
