package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type Order struct {
	UserID string `json:"user_id"`
	Date   string `json:"ordered_at"`
	Status string `json:"status"`
	Total  string `json:"total"`
}

func main() {
	// Чтение из файла
	data, err := os.ReadFile("input.json")
	if err != nil {
		panic(err)
	}

	// Распаковка данных в Order
	var orders []Order
	if err := json.Unmarshal(data, &orders); err != nil {
		panic(err)
	}

	// Траты по месяцам
	monthlySpending := make(map[time.Month]float64)
	for _, order := range orders {
		// Парсинг даты и времени заказа
		orderTime, err := parseOrderDate(order.Date)
		if err != nil {
			fmt.Printf("Ошибка при парсинге даты и времени для заказа %s: %v\n", order.UserID, err)
			continue
		}

		// Определение завершенных заказов
		if order.Status == "COMPLETED" {
			// Получение месяца заказа
			month := orderTime.Month()

			// Парсинг суммы заказа в float
			total, err := strconv.ParseFloat(order.Total, 64)
			if err != nil {
				fmt.Printf("Ошибка при парсинге суммы заказа для заказа %s: %v\n", order.UserID, err)
				continue
			}

			// Добавление суммы заказа в соответствующий месяц
			monthlySpending[month] += total
		}
	}

	// Нахождение c наибольшими тратами
	maxSpending := 0.0
	var maxSpendingMonths []time.Month
	for month, spending := range monthlySpending {
		if spending > maxSpending {
			maxSpending = spending
			maxSpendingMonths = []time.Month{month}
		} else if spending == maxSpending {
			maxSpendingMonths = append(maxSpendingMonths, month)
		}
	}

	// Сортировка по возрастанию
	sort.Slice(maxSpendingMonths, func(i, j int) bool {
		return maxSpendingMonths[i] < maxSpendingMonths[j]
	})

	// Формирование отчета
	report := make(map[string][]string)
	for _, month := range maxSpendingMonths {
		report["months"] = append(report["months"], month.String())
	}

	// Преобразование отчета в JS
	reportJSON, err := json.Marshal(report)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(reportJSON))
}

// Функция парсера для строки, представляющей дату и время заказа
func parseOrderDate(dateStr string) (time.Time, error) {
	orderTime, err := time.Parse("2006-01-02T15:04:05.000", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return orderTime, nil
}
