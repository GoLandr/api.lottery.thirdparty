package webhandler

import (
	"fmt"
	"mycommon/encode"
	"mycommon/logs"
	"net/http"
	"os"

	"api.lottery.thirdparty/webhandler/wft"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
)

// 启动网页
func WebRun(addr string) {
	fmt.Println("-----监听：", addr, "-----------")
	m := martini.Classic()
	webRoute(m)

	/* 启动HTTPS服务 */
	err := http.ListenAndServe(addr, m)

	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
}

// 所有路由配置
func webRoute(m *martini.ClassicMartini) {

	mainHandler := MainHandler{Reqid: encode.UUID()}
	wftHandler := wft.MainHandler{Reqid: encode.UUID()}

	m.Any("/third/", mainHandler.AppHandler)
	m.Any("/api/", wftHandler.AppHandler)

	m.Use(gzip.All())
}
