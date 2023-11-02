package tools

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

const (
	apiKey      = ""
	apiEndpoint = "https://api.tavily.com/search"
)

func TavilySearch(query string) string {
	client := resty.New()
	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"api_key": apiKey,
			// "messages":            []interface{}{map[string]interface{}{"role": "user", "content": userPrompt}},
			"query":               query,
			"search_depth":        "basic",
			"include_answer":      false,
			"include_images":      true,
			"include_raw_content": false,
			"max_results":         5,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}

	body := response.Body()
	fmt.Printf("body: %s: \n", body)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
	}
	content := data["results"].([]interface{})[0].(map[string]interface{})["title"].(string)
	println("content: ", content)
	return content
}
