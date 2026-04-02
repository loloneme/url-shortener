package shortenedurls

import domain "url-shortener/internal/domain/shortenedurl"

type Model = domain.ShortenedURL

type ShortenedURL struct {
	Short    string `db:"short"`
	Original string `db:"original"`
}

func (s ShortenedURL) Values() []any {
	return []any{s.Short, s.Original}
}

func (s ShortenedURL) ToModel() Model {
	return Model{
		Short:    s.Short,
		Original: s.Original,
	}
}

func (s ShortenedURL) FromModel(model *Model) ShortenedURL {
	return ShortenedURL{
		Short:    model.Short,
		Original: model.Original,
	}
}
