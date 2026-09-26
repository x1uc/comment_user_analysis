package services

import (
	"strings"

	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/models"
)

type CommentOrderType string

const (
	OrderByPopular  CommentOrderType = "popular"
	OrderByTimeDesc CommentOrderType = "timeDesc"
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

func (s *WeiboService) GetCommentData(blogId string, commentAmount int, orderType CommentOrderType) (*models.ResultData, error) {
	resultData := models.ResultData{
		ResultComment: make([]models.WeiboComment, 0),
		ResultUsers:   make([]models.WeiboUser, 0),
	}

	var curComments []models.WeiboComment
	var curSubComments []models.WeiboComment
	curUsers := make([]models.WeiboUser, 0)
	var maxId uint64
	var err error

	for len(curComments) < commentAmount {
		if orderType == OrderByPopular {
			curComments, maxId, err = s.agent.GetHotComments(blogId, "", maxId)
		} else {
			curComments, maxId, err = s.agent.GetNewComments(blogId, "", maxId)
		}

		if err != nil {
			return nil, err
		}

		for _, comment := range curComments {
			curUsers = append(curUsers, comment.User)
			curSubComments = append(curSubComments, comment.SubComments...)
			for _, subComment := range comment.SubComments {
				curUsers = append(curUsers, subComment.User)
			}
		}

		if maxId == 0 {
			break
		}
	}
	resultData.ResultComment = append(resultData.ResultComment, curComments...)
	resultData.ResultComment = append(resultData.ResultComment, curSubComments...)
	resultData.ResultUsers = append(resultData.ResultUsers, curUsers...)

	return &resultData, nil
}

func (s *WeiboService) GetUserPhoneType(user models.WeiboUser) (*models.UserPhoneInfo, error) {
	blogs, err := s.agent.GetUserBlogs(user.IDStr, 1)
	if err != nil {
		return nil, err
	}
	if len(blogs) == 0 {
		return nil, nil
	}
	phoneType := ""
	brand := ""
	var blog models.WeiboBlog
	for _, curBlog := range blogs {
		blog = curBlog
		phoneType = strings.TrimSpace(strings.ToLower(blog.Source))
		for deviceName, curBrand := range models.BrandMap {
			if strings.Contains(phoneType, strings.ToLower(deviceName)) {
				brand = curBrand
				break
			}
		}
	}
	return &models.UserPhoneInfo{
		User:       user,
		Blog:       blog,
		PhoneType:  phoneType,
		PhoneBrand: brand,
	}, nil
}
