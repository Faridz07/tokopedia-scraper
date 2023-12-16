package productrepository

import (
	"tokopedia-scraper/infrastructure/database"
	"tokopedia-scraper/internal/product/model"
	"tokopedia-scraper/pkg/session"
)

type ProductRepository interface {
	Create(
		sess *session.Session,
		product model.ProductDAO,
	) (bool, error)
}

type productRepository struct {
	db database.Database
}

func NewProductRepository(db database.Database) ProductRepository {
	return &productRepository{db: db}
}
