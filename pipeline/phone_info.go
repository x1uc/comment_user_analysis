package pipeline

import (
	"context"
	"fmt"

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
	return "phone_info"
}

func (p *PhoneInfoPlugin) Run(ctx context.Context, mem *Memory) error {
	for _, user := range mem.Users {
		if err := ctx.Err(); err != nil {
			return err
		}
		phoneInfo, err := p.service.GetUserPhoneType(user)
		if err != nil {
			fmt.Printf("Error fetching phone type for user %s: %v\n", user.IDStr, err)
			continue
		}
		if phoneInfo == nil {
			fmt.Printf("No phone info for user %s\n", user.IDStr)
			continue
		}
		userDetail, err := p.agent.GetUserDetailInfo(user.IDStr)
		if err != nil {
			fmt.Printf("Error fetching phone type for user %s: %v\n", user.IDStr, err)
			continue
		}
		phoneInfo.Detail = *userDetail
		mem.PhoneInfos = append(mem.PhoneInfos, *phoneInfo)
	}
	return nil
}
