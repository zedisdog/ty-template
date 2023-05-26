package account

import (
	"fback/internal/modules/account/handlers"
	"fback/internal/modules/account/migration"
	"fback/internal/modules/account/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zedisdog/ty/application"
	"github.com/zedisdog/ty/database/migrate"
	"github.com/zedisdog/ty/sdk/net/http/middlewares"
)

type Module struct {
	*services.Account
}

func (d Module) Name() string {
	return "account"
}

func (d *Module) Register() error {
	d.Account = services.NewAccount(application.GetInstance())

	authMiddleware := middlewares.NewAuthBuilder().WithUserExistsFunc(d.Account.UserExists).WithOnPass(func(claims jwt.MapClaims, ctx *gin.Context) error {
		if openID, ok := claims["open_id"]; ok {
			ctx.Set("open_id", openID)
		}
		return nil
	}).Build()
	application.RegisterComponent("authMiddleware", authMiddleware)

	application.Migrator().GetSourceInstance().(*migrate.EmbedDriver).Add(&migration.Migration)

	return nil
}

func (d Module) Boot() error {
	application.HttpServer[*gin.Engine]().RegisterRoutes(func(r *gin.Engine) error {
		accountHandler := handlers.NewAccount(d.Account)

		api := r.Group("api")
		api.POST("login-by-code", accountHandler.LoginByMiniCode)
		api.GET("code-login", accountHandler.LoginByCode)
		api.Any("wechat-msg", accountHandler.ReceiveWechatMessage)

		auth := api.Group("", application.Component[gin.HandlerFunc]("authMiddleware"))
		auth.GET("self", accountHandler.Self)

		return nil
	})
	return nil
}
