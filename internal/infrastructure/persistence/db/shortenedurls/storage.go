package shortenedurls

import (
	"context"
	"fmt"
	persistence "url-shortener/internal/infrastructure/persistence/db"

	sq "github.com/Masterminds/squirrel"
)

var st = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type FindSpecification interface {
	GetRule(builder sq.SelectBuilder) sq.SelectBuilder
}

func (r *Repository) Save(ctx context.Context, model *Model) error {
	entity := ShortenedURL{}.FromModel(model)

	query, args, err := st.Insert(tableName).
		Columns(writableColumns...).
		Values(entity.Values()...).
		Suffix("ON CONFLICT (original) DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if mappedErr := persistence.MapError(err); mappedErr != nil {
		return fmt.Errorf("failed to save shortened url: %w", mappedErr)
	}
	return nil
}

func (r *Repository) findBy(ctx context.Context, column string, value any) (*Model, error) {
	var entity ShortenedURL

	query, args, err := st.Select(readableColumns...).
		From(tableName).
		Where(sq.Eq{column: value}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.GetContext(ctx, &entity, query, args...)
	if mappedErr := persistence.MapError(err); mappedErr != nil {
		return nil, fmt.Errorf("failed to find shortened url: %w", mappedErr)
	}

	model := entity.ToModel()
	return &model, nil
}

func (r *Repository) GetByOriginal(ctx context.Context, original string) (*Model, error) {
	return r.findBy(ctx, "original", original)
}

func (r *Repository) GetByShort(ctx context.Context, short string) (*Model, error) {
	return r.findBy(ctx, "short", short)
}
