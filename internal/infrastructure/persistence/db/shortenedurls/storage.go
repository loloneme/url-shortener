package shortenedurls

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/domain/shortgen"
	persistence "url-shortener/internal/infrastructure/persistence/db"

	sq "github.com/Masterminds/squirrel"
)

var st = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (r *Repository) Save(ctx context.Context, original string) (*Model, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query, args, err := st.Insert(tableName).
		Columns("original").
		Values(original).
		Suffix("ON CONFLICT (original) DO NOTHING RETURNING id").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var id uint64
	err = tx.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return r.GetByOriginal(ctx, original)
		}
		return nil, fmt.Errorf("failed to insert shortened url: %w", err)
	}

	short, err := shortgen.EncodeID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to encode id: %w", err)
	}

	if _, err = tx.ExecContext(ctx, "UPDATE shortened_urls SET short = $1 WHERE id = $2", short, id); err != nil {
		return nil, fmt.Errorf("failed to update short: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &Model{ID: id, Short: short, Original: original}, nil
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
