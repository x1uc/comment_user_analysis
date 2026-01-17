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
	userDetail, err := weiboAgent.GetUserDetailInfo("2607719317")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("%+v", userDetail)
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing required environment variable: %s", key))
	}
	return val
}
