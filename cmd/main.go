package main

import (
	"fmt"
	"log"

	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/client"
	"github.com/x1uc/comment_user_analysis/config"
	"github.com/x1uc/comment_user_analysis/models"
	"github.com/x1uc/comment_user_analysis/services"
	"github.com/x1uc/comment_user_analysis/store"
)

func main() {
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	rateLimit, err := cfg.RateLimitDuration()
	if err != nil {
		log.Fatal(err)
	}

	httpClient := client.NewClient(cfg.Cookie, rateLimit)

	store1, err := store.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store1.Close()

	weiboAgent := agent.NewService(httpClient)
	weiboService := services.NewWeiboService(weiboAgent)

	resultData := models.ResultData{
		ResultComment:    make([]models.WeiboComment, 0),
		ResultUsers:      make([]models.WeiboUser, 0),
		ResultPhoneInfos: make([]models.UserPhoneInfo, 0),
	}
	for _, blog := range cfg.BlogsInfos {
		curData, err := weiboService.GetUsers(blog.BlogId, blog.CommentAmount, services.CommentOrderType(blog.OrderType))

		if err != nil {
			log.Fatalf("Failed to fetch comments for blog %s: %v", blog.BlogId, err)
		}
		resultData.ResultComment = append(resultData.ResultComment, curData.ResultComment...)
		resultData.ResultUsers = append(resultData.ResultUsers, curData.ResultUsers...)
	}

	phoneInfoList := make([]models.UserPhoneInfo, 0)

	for _, user := range resultData.ResultUsers {
		phoneInfo, err := weiboService.GetUserPhoneType(user)
		if err != nil {
			fmt.Printf("Error fetching phone type for user %s: %v\n", user.IDStr, err)
			continue
		}
		if phoneInfo == nil {
			fmt.Printf("No phone info for user %s\n", user.IDStr)
			continue
		}
		userDetail, err := weiboAgent.GetUserDetailInfo(user.IDStr)
		if err != nil {
			fmt.Printf("Error fetching phone type for user %s: %v\n", user.IDStr, err)
			continue
		}
		phoneInfo.Detail = *userDetail
		phoneInfoList = append(phoneInfoList, *phoneInfo)
	}
	resultData.ResultPhoneInfos = append(resultData.ResultPhoneInfos, phoneInfoList...)
}
