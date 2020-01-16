package main

import (
	"fmt"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"
)

func main() {
	song := domain.Song{
		Title:  "Всё идёт по плану",
		Author: "Е.Летов",
		Album:  "Всё идёт по плану",
		Verses: []domain.Verse{
			domain.Verse{
				Quotes: []domain.Quote{
					domain.Quote{
						Phrase: "Границы ключ переломлен пополам",
					},
					domain.Quote{
						Phrase: "А наш батюшка Ленин совсем усоп",
					},
					domain.Quote{
						Phrase: "А перестройка всё идёт и идёт по плану",
					},
					domain.Quote{
						Phrase: "А вся грязь превратилась в голый лёд",
					},
				},
			},
		},
	}
	fmt.Println("My first song is", song.GetQuotes())
}
