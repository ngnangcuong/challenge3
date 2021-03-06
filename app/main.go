package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	
	"challenge3/api/user"
	"challenge3/api/post"
	"challenge3/database"
	"challenge3/middleware"
)

func NewOpenAPIMiddleware() gin.HandlerFunc {
	validator := middleware.OpenapiInputValidator("./openapi.yaml")
	return validator
}

func InitRoute(router *gin.Engine) {
	validator := NewOpenAPIMiddleware()
	
	userRoute := router.Group("/user")
	{
		userRoute.Use(validator)
		userRoute.Use(middleware.Authorized())

		userRoute.POST("/login", user.LogIn)
		userRoute.GET("/logout", user.LogOut)
		userRoute.POST("/register", user.Register)
		userRoute.POST("/create-user", middleware.NeedPermission("c"), user.CreateUser)
		userRoute.DELETE("/delete-user/:userEmail", middleware.NeedPermission("d"), user.DeleteUser)
		userRoute.PATCH("/update-user/:userEmail", middleware.NeedPermission("u"), user.UpdateUser)
		userRoute.PUT("/change-role", middleware.NeedRole("admin"), user.ChangeRole)
		userRoute.POST("/new-role", middleware.NeedRole("admin"), user.NewRole)
		userRoute.GET("/", user.GetListUser)
	}

	postRoute := router.Group("/post")
	{
		postRoute.Use(validator)
		postRoute.Use(middleware.Authorized())

		postRoute.POST("/create", post.CreatePost)
		postRoute.DELETE("/delete/:postID", post.DeletePost)
		postRoute.PUT("/update/:postID", post.UpdatePost)
		postRoute.GET("/", post.GetListPost)
	}
}

func InitAPI() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	InitRoute(router)
	router.Run(":3000")
}

func main() {
	database.InitMigration()
	InitAPI()
}

