package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/user_service/internal/data"
	"github.com/suv-900/kl/user_service/internal/utils"
	"gorm.io/gorm"
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

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			//invalid request header send 401 response
			app.invalidAuthenticationHeader(w)
			return
		}

		//expired malformed token or empty string
		userid, err := VerifyToken(headerParts[1])
		if err != nil {
			app.invalidToken(w)
			return
		}

		//not accounted errors yet
		user, err := Models.Users.GetUser(userid)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				app.invalidTokenDeletedUser(w)
			default:
				app.internalServerError(w)
			}
		}

		app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
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
