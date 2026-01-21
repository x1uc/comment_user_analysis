package main

import (
	"fmt"

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
	//serviceGetStaticBlogTest()
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
