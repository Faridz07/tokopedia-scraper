package scraperusecase

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
	"tokopedia-scraper/internal/product/model"
	"tokopedia-scraper/pkg/session"
)

func (u scraperUsecase) worker(
	id int,
	sess *session.Session,
	jobs <-chan *model.CreateProductRequest,
	results chan<- *model.CreateProductRequest,
	errors chan<- error,
	counter *uint64,
) {
	for job := range jobs {
		log.Printf("Worker %d started job for product: %s\n", id, job.Name)

		err := u.ScrapeDetails(sess, job, job.URL)

		if err != nil {
			log.Printf("Worker %d encountered error for product %s: %v\n", id, job.Name, err)
			errors <- err
			continue
		}

		atomic.AddUint64(counter, 1)
		log.Printf("Worker %d: Completed product %s. Total completed: %d\n", id, job.Name, atomic.LoadUint64(counter))
		results <- job
	}
	log.Printf("Worker %d has finished processing jobs\n", id)
}

func (u scraperUsecase) writeProductsToCSV(products []*model.CreateProductRequest) error {
	if err := os.MkdirAll(u.cfg.Output.Path, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	currentDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf(u.cfg.Output.Filename, u.cfg.Output.Path, currentDate)

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Name of Product", "Description", "Image Link", "Price", "Rating", "Name of Store"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header to CSV file: %w", err)
	}

	for _, product := range products {
		record := []string{
			product.Name,
			product.Description,
			product.ImageURL,
			product.Price,
			product.Rating,
			product.StoreName,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record to CSV file: %w", err)
		}
	}

	return nil
}
