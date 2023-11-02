package prompt

import "time"
import "fmt"

var code_example = `
package foo

import (
	"fmt"
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

// PrimesUnder10 returns a string representation of prime numbers under 10
func PrimesUnder10() string {
    var primes []string
    for i := 2; i < 10; i++ {
        if IsPrime(i) {
            primes = append(primes, fmt.Sprintf("%d", i))
        }
    }
    return strings.Join(primes, ", ")
}
`

var PROMPT_TEMPLATE string = `Today is %s and you can use tools to get new information. Answer the question as best as you can using the following tools: 

%s

Use the following format:

Question: the input question you must answer
Thought: comment on what you want to do next
Action: the action to take, exactly one element of [{tool_names}]
Input Type: code snippet or plain text.
Action Input: the input to the action, Note: please provide the code in plain text without any special formatting, no backtick.
Observation: the result of the action
... (this Thought/Action/Action Input/Observation repeats N times, use it until you are sure of the answer)
Thought: I now know the final answer
Final Answer: your final answer to the original input question


Begin!

Question: %s
Thought: %s
`

func GeneratePrompt(toolsDesc, question, thought string) string {
	today := time.Now().Format("2006-Jan-02")
	return fmt.Sprintf(PROMPT_TEMPLATE, today, toolsDesc, question, thought)
}
