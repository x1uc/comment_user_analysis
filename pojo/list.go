package pojo

type BlogList struct {
	PhoneType string `json:"source"`
	MblogId   string `json:"mblogid"`
	User      User   `json:"user"`
	Id        string `json:"idstr"`
}

type User struct {
	Uid string `json:"idstr"`
}

type CommentData struct {
	Data   []CommentDataList `json:"data"`
	Max_id uint64            `json:"max_id"`
}

type CommentDataList struct {
	User CommentUser `json:"user"`
}

type CommentUser struct {
	Idstr       string `json:"idstr"`
	Profile_url string `json:"profile_url"`
}

type StatisticsData struct {
	PhoneType string
	Count     int
}
