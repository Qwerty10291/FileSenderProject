package auth

import (
	"github.com/Qwerty10291/FileSenderProject/internal/user/db"
	"github.com/gin-gonic/gin"
)

// func NewSessionAuthService(userStorage *db.UserStorage) AuthService {
// 	return &SessionAuthService{
// 	}
// }

type SessionAuthService struct {
	userStorage *db.UserStorage
}

// Create implements AuthService
func (*SessionAuthService) Create(ctx *gin.Context) {
	
}

// Delete implements AuthService
func (*SessionAuthService) Delete(ctx *gin.Context) {
	panic("unimplemented")
}

// SignIn implements AuthService
func (*SessionAuthService) SignIn(ctx *gin.Context) {
	panic("unimplemented")
}

// SignOut implements AuthService
func (*SessionAuthService) SignOut(ctx *gin.Context) {
	panic("unimplemented")
}
