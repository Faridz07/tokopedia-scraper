package scraperusecase

import (
	"tokopedia-scraper/config"
	productusecase "tokopedia-scraper/internal/product/usecase"
	"tokopedia-scraper/pkg/session"
)

type ScraperUsecase interface {
	Scrape(
		sess *session.Session,
	) (err error)
}

type scraperUsecase struct {
	cfg            config.ScrapeConfig
	productUsecase productusecase.ProductUsecase
}

func NewScraperUsecase(
	cfg config.ScrapeConfig,
	productUsecase productusecase.ProductUsecase,
) ScraperUsecase {
	return &scraperUsecase{
		cfg:            cfg,
		productUsecase: productUsecase,
	}
}
