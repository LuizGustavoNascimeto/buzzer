package middleware

// import (
// 	"backend-go/pkg/auth"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func RequireAuth(validator *auth.Validator) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := c.GetHeader("Authorization")
// 		claims, err := validator.ValidateToken(token)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 			return
// 		}

// 		// injeta no contexto — o handler pega daqui
// 		c.Set("user_id", claims.)
// 		c.Next()
// 	}
// }
