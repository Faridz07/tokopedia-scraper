package model

type (
	CreateProductRequest struct {
		URL         string `json:"url"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ImageURL    string `json:"imageUrl"`
		Price       string `json:"price"`
		Rating      string `json:"rating"`
		TotalRating string `json:"totalRating"`
		StoreName   string `json:"storeName"`
	}

	CreateProductResponse struct {
		Created bool   `json:"created"`
		Name    string `json:"productName"`
	}
)
