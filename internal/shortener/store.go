package shortener

type URLStore interface {
	Create(url ShortenModel) (ShortenModel, error)
	GetByShortID(shortID string) (ShortenModel, error)
	GetAllByUser(userID uint) ([]ShortenModel, error)
	GetAll() ([]ShortenModel, error)        
	Delete(shortID string) error            
}
