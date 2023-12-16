package scraperusecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"tokopedia-scraper/internal/product/model"
	"tokopedia-scraper/internal/scraper/helper"
	"tokopedia-scraper/pkg/constant"
	"tokopedia-scraper/pkg/session"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func (u scraperUsecase) Scrape(sess *session.Session) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	sess.SetContext(ctx)

	var products []*model.CreateProductRequest
	uniqueUrls := make(map[string]bool)

	i := 1
	for len(products) < 100 {
		url := helper.GetUrl(u.cfg.Web.BaseUrl, u.cfg.Web.PathToScrap, i)
		log.Printf("Starting Scrape URL: %s\n", url)

		tempProducts, err := u.ScrapeMainPage(sess, url)
		if err != nil {
			return err
		}

		for _, product := range tempProducts {
			if _, exists := uniqueUrls[product.URL]; !exists {
				uniqueUrls[product.URL] = true
				products = append(products, product)

				log.Printf("Save data product: %s into database\n", product.Name)
				_, err = u.productUsecase.CreateProduct(sess, product)
				if err != nil {
					return err
				}
			}
		}

		if len(products) >= 100 {
			products = products[:100]
			break
		}

		log.Printf("Finished Scrape URL: %s, we got %d data product\n", url, len(products))
		i++
	}

	if err := u.writeProductsToCSV(products); err != nil {
		err = fmt.Errorf("Error writing products to CSV: %w\n", err)
		sess.SetError(constant.ErrInternal, err)
		return err
	}

	log.Println("Products successfully written to CSV file")
	return nil
}

func (u scraperUsecase) ScrapeMainPage(
	sess *session.Session,
	url string,
) ([]*model.CreateProductRequest, error) {
	if err := chromedp.Run(sess.Ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers{
			"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"Sec-Ch-Ua":          `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": `"macOS"`,
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Site":     "same-site",
		}),
		chromedp.Navigate(url),
	); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		if err := chromedp.Run(sess.Ctx,
			chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
			chromedp.Sleep(2*time.Second),
		); err != nil {
			log.Fatal(err)
		}
	}

	var res []byte
	err := chromedp.Run(sess.Ctx,
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('div[data-testid="divProductWrapper"]')).map(div => {
				const anchorElement = div.closest('a.css-54k5sq');
				const url = anchorElement?.getAttribute('href') || '';
				if (url == '' || url.includes('https://ta.tokopedia.com/promo/')) {
					return null;
				}
				const name = div.querySelector('span.css-20kt3o')?.textContent.trim() || '';
				const imageElement = div.querySelector('div.css-1g5og91 > img');
				const imageURL = imageElement?.src || '';
				const price = div.querySelector('div.css-pp6b3e > span.css-o5uqvq')?.textContent.trim() || '';
				const ratingElement = div.querySelector('div.css-1riykrk > div > span');
				const rating = ratingElement ? ratingElement.textContent.match(/\d+/g).join('.') : '';
				const storeElements = div.querySelectorAll('div.css-tpww51 > div.css-vbihp9 > span.css-ywdpwd');
				const storeName = storeElements && storeElements.length > 1 ? storeElements[1].textContent.trim() : '';

				return { url, name, description: '', imageURL, price, rating, storeName };
			}).filter(product => product !== null)`, &res),
	)
	if err != nil {
		sess.SetError(constant.ErrInternal, err)
		return nil, err
	}

	var products []*model.CreateProductRequest
	err = json.Unmarshal(res, &products)
	if err != nil {
		sess.SetError(constant.ErrInternal, err)
		return nil, err
	}

	jobs := make(chan *model.CreateProductRequest, len(products))
	results := make(chan *model.CreateProductRequest, len(products))
	errors := make(chan error)
	var counter uint64

	for w := 0; w < u.cfg.Worker; w++ {
		go u.worker(w, sess, jobs, results, errors, &counter)
	}

	for _, product := range products {
		jobs <- product
	}
	close(jobs)

	for i := 0; i < len(products); i++ {
		select {
		case product := <-results:
			products[i] = product
		case err := <-errors:
			sess.SetError(constant.ErrInternal, err)
		}
	}

	return products, nil
}

func (u scraperUsecase) ScrapeDetails(
	sess *session.Session,
	req *model.CreateProductRequest,
	url string,
) error {
	log.Printf("Starting ScrapeDetails for product: %s\n", req.Name)

	ctx, cancel := chromedp.NewContext(sess.Ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers{
			"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"Sec-Ch-Ua":          `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": `"macOS"`,
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Site":     "same-site",
		}),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Evaluate(`document.querySelector('[data-testid="lblPDPDetailProductRatingNumber"]')?.textContent.trim() || ''`, &req.Rating),
		chromedp.Evaluate(`document.querySelector('[data-testid="lblPDPDetailProductRatingCounter"]')?.textContent.trim() || ''`, &req.TotalRating),
		chromedp.Evaluate(`document.querySelector('img[data-testid="PDPMainImage"]')?.getAttribute('src') || ''`, &req.ImageURL),
		chromedp.Evaluate(`document.querySelector('[data-testid="lblPDPDescriptionProduk"]')?.innerText || ''`, &req.Description),
	)

	if err != nil {
		log.Printf("Error in ScrapeDetails for product %s: %v\n", req.Name, err)
		sess.SetError(constant.ErrInternal, err)
		return err
	}

	log.Printf("Completed ScrapeDetails for product: %s\n", req.Name)
	return nil
}
