package converter

import (
	"github.com/manikandareas/go-clean-architecture/internal/entity"
	"github.com/manikandareas/go-clean-architecture/internal/model"
)

func BooksToResponse(books *[]entity.Book) []model.BookResponse {
	var booksResponse []model.BookResponse
	for _, book := range *books {
		booksResponse = append(booksResponse, *BookToResponse(&book))
	}
	return booksResponse
}

func BookToResponse(book *entity.Book) *model.BookResponse {
	return &model.BookResponse{
		ID:       book.ID,
		Title:    book.Title,
		AuthorId: book.AuthorId,
	}
}
