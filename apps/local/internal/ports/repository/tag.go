package repository

import (
	"context"
	"errors"

	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
)

var (
	ErrTagNotFound = errors.New("tag does not exist")
)

type Tag interface {
	ListForResource(ctx context.Context, rid string) ([]*tag.Tag, error)
	AddToResource(ctx context.Context, t *tag.Tag, rid string, ty tag.ResourceType) error
	DeleteTagFromResource(ctx context.Context, key, rid string) error
}
