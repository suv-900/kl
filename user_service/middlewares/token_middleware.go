package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/user_service/utils"
)

func TokenMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.GetHeader("token")

		if len(token) == 0 {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		tokenExpired, userid, tokenInvalid := utils.ValidateToken(token)

		if tokenExpired {
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		}
		if tokenInvalid {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Set("userid", userid)
		c.Next()
	}
}
