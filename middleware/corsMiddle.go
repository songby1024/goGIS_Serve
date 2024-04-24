package middleware

import "github.com/gin-gonic/gin"
import "github.com/gin-contrib/cors"

func CorsMiddle() gin.HandlerFunc {
	c := cors.DefaultConfig()
	c.AllowMethods = []string{"POST", "GET", "DELETE", "OPTIONS"}
	c.AllowAllOrigins = true
	return cors.New(c)
}
