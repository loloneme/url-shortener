package persistence

import (
	"database/sql"
	"errors"
	domain "url-shortener/internal/domain/shortenedurl"

	"github.com/jackc/pgx/v5/pgconn"
)

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func isNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func MapError(err error) error {
	if isUniqueViolation(err) {
		return domain.ErrDuplicate
	}
	if isNotFound(err) {
		return domain.ErrNotFound
	}
	return err
}
