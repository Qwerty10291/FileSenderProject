package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Qwerty10291/FileSenderProject/internal/user/db"
	"github.com/Qwerty10291/FileSenderProject/package/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type jwtSignInResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	TTL          int    `json:"ttl"`
}

type JwtAuthService struct {
	userStorage *db.UserStorage
	jwtStorage  *db.JwtStorage
	JwtAuthConfig
}

type JwtAuthConfig struct {
	SecretKey          string
	RefreshTokenLength int
	TTL                time.Duration
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type errorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func errorUnknown() errorResponse {
	return errorResponse{
		Status: false,
		Error:  "unknown server error",
	}
}

func errorTemplate(err string, args ...interface{}) errorResponse {
	return errorResponse{
		Status: false,
		Error:  fmt.Sprintf(err, args...),
	}
}

func errorNo() errorResponse {
	return errorResponse{
		Status: true,
		Error:  "",
	}
}

// func NewJWTAuthService(userStorage *db.UserStorage, jwtStorage *db.JwtStorage, config JwtAuthConfig) AuthService {
// 	return &JwtAuthService{
// 		userStorage:   userStorage,
// 		jwtStorage:    jwtStorage,
// 		JwtAuthConfig: config,
// 	}
// }

func (s *JwtAuthService) Create(ctx *gin.Context) {
	credentials := new(Credentials)
	err := ctx.BindJSON(credentials)
	if err != nil {
		log.Printf("bad json signature in jwt create user: %s", err)
		ctx.JSON(http.StatusBadRequest, errorTemplate("bad json format"))
		return
	}

	user, err := s.userStorage.GetByLogin(credentials.Login)
	if err != nil {
		log.Printf("user get by login failed: %s\n", err)
		ctx.JSON(http.StatusInternalServerError, errorUnknown())
		return
	}
	if user != nil {
		ctx.JSON(http.StatusForbidden, errorTemplate("user with login %s already exists", credentials.Login))
		return
	}
	hashedPassword, err := utils.HashPassword(credentials.Password)
	if err != nil {
		log.Printf("failed to hash password: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorUnknown())
		return
	}

	newUser, err := s.userStorage.Create(credentials.Login, hashedPassword)
	if err != nil {
		log.Printf("user creation error: %s\n", err)
		ctx.JSON(http.StatusInternalServerError, errorUnknown())
		return
	}

	token := db.Jwt{
		UserId:       newUser.Id,
		Token:        s.generateToken(newUser.Login),
		RefreshToken: utils.RandomString(s.RefreshTokenLength),
		Expires:      time.Now().Add(s.TTL),
	}

	err = s.jwtStorage.Create(token)
	if err != nil {
		log.Printf("failed to save jwt token: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorUnknown())
		return
	}

	ctx.JSON(http.StatusCreated, jwtSignInResponse{
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		TTL:          int(s.TTL.Seconds()),
	})
}

func (s *JwtAuthService) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("bad id")
		ctx.JSON(http.StatusBadRequest, errorTemplate("user id must be int"))
		return
	}
	ok, err := s.userStorage.Delete(id)
	if err != nil {
		log.Printf("unknown error when deleting user: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorUnknown())
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, errorTemplate("user with id %d does not exists", id))
		return
	}
	ctx.JSON(http.StatusOK, errorNo())
}

func (s *JwtAuthService) SignIn(ctx *gin.Context) {
	credentials := new(Credentials)
	err := ctx.BindJSON(credentials)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorTemplate("unknown json structure for sign in"))
		return
	}
	user, err := s.userStorage.GetByLogin(credentials.Login)
	if err != nil {
		log.Printf("unknown error when getting user by login: %s\n", credentials.Login)
		ctx.JSON(http.StatusInternalServerError, errorUnknown())
		return
	}
	if user == nil {
		ctx.JSON(http.StatusForbidden, errorTemplate("wrong username or password"))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		ctx.JSON(http.StatusForbidden, errorTemplate("wrong username or password"))
		return
	}
	jwtData, err := s.jwtStorage.FromUser(user)
	if err != nil {
		log.Printf("unknown error when getting jwt data from user: %s", err)
		ctx.JSON(http.StatusInternalServerError, errorUnknown())
		return
	}

	ctx.JSON(http.StatusOK, jwtSignInResponse{
		Token:        jwtData.Token,
		RefreshToken: jwtData.RefreshToken,
		TTL:          int(s.TTL.Seconds()),
	})
}

func (s *JwtAuthService) SignOut(ctx *gin.Context) {

}

func (service *JwtAuthService) generateToken(login string) string {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(service.TTL).Unix(),
		Subject:   login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(service.SecretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *JwtAuthService) validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token: %s", token.Header["alg"])

		}
		return []byte(service.SecretKey), nil
	})
}
