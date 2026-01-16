package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/client"
)

func main() {
	godotenv.Load()
	cookie := requireEnv("COOKIE")
	fmt.Printf("%s", cookie)
	weiboAgent := agent.NewService(client.NewClient(cookie))
	comments, max_id, err := weiboAgent.GetHotComments("5254618772671209", "2607719317", 0)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print("max_id:", max_id, "\n")
	fmt.Printf("arr len: %d", len(comments))
	sum := 0
	for _, comment := range comments {
		sum++
		sum += len(comment.Comments)
	}
	fmt.Print("total comments:", sum)
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing required environment variable: %s", key))
	}
	return val
}
