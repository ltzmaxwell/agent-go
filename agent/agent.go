package agent

import (
	"demo/gogpt/client"
	"demo/gogpt/prompt"
	"demo/gogpt/tools"

	// tools "demo/gogpt/tools"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var code_example = `
package tool

import (
	"fmt"
	"strconv"
	"strings"
)

// IsPrime checks if a number is prime or not
func IsPrime(n int) bool {
    if n <= 1 {
        return false
    }
    for i := 2; i*i <= n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}

// PrimesUnderN takes a string representation of a number n and returns a string representation of prime numbers under n
func Do() string {
	limit := 50
    var primes []string
    for i := 2; i < limit; i++ {
        if IsPrime(i) {
            primes = append(primes, fmt.Sprintf("%d", i))
        }
    }
    return strings.Join(primes, ", ")
}
`

var OBSERVATION_TOKEN = "Observation:"
var systemPrompt = "If you are returning golang code, (e.g. I want to calculate prime numbers under 50), remember, give something like: " + code_example + ", which means that 1. there is no `main` func! instead, we have  another caller whose name is `Do`, with no input and string type return value(in this case: `Do`); 2. package name should be `tool`; 3. the caller function should return string only, its signature is something like func()string, I can handle string only. why I need this is that I have an interpreter to evaluate these function and get the results, thanks"

func Start(task string) {
	var thought string
	var round int
	for {
		println("round: ", round)

		prompt := prompt.GeneratePrompt("Golang Interpreter: A Golang interpreter. Use this to execute golang code. Input should be a valid golang program,  Note: Please provide the code without markdown formatting. \nGoogle Search: Get specific information from a search query. Input should be a  question like 'How to add number in Clojure?'. Result will be the answer to the question.", task, thought)
		println("prompt is: ", prompt)

		r := client.RequestLLM(0.0, OBSERVATION_TOKEN, systemPrompt, prompt, "gpt-4", 300)
		println("response start------------------")
		fmt.Printf("%s \n", r)

		println("response over------------------")
		// plan based on this repsonse
		// initiate tool to get real answer instead of hallication
		FinalPattern := `Final Answer:\s*(.*?)(?:\n|$)`
		err, finalAnswer := parse(r, FinalPattern)
		// final answer
		if err == nil {
			println("------Final answer is---------: ", finalAnswer)
			break
		}

		actionPattern := `Action:\s*(.*?)(?:\n|$)`
		err, action := parse(r, actionPattern)
		if err != nil {
			panic("action not match")
		}
		println("Action is: ", action)

		inputTypePattern := `Input Type:\s*(.*?)(?:\n|$)`
		err, inputType := parse(r, inputTypePattern)
		if err != nil {
			panic("input type not match")
		}
		println("Input type is: ", inputType)

		// actionInputCodePattern := "(?s)Action Input:[ \t]*\n```(golang|go)\n(.*?)\n```"
		actionInputCodePattern := "(?s)Action Input:([ \t]*\n```(golang|go)\n(.*?)\n```)"
		actionInputTextPattern := "(?s)Action Input:s*(.*?)(?:\n|$)"

		var actionInput, observation string
		if inputType == "code snippet" {
			err, actionInput = parse(r, actionInputCodePattern)
			if err != nil {
				panic("action input not match")
			}
			println("action Input is: ", actionInput)
			code := extractCodeFromMarkdown(actionInput)
			println("code start:")
			println(code)
			println("code end:")
			observation = tools.Evaluate(code)
		} else {
			err, actionInput = parse(r, actionInputTextPattern)
			if err != nil {
				panic("not match")
			}
			println("action Input is: ", actionInput)
			observation = tools.TavilySearch(actionInput)
			println("search result: ")
			println(observation)
		}
		thought = generateThought(action, inputType, actionInput, observation)
		println(("---thought---:"))
		println(thought)
		round++
	}
}

func parse(str, pattern string) (error, string) {
	exp := regexp.MustCompile(pattern)
	match := exp.FindStringSubmatch(str)

	if len(match) == 0 {
		return errors.New("not match"), ""
	}
	return nil, match[1]
}

func extractCodeFromMarkdown(markdown string) string {
	// Split by triple backticks
	parts := strings.Split(markdown, "```")
	if len(parts) < 3 {
		// Not enough parts for a valid markdown code block
		return ""
	}

	// The second part should be the code with or without a language hint
	codeWithHint := strings.TrimSpace(parts[1])

	// Split by newline to separate potential hint from actual code
	codeLines := strings.Split(codeWithHint, "\n")

	if len(codeLines) < 2 {
		// No language hint, just return the code
		return codeWithHint
	}

	// Return code without the first line (the language hint)
	return strings.Join(codeLines[1:], "\n")
}

func generateThought(action, inputType, actionInput, observation string) string {
	return fmt.Sprintf("Action: %s\nInputType: %s\nActionInput: %s\nObservation: %s\n", action, inputType, actionInput, observation)
}
