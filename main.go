package main

import (
	"log"
	"tokopedia-scraper/container"
	"tokopedia-scraper/pkg/session"
)

func main() {
	cont := container.New()

	session := session.NewSession(cont.Log)
	err := cont.ScraperUsecase.Scrape(session)
	if err != nil {
		log.Printf("Scrape failed, err: %v", err)
		return
	}
}
