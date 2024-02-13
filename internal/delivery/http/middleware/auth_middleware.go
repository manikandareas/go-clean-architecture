package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manikandareas/go-clean-architecture/internal/model"
	"github.com/manikandareas/go-clean-architecture/internal/usecase"
	"github.com/manikandareas/go-clean-architecture/pkg"
	"strings"
)

func NewAuth(userUseCase *usecase.UserUseCase) fiber.Handler {
	/*
		Duty
		Ensure user have valid access token and add information user to auth.local
	*/
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			userUseCase.Log.Warnf("Invalid token")
			return fiber.ErrUnauthorized
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		userUseCase.Log.Debugf("Authorization : %s", tokenString)
		// Decode token to extract information user
		claims, err := userUseCase.JwtService.DecodeJwtToken(tokenString, pkg.ACCESS_TOKEN_KEY)
		if err != nil {
			userUseCase.Log.Warnf("Failed to decode token : %+v", err)
			return fiber.ErrUnauthorized
		}

		user := claims["User"].(map[string]interface{})
		auth := &model.Auth{
			ID:    user["user_id"].(string),
			Email: user["email"].(string),
			Name:  user["name"].(string),
		}
		// search user from db and count, return err if count < 1
		err = userUseCase.Verify(ctx.Context(), auth)
		if err != nil {
			userUseCase.Log.Warnf("Failed find user by id : %+v", err)
			return fiber.ErrUnauthorized
		}
		userUseCase.Log.Debugf("User : %+v", auth.ID)
		// inject auth information to local var
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func NewRefreshToken(userUseCase *usecase.UserUseCase) fiber.Handler {
	/*
		Duty
		Ensure user have valid refresh token and add information user to auth.local
	*/
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Refresh") {
			userUseCase.Log.Warnf("Invalid token")
			return fiber.ErrUnauthorized
		}
		tokenString := strings.Replace(authorizationHeader, "Refresh ", "", -1)
		userUseCase.Log.Debugf("Refresh Token : %s", tokenString)

		// Decode token to extract information user
		claims, err := userUseCase.JwtService.DecodeJwtToken(tokenString, pkg.REFRESH_TOKEN_KEY)
		if err != nil {
			userUseCase.Log.Warnf("Failed to decode token : %+v", err)
			return fiber.ErrUnauthorized
		}
		user := claims["User"].(map[string]interface{})
		auth := &model.Auth{
			ID:    user["user_id"].(string),
			Email: user["email"].(string),
			Name:  user["name"].(string),
		}
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
