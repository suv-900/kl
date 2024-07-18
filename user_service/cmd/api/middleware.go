package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/user_service/internal/data"
	"github.com/suv-900/kl/user_service/internal/utils"
)

// adds user with the request
func (app *application) authenticator(next http.HandlerFunc) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get("Authorization")

		if len(authorizationHeader) == 0 {
			app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		//if my header is incomplete should i response with bad request or
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			//unknown request send 401 response
			app.invalidAuthenticationHeader(w)
		}

		//expired malformed token or empty string
		userid, err := VerifyToken(headerParts[1])
		if err != nil {
			app.invalidToken(w)	
		}

		user,err:=data.
	}

	return fn
}

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
