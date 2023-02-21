package main

import (
	"log"
)

func main() {
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
