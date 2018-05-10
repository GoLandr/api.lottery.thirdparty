package base

import (
	"fmt"
	"log"
	"mycommon/logs"
	"mycommon/mathstr"
	"mycommon/utils"
	"net/url"
	"runtime"
	"strings"
	"time"

	"api.lottery.thirdparty/common"
	"api.lottery.thirdparty/model"
	lotteryutils "api.lottery.thirdparty/utils"
	"github.com/robfig/cron"
)

type Spider struct {
}

func (this *Spider) SpiderCron() {
	InitConfigs()
	LoardAPI()
	timer := time.NewTimer(time.Duration(0) * time.Second)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				err, ok := e.(error)
				if ok {
					// 日志记录
					for i := 2; i <= 8; i++ {
						_, f, line, ok := runtime.Caller(i)
						if !ok {
							continue
						}
						if i == 2 {
							logs.Error(i, "__err:[", err, "]__fname:[", f, "]__line:[", line, "]")
						} else {
							logs.Error(i, "__fname:[", f, "]__line:[", line, "]")
						}
					}
				}
			}
		}()
		//等触发时的信号
		<-timer.C
		this.LoardSpider(STATUS_YES)
		timer.Stop()

		//进入第二次执行时间
		i := 0
		c := cron.New()
		c.AddFunc(fmt.Sprint("@every ", 60, "s"), func() {
			this.LoardSpider(STATUS_NO)
			//			pj.Net_Cqssc()
			i++

		})
		c.Start()
	}()

}
func (this *Spider) LoardSpider(lordinit int) {
	if lotteryutils.JudgeTime(GLotteryAPI.CQSSC.StartTime, GLotteryAPI.CQSSC.EndTime) {
		if GLotteryAPI.CQSSC.Mode == CQSSC_API_PJ {
			Pj_SSC(PJ_CQSSC, T_CQSSC, CQSSC_TYPE, lordinit)
			//		logs.Debug("LoardSpider_err", mathstr.GetJsonStr(err))
		} else if GLotteryAPI.CQSSC.Mode == CQSSC_API_OFFICIAL {
			Official_SSC(OFFICIAL_CQSSC, T_CQSSC, CQSSC_TYPE, lordinit)
		}
	}
	if lotteryutils.JudgeTime(GLotteryAPI.XJSSC.StartTime, GLotteryAPI.XJSSC.EndTime) {
		if GLotteryAPI.XJSSC.Mode == XJSSC_API_PJ {
			//		logs.Debug("XJSSC")
			Pj_SSC(PJ_XJSSC, T_XJSSC, XJSSC_TYPE, lordinit)
		} else if GLotteryAPI.XJSSC.Mode == XJSSC_API_OFFICIAL {
			Official_SSC(OFFICIAL_XJSSC, T_XJSSC, XJSSC_TYPE, lordinit)
		}
	} else {
		logs.Debug("no in curtime")
	}

	if lordinit == STATUS_YES {
		GLotteryMgr.Cqssc.LordInit(T_CQSSC, CQSSC_NAME, CQSSC_TYPE)
		GLotteryMgr.Xjssc.LordInit(T_XJSSC, XJSSC_NAME, XJSSC_TYPE)
	}
}

type SSCModel struct {
	Expect        string `json:"expect"`
	Opencode      string `json:"opencode"`
	Opentime      string `json:"opentime"`
	Opentimestamp int    `json:"opentimestamp"`
}

func Official_SSC(urlstr string, tablename string, mode int, lordinit int) error {
	log.Println("visit Official_SSC")
	defer func() {
		if e := recover(); e != nil {
			logs.Debug("Fail to collect and replace the source")
			ChangeLotteryAPI(mode)
			logs.Error(e)
		}

	}()
	param := make(url.Values)
	result, err := common.Httppost(urlstr, param)
	if err != nil {
		return err
	}
	var redata map[string]interface{}
	mathstr.JsonUnmarsh(result, &redata)
	//	log.Println(result)
	var hmlist []SSCModel
	mathstr.JsonUnmarsh(mathstr.GetJsonPlainStr(redata["data"]), &hmlist)
	//	fmt.Println("hmlist_2_", hmlist)
	for _, v := range hmlist {
		havecount := CheckLottery(tablename, v.Expect)
		if havecount == 0 {
			log.Println("save_", tablename, "_Expect", v.Expect, "_time_", v.Opentime)
			//保存
			ssc := model.SSC{}
			ssc.Flowid = mathstr.Math2intDefault0(v.Expect)
			ball := strings.Split(v.Opencode, ",")
			ssc.One_ball = mathstr.Math2intDefault0(ball[0])
			ssc.Two_ball = mathstr.Math2intDefault0(ball[1])
			ssc.Third_ball = mathstr.Math2intDefault0(ball[2])
			ssc.Four_ball = mathstr.Math2intDefault0(ball[3])
			ssc.Five_ball = mathstr.Math2intDefault0(ball[4])
			ssc.Periods = v.Expect[8:len(v.Expect)]
			ssc.Update_date = utils.Now()
			ssc.Lottery_date = v.Opentime[0:10]
			ssc.Lottery_time = v.Opentime[11:len(v.Opentime)]
			//			log.Println(ssc.Lottery_date, "_", ssc.Lottery_time)
			if lordinit == STATUS_YES {
				SaveLottery(ssc, mode, STATUS_NO, tablename)
			} else {
				SaveLottery(ssc, mode, STATUS_YES, tablename)
			}
		}
	}
	return nil
}

//PJSSC
func Pj_SSC(urlstr string, tablename string, mode int, lordinit int) error {
	//	log.Println("visit Pj_SSC_mode", mode, "_", tablename)
	defer func() {
		logs.Debug("visit defer")
		if e := recover(); e != nil {
			logs.Debug("Fail to collect and replace the source")
			ChangeLotteryAPI(mode)
			logs.Error(e)
		}

	}()
	param := make(url.Values)
	result, err := common.Httppost(urlstr, param)
	if err != nil {
		return err
	}
	var redata map[string]interface{}
	mathstr.JsonUnmarsh(result, &redata)
	//	fmt.Println(redata)
	//	fmt.Println("hmlist_", mathstr.GetJsonPlainStr(redata["hmlist"]))
	if lordinit == STATUS_YES {
		//查看历史记录是否保存
		var hmlist map[string]string
		mathstr.JsonUnmarsh(mathstr.GetJsonPlainStr(redata["hmlist"]), &hmlist)
		//		fmt.Println("hmlist_2_", hmlist)
		for k, v := range hmlist {
			//			fmt.Println("k_", k, "ball_", strings.Split(v, ","))
			havecount := CheckLottery(tablename, k)
			//			logs.Debug("havecount_", havecount)
			if havecount == 0 {
				//保存
				ssc := model.SSC{}
				ssc.Flowid = mathstr.Math2intDefault0(k)
				ball := strings.Split(v, ",")
				ssc.One_ball = mathstr.Math2intDefault0(ball[0])
				ssc.Two_ball = mathstr.Math2intDefault0(ball[1])
				ssc.Third_ball = mathstr.Math2intDefault0(ball[2])
				ssc.Four_ball = mathstr.Math2intDefault0(ball[3])
				ssc.Five_ball = mathstr.Math2intDefault0(ball[4])
				ssc.Periods = k[8:len(k)]
				ssc.Update_date = utils.Now()
				ssc.Lottery_date = k[0:8]
				//				SaveSSC(ssc)
				//				logs.Debug("ssc_", ssc)
				SaveLottery(ssc, mode, STATUS_NO, tablename)
			}
		}
	} else {
		//最新记录
		ssc := model.SSC{}
		ssc.Flowid = mathstr.Math2intDefault0(redata["numbers"])
		var ball []string
		mathstr.JsonUnmarsh(fmt.Sprint(mathstr.GetJsonPlainStr(redata["hm"])), &ball)
		ssc.One_ball = mathstr.Math2intDefault0(ball[0])
		ssc.Two_ball = mathstr.Math2intDefault0(ball[1])
		ssc.Third_ball = mathstr.Math2intDefault0(ball[2])
		ssc.Four_ball = mathstr.Math2intDefault0(ball[3])
		ssc.Five_ball = mathstr.Math2intDefault0(ball[4])
		t_flowid := fmt.Sprint(redata["numbers"])
		ssc.Periods = t_flowid[8:len(t_flowid)]
		ssc.Update_date = utils.Now()
		ssc.Lottery_date = t_flowid[0:8]
		havecount := CheckLottery(tablename, t_flowid)
		if havecount == 0 {
			log.Println("save_", tablename, "_Expect", ssc.Flowid)
			//			SaveSSC(ssc)
			//			GLotteryMgr.Cqssc.AddRecord(ssc)
			SaveLottery(ssc, mode, STATUS_YES, tablename)
		}
	}
	return nil
}

func Pj_BJPK(urlstr string, tablename string, mode int, lordinit int) error {
	param := make(url.Values)
	result, err := common.Httppost(urlstr, param)
	if err != nil {
		return err
	}
	var redata map[string]interface{}
	mathstr.JsonUnmarsh(result, &redata)
	fmt.Println(redata)
	//	fmt.Println("hmlist_", mathstr.GetJsonPlainStr(redata["hmlist"]))
	if lordinit == STATUS_YES {
		//查看历史记录是否保存
		var hmlist map[string]string
		mathstr.JsonUnmarsh(mathstr.GetJsonPlainStr(redata["hmlist"]), &hmlist)
		fmt.Println("hmlist_2_", hmlist)
		for k, v := range hmlist {
			//			fmt.Println("k_", k, "ball_", strings.Split(v, ","))
			havecount := CheckLottery(tablename, k)
			if havecount == 0 {
				v = strings.Replace(v, "<br>", ",", 1)
				logs.Debug("v_", v)
				//保存
				bjpk := model.BJPK{}
				bjpk.Flowid = mathstr.Math2intDefault0(k)
				ball := strings.Split(v, ",")
				bjpk.One_ball = mathstr.Math2intDefault0(ball[0])
				bjpk.Two_ball = mathstr.Math2intDefault0(ball[1])
				bjpk.Third_ball = mathstr.Math2intDefault0(ball[2])
				bjpk.Four_ball = mathstr.Math2intDefault0(ball[3])
				bjpk.Five_ball = mathstr.Math2intDefault0(ball[4])
				bjpk.Six_ball = mathstr.Math2intDefault0(ball[5])
				bjpk.Seven_ball = mathstr.Math2intDefault0(ball[6])
				bjpk.Eight_ball = mathstr.Math2intDefault0(ball[7])
				bjpk.Ninth_ball = mathstr.Math2intDefault0(ball[8])
				bjpk.Ten_ball = mathstr.Math2intDefault0(ball[9])
				bjpk.Periods = k
				bjpk.Update_date = utils.Now()
				//				ssc.Lottery_date = k[0:8]
				//				SaveSSC(ssc)
				SaveLottery(bjpk, mode, STATUS_NO, tablename)
			}
		}
	} else {
		//最新记录
		bjpk := model.BJPK{}
		bjpk.Flowid = mathstr.Math2intDefault0(redata["numbers"])
		var ball []string
		mathstr.JsonUnmarsh(fmt.Sprint(mathstr.GetJsonPlainStr(redata["hm"])), &ball)
		bjpk.One_ball = mathstr.Math2intDefault0(ball[0])
		bjpk.Third_ball = mathstr.Math2intDefault0(ball[2])
		bjpk.Four_ball = mathstr.Math2intDefault0(ball[3])
		bjpk.Five_ball = mathstr.Math2intDefault0(ball[4])
		bjpk.Six_ball = mathstr.Math2intDefault0(ball[5])
		bjpk.Seven_ball = mathstr.Math2intDefault0(ball[6])
		bjpk.Eight_ball = mathstr.Math2intDefault0(ball[7])
		bjpk.Ninth_ball = mathstr.Math2intDefault0(ball[8])
		bjpk.Ten_ball = mathstr.Math2intDefault0(ball[9])
		bjpk.Periods = fmt.Sprint(redata["numbers"])
		bjpk.Update_date = utils.Now()
		//		ssc.Lottery_date = t_flowid[0:8]
		havecount := CheckLottery(tablename, bjpk.Periods)
		if havecount == 0 {
			//			SaveSSC(ssc)
			//			GLotteryMgr.Cqssc.AddRecord(ssc)
			SaveLottery(bjpk, mode, STATUS_YES, tablename)
		}
	}
	return nil
}

func SaveLottery(lottery interface{}, mode int, loadRecord int, tablename string) {
	if mode == CQSSC_TYPE {
		ssc := lottery.(model.SSC)
		SaveSSC(tablename, ssc, mode)
		if loadRecord == STATUS_YES {
			GLotteryMgr.Cqssc.AddRecord(ssc)
		}
	} else if mode == XJSSC_TYPE {
		logs.Debug("save_xj")
		ssc := lottery.(model.SSC)
		SaveSSC(tablename, ssc, mode)
		if loadRecord == STATUS_YES {
			GLotteryMgr.Xjssc.AddRecord(ssc)
		}
	}
}

//func (this *Lottery) SpiderUrl(url string) error {
//	doc, err := goquery.NewDocument(url)
//	if err != nil {
//		return nil
//	}
//	bookname := common.GbkToUtf8(doc.Find("#info h1").Text())
//	fmt.Println("doc_", mathstr.GetJsonPlainStr(bookname))
//	return nil
//}
