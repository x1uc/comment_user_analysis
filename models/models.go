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

type UserDetail struct {
	Birthday    string `json:"birthday"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
	Gender      string `json:"gender"`
	IPLocation  string `json:"ip_location"`
	RealAuth    bool   `json:"real_auth"`
	DescText    string `json:"desc_text"`
	VerifiedURL string `json:"verified_url"`
	FriendInfo  string `json:"friend_info"`
}

type WeiboComment struct {
	CreatedAt          string         `json:"created_at"`
	ID                 int64          `json:"id"`
	RootID             int64          `json:"rootid"`
	RootIDStr          string         `json:"rootidstr"`
	FloorNumber        int            `json:"floor_number"`
	Text               string         `json:"text"`
	DisableReply       int            `json:"disable_reply"`
	RestrictOperate    int            `json:"restrictOperate"`
	SourceAllowClick   int            `json:"source_allowclick"`
	SourceType         int            `json:"source_type"`
	Source             string         `json:"source"`
	MID                string         `json:"mid"`
	IDStr              string         `json:"idstr"`
	Liked              bool           `json:"liked"`
	PicNum             int            `json:"pic_num"`
	ReadTimeType       string         `json:"readtimetype"`
	AnalysisExtra      string         `json:"analysis_extra"`
	CmtExt             string         `json:"cmt_ext"`
	MatchAIPlayPicture bool           `json:"match_ai_play_picture"`
	RID                string         `json:"rid"`
	AllowFollow        bool           `json:"allow_follow"`
	ItemCategory       string         `json:"item_category"`
	DegradeType        string         `json:"degrade_type"`
	ReportScheme       string         `json:"report_scheme"`
	BizMarkType        int            `json:"biz_mark_type"`
	Comments           []WeiboComment `json:"comments"`
	HasVisible         bool           `json:"hasvisible"`
	PreviousCursor     int            `json:"previous_cursor"`
	NextCursor         int            `json:"next_cursor"`
	TotalNumber        int            `json:"total_number"`
	SinceID            int            `json:"since_id"`
	MaxID              int            `json:"max_id"`
	LikeCounts         int            `json:"like_counts"`
	TextRaw            string         `json:"text_raw"`
	IsExpand           bool           `json:"isExpand"`
	User               WeiboUser      `json:"user"`
}

// WeiboStatus represents the root object
type WeiboBlog struct {
	CreatedAt      string    `json:"created_at"`
	ID             int64     `json:"id"`
	IDStr          string    `json:"idstr"`
	Mid            string    `json:"mid"`
	MblogID        string    `json:"mblogid"`
	User           WeiboUser `json:"user"`
	CanEdit        bool      `json:"can_edit"`
	TextLength     int       `json:"textLength"`
	Source         string    `json:"source"`
	Favorited      bool      `json:"favorited"`
	PicIDs         []string  `json:"pic_ids"`
	PicNum         int       `json:"pic_num"`
	IsPaid         bool      `json:"is_paid"`
	PicBgNew       string    `json:"pic_bg_new"`
	RepostsCount   int       `json:"reposts_count"`
	CommentsCount  int       `json:"comments_count"`
	AttitudesCount int       `json:"attitudes_count"`
	IsLongText     bool      `json:"isLongText"`
	Text           string    `json:"text"`
	TextRaw        string    `json:"text_raw"`
	RegionName     string    `json:"region_name"`

	// Additional nested objects
}
