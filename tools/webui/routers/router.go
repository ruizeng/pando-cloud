package routers

import (
	"github.com/PandoCloud/pando-cloud/tools/webui/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
