package main

import (
	"time"

	"github.com/Jolek/be-dashboard/config"
	"github.com/Jolek/be-dashboard/controller"
	"github.com/Jolek/be-dashboard/middleware"
	"github.com/Jolek/be-dashboard/repository"
	"github.com/Jolek/be-dashboard/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB                  = config.SetupDatabaseConnection()
	repo    repository.UserRepository = repository.NewUserRepository(db)
	tokenUC usecase.Token             = usecase.NewTokenUc()
	authUC  usecase.Auth              = usecase.NewAuthUC(repo)
	ac      controller.AuthController = controller.NewAuthController(authUC, tokenUC)
)

func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/register", ac.Register)
		authRoutes.POST("/login", ac.Login)
		authRoutes.POST("/logout", ac.Logout, middleware.AuthorizeJWT(tokenUC))
	}

	r.Run()
}
