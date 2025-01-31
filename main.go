package main

import (
	"comment_phone_analyse/crawler"
	"fmt"
	"time"
)

var uid string = "6290114447"
var cookie string = ""
var limit int = 1000

func main() {
	userSet := make(map[string]bool)
	phoneMap := make(map[string]int)
	var cnt = 0
	var page = 1
	for cnt < limit {
		// get blog
		result := crawler.GetBlog(uid, cookie, page)
		if result == nil || len(result.Data.List) == 0 {
			fmt.Println("No more blog")
			break
		}

		// get comment user
		for _, blog := range result.Data.List {
			if blog.User.Uid != uid {
				continue
			}

			userUrls := crawler.GetCommentUser(cookie, blog.MblogId, uid)

			for _, userUrl := range userUrls.Data {
				// filter duplicate user
				if _, ok := userSet[userUrl.User.Idstr]; ok {
					continue
				}
				userSet[userUrl.User.Idstr] = true

				time.Sleep(1 * time.Second)
				// get phone type
				phoneType := crawler.GetSource(userUrl.User.Idstr, cookie)
				phoneMap[phoneType]++
			}
			cnt = 0
			for phoneType, num := range phoneMap {
				if phoneType == "未知" {
					continue
				}
				cnt = cnt + num
				fmt.Printf("PhoneType: %s, Num: %d\n", phoneType, num)
			}
			fmt.Println("=====================================")
			fmt.Printf("cnt: %d\n", cnt)
			if cnt >= limit {
				break
			}
		}
		time.Sleep(3 * time.Second)
		page++
	}
}
