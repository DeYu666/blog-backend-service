package middleware

import (
	"github.com/DeYu666/blog-backend-service/delivery/response"
	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/service"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		log := blog.Extract(c)

		tokenStr := c.Request.Header.Get("Authorization")

		log.Sugar().Debug("JWTAuth", tokenStr)

		token, err := service.JwtService.ValidateToken(tokenStr)

		log.Sugar().Debug("JWTAuth", token, err)

		if token == nil || err != nil {
			response.TokenFail(c)
			c.Abort()
			return
		}

		c.Set("author_id", token.Claims.(*model.CustomClaims).Id)

	}
}
