package reprint

type ParamCSDN struct {
	ArticleId         string   `json:"article_id"`        // 文章ID
	AuthorizedStatus  bool     `json:"authorized_status"` // default:"false"
	Categories        string   `json:"categories"`        // 文章分类
	CheckOriginal     bool     `json:"check_original"`    // default:"false"
	Content           string   `json:"content"`           // 文章内容，html 形式的
	CoverImages       []string `json:"cover_images"`      // 封面图片，不知道放什么具体内容，只知道是数组
	CoverType         int      `json:"cover_type"`        // 封面属性： 0：无封面；1：单图；   default:"0"
	CreatorActivityId string   `json:"creator_activity_id"`
	Description       string   `json:"description"` // 摘要
	IsNew             int      `json:"is_new"`
	Level             string   `json:"level"`          //?
	NotAutoSaved      int      `json:"not_auto_saved"` //  default:"1"
	OriginalLink      string   `json:"original_link"`
	ReadType          string   `json:"read_type"`
	Reason            string   `json:"reason"`
	ResourceUrl       string   `json:"resource_url"`
	ScheduledTime     int      `json:"scheduled_time"` // default:"0"
	Source            string   `json:"source"`
	Status            string   `json:"status" ` // 发布形式 ？可能是 default:"0"
	Tags              string   `json:"tags"`    // 文章标签
	Title             string   `json:"title"`   // 文章标题
	Type              string   `json:"type"`
	VoteId            int      `json:"vote_id"` //  default:"0"
}

type ReceiveCSDN struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
