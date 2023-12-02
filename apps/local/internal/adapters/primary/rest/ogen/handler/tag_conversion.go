package handler

import (
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/internal/pkg/errors"
	camplocal "github.com/awlsring/camp/packages/camp_local"
)

func tagsToDomain(tagList []camplocal.Tag) ([]*tag.Tag, error) {
	tags := make([]*tag.Tag, len(tagList))
	for i, t := range tagList {

		key, err := tag.TagKeyFromString(t.Key)
		if err != nil {
			return nil, errors.New(errors.ErrValidation, err)
		}

		value, err := tag.TagValueFromString(t.Value)
		if err != nil {
			return nil, errors.New(errors.ErrValidation, err)
		}

		tags[i] = &tag.Tag{
			Key:   key,
			Value: value,
		}
	}
	return tags, nil
}
