package main

import (
	mysql "e-course-management/pkg/db/mysql"
	"github.com/gin-gonic/gin"

	forgotPassword "e-course-management/internal/forgot_password/injector"
	oauth "e-course-management/internal/oauth/injector"
	register "e-course-management/internal/register/injector"
)

func main() {
	r := gin.Default()
	db := mysql.DB()

	forgotPassword.InitializedService(db).Route(&r.RouterGroup)
	oauth.InitializedService(db).Route(&r.RouterGroup)
	register.InitializedService(db).Route(&r.RouterGroup)

	r.Run()
}
