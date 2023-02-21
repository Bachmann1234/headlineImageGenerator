package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func writeResponse(body []byte) {
	file, err := os.Create("headline.webp")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	_, err = file.Write(body)
	if err != nil {
		log.Fatalln(err)
	}
}

func generateImage(prompt string) error {
	client := &http.Client{}
	postBody, _ := json.Marshal(map[string]any{
		"prompt":          prompt,
		"n":               1,
		"size":            "1024x1024",
		"response_format": "b64_json",
	})
	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/images/generations",
		requestBody,
	)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
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
		panic(err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	//['data'][0]['b64_json']
	if data["error"] != nil {
		// Headline either too violent or contains a celebrity
		return errors.New("Dalle Rejected headline")
	}
	item := data["data"].([]interface{})[0].(map[string]interface{})
	encodedImage := item["b64_json"].(string)
	imageBytes, err := base64.StdEncoding.DecodeString(encodedImage)
	if err != nil {
		log.Fatal("error:", err)
	}
	writeResponse(imageBytes)
	return nil
}

func censorHeadline(headline string) string {
	client := &http.Client{}
	postBody, _ := json.Marshal(map[string]any{
		"model":             "text-davinci-003",
		"prompt":            "Make this NYT headline so it is PG rated, contains no names, and retains its original meaning: \"" + headline + "\"",
		"temperature":       0.7,
		"max_tokens":        256,
		"top_p":             1,
		"frequency_penalty": 0,
		"presence_penalty":  0,
	})
	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/completions",
		requestBody,
	)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
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
		panic(err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		// Headline either too violent or contains a celebrity
		panic(err)
	}
	quotedHeadline := strings.TrimSpace(data["choices"].([]interface{})[0].(map[string]interface{})["text"].(string))
	return strings.TrimPrefix(strings.TrimSuffix(quotedHeadline, "\""), "\"")
}
