package user

import (
	"github.com/Qwerty10291/FileSenderProject/internal/config"
	"github.com/Qwerty10291/FileSenderProject/internal/user/auth_service"
	"github.com/gin-gonic/gin"
)

func NewUserService(config config.Config) *UserService {
	return &UserService{}
}

type UserService struct {
	auth.AuthService
}

func (s *UserService) Register(engine *gin.Engine) {
	
}

func (s *UserService) userPage() {

}