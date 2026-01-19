package client

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	httpClient  *http.Client
	cookie      string
	rateLimit   time.Duration
	lastReqTime time.Time
	mu          sync.Mutex
}

func NewClient(cookie string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cookie: cookie,
	}
}

func (c *Client) SetRateLimit(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rateLimit = interval
}

func (c *Client) Get(url string) ([]byte, error) {
	c.mu.Lock()
	if c.rateLimit > 0 {
		elapsed := time.Since(c.lastReqTime)
		if elapsed < c.rateLimit {
			time.Sleep(c.rateLimit - elapsed)
		}
		c.lastReqTime = time.Now()
	}
	c.mu.Unlock()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code error: %d", resp.StatusCode)
	}

	// handle compressed response
	reader := c.getReader(resp.Body, resp.Header.Get("Content-Encoding"))

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

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

func (c *Client) getReader(body io.ReadCloser, encoding string) io.Reader {
	if strings.Contains(encoding, "gzip") {
		gzipReader, err := gzip.NewReader(body)
		if err != nil {
			return body
		}
		return gzipReader
	}
	return body
}
