package main

import (
	"comment_phone_analyse/crawler"
	"comment_phone_analyse/picexport"
	"comment_phone_analyse/pojo"
	"flag"
	"fmt"
	"sort"
	"time"
)

var limit int = 1000

var brandDict = map[string]bool{
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

func main() {

	uid := flag.String("uid", "Default", "User ID to be counted")
	cookie := flag.String("cookie", "Defaulta", "Weibo Cookie")
	limit := flag.Int("limit", limit, "Number of counted，default 1000")

	flag.Parse()
	if *uid == "Default" || *cookie == "Default" {
		fmt.Println("请输入待统计用户的uid or cookie ，go run main.go -uid xxxxxxx -cookie \"xxxxxx\"")
		return
	}
	fmt.Println("统计中...")

	unknowDevice := make(map[string]int)
	userSet := make(map[string]bool)
	phoneMap := make(map[string]int)

	// total users
	var cnt = 0
	// blog page
	var page = 1

	for cnt < *limit {
		// get blog
		result := crawler.GetBlog(*uid, *cookie, page)
		if result == nil || len(result.Data.List) == 0 {
			fmt.Println("No more blog")
			break
		}

		// get comment user
		for _, blog := range result.Data.List {

			// filter others blog , example: like, forward
			if blog.User.Uid != *uid {
				continue
			}

			userUrls := crawler.GetCommentUser(*cookie, blog.MblogId, *uid)
			// statistics
			for _, userUrl := range userUrls.Data {

				// filter duplicate user
				if _, ok := userSet[userUrl.User.Idstr]; ok {
					continue
				}
				userSet[userUrl.User.Idstr] = true

				time.Sleep(1 * time.Second)
				// get phone type
				phoneType := crawler.GetSource(userUrl.User.Idstr, *cookie)
				phoneMap[phoneType]++

				if _, exists := brandDict[phoneType]; !exists {
					unknowDevice[phoneType]++
					continue
				}
			}
			// print phone type
			cnt = 0
			for phoneType, num := range phoneMap {
				if _, exists := brandDict[phoneType]; exists {
					cnt = cnt + num
					fmt.Printf("PhoneType: %s, Num: %d\n", phoneType, phoneMap[phoneType])
				}
			}
			fmt.Printf("cnt: %d\n", cnt)
			fmt.Println("=====================================")
			fmt.Println("统计中...")
			if cnt >= *limit {
				break
			}
		}
		time.Sleep(3 * time.Second)
		page++
	}

	// sort by value
	sortByValue := func(m map[string]int) []string {
		type kv struct {
			Key   string
			Value int
		}
		var sortedSlice []kv
		for k, v := range m {
			sortedSlice = append(sortedSlice, kv{k, v})
		}

		sort.Slice(sortedSlice, func(i, j int) bool {
			return sortedSlice[i].Value > sortedSlice[j].Value
		})

		sortedKeys := make([]string, len(sortedSlice))
		for i, kv := range sortedSlice {
			sortedKeys[i] = kv.Key
		}
		return sortedKeys
	}

	unknownKeys := sortByValue(unknowDevice)
	knownKeys := sortByValue(phoneMap)

	fmt.Println("==========================最终统计结果：未知机型============================")
	for _, phoneType := range unknownKeys {
		if _, exists := brandDict[phoneType]; !exists {
			fmt.Printf("PhoneType: %s, Num: %d\n", phoneType, unknowDevice[phoneType])
		}
	}

	fmt.Println("==========================最终统计结果：已知机型============================")
	dataArr := []pojo.StatisticsData{}
	for _, phoneType := range knownKeys {
		if _, exists := brandDict[phoneType]; exists {
			fmt.Printf("PhoneType: %s, Num: %d\n", phoneType, phoneMap[phoneType])
			dataArr = append(dataArr, pojo.StatisticsData{
				PhoneType: phoneType,
				Count:     phoneMap[phoneType],
			})
		}
	}

	picexport.Export(*uid, dataArr)
	picexport.ExportPieChart(*uid, dataArr)
}
