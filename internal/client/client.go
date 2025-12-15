package client

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client 微博API客户端
type Client struct {
	httpClient *http.Client
	cookie     string
}

// NewClient 创建新的微博客户端
func NewClient(cookie string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cookie: cookie,
	}
}

// Get 发送GET请求并处理gzip压缩
func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 处理压缩响应
	reader := c.getReader(resp.Body, resp.Header.Get("Content-Encoding"))

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return body, nil
}

// setHeaders 设置请求头
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Client-Version", "v2.47.130")
	req.Header.Set("Cookie", c.cookie)
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://weibo.com")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Server-Version", "v2025.10.31.1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-XSRF-TOKEN", "BRrWyvo1P2k853zXiY0UL7xp")
}

// getReader 根据Content-Encoding返回相应的Reader
func (c *Client) getReader(body io.ReadCloser, encoding string) io.Reader {
	if strings.Contains(encoding, "gzip") {
		gzipReader, err := gzip.NewReader(body)
		if err != nil {
			// 如果解压失败，返回原始body
			return body
		}
		return gzipReader
	}
	return body
}
