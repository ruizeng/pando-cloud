package main

import (
	_ "github.com/PandoCloud/pando-cloud/tools/webui/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/static","static")
	beego.DirectoryIndex=true
	beego.Run()
}

