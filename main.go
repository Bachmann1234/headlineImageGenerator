package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	headline := getRandomPopularNYTArticle()
	log.Print(headline)
	err := generateImage(headline)
	if err != nil {
		log.Print("Headline denied by DALLE2")
		censoredHeadline := censorHeadline(headline)
		log.Print(censoredHeadline)
		err := generateImage(censoredHeadline)
		if err != nil {
			panic(err)
		}
	}
}
