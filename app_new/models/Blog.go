package models

// BlogGeneralCategories 是博客分类总类
type BlogGeneralCategories struct {
	ID
	Name           string           `json:"name"`
	BlogCategories []BlogCategories `gorm:"foreignKey:GeneralID" json:"cub_cate" json:"-"`
}

func (BlogGeneralCategories) GetTableName() string {
	return "blog_general_categories"
}

// BlogCategories 是博客文章分类
type BlogCategories struct {
	ID
	Name      string                `json:"name"`
	GeneralID int                   `json:"general_id" json:"-"`
	General   BlogGeneralCategories `gorm:"foreignKey:GeneralID"  json:"general"`
}

func (BlogCategories) GetTableName() string {
	return "blog_categories"
}

// BlogTag 是博客标签类
type BlogTag struct {
	ID
	Name string     `json:"name"`
	Post []BlogPost `gorm:"many2many:blog_post_tags" json:"-"`
}

func (BlogTag) GetTableName() string {
	return "blog_tags"
}

// BlogPost 是博客信息
type BlogPost struct {
	ID
	Timestamps
	Title      string         `json:"title"`
	Body       string         `gorm:"type:text" json:"content" json:"-"`
	Excerpt    string         `json:"abstract"`
	CategoryID int            `json:"category_id" json:"-"`
	Category   BlogCategories `gorm:"foreignKey:CategoryID" json:"category"`
	Tag        []BlogTag      `gorm:"many2many:blog_post_tags" json:"tags"`
	AuthorID   int            `json:"author_id"`
	Author     AuthUser       `gorm:"foreignKey:AuthorID" json:"author"`
	Views      int            `json:"views"`
	Likes      int            `json:"likes"`
	CoverURL   string         `json:"cover_url"`
	TitleURL   string         `json:"img_src"`
	IsOpen     bool           `json:"is_open"`
}

func (BlogPost) GetTableName() string {
	return "blog_posts"
}

// ChickenSoup 鸡汤
type ChickenSoup struct {
	ID
	Sentence string `json:"sentence"`
}

func (ChickenSoup) GetTableName() string {
	return "chicken_soups"
}

//type BlogLinkInfo struct {
//	ID
//	Name string `json:"name"`
//	URL  string `json:"url"`
//}

type BlogPostPs struct {
	ID
	Password string `json:"password"`
}

func (BlogPostPs) GetTableName() string {
	return "blog_post_ps"
}
