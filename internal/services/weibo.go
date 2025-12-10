package services

import (
	"encoding/json"
	"fmt"
	"single_analysis/internal/client"
	"single_analysis/internal/models"
	"single_analysis/internal/utils"
)

// WeiboService 微博服务
type WeiboService struct {
	client       *client.Client
	phoneMapping models.PhoneBrandMapping
}

// NewWeiboService 创建微博服务
func NewWeiboService(cookie string) *WeiboService {
	return &WeiboService{
		client:       client.NewClient(cookie),
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
func (w *WeiboService) GetComments(blogID string, uid string, max_id uint64) ([]models.CommentData, uint64, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/buildComments?flow=0&is_reload=1&id=%s&is_show_bulletin=2&is_mix=0&max_id=%d&count=20&uid=%s&fetch_level=0&locale=zh-CN", blogID, max_id, uid)

	body, err := w.client.Get(url)
	if err != nil {
		return nil, 0, utils.NewNetworkError("获取评论列表失败", err)
	}

	var response models.CommentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, 0, utils.NewParseError("解析评论数据失败", err)
	}

	return response.Data, response.MaxID, nil
}

// GetUserPhoneType 获取用户手机类型
func (w *WeiboService) GetUserPhoneType(uid string) (string, error) {
	blogs, err := w.GetBlogs(uid, 1)
	if err != nil {
		return "", fmt.Errorf("获取用户博客失败: %w", err)
	}

	for _, blog := range blogs {
		if blog.User.ID == uid && blog.PhoneType != "" {
			return w.phoneMapping.GetBrand(blog.PhoneType), nil
		}
	}

	return "未知设备", nil
}

// GetUserBlogsAndComments 获取用户博客和评论用户（分页处理）
func (w *WeiboService) GetUserBlogsAndComments(uid string, blog_list []string, interval int, callback func([]models.CommentData)) {
	totalProcessed := 0
	userSet := make(map[string]bool)

	flag := true
	max_id := uint64(0)

	for _, blog_id := range blog_list {
		for flag || max_id != 0 {
			flag = false
			// 获取评论用户
			comments, cur_max_id, err := w.GetComments(blog_id, uid, max_id)
			max_id = cur_max_id
			if err != nil {
				fmt.Printf("获取评论失败: %v, 跳过博客 %s\n", err, blog_id)
				break
			}

			// 评论获得完了
			if max_id == 0 {
				fmt.Printf("没有更多评论\n")
				break
			}

			newComments := make([]models.CommentData, 0)
			for _, comment := range comments {
				if !userSet[comment.User.ID] {
					userSet[comment.User.ID] = true
				}
				newComments = append(newComments, comment)
			}

			// 如果有新用户，调用回调函数
			if len(newComments) > 0 {
				callback(newComments)
				totalProcessed += len(newComments)
				fmt.Printf("已处理 %d 个用户\n", totalProcessed)
			}
		}
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
