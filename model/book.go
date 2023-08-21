package model

// BooksList 书籍列表
type BooksList struct {
	ID
	BookName   string `json:"book_name"`
	BookStatus string `json:"book_status"`
	Abstract   string `json:"abstract"`
	Timestamps
}

func (BooksList) GetTableName() string {
	return "books_lists"
}

// BookContent 书籍内容
type BookContent struct {
	ID
	BookContent string    `gorm:"type:text" json:"book_content"`
	BookId      int       `json:"book_id"`
	Book        BooksList `gorm:"foreignKey:BookId"  json:"book"`
	Timestamps
}

func (BookContent) GetTableName() string {
	return "book_contents"
}

type ShelfArrWithCount struct {
	ShelfArr []BooksList `json:"shelf_arr"`
	Count    int64       `json:"count"`
}
