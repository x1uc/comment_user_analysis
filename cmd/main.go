package main

import (
	"context"
	"log"

	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/client"
	"github.com/x1uc/comment_user_analysis/config"
	"github.com/x1uc/comment_user_analysis/pipeline"
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

	err = pipeline.Run(context.Background(), []pipeline.Plugin{
		pipeline.NewCommentUserPlugin(cfg, weiboService),
		pipeline.NewPhoneInfoPlugin(weiboService, weiboAgent),
		pipeline.NewStorePlugin(store1),
	})
	if err != nil {
		log.Fatal(err)
	}
}
