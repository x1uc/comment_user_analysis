package main

import (
	"fmt"
	"time"

	"strings"

	"github.com/joho/godotenv"
	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/client"
	"github.com/x1uc/comment_user_analysis/services"
	"github.com/x1uc/comment_user_analysis/utils"
)

func main() {
	godotenv.Load()
	//serviceGetUserTest()
	//serviceGetBlogTest()
	//serviceGetStaticBlogTest()
	//serviceGetBlogTest()
	// serviceGetStaticBlogTest()
	serviceGetUserBlogTest()
}

func serviceGetUserTest() {
	cookie := utils.RequireEnv("COOKIE")
	weiboAgent := agent.NewService(client.NewClient(cookie))
	weiboService := services.NewWeiboService(weiboAgent)
	users, comments, err := weiboService.GetUsers("5254618772671209", "2607719317", 50, true)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Fetched %d users and %d comments\n", len(users), len(comments))
	// for _, comment := range comments {
	// 	fmt.Printf("Comment: %+v\n", comment)
	// }
}

func serviceGetBlogTest() {
	cookie := utils.RequireEnv("COOKIE")
	weiboAgent := agent.NewService(client.NewClient(cookie))
	blogProvider := &services.AgentBlogProvider{
		Agent:  weiboAgent,
		UID:    "2607719317",
		Number: 100,
	}
	blogs, err := blogProvider.GetBlogs()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Fetched %d blogs\n", len(blogs))
	for _, blogID := range blogs {
		fmt.Printf("Blog ID: %s\n", blogID)
	}

}

func serviceGetStaticBlogTest() {
	blogIDs := utils.RequireEnv("BLOG_IDS")
	blogProvider := &services.StaticBlogProvider{
		BlogIDs: strings.Split(blogIDs, ","),
	}
	blogs, err := blogProvider.GetBlogs()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Fetched %d blogs from static provider\n", len(blogs))
	for _, blogID := range blogs {
		fmt.Printf("Blog ID: %s\n", blogID)
	}
}

func serviceGetUserBlogTest() {
	cookie := utils.RequireEnv("COOKIE")
	client := client.NewClient(cookie)
	rate := 2000 * time.Millisecond
	client.SetRateLimit(rate)
	weiboAgent := agent.NewService(client)
	weiboService := services.NewWeiboService(weiboAgent)
	blogProvider := &services.StaticBlogProvider{
		BlogIDs: []string{"5254998191509253"},
	}

	comments, err := weiboService.GetCommentsForBlogs(blogProvider, 20)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Fetched %d comments\n", len(comments))

	for _, comment := range comments {
		commentUserInfo, err := weiboService.GetUserPhoneType(comment)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("Fetched comment user info: %+v\n", commentUserInfo)
	}
}
