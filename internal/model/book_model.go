package model

type BookResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	AuthorId string `json:"author_id"`
}

type BookRequest struct {
	Title    string `json:"title"  validate:"required"`
	AuthorId string `json:"author_id"  validate:"required"`
}
