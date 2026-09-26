package pipeline

import (
	"context"
	"fmt"

	"github.com/x1uc/comment_user_analysis/models"
)

type Memory struct {
	Comments   []models.WeiboComment
	Users      []models.WeiboUser
	PhoneInfos []models.UserPhoneInfo
}

type Plugin interface {
	Name() string
	Run(ctx context.Context, mem *Memory) error
}

func Run(ctx context.Context, plugins []Plugin) error {
	mem := &Memory{
		Comments:   make([]models.WeiboComment, 0),
		Users:      make([]models.WeiboUser, 0),
		PhoneInfos: make([]models.UserPhoneInfo, 0),
	}
	for _, plugin := range plugins {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := plugin.Run(ctx, mem); err != nil {
			return fmt.Errorf("%s: %w", plugin.Name(), err)
		}
	}
	return nil
}
