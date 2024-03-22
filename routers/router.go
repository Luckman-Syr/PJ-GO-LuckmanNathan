package routers

import (
	"project-akhir/controllers"
	"project-akhir/middleware"

	"github.com/gin-gonic/gin"
)

func StartRouter() *gin.Engine {
	r := gin.Default()
	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
		userRouter.PUT("/update/:userId", middleware.Authentication(), controllers.UserUpdate)
		userRouter.DELETE("/delete/:userId", middleware.Authentication(), controllers.UserDelete)
	}

	photoRouter := r.Group("/photos")
	{
		photoRouter.POST("/", middleware.Authentication(), middleware.FileUploadMiddleware(), controllers.PhotoUpload)
		photoRouter.GET("/", middleware.Authentication(), controllers.PhotoGetAll)
		photoRouter.PUT("/:photoId", middleware.Authentication(), middleware.FileUploadMiddleware(), controllers.PhotoUpdate)
		photoRouter.DELETE("/:photoId", middleware.Authentication(), controllers.PhotoDelete)
	}

	commentRouter := r.Group("/comments")
	{
		commentRouter.POST("/", middleware.Authentication(), controllers.CommentCreate)
		commentRouter.GET("/", middleware.Authentication(), controllers.CommentList)
		commentRouter.PUT("/:commentId", middleware.Authentication(), controllers.CommentUpdate)
		commentRouter.DELETE("/:commentId", middleware.Authentication(), controllers.CommentDelete)
	}

	SocialMediaRouter := r.Group("/socialmedias")
	{
		SocialMediaRouter.POST("/", middleware.Authentication(), controllers.SocialMediaCreate)
		SocialMediaRouter.GET("/", middleware.Authentication(), controllers.SocialMediaList)
		SocialMediaRouter.PUT("/:socialMediaId", middleware.Authentication(), controllers.SocialMediaUpdate)
		SocialMediaRouter.DELETE("/:socialMediaId", middleware.Authentication(), controllers.SocialMediaDelete)
	}

	return r
}
