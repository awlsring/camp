package repository

import (
	"context"
	"errors"

	"github.com/awlsring/camp/internal/pkg/domain/tag"
)

var (
	ErrTagNotFound          = errors.New("tag does not exist")
	ErrDuplicateResourceTag = errors.New("duplicate resource tag")
	ErrUnknownFailure       = errors.New("unknown error")
)

type Tag interface {
	ListForResource(ctx context.Context, rid string) ([]*tag.Tag, error)
	AddToResource(ctx context.Context, t *tag.Tag, rid string, ty tag.ResourceType) error
	DeleteTagFromResource(ctx context.Context, key tag.TagKey, rid string) error
}
