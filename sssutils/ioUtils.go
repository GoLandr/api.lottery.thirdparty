package sssutils

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/sidbusy/weixinmp"

	"encoding/json"
	"mycommon/logs"
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"
	"net/http"
	"strings"

	"github.com/go-martini/martini"
)

type NormalResponse struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

/*
 * 过滤器
 */
func CurFilter(c martini.Context, w http.ResponseWriter, r *http.Request) {
	c.Next()
}

/*
 * 输出错误
 */
func WriteError(w http.ResponseWriter, code int, msg ...interface{}) {
	// 设置header
	w.Header().Set("Content-Type", "application/json")
	// 设置返回
	ret := &NormalResponse{}
	ret.Ret = code
	ret.Msg = fmt.Sprint(msg...)

	str, err := json.Marshal(ret)
	if nil != err {
		logs.Info(err)
		w.Write([]byte("{\"ret\": -1, \"msg\":\"internal error\"}"))
	} else {
		w.Write(str)
	}
}

/*
 * 读取json
 */
func ReadParam(rBody io.ReadCloser, realUrl string) (weixinmp.WxRequestMsg, error) {
	if nil == rBody {
		return weixinmp.WxRequestMsg{}, nil
	}
	// 读取数据区
	body, err := ioutil.ReadAll(rBody)
	if nil != err {
		logs.Error(err)
		return weixinmp.WxRequestMsg{}, err
	}

	var errShow error
	var bodyStr string

	if len(body) > 0 {
		bodyStr, errShow = url.QueryUnescape(string(body))
		if nil != errShow {
			logs.Error(errShow)
		}
	}

	bodyStr = strings.TrimSpace(bodyStr)
	logs.Info("---wechat bodystr: ", bodyStr)
	if "" == bodyStr {
		logs.Info("空的微信post body")
		return weixinmp.WxRequestMsg{}, nil
	}

	var wxReqMsg weixinmp.WxRequestMsg
	var wxHeader weixinmp.MsgHeader
	err = xml.Unmarshal([]byte(bodyStr), &wxHeader)
	if nil != err {
		logs.Error(err)
		return weixinmp.WxRequestMsg{}, err
	}

	wxReqMsg.WxHeader = wxHeader

	switch wxHeader.MsgType {
	case weixinmp.MsgTypeEvent:
		wxReqMsg.IsEvent = true
		xml.Unmarshal([]byte(bodyStr), &wxReqMsg.EventMsg)
		break
	case weixinmp.MsgTypeImage:
		wxReqMsg.IsImage = true
		xml.Unmarshal([]byte(bodyStr), &wxReqMsg.ImageMsg)
		break
	case weixinmp.MsgTypeLink:
		logs.Info(" MsgTypeLink type")
		break
	case weixinmp.MsgTypeLocation:
		logs.Info(" MsgTypeLocation type")
		break
	case weixinmp.MsgTypeText:
		wxReqMsg.IsText = true
		xml.Unmarshal([]byte(bodyStr), &wxReqMsg.TextMsg)
		break
	case weixinmp.MsgTypeVideo:
		wxReqMsg.IsVedio = true
		xml.Unmarshal([]byte(bodyStr), &wxReqMsg.VideoMsg)
		break
	case weixinmp.MsgTypeVoice:
		wxReqMsg.IsVoice = true
		xml.Unmarshal([]byte(bodyStr), &wxReqMsg.VoiceMsg)
		break
	default:
		err := errors.New("错误的微信请求类型：" + wxHeader.MsgType)
		panic(err)
	}

	return wxReqMsg, err

}

/*
 * 读取json
 */
func ReadJson(r *http.Request, obj interface{}, uniquechar string) error {

	reqid := mathstr.RandChar(10)

	body, err := ioutil.ReadAll(r.Body)
	// logs.Debug(reqid, "___err:", err)
	// logs.Debug(reqid, "___body:", err)
	if nil != err {
		return err
	}
	realUrl := r.URL.RawQuery
	var errShow error
	var bodyStr string

	// logs.Debug(reqid, `__realurl:`, realUrl)
	// logs.Debug(reqid, `_____body:`, string(body))

	if "" != realUrl {
		bodyStr, errShow = url.QueryUnescape(realUrl)
	} else {
		bodyStr, errShow = url.QueryUnescape(string(body))
		// logs.Debug(reqid, `_____bodyStr:`, bodyStr)
	}
	if nil != errShow {

		hasErrChar := false
		for _, ch := range body {

			switch {
			case ch > '~':
				hasErrChar = true
				break
			case ch == '\r':
			case ch == '\n':
			case ch == '\t':
			case ch < ' ':
				hasErrChar = true
				break
			case ch == '%':
				hasErrChar = true
				break
			}
		}

		bodyStr, errShow = url.QueryUnescape(string(body))
		if nil != errShow || hasErrChar {
			utils.ThrowErrorStr(mathstr.S_SFT(`含非法字符[{0}]`, errShow.Error()))
		}
	}

	// 新方式--beego
	//	bodyStr := getJsonParam(r.Form)

	logs.Info(reqid, "<<<url[", r.URL, "]<<<content[", uniquechar, "]<<<<<:", bodyStr)

	if "" == bodyStr {
		//		err := errors.New("参数为空")
		//		return err
		return nil
	}

	return json.Unmarshal([]byte(bodyStr), obj)
}

/*
 * 读取form
 */
func ReadForm(r *http.Request, uniquechar string) (map[string]string, error) {

	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		return nil, err
	}
	realUrl := r.URL.RawQuery
	var errShow error
	var bodyStr string
	if "" != realUrl {
		bodyStr, errShow = url.QueryUnescape(string(realUrl))
	} else {
		bodyStr, errShow = url.QueryUnescape(string(body))
	}
	if nil != errShow {
		utils.ThrowError(errShow)
	}

	// 新方式--beego
	//	bodyStr := getJsonParam(r.Form)

	logs.Info("<<<url[", r.URL, "]<<<readform[", uniquechar, "]<<<<<:", bodyStr)

	if "" == bodyStr {
		err := errors.New("参数为空")
		return nil, err
	}

	kvsr := strings.Split(bodyStr, "&")
	resmap := make(map[string]string)
	for _, v := range kvsr {
		str12 := strings.Split(v, "=")
		if len(str12) < 2 {
			continue
		}
		if "" == str12[0] {
			continue
		}

		resmap[str12[0]] = str12[1]
	}

	return resmap, nil
}

/*
 * 读取xml
 */
func ReadXml(r *http.Request, obj interface{}, uniquechar string) error {

	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		return err
	}
	realUrl := r.URL.RawQuery
	var errShow error
	var bodyStr string
	if "" != realUrl {
		bodyStr, errShow = url.QueryUnescape(string(realUrl))
	} else {
		bodyStr, errShow = url.QueryUnescape(string(body))
	}
	if nil != errShow {
		utils.ThrowError(errShow)
	}

	// 新方式--beego
	//	bodyStr := getJsonParam(r.Form)

	logs.Info("<<<url[", r.URL, "]<<<readxml[", uniquechar, "]<<<<<:", bodyStr)

	if "" == bodyStr {
		err := errors.New("参数为空")
		return err
	}

	return xml.Unmarshal([]byte(bodyStr), obj)
}

// 获取get /型的参数
func GetSpecialGetParam(urlPath string, patternStr string) []string {
	// 以 "/"为分隔符，进行切割
	str2 := strings.Replace(urlPath, patternStr, "", -1)
	params := strings.Split(str2, "/")
	var response []string
	for _, v := range params {
		v = strings.TrimSpace(v)
		if "" == v {
			continue
		}
		response = append(response, v)
	}
	return response
}

/*
 * 输出json
 */
func WriteJson(w http.ResponseWriter, obj interface{}, uniquechar string) {
	// 设置header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Authorization,Token,Back-Token")

	str := GetJsonStr(obj)
	logs.Info(">>>>>>END[", uniquechar, "]>>>>>>:", str)

	w.Write([]byte(str))
}

/*
 * 输出json
 */
func WriteJsonStr(w http.ResponseWriter, str string, uniquechar string) {
	// 设置header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Authorization,Token,Back-Token")

	logs.Info(">>>>>>END[", uniquechar, "]>>>>>>:", str)

	w.Write([]byte(str))
}

/*
 * 输出string
 */
func WriteText(w http.ResponseWriter, obj string, uniquechar string) {
	// 设置header
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Authorization,Token,Back-Token")

	str := GetJsonStr(obj)
	logs.Info(">>>>>>END[", uniquechar, "]>>>>>>:", str)

	w.Write([]byte(obj))
}

/*
 * 输出html
 */
func WriteHtml(w http.ResponseWriter, obj string, uniquechar string) {
	// 设置header
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Authorization,Token,Back-Token")

	str := GetJsonStr(obj)
	logs.Info(">>>>>>END[", uniquechar, "]>>>>>>:", str)

	w.Write([]byte(obj))
}

/*
 * 输出json
 */
func WriteXml(w http.ResponseWriter, obj interface{}, uniquechar string) {
	// 设置header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Authorization,Token,Back-Token")

	str := GetXmlStr(obj)
	logs.Info(">>>>>>END[", uniquechar, "]>>>>>>:", str)

	w.Write([]byte(str))
}

func GetJsonStr(obj interface{}) string {

	btes, err := json.Marshal(obj)

	if nil != err {
		logs.Error(err)
	}

	return string(btes)

}

func GetXmlStr(obj interface{}) string {
	btes, err := xml.Marshal(obj)
	if nil != err {
		logs.Error(err)
	}

	return string(btes)

}

// 获取参数 -- 只适用于beego
func getJsonParam(rVals map[string][]string) string {

	var bodyStr string
	for k, _ := range rVals {
		if "" != k {
			bodyStr = k
			break
		}
	}

	return bodyStr
}
