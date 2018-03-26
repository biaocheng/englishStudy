package routers

import (
	"github.com/biaocheng/englishStudy/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
