package pipeline

import (
	"context"
	"log/slog"

	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/services"
)

type PhoneInfoPlugin struct {
	service *services.WeiboService
	agent   *agent.WeiboAgent
}

func NewPhoneInfoPlugin(service *services.WeiboService, weiboAgent *agent.WeiboAgent) *PhoneInfoPlugin {
	return &PhoneInfoPlugin{service: service, agent: weiboAgent}
}

func (p *PhoneInfoPlugin) Name() string {
	return "手机信息"
}

func (p *PhoneInfoPlugin) Run(ctx context.Context, mem *Memory) error {
	for _, user := range mem.Users {
		if err := ctx.Err(); err != nil {
			return err
		}
		phoneInfo, err := p.service.GetUserPhoneType(user)
		if err != nil {
			slog.Error("获取用户手机信息失败", "user_id", user.IDStr, "error", err)
			continue
		}
		if phoneInfo == nil {
			slog.Info("用户没有手机信息", "user_id", user.IDStr)
			continue
		}
		userDetail, err := p.agent.GetUserDetailInfo(user.IDStr)
		if err != nil {
			slog.Error("获取用户详情失败", "user_id", user.IDStr, "error", err)
			continue
		}
		phoneInfo.Detail = *userDetail
		mem.PhoneInfos = append(mem.PhoneInfos, *phoneInfo)
		slog.Info("获取用户手机信息成功", "user_id", user.IDStr, "user", user.ScreenName, "brand", phoneInfo.PhoneBrand, "phone", phoneInfo.PhoneType)
	}
	return nil
}
