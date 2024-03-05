package routers

import (
	"Dentify-X/app/handlers"
	"Dentify-X/app/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
	"gorm.io/gorm"
)

func Rout(db *gorm.DB) *gin.Engine {
	middleware.SaveLogs()
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), gindump.Dump())
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/psignup", func(c *gin.Context) {
		handlers.PsignupHandler(db, c)
	})

	r.POST("/dsignupreq", func(c *gin.Context) {
		handlers.DsignupHandler(db, c)
	})

	r.POST("/plogin", func(c *gin.Context) {
		handlers.Ploginhandler(db, c)
	})

	return r
}
