package crawler

import (
	"comment_phone_analyse/pojo"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 响应体结构体
type Response struct {
	Data ResponseData `json:"data"`
}

// 解析 res.data
type ResponseData struct {
	List []pojo.BlogList `json:"list"`
}

// 创建一个辅助函数来处理错误
func handleError(msg string, err error) *Response {
	if err != nil {
		fmt.Println(msg, err)
	}
	return nil
}

// 设置请求头的函数
func SetHeaders(req *http.Request, cookie string) {
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Client-Version", "v2.47.25")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-CH-UA", `"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Server-Version", "v2025.01.23.1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

// 获取博客数据的函数
func GetBlog(uid string, cookie string, page int) *Response {
	url := fmt.Sprintf("%s%s%s%d%s", "https://weibo.com/ajax/statuses/mymblog?uid=", uid, "&page=", page, "&feature=0")

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return handleError("Error creating request:", err)
	}

	// 设置请求头
	SetHeaders(req, cookie)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return handleError("Request error:", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return handleError("Error reading response:", err)
	}

	// 解析 JSON 数据
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		fmt.Println("Response Body:", string(body)) // 打印原始 JSON 以便调试
		return nil
	}

	return &result
}
