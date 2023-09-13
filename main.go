package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type ArithmeticResponse struct {
	Result int `json:"result"`
}

func handleArithmetic(w http.ResponseWriter, r *http.Request) {
	// Проверяем наличие прав доступа в заголовке User-Access
	access := r.Header.Get("User-Access")
	if access != "superuser" {
		// Выводим сообщение в консоли
		log.Println("Access denied")

		// Отправляем соответствующий ответ клиенту
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Access denied"})
		return
	}

	// Получаем выражение из параметра запроса
	expression := r.URL.Query().Get("expression")
	terms := strings.Split(expression, " ")

	if len(terms) < 3 {
		// Некорректное выражение
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid expression"})
		return
	}

	result, err := strconv.Atoi(terms[0])
	if err != nil {
		// Некорректный операнд
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid operand"})
		return
	}

	for i := 1; i < len(terms); i += 2 {
		operator := terms[i]
		operand, err := strconv.Atoi(terms[i+1])
		if err != nil {
			// Некорректный операнд
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid operand"})
			return
		}

		switch operator {
		case "+":
			result += operand
		case "-":
			result -= operand
		default:
			// Некорректный оператор
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid operator"})
			return
		}
	}

	response := ArithmeticResponse{
		Result: result,
	}

	// Отправляем JSON-ответ со статусом 200
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", handleArithmetic)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
