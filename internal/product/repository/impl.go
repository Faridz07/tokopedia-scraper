package productrepository

import (
	"tokopedia-scraper/internal/product/model"
	"tokopedia-scraper/pkg/constant"
	"tokopedia-scraper/pkg/session"
)

func (r productRepository) Create(
	sess *session.Session,
	product model.ProductDAO,
) (bool, error) {
	query := `
	INSERT INTO products (url, name, description, image_link, price, rating, total_rating, store_name, scraped_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP)
	ON CONFLICT (url) 
	DO UPDATE SET
		name = EXCLUDED.name,
		description = EXCLUDED.description,
		image_link = EXCLUDED.image_link,
		price = EXCLUDED.price,
		rating = EXCLUDED.rating,
		total_rating = EXCLUDED.total_rating,
		store_name = EXCLUDED.store_name,
		updated_at = CURRENT_TIMESTAMP;
	`

	rows, err := r.db.ExecContext(sess.Ctx,
		query,
		product.URL,
		product.Name,
		product.Description,
		product.ImageLink,
		product.Price,
		product.Rating,
		product.TotalRating,
		product.StoreName,
		product.ScrapedAt)
	if err != nil {
		sess.SetError(constant.ErrDatabase, err)
		return false, err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		sess.SetError(constant.ErrDatabase, err)
		return false, err
	}

	return affected != 0, nil
}
