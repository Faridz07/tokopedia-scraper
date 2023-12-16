package helper

import (
	"strconv"
	"strings"
)

func ParsePrice(str string) (float64, error) {
	if str == "" {
		return 0, nil
	}

	str = strings.ReplaceAll(str, "Rp", "")
	str = strings.ReplaceAll(str, ".", "")

	price, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func ParseRating(str string) (float64, error) {
	if str == "" {
		return 0, nil
	}

	price, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func ParseTotalRating(str string) (int, error) {
	if str == "" {
		return 0, nil
	}

	str = strings.ReplaceAll(str, "(", "")
	str = strings.ReplaceAll(str, ")", "")
	str = strings.ReplaceAll(strings.ToLower(str), "rating", "")
	str = strings.ReplaceAll(str, ".", "")
	str = strings.TrimSpace(str)

	result, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return result, nil
}
