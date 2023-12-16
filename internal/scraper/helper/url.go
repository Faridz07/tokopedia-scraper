package helper

import "fmt"

func GetUrl(baseUrl, path string, index int) string {
	path = fmt.Sprintf(path, index)
	return fmt.Sprintf("%s%s", baseUrl, path)
}
