package router

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/sso/global"
	"github.com/beewit/sso/handler"

	"github.com/labstack/echo"

	"fmt"
)

func Start() {
	fmt.Printf("登陆授权系统启动")

	e := echo.New()
	e.Static("/static", "static")
	e.Static("/page", "page")
	e.File("/", "page/login.html")

	e.POST("/pass/login", handler.Login)
	e.POST("/pass/register", handler.Register)
	e.POST("/pass/regSendSms", handler.RegSendSms)
	e.POST("/pass/checkRegMobile", handler.CheckRegMobile)

	e.GET("/img/code", handler.ImgCode)

	utils.Open(global.Host)

	e.Logger.Fatal(e.Start(":8080"))
}
