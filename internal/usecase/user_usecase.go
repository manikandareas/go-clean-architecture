package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/manikandareas/go-clean-architecture/internal/entity"
	"github.com/manikandareas/go-clean-architecture/internal/model"
	"github.com/manikandareas/go-clean-architecture/internal/model/converter"
	"github.com/manikandareas/go-clean-architecture/internal/repository"
	"github.com/manikandareas/go-clean-architecture/pkg"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	JwtService     *pkg.JwtService
}

func NewUserUseCase(DB *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, jwtService *pkg.JwtService) *UserUseCase {
	return &UserUseCase{DB: DB, Log: log, Validate: validate, UserRepository: userRepository, JwtService: jwtService}
}

// TODO: Refactor Verify to unused token from db

// Verify verifies the user authentication request.
//
// ctx - The context for the request.
// request - The authentication request to be verified.
// Duties - count user by id from database to ensure request from a valid user
// error - An error, if any.

func (c *UserUseCase) Verify(ctx context.Context, request *model.Auth) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return fiber.ErrBadRequest
	}

	count, err := c.UserRepository.CountById(tx, request.ID)
	if err != nil {
		c.Log.Warnf("Failed find by user id : %+v", err)
		return fiber.ErrNotFound
	}

	if !(count > 0) {
		c.Log.Warnf("User not found : %+v", err)
		return fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}
	return nil
}

func (c *UserUseCase) RefreshToken(ctx context.Context, request *model.Auth) (*model.BackendTokens, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	expireTime := jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
	claimsAccessToken := &model.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expireTime,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		User: map[string]string{
			"user_id": request.ID,
			"email":   request.Email,
			"name":    request.Name,
		},
	}
	accessToken, err := c.JwtService.GenerateJwtToken(claimsAccessToken, pkg.ACCESS_TOKEN_KEY)
	if err != nil {
		c.Log.Warnf("Failed to generate jwt token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	claimsRefreshToken := &model.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		User: map[string]string{
			"user_id": request.ID,
			"email":   request.Email,
			"name":    request.Name,
		},
	}
	refreshToken, err := c.JwtService.GenerateJwtToken(claimsRefreshToken, pkg.REFRESH_TOKEN_KEY)
	if err != nil {
		c.Log.Warnf("Failed to generate jwt token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	backendToken := &model.BackendTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expireTime.Unix(),
	}
	return backendToken, nil
}

func (c *UserUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	registeredUser, err := c.UserRepository.FindByEmail(tx, request.Email)
	if registeredUser.Email != "" {
		c.Log.Warnf("User already exists : %+v", err)
		return nil, fiber.ErrConflict
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to hash password : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		ID:       uuid.NewString(),
		Password: string(password),
		Email:    request.Email,
		Name:     request.Name,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed to create user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginUserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user, err := c.UserRepository.FindByEmail(tx, request.Email)
	if err != nil {
		c.Log.Warnf("Failed find by user email : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare user password with bcrypt hash : %+v", err)
		return nil, fiber.ErrUnauthorized
	}
	expireTime := jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
	claimsAccessToken := &model.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expireTime,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		User: map[string]string{
			"user_id": user.ID,
			"email":   user.Email,
			"name":    user.Name,
		},
	}

	accessToken, err := c.JwtService.GenerateJwtToken(claimsAccessToken, pkg.ACCESS_TOKEN_KEY)
	if err != nil {
		c.Log.Warnf("Failed to generate jwt token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	user.Token = accessToken
	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	claimsRefreshToken := &model.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		User: map[string]string{
			"user_id": user.ID,
			"email":   user.Email,
			"name":    user.Name,
		},
	}

	refreshToken, err := c.JwtService.GenerateJwtToken(claimsRefreshToken, pkg.REFRESH_TOKEN_KEY)
	if err != nil {
		c.Log.Warnf("Failed to generate jwt token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToLoginResponse(user, expireTime, accessToken, refreshToken), nil
}
