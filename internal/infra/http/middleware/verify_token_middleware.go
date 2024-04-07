package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/silasstoffel/account-service/internal/domain/account"
	"github.com/silasstoffel/account-service/internal/infra/service/token"
)

type VerifyTokenMiddlewareParams struct {
	TokenManagerService         token.TokenService
	AccountPermissionRepository domain.AccountPermissionRepository
}

func NewVerifyTokenMiddleware(
	tokenManagerService token.TokenService,
	accountPermissionRepository domain.AccountPermissionRepository,
) *VerifyTokenMiddlewareParams {
	return &VerifyTokenMiddlewareParams{
		TokenManagerService:         tokenManagerService,
		AccountPermissionRepository: accountPermissionRepository,
	}
}

func (ref *VerifyTokenMiddlewareParams) VerifyTokenMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 8 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	lp := "[verify-token-middleware]"
	token := authHeader[7:]
	data, err := ref.TokenManagerService.VerifyToken(token)
	response := &gin.H{"code": "INVALID_CREDENTIALS", "message": "Invalid credentials"}
	if err != nil {
		log.Println(lp, "failure o verify token", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	items, err := ref.AccountPermissionRepository.FindByAccountId(data.Sub)
	if err != nil {
		log.Println(lp, "failure to get account permissions", err)
		c.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}
	var permissions []string
	for _, p := range items {
		permissions = append(permissions, p.Scope)
	}
	c.Set("accountId", data.Sub)
	c.Set("accountPermissions", permissions)
	c.Next()
}
