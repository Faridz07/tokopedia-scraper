package util

import (
	"time"
)

func GetJakartaTime() time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now()
	}

	return time.Now().In(loc)
}
