package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/client"
	"github.com/x1uc/comment_user_analysis/services"
)

func main() {
	godotenv.Load()
	serviceGetUserTest()
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing required environment variable: %s", key))
	}
	return val
}

func serviceGetUserTest() {
	cookie := requireEnv("COOKIE")
	weiboAgent := agent.NewService(client.NewClient(cookie))
	weiboService := services.NewWeiboService(weiboAgent)
	users, comments, err := weiboService.GetUsers("5254618772671209", "2607719317", 10, true)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Fetched %d users and %d comments\n", len(users), len(comments))
	// for _, user := range users {
	// 	fmt.Printf("User: %+v\n", user)
	// }
	// for _, comment := range comments {
	// 	fmt.Printf("Comment: %+v\n", comment)
	// }
}
