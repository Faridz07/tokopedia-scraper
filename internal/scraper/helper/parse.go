package helper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tokopedia-scraper/pkg/constant"
	"tokopedia-scraper/pkg/session"
)

func ParsePrice(
	sess *session.Session,
	priceStr string,
) float64 {
	re := regexp.MustCompile(`[^\d,.-]+`)
	priceStr = re.ReplaceAllString(priceStr, "")
	priceStr = strings.ReplaceAll(priceStr, ".", "")

	priceStr = strings.Replace(priceStr, ",", ".", -1)

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		sess.SetError(constant.ErrInternal,
			fmt.Errorf("Error parsing price: %w", err))
		return 0.0
	}
	return price
}

func ParseRating(
	sess *session.Session,
	ratingStr string,
) float64 {
	ratingStr = strings.TrimLeft(ratingStr, "(")
	ratingStr = strings.TrimRight(ratingStr, ")")
	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		sess.SetError(constant.ErrInternal,
			fmt.Errorf("Error parsing rating: %w", err))
		return 0.0
	}
	return rating
}
