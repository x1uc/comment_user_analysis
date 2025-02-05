package crawler

import (
	"strings"
)

// 判断手机品牌映射
var phoneBrands = map[string]string{
	"Huawei":    "华为",
	"华为":        "华为",
	"nova":      "华为",
	"HarmonyOS": "华为",
	"Xiaomi":    "小米",
	"小米":        "小米",
	"OPPO":      "OPPO",
	"Vivo":      "Vivo",
	"iPhone":    "苹果",
	"Samsung":   "三星",
	"三星":        "三星",
	"Meizu":     "魅族",
	"魅族":        "魅族",
	"realme":    "真我",
	"真我":        "真我",
	"redmi":     "红米",
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
	"Android":   "未知Android设备",
}

// GetSource 获取手机品牌
func GetSource(uid string, cookie string) string {
	blogList := GetBlog(uid, cookie, 1)

	// 遍历 blogList 来查找用户信息
	for _, blog := range blogList.Data.List {
		// 保证是本人发的微博
		if blog.User.Uid == uid {
			// 判断 PhoneType 是否包含某个品牌名
			brand := strings.TrimSpace(blog.PhoneType) // 去除空格
			for key, value := range phoneBrands {
				// 如果 PhoneType 包含品牌名称，返回对应的中文名称
				if strings.Contains(strings.ToLower(brand), strings.ToLower(key)) {
					// if key == "Android" {
					// 	file, err := os.OpenFile("out-android-device.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					// 	if err != nil {
					// 		fmt.Println("Error opening file:", err)
					// 	}
					// 	defer file.Close()
					// 	if _, err := file.WriteString(fmt.Sprintf("%s %s\n", uid, brand)); err != nil {
					// 		fmt.Println("Error writing to file:", err)
					// 	}
					// }
					return value
				}
			}
		}
	}

	if len(blogList.Data.List) == 0 {
		return "WARNING: No blog found"
	} else {
		for _, blog := range blogList.Data.List {
			// 保证是本人发的微博
			if blog.User.Uid == uid {
				// 判断 PhoneType 是否包含某个品牌名
				brand := strings.TrimSpace(blog.PhoneType) // 去除空格
				return brand
			}
		}
		return "WARNING: No blog found"
	}
}
