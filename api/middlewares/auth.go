package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	"ticken-validator-service/security/jwt"
	"ticken-validator-service/utils"
)

type AuthMiddleware struct {
	validator   *validator.Validate
	jwtVerifier jwt.Verifier
	apiPrefix   string
}

func NewAuthMiddleware(jwtVerifier jwt.Verifier, apiPrefix string) *AuthMiddleware {
	middleware := new(AuthMiddleware)

	middleware.validator = validator.New()
	middleware.jwtVerifier = jwtVerifier
	middleware.apiPrefix = apiPrefix

	return middleware
}

func (middleware *AuthMiddleware) Setup(router gin.IRouter) {
	router.Use(middleware.isJWTAuthorized())
}

func (middleware *AuthMiddleware) isFreeURI(uri string) bool {
	uri = strings.Replace(uri, middleware.apiPrefix, "", 1)
	return uri == "/healthz"
}

func (middleware *AuthMiddleware) isJWTAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if middleware.isFreeURI(c.Request.URL.Path) {
			return
		}

		rawAccessToken := c.GetHeader("Authorization")

		token, err := middleware.jwtVerifier.Verify(rawAccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.HttpResponse{Message: "authorisation failed while verifying the token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("jwt", token)
	}
}
