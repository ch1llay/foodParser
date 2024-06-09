package models

// Product структура для хранения информации о продукте
type Product struct {
	ImageUrl     string `json:"image_url"`
	Name         string `json:"name"`
	Protein      string `json:"protein"`
	Fat          string `json:"fat"`
	Carbohydrate string `json:"carbohydrate"`
	Calories     string `json:"calories"`
}
