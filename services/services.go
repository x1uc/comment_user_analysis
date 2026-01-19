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

func (s *WeiboService) GetUsers(blogId string, uid string, num int, popular bool) ([]models.WeiboUser, []models.WeiboComment, error) {
	var users []models.WeiboUser
	var summary_comments []models.WeiboComment
	var max_id uint64
	var err error

	for len(users) < num {
		var comments []models.WeiboComment
		if popular {
			comments, max_id, err = s.agent.GetHotComments(blogId, uid, max_id)
		} else {
			comments, max_id, err = s.agent.GetNewComments(blogId, uid, max_id)
		}

		if err != nil {
			return nil, nil, err
		}

		for _, comment := range comments {
			users = append(users, comment.User)
			summary_comments = append(summary_comments, comment)
			for _, reply_comment := range comment.Comments {
				users = append(users, reply_comment.User)
				summary_comments = append(summary_comments, reply_comment)
			}
		}

		if max_id == 0 {
			break
		}
	}

	return users, summary_comments, nil
}

func (s *WeiboService) GetUserPhoneType(comment models.WeiboComment) (*models.CommentUserInfo, error) {
	blogs, err := s.agent.GetUserBlogs(comment.User.IDStr, 1)
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
	return &models.CommentUserInfo{
		Comment:    comment,
		Blog:       blog,
		PhoneType:  phone_type,
		PhoneBrand: brand,
	}, nil
}

func (s *WeiboService) GetCommentsForBlogs(blog_provider BlogProvider, pre_num int) ([]models.WeiboComment, error) {
	blog_ids, err := blog_provider.GetBlogs()
	if err != nil {
		return nil, err
	}

	var all_comments []models.WeiboComment
	for _, blog_id := range blog_ids {
		comments, err := s.GetComments(blog_id, pre_num)
		if err != nil {
			return nil, err
		}
		all_comments = append(all_comments, comments...)
	}

	return all_comments, nil
}

func (s *WeiboService) GetComments(blog_id string, num int) ([]models.WeiboComment, error) {
	var all_comments []models.WeiboComment
	var max_id uint64
	for len(all_comments) < num {
		comments, new_max_id, err := s.agent.GetNewComments(blog_id, "", max_id)
		if err != nil {
			return nil, err
		}
		all_comments = append(all_comments, comments...)
		if new_max_id == 0 {
			break
		}
		max_id = new_max_id
	}
	return all_comments, nil
}
