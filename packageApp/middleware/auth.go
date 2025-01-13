package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// IsAuthenticated middleware checks if a user is logged in
func IsAuthenticated(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}
	c.Next()
}
