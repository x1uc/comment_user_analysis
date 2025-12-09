package main

import (
	"single_analysis/internal/llm"
)

func main() {
	test := llm.DeepSeek{
		Api_key: "123",
	}
	test.GetCommentLevel()
}
