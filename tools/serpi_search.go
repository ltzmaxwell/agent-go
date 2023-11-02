package tools

import (
	"fmt"
	g "github.com/serpapi/google-search-results-golang"
)

func Search(query string) {
	parameter := map[string]string{
		"q":             query,
		"location":      "Austin, Texas, United States",
		"hl":            "en",
		"gl":            "us",
		"google_domain": "google.com",
		"api_key":       "4c8f07c0c2204563d6a4be12962555b7cd8a02093c3d394b8619ae017a093201",
	}

	search := g.NewGoogleSearch(parameter, "4c8f07c0c2204563d6a4be12962555b7cd8a02093c3d394b8619ae017a093201")
	results, err := search.GetJSON()
	if err != nil {
		println(err.Error())
	}

	fmt.Printf("result: %T \n", results)
	fmt.Printf("result: %s \n", results)

	r := results["organic_results"].([]interface{})
	first_result := r[0].(map[string]interface{})
	fmt.Println(first_result["title"].(string))
}

// func processResponse(results g.SearchResult) {

// }
