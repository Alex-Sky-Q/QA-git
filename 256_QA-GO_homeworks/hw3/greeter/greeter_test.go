package greeter

import (
	"fmt"
	"testing"
)

func TestGreet(t *testing.T) {
	type TestCase struct {
		name string
		hour int
		exp  string
	}
	// Посчитайте необходимое количество тест кейсов для проверки всех возможных состояний функции Greet
	// Временные интервалы: 0-6, 6-12, 12-18, 18-22, 22-0 + <0 + >=24
	// По классам эквивалентности: 5 тест-кейсов (по одному в каждом интервале) + 2 тест-кейса <0 и >=24
	// По граничным значениям: 13 тест-кейсов (граничное значение и два его соседа) + 2 тест-кейса <0 и >=24
	// UPD. Выше посчитал видимо по методу черного ящика. Если считать по методу белого, то все состояния можно
	// проверить за 5 тест-кейсов: <0, 0-6, 6-12, 12-18, 18-22

	tcs := []TestCase{
		// Boundary values and equivalence classes
		{"alli", 0, "Good night Alli!"},
		{"Alli", 5, "Good night Alli!"},
		{"Alli", 6, "Good morning Alli!"},
		{"Alli", 7, "Good morning Alli!"},
		{"Alli", 11, "Good morning Alli!"},
		{"Alli", 12, "Hello Alli!"},
		{"Alli", 13, "Hello Alli!"},
		{"Alli", 17, "Hello Alli!"},
		{"Alli", 18, "Good evening Alli!"},
		{"Alli", 19, "Good evening Alli!"},
		{"Alli", 21, "Good evening Alli!"},
		{"Alli", 22, "Good night Alli!"},
		{"Alli", 23, "Good night Alli!"},
		// Trimmed spaces
		{" ted ", 23, "Good night Ted!"},
		{"ted ", 23, "Good night Ted!"},
		{" ted", 23, "Good night Ted!"},
		{" Ted", 23, "Good night Ted!"},
		{"  Ted", 23, "Good night Ted!"},
		// Different languages
		{"саша", 23, "Good night Саша!"},
		{"léa", 23, "Good night Léa!"},
		// Compound names
		{"Anna Maria", 23, "Good night Anna Maria!"},
		{"Anna-Maria", 23, "Good night Anna-Maria!"},
		{" Anna Maria ", 23, "Good night Anna Maria!"},
		{"anna-maria", 23, "Good night Anna-Maria!"},
		{"anna maria", 23, "Good night Anna Maria!"},
		// Language specific rules
		{"patrick o'brian", 23, "Good night Patrick O'Brian!"},
		{"leonardo da vinci", 23, "Good night Leonardo da Vinci!"},
		// Special symbols
		{"anna maria", 23, "Good night Anna Maria!"},
		{"anna\u00a0maria", 23, "Good night Anna\u00a0Maria!"},
		{"dr. smith", 23, "Good night Dr. Smith!"},
		// Possible negative scenarios
		{"Alli", 24, "Hour is out of range, should be between 0 and 23"},
		{"Alli", -1, "Hour is out of range, should be between 0 and 23"},
	}

	for _, tc := range tcs {
		t.Run("TestGreet "+fmt.Sprint(tc.hour), func(t *testing.T) {
			res := Greet(tc.name, tc.hour)
			if res != tc.exp {
				t.Fatalf("Got: '%s'. Want: '%s'", res, tc.exp)
			}
		})
	}

}
