package main

import (
	_ "beego-crud-api/routers"
	servicesDb "beego-crud-api/services/db"
	"github.com/astaxie/beego"
	_ "github.com/lib/pq"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	servicesDb.InitSql()
	beego.Run()
}
