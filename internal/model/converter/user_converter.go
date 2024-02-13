package converter

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/manikandareas/go-clean-architecture/internal/entity"
	"github.com/manikandareas/go-clean-architecture/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

func UserToLoginResponse(user *entity.User, expiresIn *jwt.NumericDate, token ...string) *model.LoginUserResponse {
	return &model.LoginUserResponse{
		User: UserToResponse(user),
		BackendTokens: model.BackendTokens{
			AccessToken:  token[0],
			RefreshToken: token[1],
			ExpiresIn:    expiresIn.Unix(),
		},
	}
}

//func ClaimsToBackendTokens(claims *model.BackendTokens) *model.BackendTokens {
//	return &model.BackendTokens{
//		AccessToken:  claims.AccessToken,
//		RefreshToken: claims.RefreshToken,
//		ExpiresIn:    claims.ExpiresIn,
//	}
//}
