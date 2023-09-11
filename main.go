package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Product struct {
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Rating float64 `json:"rating"`
}

func main() {
	filePath := os.Args[1] // Путь к файлу передается в качестве аргумента командной строки

	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Создаем срез для хранения данных о продуктах
	var products []Product

	// Определяем формат файла на основе его расширения
	fileExt := filepath.Ext(filePath)

	// Читаем данные из файла в зависимости от формата
	switch fileExt {
	case ".csv":
		products, err = readCSV(file)
		if err != nil {
			log.Fatal(err)
		}
	case ".json":
		products, err = readJSON(file)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Неподдерживаемый формат файла: %s", fileExt)
	}

	// Находим "самый дорогой продукт"
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price > products[j].Price
	})
	mostExpensive := products[0]

	// Находим "самый высокий рейтинг"
	sort.Slice(products, func(i, j int) bool {
		return products[i].Rating > products[j].Rating
	})
	highestRating := products[0]

	// Выводим результаты
	fmt.Printf("Самый дорогой продукт: %s (цена: %.2f)\n", mostExpensive.Name, mostExpensive.Price)
	fmt.Printf("Самый высокий рейтинг: %s (рейтинг: %.2f)\n", highestRating.Name, highestRating.Rating)
}

// Функция для чтения данных из CSV-файла
func readCSV(file io.Reader) ([]Product, error) {
	reader := csv.NewReader(bufio.NewReader(file))
	var products []Product

	// Пропускаем заголовок CSV-файла
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Читаем данные из CSV-файла
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		price, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}

		rating, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}

		product := Product{
			Name:   record[0],
			Price:  price,
			Rating: rating,
		}
		products = append(products, product)
	}

	return products, nil
}

// Функция для чтения данных из JSON-файла
func readJSON(file io.Reader) ([]Product, error) {
	decoder := json.NewDecoder(bufio.NewReader(file))
	var products []Product

	// Читаем данные из JSON-файла
	err := decoder.Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}
