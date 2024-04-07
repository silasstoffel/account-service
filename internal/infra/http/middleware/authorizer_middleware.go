package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthorizerMiddlewareInstance struct {
	// map[method|path] = scope
	// map[POST|/v1/accounts/] = "scope1,scope2,scopeN"
	Permissions map[string]string
}

func NewAuthorizerMiddleware(permissions map[string]string) *AuthorizerMiddlewareInstance {
	return &AuthorizerMiddlewareInstance{Permissions: permissions}
}

func (ref *AuthorizerMiddlewareInstance) AuthorizerMiddleware(c *gin.Context) {
	scopes := c.GetString("accountPermissions")
	response := gin.H{"code": "INVALID_PERMISSIONS", "message": "Invalid permission"}
	if len(scopes) == 0 || len(ref.Permissions) == 0 {
		log.Println("no permissions", scopes)
		c.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	if !ref.isAllowed(c) {
		c.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}
	c.Next()
}

func (ref *AuthorizerMiddlewareInstance) isAllowed(c *gin.Context) bool {
	m := c.Request.Method
	p := c.FullPath()
	key := m + "|" + p
	p, ok := ref.Permissions[key]
	if !ok {
		return false
	}

	scope := c.GetString("accountPermissions")
	permissions := strings.Split(p, ",")
	scopes := strings.Split(scope, ",")
	for i := range permissions {
		for j := range scopes {
			if scopes[j] == permissions[i] {
				return true
			}
		}
	}
	return false
}
