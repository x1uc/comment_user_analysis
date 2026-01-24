package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/client"
	"github.com/x1uc/comment_user_analysis/models"
	"github.com/x1uc/comment_user_analysis/services"
	"github.com/x1uc/comment_user_analysis/store"
	"github.com/x1uc/comment_user_analysis/utils"
)

func main() {
	godotenv.Load()
	cookie := utils.RequireEnv("COOKIE")
	client := client.NewClient(cookie)
	rate := 2000 * time.Millisecond
	client.SetRateLimit(rate)

	store1 := store.NewStore("test.db")

	weibo_agent := agent.NewService(client)
	weibo_service := services.NewWeiboService(weibo_agent)
	blog_provider := &services.StaticBlogProvider{
		BlogIDs: []string{"5254998191509253"},
	}

	users, comments, err := weibo_service.GetUsers(blog_provider, 20, false)

	phone_info_list := make([]models.UserPhoneInfo, 0)

	for _, user := range users {
		phone_info, err := weibo_service.GetUserPhoneType(user)
		if err != nil {
			fmt.Printf("Error fetching phone type for user %s: %v\n", user.IDStr, err)
			continue
		}
		if phone_info == nil {
			fmt.Printf("No phone info for user %s\n", user.IDStr)
			continue
		}
		user_detail, err := weibo_agent.GetUserDetailInfo(user.IDStr)
		if err != nil {
			fmt.Printf("Error fetching phone type for user %s: %v\n", user.IDStr, err)
			continue
		}
		phone_info.Detail = *user_detail
		phone_info_list = append(phone_info_list, *phone_info)
	}
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Fetched %d comments\n", len(comments))
	fmt.Printf("USER LEN %d", len(users))
	for _, info := range phone_info_list {
		if err := store1.InsertInfo(info); err != nil {
			log.Printf("Failed to insert info for user %s: %v", info.User.ScreenName, err)
		}
		fmt.Printf("User: %s, Phone Type: %s, Brand: %s, Blog Count: %s\n", info.User.ScreenName, info.PhoneType, info.PhoneBrand, info.Detail.IPLocation)
	}
}
