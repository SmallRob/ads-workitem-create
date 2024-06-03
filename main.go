package main

import (
	//"adswork/cmd"

	"adswork/bootstrap"
	"adswork/cmd"
	"adswork/common"
	_ "embed"
)

func main() {
	//
	defer func() {
		if err := recover(); err != nil {
			common.JiaLog.Fatal(err)
		}
	}()

	// print version
	//versionInfo := formatVerion()

	//common.JiaLog.Console(versionInfo)
	//common.JiaLog.Console("--------开始执行")

	// command line startup
	//cmd.Execute(versionInfo)
	cmd.Execute()
}

// running before main
func init() {

	// 初始化配置
	bootstrap.InitializeConfig()
	// 初始化日志
	//global.App.Log = bootstrap.InitializeLog()

	//日志初始化
	common.JiaLog.InitLogger("error")
}
