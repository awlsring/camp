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
	return &tag.Tag{
		Key:   t.Key,
		Value: t.Value,
	}, nil
}
