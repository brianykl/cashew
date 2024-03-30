// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"cashew-api/conf"
	"cashew-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
	)
	beego.AddNamespace(ns)

	oauthConfig := conf.LoadOAuthConfig()
	oauthController := controllers.NewOAuthController(*oauthConfig)

	beego.Router("/oauth/redirect", oauthController, "get:Redirect")
	beego.Router("/oauth/callback", oauthController, "get:Callback")
}
