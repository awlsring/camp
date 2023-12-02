package tag_repository

import "github.com/awlsring/camp/apps/local/internal/core/domain/tag"

type TagSql struct {
	Identifier         int    `db:"id"`
	ResourceIdentifier string `db:"resource_identifier"`
	ResourceType       string `db:"resource_type"`
	Key                string `db:"tag_key"`
	Value              string `db:"tag_value"`
}

func (t *TagSql) ToModel() (*tag.Tag, error) {
	key, err := tag.TagKeyFromString(t.Key)
	if err != nil {
		return nil, err
	}

	value, err := tag.TagValueFromString(t.Value)
	if err != nil {
		return nil, err
	}

	return &tag.Tag{
		Key:   key,
		Value: value,
	}, nil
}
