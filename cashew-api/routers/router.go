// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"cashew-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/v1/user", &controllers.UserController{}, "post:Create")
	beego.Router("/v1/user/login", &controllers.UserController{}, "post:Login")
	beego.Router("/v1/protected/plaid/link/create", &controllers.PlaidController{}, "get:GenerateLinkToken")
	beego.Router("/v1/plaid/link/create", &controllers.PlaidController{}, "get:GenerateLinkToken")
}
