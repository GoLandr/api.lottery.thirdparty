// bootstrap
package main

import (
	"log"
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"os"

	"net/http"

	"api.lottery.thirdparty/config"
	"api.lottery.thirdparty/handler/base"
	"api.lottery.thirdparty/webhandler"
	//	_ "net/http/pprof"
)

const current_version = "0.0.00"

func main() {

	//	if *cpuprofile != "" {
	//		f, err := os.Create(*cpuprofile)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		pprof.StartCPUProfile(f)
	//		defer pprof.StopCPUProfile()
	//	}

	//	go ppmem()

	args := os.Args
	argsLen := len(args)
	if argsLen <= 1 {
		webhandler.WebRun(fmt.Sprint(":", config.GetListen()))
		return
	}

	mathstr.RouteArgs(args, current_version)

}

func ppmem() {
	log.Println(123456789)
	log.Println(http.ListenAndServe(":6111", nil))
}
func init() {
	log.Println("spider_init_____")
	spider := new(base.Spider)
	spider.SpiderCron()
}
