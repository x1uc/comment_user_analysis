package services

import (
	"github.com/x1uc/comment_user_analysis/agent"
	"github.com/x1uc/comment_user_analysis/models"
)

type WeiboService struct {
	agent *agent.WeiboAgent
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
