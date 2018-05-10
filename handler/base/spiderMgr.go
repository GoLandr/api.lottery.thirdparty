package base

var GLotteryAPI LotteryAPI //全局API管理
var CQSSC_API []API
var XJSSC_API []API
var TJSSC_API []API
var YNSSC_API []API

type LotteryAPI struct {
	CQSSC APIMgr
	XJSSC APIMgr
	TJSSC APIMgr
	YNSSC APIMgr
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

	TJSSC_API = append(TJSSC_API, API{PJ_TJSSC, TJSSC_API_PJ})
	TJSSC_API = append(TJSSC_API, API{OFFICIAL_TJSSC, TJSSC_API_OFFICIAL})

	YNSSC_API = append(YNSSC_API, API{PJ_YNSSC, YNSSC_API_PJ})
	YNSSC_API = append(YNSSC_API, API{OFFICIAL_YNSSC, YNSSC_API_OFFICIAL})

	cqssc := APIMgr{CQSSC_API[0].URL, CQSSC_START_TIME, CQSSC_END_TIME, CQSSC_API[0].Mode, CQSSC_API[0].Mode}
	GLotteryAPI.CQSSC = cqssc

	xjssc := APIMgr{XJSSC_API[0].URL, XJSSC_START_TIME, XJSSC_END_TIME, XJSSC_API[0].Mode, XJSSC_API[0].Mode}
	GLotteryAPI.XJSSC = xjssc

	tjssc := APIMgr{TJSSC_API[0].URL, TJSSC_START_TIME, TJSSC_END_TIME, TJSSC_API[0].Mode, TJSSC_API[0].Mode}
	GLotteryAPI.TJSSC = tjssc

	ynssc := APIMgr{YNSSC_API[0].URL, YNSSC_START_TIME, YNSSC_END_TIME, YNSSC_API[0].Mode, YNSSC_API[0].Mode}
	GLotteryAPI.YNSSC = ynssc
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
			GLotteryAPI.XJSSC.URL = XJSSC_API[GLotteryAPI.XJSSC.Index].URL
			GLotteryAPI.XJSSC.Mode = XJSSC_API[GLotteryAPI.XJSSC.Index].Mode
			GLotteryAPI.XJSSC.Index = XJSSC_API[GLotteryAPI.XJSSC.Index].Mode
		}
	}

	if mode == TJSSC_TYPE {
		if GLotteryAPI.TJSSC.Index == len(TJSSC_API) {
			GLotteryAPI.TJSSC.Mode = TJSSC_API[0].Mode
			GLotteryAPI.TJSSC.Index = TJSSC_API[0].Mode
			GLotteryAPI.TJSSC.URL = TJSSC_API[0].URL
		} else {
			GLotteryAPI.TJSSC.URL = TJSSC_API[GLotteryAPI.TJSSC.Index].URL
			GLotteryAPI.TJSSC.Mode = TJSSC_API[GLotteryAPI.TJSSC.Index].Mode
			GLotteryAPI.TJSSC.Index = TJSSC_API[GLotteryAPI.TJSSC.Index].Mode
		}
	}

	if mode == YNSSC_TYPE {
		if GLotteryAPI.YNSSC.Index == len(YNSSC_API) {
			GLotteryAPI.YNSSC.Mode = YNSSC_API[0].Mode
			GLotteryAPI.YNSSC.Index = YNSSC_API[0].Mode
			GLotteryAPI.YNSSC.URL = YNSSC_API[0].URL
		} else {
			GLotteryAPI.YNSSC.URL = YNSSC_API[GLotteryAPI.YNSSC.Index].URL
			GLotteryAPI.YNSSC.Mode = YNSSC_API[GLotteryAPI.YNSSC.Index].Mode
			GLotteryAPI.YNSSC.Index = YNSSC_API[GLotteryAPI.YNSSC.Index].Mode
		}
	}
}
