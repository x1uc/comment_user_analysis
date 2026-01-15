package main

import (
	"fmt"

	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/client"
)

func main() {
	cookie := ""
	weiboAgent := agent.NewService(client.NewClient(cookie))
	user, err := weiboAgent.GetUserInfo("2607719317")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("%w", user)
}
