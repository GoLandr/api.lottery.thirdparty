package spider

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
	"api.lottery.thirdparty/handler/base"
	"api.lottery.thirdparty/model"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
)

type Lottery struct {
}

//澳门新葡京
type PuJing struct {
	Url     string
	IsFirst bool //是否第一次调用
}
type Spider struct {
}

func (this *Spider) SpiderCron() {
	base.InitConfigs()
	pj := new(PuJing)
	pj.IsFirst = true
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
		pj.Pj_SSC("http://www.j0024.com/lottery/getAllCqsscAutoList", "lottery_cqssc", base.CQSSC_TYPE)
		pj.Pj_SSC("http://www.j0024.com/lottery/getAllXjsscAutoList", "lottery_xjssc", base.XJSSC_TYPE)
		//		pj.Net_Cqssc()
		pj.IsFirst = false
		timer.Stop()
		logs.Debug("SSC加载初始化")
		base.GLotteryMgr.Cqssc.LordInit("lottery_cqssc", "重庆时时彩")
		base.GLotteryMgr.Xjssc.LordInit("lottery_xjssc", "新疆时时彩")
		//进入第二次执行时间
		i := 0
		c := cron.New()
		c.AddFunc(fmt.Sprint("@every ", 60, "s"), func() {
			pj.Pj_SSC("http://www.j0024.com/lottery/getAllCqsscAutoList", "lottery_cqssc", base.CQSSC_TYPE)
			pj.Pj_SSC("http://www.j0024.com/lottery/getAllXjsscAutoList", "lottery_xjssc", base.XJSSC_TYPE)
			//			pj.Net_Cqssc()
			i++

		})
		c.Start()
	}()

}

type SSCModel struct {
	Expect        string `json:"expect"`
	Opencode      string `json:"opencode"`
	Opentime      string `json:"opentime"`
	Opentimestamp int    `json:"opentimestamp"`
}

func (this *PuJing) Net_Cqssc() error {
	urlstr := "http://f.apiplus.net/cqssc-20.json"
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
		havecount := base.CheckLottery("lottery_cqssc", v.Expect)
		if havecount == 0 {
			log.Println("save_Expect", v.Expect, "_time_", v.Opentime)
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
			//			base.SaveSSC(ssc)
			SaveLottery(ssc, base.CQSSC_TYPE, base.STATUS_NO, "lottery_cqssc")
			if !this.IsFirst {
				//				base.GLotteryMgr.Cqssc.AddRecord(ssc)
				SaveLottery(ssc, base.CQSSC_TYPE, base.STATUS_YES, "lottery_cqssc")
			}
		}
	}
	return nil
}

//葡京时时彩
func (this *PuJing) Pj_SSC(urlstr string, tablename string, mode int) error {
	//	urlstr := "http://www.j0024.com/lottery/getAllCqsscAutoList"
	param := make(url.Values)
	result, err := common.Httppost(urlstr, param)
	if err != nil {
		return err
	}
	var redata map[string]interface{}
	mathstr.JsonUnmarsh(result, &redata)
	fmt.Println(redata)
	fmt.Println("hmlist_", mathstr.GetJsonPlainStr(redata["hmlist"]))
	if this.IsFirst {
		//查看历史记录是否保存
		var hmlist map[string]string
		mathstr.JsonUnmarsh(mathstr.GetJsonPlainStr(redata["hmlist"]), &hmlist)
		fmt.Println("hmlist_2_", hmlist)
		for k, v := range hmlist {
			//			fmt.Println("k_", k, "ball_", strings.Split(v, ","))
			havecount := base.CheckLottery(tablename, k)
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
				//				base.SaveSSC(ssc)
				SaveLottery(ssc, mode, base.STATUS_NO, tablename)
			}
		}
		this.IsFirst = false
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
		havecount := base.CheckLottery(tablename, t_flowid)
		if havecount == 0 {
			//			base.SaveSSC(ssc)
			//			base.GLotteryMgr.Cqssc.AddRecord(ssc)
			SaveLottery(ssc, mode, base.STATUS_YES, tablename)
		}
	}
	return nil
}

func SaveLottery(lottery interface{}, mode int, loadRecord int, tablename string) {
	if mode == base.CQSSC_TYPE {
		ssc := lottery.(model.SSC)
		base.SaveSSC(tablename, ssc, mode)
		if loadRecord == base.STATUS_YES {
			base.GLotteryMgr.Cqssc.AddRecord(ssc)
		}
	} else if mode == base.XJSSC_TYPE {
		ssc := lottery.(model.SSC)
		base.SaveSSC(tablename, ssc, mode)
		if loadRecord == base.STATUS_YES {
			base.GLotteryMgr.Xjssc.AddRecord(ssc)
		}
	}
}

func (this *Lottery) SpiderUrl(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil
	}
	bookname := common.GbkToUtf8(doc.Find("#info h1").Text())
	fmt.Println("doc_", mathstr.GetJsonPlainStr(bookname))
	return nil
}
