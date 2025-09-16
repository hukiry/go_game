package main

import (
	"fmt"
	"game/server"
	"game/server/logic"
	"game/util"
	"os"
	"strings"
	"time"
)

var tcpServer *server.TcpServer

var TIPS_HELP string = `Help:
	-> start 启动服务器
	-> stop  关闭服务器
	-> exit  退出程序
	-> ...
`

func main() {
	//EditorProtoUtil.ExportProtocolBinary()
	//return
	fmt.Println(util.GetHUKIRY_TEXT())
	var name string
	for {
		util.LogInput(" ")
		fmt.Scan(&name)
		switch strings.ToLower(name) {
		case "start":
			EnableServer()
			break
		case "stop":
			if tcpServer != nil {
				tcpServer.Stop()
				tcpServer = nil
				util.Log("服务器停止...")
			} else {
				util.LogError("服务器未启动， 请输入 start 创建服务器...")
				time.Sleep(time.Second / 5)
			}
			break
		case "exit":
			util.Log("退出程序...")
			os.Exit(0)
		default:
			fmt.Print(TIPS_HELP)
			util.LogError("输入不正确...")
			time.Sleep(time.Second / 5)
			break
		}
	}
}

func EnableServer() {
	logic.InitMessage()
	tcpServer = server.NewServer()
	err := tcpServer.Start("0.0.0.0:9301")
	if err != nil {
		util.LogError("启动服务器失败: %v\n", err)
	}
}

////time.Sleep(time.Second / 5)
//func test() {
//	err := fmt.Errorf("error occurred at: %v", time.Now())
//	defer func() {
//		//延时调用函数，恢复到正常
//		if err := recover(); err != nil {
//			util.Log("panic occurred:", err)
//		}
//	}()
//	panic(err) //抛出异常，并中断后面语句
//}
