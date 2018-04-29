package webhandler

import (
	"fmt"
	"mycommon/logs"
	"mycommon/mathstr"
	"net/http"
	"runtime"

	"api.lottery.thirdparty/model"
	"api.lottery.thirdparty/sssutils"
)

type MainHandler struct {
	Reqid string
}

func (this *MainHandler) AppHandler(w http.ResponseWriter, r *http.Request) {

	reqid := this.Reqid
	response := &model.FlagObj{}
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err, ok := e.(error)
			if ok {

				response.Flag = 10
				response.Val = e
				response.Msg = "inner err"

				// 日志记录
				for i := 3; i <= 7; i++ {
					_, f, line, ok := runtime.Caller(i)
					if !ok {
						continue
					}
					if i == 3 {
						logs.Error("__err:[", err, i, "]__fname:[", f, "]__line:[", line, "]")
					} else {
						logs.Error("__fname:[", f, "]__line:[", line, "]")
					}
				}

			}
		}

		sssutils.WriteJson(w, response, reqid)

		runtime.GC()
	}()

	specialParams := sssutils.GetSpecialGetParam(r.URL.String(), "/third/")
	paramLen := len(specialParams)

	var data map[string]interface{}
	// 读取参数
	err := sssutils.ReadJson(r, &data, reqid)
	if nil != err {
		if `参数为空` == err.Error() {
			response.Flag = 2
			response.Val = err.Error()
			return
		}
	}

	fmt.Println("__params:", mathstr.GetJsonPlainStr(specialParams))

	switch paramLen {

	case 1:
		this.route1(specialParams[0], data, response, w, r)
		break

	case 2:
		this.route2(specialParams[0], specialParams[1], data, response, w, r)
		break

	default:
		response.Flag = 10
		response.Val = "error router"
	}

}

// 1 route,eg: /third/oninfo
func (this *MainHandler) route1(funcName string, data map[string]interface{}, response *model.FlagObj, w http.ResponseWriter, r *http.Request) {

	responseJson, errCode := new(Statistics).Handler(funcName, data, w, r)
	response.Flag = errCode
	response.Val = responseJson
	if -1 == response.Flag {
		response.Msg = "success"
	}

}

// 2 route,eg: /third/user/select
func (this *MainHandler) route2(entityName string, funcName string, data map[string]interface{}, response *model.FlagObj, w http.ResponseWriter, r *http.Request) {

	var responseJson interface{}
	var errCode int = 0

	switch entityName {

	case "user":
		responseJson, errCode = new(Users).Handler(funcName, data, w, r)

	case "cut":
		responseJson, errCode = new(Cutfunc).Handler(funcName, data, w, r)

	case "probability":
		responseJson, errCode = new(Probability).Handler(funcName, data, w, r)

	case "customer":
		responseJson, errCode = new(Customer).Handler(funcName, data, w, r)

	case "rank":
		responseJson, errCode = new(RankPrize).Handler(funcName, data, w, r)

	case "record":
		responseJson, errCode = new(Record).Handler(funcName, data, w, r)

	case "syscnf":
		responseJson, errCode = new(SysCnf).Handler(funcName, data, w, r)

	default:
		responseJson = "invalid entityname [ " + entityName + " ]"
		errCode = 777
	}

	response.Flag = errCode
	response.Val = responseJson
	if -1 == response.Flag {
		response.Msg = "success"
	}

}
