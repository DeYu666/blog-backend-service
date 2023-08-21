package service

import (
	"context"
	"time"

	"github.com/DeYu666/blog-backend-service/model"
	"github.com/dgrijalva/jwt-go"
)

type JwtConfig struct {
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	JwtTtl int64  `mapstructure:"jwt_ttl" json:"jwt_ttl" yaml:"jwt_ttl"` // token 有效期（秒）
}

// JwtUser 所有需要颁发 token 的用户模型必须实现这个接口
type JwtUser interface {
	GetUid() string
}

var JwtService jwtService

type jwtService struct {
	secret string
}

func (j *jwtService) Init(ctx context.Context, secret string) {
	j.secret = secret
}

// 没有设置过期时间
func (j *jwtService) GenerateToken(user JwtUser) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		model.CustomClaims{
			StandardClaims: jwt.StandardClaims{
				Id:        user.GetUid(),
				NotBefore: time.Now().Unix() - 1000,
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(j.secret))

	return tokenStr, err
}

func (j *jwtService) ValidateToken(tokenStr string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
