package pipeline

import (
	"context"
	"fmt"

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
	return "comment_user"
}

func (p *CommentUserPlugin) Run(ctx context.Context, mem *Memory) error {
	for _, blog := range p.cfg.BlogsInfos {
		if err := ctx.Err(); err != nil {
			return err
		}
		curData, err := p.service.GetUsers(blog.BlogId, blog.CommentAmount, services.CommentOrderType(blog.OrderType))
		if err != nil {
			return fmt.Errorf("fetch comments for blog %s: %w", blog.BlogId, err)
		}
		mem.Comments = append(mem.Comments, curData.ResultComment...)
		mem.Users = append(mem.Users, curData.ResultUsers...)
	}
	return nil
}
