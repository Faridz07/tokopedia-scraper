package model

import "time"

type ProductDAO struct {
	ID          int
	URL         string
	Name        string
	Description string
	ImageLink   string
	Price       float64
	Rating      float64
	TotalRating int
	StoreName   string
	ScrapedAt   time.Time
	Updated     time.Time
}
