package crawler

import (
	"strings"
)

// 判断手机品牌映射
var phoneBrands = map[string]string{
	"Huawei":    "华为",
	"HUAWEI":    "华为",
	"华为":        "华为",
	"Nova":      "华为",
	"nova":      "华为",
	"HarmonyOS": "华为",
	"Xiaomi":    "小米",
	"小米":        "小米",
	"OPPO":      "OPPO",
	"oppo":      "OPPO",
	"Vivo":      "Vivo",
	"vivo":      "Vivo",
	"iPhone":    "苹果",
	"Samsung":   "三星",
	"三星":        "三星",
	"Meizu":     "魅族",
	"魅族":        "魅族",
	"Realme":    "真我",
	"realme":    "真我",
	"真我":        "真我",
	"Redmi":     "红米",
	"redmi":     "红米",
	"一加":        "一加",
	"OnePlus":   "一加",
	"oneplus":   "一加",
	"荣耀":        "荣耀",
	"Honor":     "荣耀",
	"honor":     "荣耀",
	"ZTE":       "中兴",
	"中兴":        "中兴",
	"Nubia":     "努比亚",
	"努比亚":       "努比亚",
	"iqoo":      "IQOO",
	"IQOO":      "IQOO",
}

// GetSource 获取手机品牌
func GetSource(uid string, cookie string) string {
	blogList := GetBlog(uid, cookie, 1)

	// 遍历 blogList 来查找用户信息
	for _, blog := range blogList.Data.List {
		if blog.User.Uid == uid {
			// 判断 PhoneType 是否包含某个品牌名
			brand := strings.TrimSpace(blog.PhoneType) // 去除空格
			for key, value := range phoneBrands {
				// 如果 PhoneType 包含品牌名称，返回对应的中文名称
				if strings.Contains(strings.ToLower(brand), strings.ToLower(key)) {
					return value
				}
			}
		}
	}
	// 如果没有找到匹配的品牌或列表为空
	return "未知"
}
