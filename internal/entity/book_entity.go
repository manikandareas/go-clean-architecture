package entity

type Book struct {
	ID       string `gorm:"column:id;primaryKey"`
	Title    string `gorm:"column:title"`
	AuthorId string `gorm:"column:author_id"`
}

func (b *Book) TableName() string {
	return "books"
}
