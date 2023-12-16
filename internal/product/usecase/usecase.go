package productusecase

import (
	"tokopedia-scraper/internal/product/model"
	productrepository "tokopedia-scraper/internal/product/repository"
	"tokopedia-scraper/pkg/session"
)

type ProductUsecase interface {
	CreateProduct(
		sess *session.Session,
		req *model.CreateProductRequest,
	) (*model.CreateProductResponse, error)
}

type productUsecase struct {
	productRepository productrepository.ProductRepository
}

func NewProductUsecase(productRepository productrepository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepository: productRepository}
}
