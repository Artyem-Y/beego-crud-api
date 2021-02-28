// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"beego-crud-api/controllers"
	"beego-crud-api/middlewares"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type",
			"X-ACCESS-TOKEN", "X-ACCESS-PRODUCT"},
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin"},
	}))

	beego.Router("/", &controllers.BaseController{})
	beego.Router("/active/:code([0-9a-zA-Z]+)", &controllers.AuthController{}, "get:ConfirmEmail")

	nsAuth := beego.NewNamespace("/auth",
		beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
		beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		beego.NSRouter("/:id/access_token", &controllers.AuthController{}, "post:CheckAccessToken"),
		beego.NSRouter("/:id/logout", &controllers.AuthController{}, "get:Logout"),
	)

	nsApi := beego.NewNamespace("/api",

		beego.NSBefore(middlewares.Jwt),
		beego.NSRouter("/:id/me", &controllers.UsersController{}, "get:GetCurrent"),

		beego.NSBefore(middlewares.CheckEmailIsValid),
		beego.NSRouter("/:id/email", &controllers.UsersController{}, "put:PutEmail"),
		beego.NSRouter("/:id/valid_email", &controllers.UsersController{}, "get:CheckEmail"),
		beego.NSRouter("/:id/password", &controllers.UsersController{}, "put:PutPassword"),
		beego.NSRouter("/:id/validate_email", &controllers.UsersController{}, "get:ValidateEmail"),
	)

	nsAdmin := beego.NewNamespace("/admin",
		beego.NSBefore(middlewares.CheckIfRoleIsAdmin),
		beego.NSRouter("/users", &controllers.UsersController{}, "get:GetUsers"),
	)

	beego.AddNamespace(nsAuth, nsApi, nsAdmin)
}
