package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"foodParser/models"
	"github.com/google/uuid"
	"os"
	"strconv"
)

func getData() []models.Product {
	f, err := os.Open("calorizator_products.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Чтение файла с ридером
	wr := bytes.Buffer{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}

	if err := sc.Err(); err != nil {
		panic(err)
	}

	var products []models.Product
	if err := json.Unmarshal(wr.Bytes(), &products); err != nil {
		panic(err)
	}

	return products
}

func createSql(products []models.Product) {
	f, err := os.OpenFile("data.sql", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_, err = f.WriteString("insert into products(id, name,protein,fat, carbohydrate, calories) values(")
	if err != nil {
		return
	}

	for _, product := range products {
		proteins, err := strconv.ParseFloat(product.Protein, 32)
		fat, err := strconv.ParseFloat(product.Fat, 32)
		carbohydrate, err := strconv.ParseFloat(product.Carbohydrate, 32)
		calories, err := strconv.ParseInt(product.Calories, 10, 32)

		if err != nil {
			fmt.Println("Ошибка:", err)
		}
		_, err = f.WriteString(fmt.Sprintf(

			"('%s', '%s', %f, %f, %f, %d),\n",
			uuid.New(),
			cleanString(product.Name),
			proteins,
			fat,
			carbohydrate,
			calories,
		))
	}

	_, err = f.WriteString(")")

}

func main() {

	products := getData()
	createSql(products)

}
