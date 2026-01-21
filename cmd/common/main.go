package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/client"
	"github.com/x1uc/comment_user_analysis/models"
	"github.com/x1uc/comment_user_analysis/services"
	"github.com/x1uc/comment_user_analysis/utils"
)

func main() {
	godotenv.Load()
	cookie := utils.RequireEnv("COOKIE")
	client := client.NewClient(cookie)
	rate := 3000 * time.Millisecond
	client.SetRateLimit(rate)

	weiboAgent := agent.NewService(client)
	weiboService := services.NewWeiboService(weiboAgent)
	blogProvider := &services.StaticBlogProvider{
		BlogIDs: []string{"5254998191509253"},
	}

	users, comments, err := weiboService.GetUsers(blogProvider, 20, false)

	phone_info_list := make([]models.UserPhoneInfo, 0)

	for _, user := range users {
		phone_info, err := weiboService.GetUserPhoneType(user)
		if err != nil {
			fmt.Printf("Error fetching phone type for user %s: %v\n", user.IDStr, err)
			continue
		}
		if phone_info == nil {
			fmt.Printf("No phone info for user %s\n", user.IDStr)
			continue
		}
		user_detail, err := weiboAgent.GetUserDetailInfo(user.IDStr)
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
		fmt.Printf("User: %s, Phone Type: %s, Brand: %s, Blog Count: %d\n", info.User.ScreenName, info.PhoneType, info.PhoneBrand, info.Detail.IPLocation)
	}
}
