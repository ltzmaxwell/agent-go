package main

import (
	"demo/gogpt/agent"
)

func main() {
	question := "赵雷的最新专辑是什么"
	// question := "tell me all prime numbers under 50"
	agent.Start(question)
}
