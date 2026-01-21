package services

import (
	"strings"

	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/models"
)

type WeiboService struct {
	agent *agent.WeiboAgent
}

type BlogProvider interface {
	GetBlogs() ([]string, error)
}

type StaticBlogProvider struct {
	BlogIDs []string
}

func (p *StaticBlogProvider) GetBlogs() ([]string, error) {
	return p.BlogIDs, nil
}

type AgentBlogProvider struct {
	Agent  *agent.WeiboAgent
	UID    string
	Number int
}

func (p *AgentBlogProvider) GetBlogs() ([]string, error) {

	cur_count := 0
	page := 1
	var ids []string

	for cur_count < p.Number {
		blogs, err := p.Agent.GetUserBlogs(p.UID, page)
		if err != nil {
			return nil, err
		}

		for _, blog := range blogs {
			ids = append(ids, blog.IDStr)
			cur_count++
			if cur_count >= p.Number {
				break
			}
		}

		if len(blogs) == 0 {
			break
		}
		page++
	}
	return ids, nil
}

func NewWeiboService(agent *agent.WeiboAgent) *WeiboService {
	return &WeiboService{agent: agent}
}

// num: number of users to fetch per blog,
// popular: whether to fetch popular comments or new comments
func (s *WeiboService) GetUsers(blog_provider BlogProvider, num int, popular bool) ([]models.WeiboUser, []models.WeiboComment, error) {
	blog_ids, err := blog_provider.GetBlogs()
	if err != nil {
		return nil, nil, err
	}
	summary_users := make([]models.WeiboUser, 0)
	summary_comments := make([]models.WeiboComment, 0)
	for _, blog_id := range blog_ids {

		var users []models.WeiboUser
		var max_id uint64
		var err error

		for len(users) < num {
			var comments []models.WeiboComment
			if popular {
				comments, max_id, err = s.agent.GetHotComments(blog_id, "", max_id)
			} else {
				comments, max_id, err = s.agent.GetNewComments(blog_id, "", max_id)
			}

			if err != nil {
				return nil, nil, err
			}

			for _, comment := range comments {
				users = append(users, comment.User)
				summary_comments = append(summary_comments, comment)
				for _, reply_comment := range comment.SubComments {
					users = append(users, reply_comment.User)
					summary_comments = append(summary_comments, reply_comment)
				}
			}

			if max_id == 0 {
				break
			}
		}
		summary_users = append(summary_users, users...)
	}

	return summary_users, summary_comments, nil
}

func (s *WeiboService) GetUserPhoneType(user models.WeiboUser) (*models.UserPhoneInfo, error) {
	blogs, err := s.agent.GetUserBlogs(user.IDStr, 1)
	if err != nil {
		return nil, err
	}
	if len(blogs) == 0 {
		return nil, nil
	}
	phone_type := ""
	brand := ""
	var blog models.WeiboBlog
	for _, cur_blog := range blogs {
		blog = cur_blog
		phone_type = strings.TrimSpace(strings.ToLower(blog.Source))
		for device_name, cur_brand := range models.BrandMap {
			if strings.Contains(phone_type, strings.ToLower(device_name)) {
				brand = cur_brand
				break
			}
		}
	}
	return &models.UserPhoneInfo{
		User:       user,
		Blog:       blog,
		PhoneType:  phone_type,
		PhoneBrand: brand,
	}, nil
}
