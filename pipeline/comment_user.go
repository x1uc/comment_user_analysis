package pipeline

import (
	"context"
	"log/slog"

	"github.com/x1uc/comment_user_analysis/config"
	"github.com/x1uc/comment_user_analysis/services"
)

type CommentUserPlugin struct {
	cfg     *config.BlogCrawlInfos
	service *services.WeiboService
}

func NewCommentUserPlugin(cfg *config.BlogCrawlInfos, service *services.WeiboService) *CommentUserPlugin {
	return &CommentUserPlugin{cfg: cfg, service: service}
}

func (p *CommentUserPlugin) Name() string {
	return "拉取评论"
}

func (p *CommentUserPlugin) Run(ctx context.Context, mem *Memory) error {
	for _, blog := range p.cfg.BlogsInfos {
		if err := ctx.Err(); err != nil {
			return err
		}
		curData, err := p.service.GetCommentData(blog.BlogId, blog.CommentAmount, services.CommentOrderType(blog.OrderType))
		if err != nil {
			slog.Error("拉取微博评论失败", "blog_id", blog.BlogId, "error", err)
			return err
		}

		mem.Comments = append(mem.Comments, curData.ResultComment...)
		mem.Users = append(mem.Users, curData.ResultUsers...)
		slog.Info("拉取微博成功", "blog_id", blog.BlogId, "comments", len(curData.ResultComment), "users", len(curData.ResultUsers))
	}
	return nil
}
