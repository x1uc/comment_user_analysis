package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"

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
		slog.Error("加载配置失败", "error", err)
		os.Exit(1)
	}

	rateLimit, err := cfg.RateLimitDuration()
	if err != nil {
		slog.Error("解析限流间隔失败", "error", err)
		os.Exit(1)
	}

	httpClient := client.NewClient(cfg.Cookie, rateLimit)

	DB, err := store.NewStore(cfg.DBPath)
	if err != nil {
		fatal("初始化数据库失败", err)
	}
	defer DB.Close()

	batchID, err := uuid.NewUUID()
	if err != nil {
		fatal("生成批次号失败", err)
	}
	slog.Info("批次号", "batch_id", batchID)

	weiboAgent := agent.NewService(httpClient)
	weiboService := services.NewWeiboService(weiboAgent)

	err = pipeline.Run(context.Background(), []pipeline.Plugin{
		pipeline.NewCommentUserPlugin(cfg, weiboService),
		pipeline.NewPhoneInfoPlugin(weiboService, weiboAgent),
		pipeline.NewStorePlugin(DB, batchID.String()),
		pipeline.NewPhoneBrandChartPlugin(filepath.Join("output", "phone_brand_pie"+batchID.String()+".html")),
	})
	if err != nil {
		fatal("执行失败", err)
	}
}

func fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}
