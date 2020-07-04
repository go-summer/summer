package summer

import (
	"net"
	"net/http"
)

const (
	defaultPort = "80"
)

func init() {
	// 初始化 Controllers
	logo()
	App = Application{
		Port: defaultPort,
	}
}

// App 主程序管理器
var App Application

// 主程序初始化
type Application struct {
	Port        string
	Controllers []Controller
}

// 主程序启动
func (a Application) Run() {
	// 创建路由
	re := CreateNewRemux()

	// 注册全部的 handler
	for _, c := range App.Controllers {
		re.SetHandlerMapping(c.Path(), c.Handler)
	}

	listener, err := net.Listen("tcp", ":"+App.Port)
	if err != nil {
		print(err)
		return
	}
	svr := http.Server{Handler: re}

	err = svr.Serve(listener)
	if err != nil {
		print(err)
	}
}

func logo() {
	//http://patorjk.com/software/taag/#p=display&f=Doom&t=Summer
	println(

		" _____                                     \n" +
			"/  ___|                                    \n" +
			"\\ `--. _   _ _ __ ___  _ __ ___   ___ _ __ \n" +
			" `--. \\ | | | '_ ` _ \\| '_ ` _ \\ / _ \\ '__|\n" +
			"/\\__/ / |_| | | | | | | | | | | |  __/ |   \n" +
			"\\____/ \\__,_|_| |_| |_|_| |_| |_|\\___|_|   \n")
}
