package wft

import (
	"mycommon/logs"
	fmt "mycommon/myinherit/myfmt"
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
	response := &model.FlagObjOfWft{}
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err, ok := e.(error)
			if ok {

				response.Flag = 10
				response.Val = err.Error()
				response.Msg = "inner err"

				// 日志记录
				for i := 3; i <= 8; i++ {
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

		if response.Flag == 1100 {
			// 纯文本输出
			sssutils.WriteText(w, fmt.Sprint(response.Val), reqid)
		} else {
			sssutils.WriteJson(w, response, reqid)
		}

		runtime.GC()
	}()

	specialParams := sssutils.GetSpecialGetParam(r.URL.String(), "/api/")
	paramLen := len(specialParams)

	if paramLen <= 0 {
		response.Flag = 3
		response.Val = "invalid url"
		return
	}
	if 1 == paramLen && "buyCallback" == specialParams[0] {
		// 微信支付回调
		var data model.XmlResp
		err := sssutils.ReadXml(r, &data, reqid)
		if nil != err {
			if `参数为空` == err.Error() {
				response.Flag = 2
				response.Val = err.Error()
				return
			}
		}

		resp := new(Weifutong).BuyCallback(&data)
		response.Flag = 1100 // 1100 代表返回text
		response.Val = resp
		return
	}

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

	switch paramLen {

	case 1:
		this.route1(specialParams[0], data, response, w, r)
		break

	default:
		response.Flag = 10
		response.Val = "error router"
	}

}

// 1 route,eg: /api/getMoney
func (this *MainHandler) route1(funcName string, data map[string]interface{}, response *model.FlagObjOfWft, w http.ResponseWriter, r *http.Request) {

	responseJson, errCode := new(Weifutong).Handler(funcName, data, w, r)
	response.Flag = errCode
	response.Val = responseJson
	if -1 == response.Flag {
		response.Msg = "success"
	}
}
