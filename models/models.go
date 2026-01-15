package models

type WeiboUser struct {
	ID                int64  `json:"id"`
	IDStr             string `json:"idstr"`
	PcNew             int    `json:"pc_new"`
	ScreenName        string `json:"screen_name"`
	ProfileImageURL   string `json:"profile_image_url"`
	ProfileURL        string `json:"profile_url"`
	Verified          bool   `json:"verified"`
	VerifiedType      int    `json:"verified_type"`
	Domain            string `json:"domain"`
	Weihao            string `json:"weihao"`
	VerifiedTypeExt   int    `json:"verified_type_ext"`
	AvatarLarge       string `json:"avatar_large"`
	AvatarHD          string `json:"avatar_hd"`
	FollowMe          bool   `json:"follow_me"`
	Following         bool   `json:"following"`
	MbRank            int    `json:"mbrank"`
	MbType            int    `json:"mbtype"`
	VPlus             int    `json:"v_plus"`
	UserAbility       int    `json:"user_ability"`
	PlanetVideo       bool   `json:"planet_video"`
	VerifiedReason    string `json:"verified_reason"`
	Description       string `json:"description"`
	Location          string `json:"location"`
	Gender            string `json:"gender"`
	FollowersCount    int    `json:"followers_count"`
	FollowersCountStr string `json:"followers_count_str"`
	FriendsCount      int    `json:"friends_count"`
	StatusesCount     int    `json:"statuses_count"`
	URL               string `json:"url"`
	Svip              int    `json:"svip"`
	Vvip              int    `json:"vvip"`
	CoverImagePhone   string `json:"cover_image_phone"`
	TopUser           int    `json:"top_user"`
	UserType          int    `json:"user_type"`
	IsStar            string `json:"is_star"`
	IsMuteUser        bool   `json:"is_muteuser"`
	SpecialFollow     bool   `json:"special_follow"`
}

type UserDeteil struct {
	UID         string
	CreateAt    string `json:"created_at"`
	DescText    string `json:"desc_text"`
	Gender      string `json:"gender"`
	IPLocation  string `json:"ip_location"`
	Description string `json:"description"`
}
