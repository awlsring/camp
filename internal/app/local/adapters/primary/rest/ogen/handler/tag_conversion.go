package handler

import (
	"github.com/awlsring/camp/internal/pkg/domain/tag"
	"github.com/awlsring/camp/internal/pkg/errors"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
)

func tagToDomain(t camplocal.Tag) (*tag.Tag, error) {
	key, err := tag.TagKeyFromString(t.Key)
	if err != nil {
		return nil, errors.New(errors.ErrValidation, err)
	}

	value, err := tag.TagValueFromString(t.Value)
	if err != nil {
		return nil, errors.New(errors.ErrValidation, err)
	}

	return &tag.Tag{
		Key:   key,
		Value: value,
	}, nil
}

func tagsToDomain(tagList []camplocal.Tag) ([]*tag.Tag, error) {
	tags := make([]*tag.Tag, len(tagList))
	for i, t := range tagList {
		tag, err := tagToDomain(t)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}
	return tags, nil
}

func domainToTags(tags []*tag.Tag) []camplocal.Tag {
	tagList := make([]camplocal.Tag, len(tags))
	for i, t := range tags {
		tagList[i] = camplocal.Tag{
			Key:   t.Key.String(),
			Value: t.Value.String(),
		}
	}
	return tagList
}
