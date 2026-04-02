package shortenedurls

import "github.com/jmoiron/sqlx"

const (
	alias     = "su"
	tableName = "shortened_urls"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}
