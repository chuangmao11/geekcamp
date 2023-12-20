package main

import (
	"geekcamp/webook/internal/repository"
	dao2 "geekcamp/webook/internal/repository/dao"
	"geekcamp/webook/internal/service"
	"geekcamp/webook/internal/web"
	"geekcamp/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	server := initWebServer()
	db := initDB()
	u := initUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		// ExposeHeaders: 用于获取需要前端截获的头，如 jwt
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	//store := cookie.NewStore([]byte("secret"))
	/*store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	[]byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))*/
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		[]byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	if err != nil {
		panic(err)
	}
	//store := &sqlx_store.Store{}

	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	dao := dao2.NewUserDAO(db)
	repo := repository.NewUserRepository(dao)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao2.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
