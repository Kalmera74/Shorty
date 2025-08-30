package shortener

type ShortStore interface {
	Create(short ShortModel) (ShortModel, error)

	GetById(id uint) (ShortModel, error)

	GetByShortUrl(shortUrl string) (ShortModel, error)
	GetByLongUrl(originalUrl string) (ShortModel, error)

	GetAllByUser(userID uint) ([]ShortModel, error)

	GetAll() ([]ShortModel, error)

	Delete(shortenID uint) error
}
