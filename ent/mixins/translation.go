package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type TranslationMixin struct {
	mixin.Schema
	ParentField string
}

func (TranslationMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("locale").
			NotEmpty().
			Comment("IETF language tag, e.g. en, de, fr-FR"),
	}
}

func (m TranslationMixin) Indexes() []ent.Index {
	indexes := []ent.Index{
		index.Fields("locale"),
	}

	if m.ParentField != "" {
		indexes = append(indexes,
			index.Fields("locale", m.ParentField).Unique(),
		)
	}

	return indexes
}
