package repository

import (
	"github.com/manikandareas/go-clean-architecture/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{Log: log}
}

func (r *UserRepository) FindByToken(db *gorm.DB, user *entity.User, token string) error {
	return db.Where("token = ?", token).First(user).Error
}

func (r *UserRepository) FindByEmail(db *gorm.DB, email string) (*entity.User, error) {
	user := new(entity.User)
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}
