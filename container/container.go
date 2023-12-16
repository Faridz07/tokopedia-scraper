package container

import (
	"strings"
	"tokopedia-scraper/config"
	"tokopedia-scraper/infrastructure/database"
	"tokopedia-scraper/infrastructure/log"
	productrepository "tokopedia-scraper/internal/product/repository"
	productusecase "tokopedia-scraper/internal/product/usecase"
	scraperusecase "tokopedia-scraper/internal/scraper/usecase"
)

type Container struct {
	Config         *config.Config
	Log            log.Log
	ScraperUsecase scraperusecase.ScraperUsecase
}

func New() *Container {
	// Initiate config
	cfg := config.New()

	// Initiate database
	db := database.New(&cfg.Infrastructure.Database)

	if strings.ToLower(cfg.App.Env) == "dev" {
		database.MigrateDatabase(db, cfg.Infrastructure.Database.MigrationPath)
	}

	// Initiate log
	log := log.New(cfg.Infrastructure.Log)

	// Initiate repository
	productRepository := productrepository.NewProductRepository(db)

	// Initiate usecase
	productUsecase := productusecase.NewProductUsecase(productRepository)
	scraperUsecase := scraperusecase.NewScraperUsecase(cfg.Scrape, productUsecase)

	return &Container{
		Config:         cfg,
		Log:            log,
		ScraperUsecase: scraperUsecase,
	}
}
