package crawler

import (
	"comment_phone_analyse/pojo"
	"encoding/json"
	"fmt"
	"net/http"
)

// 处理错误的辅助函数
func HandleError(msg string, err error) *pojo.CommentData {
	if err != nil {
		fmt.Println(msg, err)
	}
	return nil
}

// 获取评论数据的函数
func GetCommentUser(cookie string, idstr string, uid string) *pojo.CommentData {
	first_url := "https://weibo.com/ajax/statuses/buildComments?flow=0&is_reload=1&id=" + idstr + "&is_show_bulletin=2&is_mix=0&count=10&uid=" + uid + "&fetch_level=0&locale=zh-CN"

	// 创建请求
	req, err := http.NewRequest("GET", first_url, nil)
	if err != nil {
		return HandleError("Error creating request:", err)
	}

	// 设置请求头
	SetHeaders(req, cookie)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return HandleError("Request error:", err)
	}
	defer resp.Body.Close()

	var res pojo.CommentData

	// 解析响应
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return HandleError("Error decoding response:", err)
	}

	// 返回解析后的评论数据
	return &res
}
