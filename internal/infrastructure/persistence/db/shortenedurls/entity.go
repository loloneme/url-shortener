package shortenedurls

import domain "url-shortener/internal/domain/shortenedurl"

type Model = domain.ShortenedURL

type ShortenedURL struct {
	ID       uint64 `db:"id"`
	Short    string `db:"short"`
	Original string `db:"original"`
}

func (s ShortenedURL) Values() []any {
	return []any{s.ID, s.Short, s.Original}
}

func (s ShortenedURL) ToModel() Model {
	return Model{
		ID:       s.ID,
		Short:    s.Short,
		Original: s.Original,
	}
}

func (s ShortenedURL) FromModel(model *Model) ShortenedURL {
	return ShortenedURL{
		ID:       model.ID,
		Short:    model.Short,
		Original: model.Original,
	}
}
