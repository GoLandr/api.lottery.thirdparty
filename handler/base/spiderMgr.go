package base

var GLotteryAPI LotteryAPI //全局API管理
var CQSSC_API []API
var XJSSC_API []API

type LotteryAPI struct {
	CQSSC APIMgr
	XJSSC APIMgr
}
type APIMgr struct {
	URL       string
	StartTime string //开始时间
	EndTime   string //结束时间
	Mode      int    //类型
	Index     int    //当前选择
}

type API struct {
	URL  string
	Mode int
}

func LoardAPI() {
	CQSSC_API = append(CQSSC_API, API{PJ_CQSSC, CQSSC_API_PJ})
	CQSSC_API = append(CQSSC_API, API{OFFICIAL_CQSSC, CQSSC_API_OFFICIAL})
	XJSSC_API = append(XJSSC_API, API{PJ_XJSSC, XJSSC_API_PJ})
	XJSSC_API = append(XJSSC_API, API{OFFICIAL_XJSSC, XJSSC_API_OFFICIAL})
	cqssc := APIMgr{CQSSC_API[0].URL, CQSSC_START_TIME, CQSSC_END_TIME, CQSSC_API[0].Mode, CQSSC_API[0].Mode}
	GLotteryAPI.CQSSC = cqssc
	xjssc := APIMgr{XJSSC_API[0].URL, XJSSC_START_TIME, XJSSC_END_TIME, XJSSC_API[0].Mode, XJSSC_API[0].Mode}
	GLotteryAPI.XJSSC = xjssc
}
func ChangeLotteryAPI(mode int) {
	if mode == CQSSC_TYPE {
		if GLotteryAPI.CQSSC.Index == len(CQSSC_API) {
			GLotteryAPI.CQSSC.Mode = CQSSC_API[0].Mode
			GLotteryAPI.CQSSC.Index = CQSSC_API[0].Mode
			GLotteryAPI.CQSSC.URL = CQSSC_API[0].URL
		} else {
			//			logs.Debug("cqsss_api", mathstr.GetJsonPlainStr(CQSSC_API))
			//			logs.Debug("cqssc_", GLotteryAPI.CQSSC)
			GLotteryAPI.CQSSC.URL = CQSSC_API[GLotteryAPI.CQSSC.Index].URL
			GLotteryAPI.CQSSC.Mode = CQSSC_API[GLotteryAPI.CQSSC.Index].Mode
			GLotteryAPI.CQSSC.Index = CQSSC_API[GLotteryAPI.CQSSC.Index].Mode
		}
	}
	if mode == XJSSC_TYPE {
		if GLotteryAPI.XJSSC.Index == len(XJSSC_API) {
			GLotteryAPI.XJSSC.Mode = XJSSC_API[0].Mode
			GLotteryAPI.XJSSC.Index = XJSSC_API[0].Mode
			GLotteryAPI.XJSSC.URL = XJSSC_API[0].URL
		} else {
			GLotteryAPI.XJSSC.URL = CQSSC_API[GLotteryAPI.XJSSC.Index].URL
			GLotteryAPI.XJSSC.Mode = CQSSC_API[GLotteryAPI.XJSSC.Index].Mode
			GLotteryAPI.XJSSC.Index = CQSSC_API[GLotteryAPI.XJSSC.Index].Mode
		}
	}
}
