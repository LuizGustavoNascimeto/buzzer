package middleware

import (
	"backend-go/pkg/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth(validator *auth.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		// suporta "Bearer <token>" ou token direto
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, err := validator.ValidateToken(token)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("user_id", claims)
		c.Next()
	}
}
func SoftAuth(validator *auth.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims, err := validator.ValidateToken(token)
		if err == nil {
			c.Set("user_id", claims)
			c.Set("authenticated", true)
		} else {
			c.Set("authenticated", false)
		}
		c.Next() // nunca aborta — deixa o handler decidir
	}
}
