package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["mailadm/controllers:MailboxController"] = append(beego.GlobalControllerRouter["mailadm/controllers:MailboxController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			Params:           nil})

	beego.GlobalControllerRouter["mailadm/controllers:MailboxController"] = append(beego.GlobalControllerRouter["mailadm/controllers:MailboxController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["mailadm/controllers:DomainController"] = append(beego.GlobalControllerRouter["mailadm/controllers:DomainController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})
}
