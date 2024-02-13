package repository

import (
	"github.com/manikandareas/go-clean-architecture/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookRepository struct {
	Repository[entity.Book]
	Log *logrus.Logger
}

func NewBookRepository(log *logrus.Logger) *BookRepository {
	return &BookRepository{
		Log: log,
	}
}

func (r *BookRepository) FindAll(tx *gorm.DB, books *[]entity.Book) error {
	if err := tx.Find(&books).Error; err != nil {
		return err
	}
	return nil
}
