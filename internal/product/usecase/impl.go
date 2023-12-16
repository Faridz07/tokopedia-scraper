package productusecase

import (
	"fmt"
	"tokopedia-scraper/internal/product/helper"
	"tokopedia-scraper/internal/product/model"
	"tokopedia-scraper/pkg/constant"
	"tokopedia-scraper/pkg/session"
	"tokopedia-scraper/pkg/util"
)

func (u productUsecase) CreateProduct(
	sess *session.Session,
	req *model.CreateProductRequest,
) (*model.CreateProductResponse, error) {
	price, err := helper.ParsePrice(req.Price)
	if err != nil {
		err = fmt.Errorf("Error converting price to float: %w", err)
		sess.SetError(constant.ErrInternal, err)
		return nil, err
	}

	rating, err := helper.ParseRating(req.Rating)
	if err != nil {
		err = fmt.Errorf("Error converting rating to float: %w", err)
		sess.SetError(constant.ErrInternal, err)
		return nil, err
	}

	totalRating, err := helper.ParseTotalRating(req.TotalRating)
	if err != nil {
		err = fmt.Errorf("Error converting total rating to int: %w", err)
		sess.SetError(constant.ErrInternal, err)
		return nil, err
	}

	product := model.ProductDAO{
		URL:         req.URL,
		Name:        req.Name,
		Description: req.Description,
		ImageLink:   req.ImageURL,
		Price:       price,
		Rating:      rating,
		TotalRating: totalRating,
		StoreName:   req.StoreName,
		ScrapedAt:   util.GetJakartaTime(),
	}

	created, err := u.productRepository.Create(sess, product)
	if err != nil {
		return nil, err
	}

	return &model.CreateProductResponse{
		Created: created,
		Name:    req.Name,
	}, nil
}
