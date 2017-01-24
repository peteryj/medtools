package routers

import (
	"github.com/peteryj/medtools/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
