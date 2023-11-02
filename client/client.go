package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

// Use your API KEY here

const (
	apiKey      string = ""
	apiEndpoint        = "https://api.openai.com/v1/chat/completions"
)

func RequestLLM(temperature float32, stop string, systemPrompt string, userPrompt string, model string, maxTokens int) string {
	client := resty.New()
	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":       model,
			"messages":    []interface{}{map[string]interface{}{"role": "user", "content": userPrompt}, map[string]interface{}{"role": "system", "content": systemPrompt}},
			"max_tokens":  maxTokens,
			"temperature": temperature,
			"stop":        stop,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}

	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return ""
	}

	// Extract the content from the JSON response
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	// fmt.Println(content)
	return content
}
