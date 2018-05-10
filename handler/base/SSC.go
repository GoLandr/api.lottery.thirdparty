package base

import (
	"fmt"
	"log"
	"mycommon/logs"
	"mycommon/utils"
	"sort"

	"api.lottery.thirdparty/global"

	"api.lottery.thirdparty/model"
	lotteryutils "api.lottery.thirdparty/utils"
)

type SSCSlice []*model.SSC

//排序
func (a SSCSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a SSCSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a SSCSlice) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return a[j].Flowid > a[i].Flowid
}

type SSC struct {
	Name        string
	RecordList  []*model.SSC
	Limit       map[int]*model.Limit
	Total_Limit model.Limit //总和
	Stars       map[int]*model.StarsLimt
	Pred_Limit  model.PredLimt
	Mode        int //类型
}

func (this *SSC) LordInit(tablename string, lotteryName string, mode int) {
	log.Println("SSC_LordInit")
	this.Name = lotteryName
	this.Mode = mode
	//	if this.RecordMap == nil {
	//		this.RecordMap = make(map[int]*model.SSC)
	//	}
	defer func() {
		if e := recover(); e != nil {
			logs.Error(e)
		}
	}()
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	var result []model.SSC
	err := orm.SetTable(tablename).OrderBy("flowid desc").Offset(0).Limit(100).FindAll(&result)
	utils.ThrowError(err)
	for k, _ := range result {
		this.RecordList = append(this.RecordList, &result[k])
		//		this.RecordMap[v.Flowid] = &result[k]
	}
	//排序
	sort.Sort(SSCSlice(this.RecordList))
	//	log.Println("sort_", mathstr.GetJsonPlainStr(this.RecordList))
	for _, v := range this.RecordList {
		this.BaseStat(5, v)
		this.StarsStat(10, v)
	}
	//	this.Print()
	this.pushMsg()
}

//添加记录
func (this *SSC) AddRecord(record model.SSC) {
	if len(this.RecordList) > 10000 {
		this.RecordList = append(this.RecordList[:1], this.RecordList[1:]...)
	}
	this.RecordList = append(this.RecordList, &record)
	this.BaseStat(5, &record)
	this.StarsStat(10, &record)
	//	this.Print()
	this.pushMsg()
}

//五星统计
func (this *SSC) StarsStat(ballSize int, record *model.SSC) {
	if this.Stars == nil {
		this.Stars = make(map[int]*model.StarsLimt)
		for i := 0; i < ballSize; i++ {
			limit := new(model.StarsLimt)
			this.Stars[i] = limit
		}
	}
	result := make([]int, 5)
	result[0] = record.One_ball
	result[1] = record.Two_ball
	result[2] = record.Third_ball
	result[3] = record.Four_ball
	result[4] = record.Five_ball
	//	log.Println("result_", result)
	for i := 0; i < ballSize; i++ {
		isOpen := false
		for j := 0; j < 5; j++ {
			if i == result[j] {
				isOpen = true
				break
			}
		}
		if isOpen {
			this.Stars[i].Open += 1
			this.Stars[i].No = 0
		} else {
			this.Stars[i].Open = 0
			this.Stars[i].No += 1
		}
	}
}
func (this *SSC) BaseStat(ballSize int, record *model.SSC) {
	if this.Limit == nil {
		this.Limit = make(map[int]*model.Limit)
		for i := 1; i <= ballSize; i++ {
			limit := new(model.Limit)
			this.Limit[i] = limit
		}
	}
	for k, v := range this.Limit {
		big := 0
		small := 0
		odd := 0
		even := 0

		if k == ONE_BALL {
			big, small = lotteryutils.GetBigSmall(record.One_ball, SSC_SPLIT)
			odd, even = lotteryutils.GetOddEven(record.One_ball)
		} else if k == TWO_BALL {
			big, small = lotteryutils.GetBigSmall(record.Two_ball, SSC_SPLIT)
			odd, even = lotteryutils.GetOddEven(record.Two_ball)
		} else if k == THRID_BALL {
			big, small = lotteryutils.GetBigSmall(record.Third_ball, SSC_SPLIT)
			odd, even = lotteryutils.GetOddEven(record.Third_ball)
		} else if k == FOUR_BALL {
			big, small = lotteryutils.GetBigSmall(record.Four_ball, SSC_SPLIT)
			odd, even = lotteryutils.GetOddEven(record.Four_ball)
		} else if k == FIVE_BALL {
			big, small = lotteryutils.GetBigSmall(record.Five_ball, SSC_SPLIT)
			odd, even = lotteryutils.GetOddEven(record.Five_ball)
		}

		if big == 1 {
			v.Big += big
			v.Small = 0
		} else {
			v.Big = 0
			v.Small += small
		}
		if odd == 1 {
			v.Odd += odd
			v.Even = 0
		} else {
			v.Even += even
			v.Odd = 0
		}
	}
	//计算总和
	total_small := 0
	total_odd := 0
	total_even := 0
	total_big := 0
	total_odd, total_even, total_big, total_small = lotteryutils.GetTotalStat(this.recordToArray(record), SSC_TOTAL_SPLIT)
	if total_big == 1 {
		this.Total_Limit.Big += total_big
		this.Total_Limit.Small = 0
	} else {
		this.Total_Limit.Big = 0
		this.Total_Limit.Small += total_small
	}
	if total_odd == 1 {
		this.Total_Limit.Odd += total_odd
		this.Total_Limit.Even = 0
	} else {
		this.Total_Limit.Even += total_even
		this.Total_Limit.Odd = 0
	}
	//计算龙虎
	dragon, tiger, draw := lotteryutils.GetPredStat(record.One_ball, record.Five_ball,
		this.Pred_Limit.Dragon, this.Pred_Limit.Tiger, this.Pred_Limit.Draw)
	this.Pred_Limit.Dragon = dragon
	this.Pred_Limit.Tiger = tiger
	this.Pred_Limit.Draw = draw
}

func (this *SSC) recordToArray(record *model.SSC) []int {
	var array []int
	array = append(array, record.One_ball)
	array = append(array, record.Two_ball)
	array = append(array, record.Third_ball)
	array = append(array, record.Four_ball)
	array = append(array, record.Five_ball)
	return array
}

func (this *SSC) Print() {
	str := fmt.Sprint(this.Name, "\n")

	for i := 1; i <= len(this.Limit); i++ {
		v, _ := this.Limit[i]
		str = fmt.Sprint(str, "第", i, "球:大已开出", v.Big, "期,小已开出",
			v.Small, "期,单已开出", v.Odd, "期,双已开出", v.Even, "期", "\n")
	}
	str = fmt.Sprint(str, "总和大已开出", this.Total_Limit.Big, "期，总和小已开出", this.Total_Limit.Small, "期\n")
	str = fmt.Sprint(str, "总和单已开出", this.Total_Limit.Odd, "期，总和双已开出", this.Total_Limit.Even, "期\n")
	str = fmt.Sprint(str, "龙已开出", this.Pred_Limit.Dragon, "期，虎已开出", this.Pred_Limit.Tiger, "期，和已开出", this.Pred_Limit.Draw, "期\n")
	for i := 0; i < len(this.Stars); i++ {
		v, _ := this.Stars[i]
		str = fmt.Sprint(str, "号码", i, "未出次数", v.No, "次 已出次数", v.Open, "次", "\n")
	}
	logs.Debug(str)
}
func (this *SSC) pushMsg() {
	BSlimit, _ := GBigSmallLimit[this.Mode]
	OElimit, _ := GOddEvenLimit[this.Mode]
	//
	//	starlimit, ok := GStarsLimit[this.Mode]
	//
	//	predlimit, ok := GPredLimit[this.Mode]
	for i := 1; i <= len(this.Limit); i++ {
		v, _ := this.Limit[i]
		BS_maxVal := lotteryutils.GetMaxValue(v.Big, v.Small)
		BS_retLst := GetPushMenber(BSlimit, BS_maxVal)
		if len(BS_retLst) > 0 {
			msg := fmt.Sprint(this.Name, ":第", i, "球:大已开出", v.Big, "期,小已开出", v.Small, "期")
			sendMsgToFriend(BS_retLst, msg)
		}
		OE_maxVal := lotteryutils.GetMaxValue(v.Odd, v.Even)
		OE_retLst := GetPushMenber(OElimit, OE_maxVal)
		if len(OE_retLst) > 0 {
			msg := fmt.Sprint(this.Name, ":第", i, "球:单已开出", v.Odd, "期,双已开出", v.Even, "期")
			sendMsgToFriend(OE_retLst, msg)
		}
	}
	total_BS_limit, tok := GTotalBSLimit[this.Mode]
	if tok {
		Total_BS_maxVal := lotteryutils.GetMaxValue(this.Total_Limit.Big, this.Total_Limit.Small)
		Total_BS_retLst := GetPushMenber(total_BS_limit, Total_BS_maxVal)
		if len(Total_BS_retLst) > 0 {
			msg := fmt.Sprint(this.Name, ":总和大已开出", this.Total_Limit.Big, "期,总和小已开出", this.Total_Limit.Small, "期")
			sendMsgToFriend(Total_BS_retLst, msg)
		}

	}
	total_OE_limit, oeOk := GTotalOELimit[this.Mode]
	if oeOk {
		Total_OE_maxVal := lotteryutils.GetMaxValue(this.Total_Limit.Odd, this.Total_Limit.Even)
		Total_OE_retLst := GetPushMenber(total_OE_limit, Total_OE_maxVal)
		if len(Total_OE_retLst) > 0 {
			msg := fmt.Sprint(this.Name, ":总和单已开出", this.Total_Limit.Odd, "期,总和双已开出", this.Total_Limit.Even, "期")
			sendMsgToFriend(Total_OE_retLst, msg)
		}
	}
	pred_limit, pok := GPredLimit[this.Mode]
	if pok {
		pred_maxVal := lotteryutils.GetMaxValue(this.Pred_Limit.Dragon, this.Pred_Limit.Draw, this.Pred_Limit.Tiger)
		pred_retLst := GetPushMenber(pred_limit, pred_maxVal)
		if len(pred_retLst) > 0 {
			msg := fmt.Sprint(this.Name, ":龙已开出", this.Pred_Limit.Dragon, "期，虎已开出", this.Pred_Limit.Tiger, "期，和已开出", this.Pred_Limit.Draw, "期")
			sendMsgToFriend(pred_retLst, msg)
		}
	}
	star_limit, sok := GStarsLimit[this.Mode]
	if sok {
		for i := 0; i < len(this.Stars); i++ {
			v, _ := this.Stars[i]
			star_maxVal := lotteryutils.GetMaxValue(v.Open, v.No)
			star_retLst := GetPushMenber(star_limit, star_maxVal)
			if len(star_retLst) > 0 {
				msg := fmt.Sprint(this.Name, ":号码", i, "未出次数", v.No, "次 已出次数", v.Open, "次")
				sendMsgToFriend(star_retLst, msg)
			}
		}
	}
}
