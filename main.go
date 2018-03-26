package main

import (
	_ "github.com/biaocheng/englishStudy/routers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"github.com/biaocheng/englishStudy/service"
)

func init(){
	beego.SetLevel(beego.LevelDebug)

	config := make(map[string]interface{})
	// config["filename"] = "e:/golang/go_pro/logs/logcollect.log"
	// config["filename"] = "e://golang//go_pro//logs//logcollect.log"
	config["filename"] = "d:\\logs\\log.txt"
	config["level"] = logs.LevelDebug

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:Ab20171221@tcp(47.52.197.97)/english_study?charset=utf8")
	orm.RunSyncdb("default", true, true)
}

func main() {
	orm.Debug = true

	service.GetWordS()

}

