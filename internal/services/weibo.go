package services

import (
	"comment_phone_analyse/config"
	"comment_phone_analyse/internal/client"
	"comment_phone_analyse/internal/models"
	"comment_phone_analyse/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// WeiboService 微博服务
type WeiboService struct {
	client       *client.Client
	phoneMapping models.PhoneBrandMapping
}

// NewWeiboService 创建微博服务
func NewWeiboService() *WeiboService {
	cfg := config.GetGlobalConfig()
	return &WeiboService{
		client:       client.NewClient(cfg.Cookie),
		phoneMapping: getDefaultPhoneMapping(),
	}
}

// GetBlogs 获取用户博客列表
func (w *WeiboService) GetBlogs(uid string, page int) ([]models.Blog, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/mymblog?uid=%s&page=%d&feature=0", uid, page)

	body, err := w.client.Get(url)
	if err != nil {
		return nil, utils.NewNetworkError("获取博客列表失败", err)
	}

	var response models.BlogResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, utils.NewParseError("解析博客数据失败", err)
	}

	if len(response.Data.List) == 0 {
		return nil, utils.ErrNoMoreData
	}

	return response.Data.List, nil
}

// GetComments 获取博客评论用户列表
func (w *WeiboService) GetComments(blogID string, uid string, max_id uint64) (*models.CommentResponse, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/buildComments?flow=0&is_reload=1&id=%s&is_show_bulletin=2&is_mix=0&count=20&uid=%s&fetch_level=0&locale=zh-CN&max_id=%v", blogID, uid, max_id)

	body, err := w.client.Get(url)
	if err != nil {
		return nil, utils.NewNetworkError("获取评论列表失败", err)
	}

	var response models.CommentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, utils.NewParseError("解析评论数据失败", err)
	}

	return &response, nil
}

// GetUserPhoneType 获取用户手机类型
func (w *WeiboService) GetUserPhoneType(uid string) (string, error) {
	blogs, err := w.GetBlogs(uid, 1)
	if err != nil {
		return "", fmt.Errorf("获取用户博客失败: %w", err)
	}
	userPhone := ""
	for _, blog := range blogs {
		if blog.User.ID == uid && blog.PhoneType != "" {
			curBrand := w.phoneMapping.GetBrand(blog.PhoneType)
			if w.IsKnownBrand(curBrand) {
				return curBrand, nil
			} else {
				userPhone = curBrand
			}
		}
	}
	if userPhone != "" {
		return userPhone, nil
	}
	return "未知设备", nil
}

// GetUserBlogsAndComments 获取用户博客和评论用户
func (w *WeiboService) GetUserBlogsAndComments(callback func([]models.CommentUser)) {
	cfg := config.GetGlobalConfig()

	page := 1
	totalProcessed := 0
	processedUsers := make(map[string]bool)

	for totalProcessed < cfg.Limit {
		// 获取博客列表
		blogs, err := w.GetBlogs(cfg.UID, page)
		if err != nil {
			if errors.Is(err, utils.ErrNoMoreData) {
				fmt.Println("没有更多博客了")
				break
			}
			fmt.Printf("获取第%d页博客失败: %v\n", page, err)
			break
		}

		// 处理每条博客的评论
		for _, blog := range blogs {
			// 只处理用户本人发布的博客
			if blog.User.ID != cfg.UID {
				continue
			}

			if totalProcessed >= cfg.Limit {
				break
			}

			singleCount := 0
			maxID := uint64(0)
			isFirstPage := true

			// 获取该博客的所有评论
			for maxID != 0 || isFirstPage {
				isFirstPage = false

				response, err := w.GetComments(blog.MblogID, cfg.UID, maxID)
				if err != nil {
					fmt.Printf("获取评论失败: %v, 跳过博客 %s\n", err, blog.MblogID)
					continue
				}

				maxID = response.MaxID

				// 过滤重复用户
				var newUsers []models.CommentUser
				for _, comment := range response.Data {
					if !processedUsers[comment.User.ID] {
						processedUsers[comment.User.ID] = true
						newUsers = append(newUsers, comment.User)
					}
				}

				// 调用回调处理新用户
				if len(newUsers) > 0 {
					callback(newUsers)
					totalProcessed += len(newUsers)
					singleCount += len(newUsers)
					fmt.Printf("已处理 %d 个用户\n", totalProcessed)
				}

				// 检查限制
				if singleCount >= cfg.SingleLimit || totalProcessed >= cfg.Limit {
					break
				}
			}
		}

		if totalProcessed >= cfg.Limit {
			break
		}

		// 页面间延迟
		time.Sleep(time.Duration(cfg.Interval) * time.Second)
		page++
	}
}

// getDefaultPhoneMapping 获取默认的手机品牌映射
func getDefaultPhoneMapping() models.PhoneBrandMapping {
	return models.PhoneBrandMapping{
		"Huawei":    "华为",
		"华为":        "华为",
		"nova":      "华为",
		"HarmonyOS": "华为",
		"Xiaomi":    "小米",
		"小米":        "小米",
		"OPPO":      "OPPO",
		"Vivo":      "Vivo",
		"iPhone":    "苹果",
		"苹果":        "苹果",
		"Samsung":   "三星",
		"三星":        "三星",
		"Meizu":     "魅族",
		"魅族":        "魅族",
		"realme":    "真我",
		"真我":        "真我",
		"redmi":     "红米",
		"红米":        "红米",
		"一加":        "一加",
		"OnePlus":   "一加",
		"荣耀":        "荣耀",
		"Honor":     "荣耀",
		"honor":     "荣耀",
		"ZTE":       "中兴",
		"中兴":        "中兴",
		"Nubia":     "努比亚",
		"努比亚":       "努比亚",
		"IQOO":      "IQOO",
		"Neo5":      "IQOO",
		"Android":   "Android设备",
	}
}

// IsKnownBrand 检查是否为已知品牌
func (w *WeiboService) IsKnownBrand(phoneType string) bool {
	knownBrands := map[string]bool{
		"华为":        true,
		"小米":        true,
		"OPPO":      true,
		"Vivo":      true,
		"苹果":        true,
		"三星":        true,
		"魅族":        true,
		"真我":        true,
		"红米":        true,
		"一加":        true,
		"荣耀":        true,
		"中兴":        true,
		"努比亚":       true,
		"IQOO":      true,
		"未知Android": true,
	}
	return knownBrands[phoneType]
}
