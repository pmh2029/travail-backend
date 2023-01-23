package utils

import (
	"travail/internal/pkg/domains/models/dtos/res"
	"travail/internal/pkg/domains/models/entities"
)

// convertUserEntityToUserResponse func
func ConvertUserEntityToUserResponse(user entities.User) res.UserResponse {
	return res.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
