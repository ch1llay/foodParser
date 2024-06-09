package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Product структура для хранения информации о продукте
type Product struct {
	ImageUrl     string `json:"image_url"`
	Name         string `json:"name"`
	Protein      string `json:"protein"`
	Fat          string `json:"fat"`
	Carbohydrate string `json:"carbohydrate"`
	Calories     string `json:"calories"`
}

func cleanString(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
}

var productsAll = make([][]Product, 84)

func getFromPage(page int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("page %d start \n", page)

	url := fmt.Sprintf("https://www.calorizator.ru/product/all?page=%d", page)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)

	// Парсинг HTML с помощью goquery
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	var products []Product
	table := doc.Find("table")
	tbody := table.Find("tbody")
	rows := tbody.Find("tr")

	// Поиск таблицы с продуктами
	rows.Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		if cells.Length() == 6 {
			product := Product{
				ImageUrl:     cleanString(cells.Eq(0).Find("a").AttrOr("href", "")),
				Name:         cleanString(cells.Eq(1).Text()),
				Protein:      cleanString(cells.Eq(2).Text()),
				Fat:          cleanString(cells.Eq(3).Text()),
				Carbohydrate: cleanString(cells.Eq(4).Text()),
				Calories:     cleanString(cells.Eq(5).Text()),
			}
			products = append(products, product)
		}
	})

	fmt.Printf("page %d complete\n", page)
	productsAll[page] = products
}

func main() {
	// URL страницы с продуктами

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		if i%5 == 0 {
			time.Sleep(time.Second * 2)
		}
		wg.Add(1)
		go getFromPage(i, &wg)
	}
	wg.Wait()

	productsAll_ := make([]Product, 0, 0)
	for _, product := range productsAll {
		productsAll_ = append(productsAll_, product...)
	}

	fmt.Printf("len %d", len(productsAll_))
	// Сохранение данных в JSON файл
	file, err := os.Create("calorizator_products.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(productsAll_)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Данные успешно спарсены и сохранены в файл calorizator_products.json")
}
