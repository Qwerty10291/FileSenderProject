package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/Qwerty10291/FileSenderProject/internal/user/db"
)

type AuthService interface{
	GetUser(*gin.Engine) *db.User
	AuthRequiredMiddleware() *gin.HandlerFunc
}

