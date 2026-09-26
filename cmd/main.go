package main

import (
	"context"
	"github.com/google/uuid"
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

	DB, err := store.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer DB.Close()

	batchID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("batch id: %s", batchID)

	weiboAgent := agent.NewService(httpClient)
	weiboService := services.NewWeiboService(weiboAgent)

	err = pipeline.Run(context.Background(), []pipeline.Plugin{
		pipeline.NewCommentUserPlugin(cfg, weiboService),
		pipeline.NewPhoneInfoPlugin(weiboService, weiboAgent),
		pipeline.NewStorePlugin(DB, batchID.String()),
	})
	if err != nil {
		log.Fatal(err)
	}
}
