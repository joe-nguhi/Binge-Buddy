package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	USER_ID_CONTEXT_KEY   = "userID"
	USER_ROLE_CONTEXT_KEY = "userRole"
)

func GetUserIdFromContext(c *gin.Context) (string, error) {
	id := c.GetString(USER_ID_CONTEXT_KEY)

	if id == "" {
		return "", errors.New("user ID not found")
	}

	return id, nil
}

func GetUserRoleFromContext(c *gin.Context) (string, error) {
	role := c.GetString(USER_ROLE_CONTEXT_KEY)

	if role == "" {
		return "", errors.New("user role not found")
	}

	return role, nil
}
