package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/manikandareas/go-clean-architecture/internal/entity"
	"github.com/manikandareas/go-clean-architecture/internal/model"
	"github.com/manikandareas/go-clean-architecture/internal/model/converter"
	"github.com/manikandareas/go-clean-architecture/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	BookRepository *repository.BookRepository
}

func NewBookUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, bookRepository *repository.BookRepository) *BookUseCase {
	return &BookUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		BookRepository: bookRepository,
	}
}

func (c *BookUseCase) List(ctx context.Context) ([]model.BookResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var err error
	books := new([]entity.Book)

	if err = c.BookRepository.FindAll(tx, books); err != nil {
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.BooksToResponse(books), nil
}

func (c *BookUseCase) Create(ctx context.Context, request *model.BookRequest) (*model.BookResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return nil, fiber.ErrBadRequest
	}

	book := &entity.Book{
		ID:       uuid.NewString(),
		Title:    request.Title,
		AuthorId: request.AuthorId,
	}
	if err := c.BookRepository.Create(tx, book); err != nil {
		c.Log.WithError(err).Error("failed to create book")
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}
	return converter.BookToResponse(book), nil
}
