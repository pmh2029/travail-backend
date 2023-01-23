package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	jwt "travail/pkg/shared/auth"
)

func CheckAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, "missing")
			c.Abort()
			return
		}
		tokenReq := strings.Replace(authorization, "Bearer ", "", -1)
		tokenReqParts := strings.Split(tokenReq, ".")
		if len(tokenReqParts) != 3 {
			c.JSON(http.StatusUnauthorized,
				gin.H{"Message": "Token is not JWT"})
			c.Abort()
			return
		}
		if !jwt.VerifyJWT(tokenReq) {
			c.JSON(http.StatusUnauthorized,
				gin.H{"Message": "Token is invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}
