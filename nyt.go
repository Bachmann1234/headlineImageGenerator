package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
)

/*
Makes a get request to a NYT api and returns the top item from the json list
*/
func getRandomPopularNYTArticle() string {
	resp, err := http.Get("https://api.nytimes.com/svc/mostpopular/v2/viewed/1.json?api-key=" + os.Getenv("NYT_API_KEY"))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}
	articles := data["results"].([]interface{})
	chosenArticle := articles[rand.Intn(len(articles))].(map[string]interface{})
	return chosenArticle["title"].(string)
}
