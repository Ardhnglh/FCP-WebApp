package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				contentType := ctx.GetHeader("Content-Type")
				if contentType == "application/json" {
					ctx.JSON(http.StatusUnauthorized, err)
					ctx.Abort()
					return
				} else {
					ctx.Redirect(http.StatusSeeOther, "user/login")
					ctx.Abort()
					return
				}
			}
			ctx.JSON(http.StatusBadRequest, err)
			ctx.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(cookie, &model.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, err
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, err)
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusBadRequest, err)
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(*model.Claims)
		if ok && token.Valid {
			ctx.Set("email", claims.Email)
		}

		ctx.Next()
	})
}
