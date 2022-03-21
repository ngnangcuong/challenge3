package middleware

import (
	"challenge3/database"
	repo "challenge3/repository"
	"strings"

	"github.com/gin-gonic/gin"
)

func NeedPermission(permit string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role")
		connection := database.GetDatabase()
		roleRepo := repo.NewRoleRepo(connection)

		permission, _ := roleRepo.Find(role.(string))

		// var permission models.Role
		// connection.Where("name = ?", role.(string)).First(&permission)

		if ok := strings.Contains(permission.Permission, permit); !ok {
			c.Set("Permission", false)
			return
			c.Next()
		}
		c.Set("Permission", true)
		c.Next()
	}
}

func NeedRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCheck := c.MustGet("role").(string)
		if role != roleCheck {
			c.Set("Permission", false)
			return
			c.Next()
		}
		c.Set("Permission", true)
		c.Next()
	}
}