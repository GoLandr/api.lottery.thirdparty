package base

import (
	"fmt"
	"log"
	"mycommon/logs"
	"mycommon/utils"
	"sort"

	"api.lottery.thirdparty/global"

	"mycommon/mathstr"

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
	RecordList []*model.SSC
	Limit      map[int]*model.Limit
	Stars      map[int]*model.StarsLimt
}

func (this *SSC) LordInit(tablename string) {
	log.Println("SSC_LordInit")
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
	log.Println("sort_", mathstr.GetJsonPlainStr(this.RecordList))
	for _, v := range this.RecordList {
		this.BaseStat(5, v)
		this.StarsStat(10, v)
	}
	this.Print()
}

//添加记录
func (this *SSC) AddRecord(record model.SSC) {
	this.RecordList = append(this.RecordList, &record)
	this.BaseStat(5, &record)
	this.StarsStat(10, &record)
	this.Print()
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
			big, small = lotteryutils.GetBigSmall(record.One_ball, 5)
			odd, even = lotteryutils.GetOddEven(record.One_ball)
		} else if k == TWO_BALL {
			big, small = lotteryutils.GetBigSmall(record.Two_ball, 5)
			odd, even = lotteryutils.GetOddEven(record.Two_ball)
		} else if k == THRID_BALL {
			big, small = lotteryutils.GetBigSmall(record.Third_ball, 5)
			odd, even = lotteryutils.GetOddEven(record.Third_ball)
		} else if k == FOUR_BALL {
			big, small = lotteryutils.GetBigSmall(record.Four_ball, 5)
			odd, even = lotteryutils.GetOddEven(record.Four_ball)
		} else if k == FIVE_BALL {
			big, small = lotteryutils.GetBigSmall(record.Five_ball, 5)
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
}

func (this *SSC) Print() {
	for i := 1; i <= len(this.Limit); i++ {
		v, _ := this.Limit[i]
		fmt.Println("第", i, "球:大已开出", v.Big, "期,小已开出",
			v.Small, "期,单已开出", v.Odd, "期,双已开出", v.Even, "期")
	}
	for i := 0; i < len(this.Stars); i++ {
		v, _ := this.Stars[i]
		fmt.Println("号码", i, "未出次数", v.No, "次 已出次数", v.Open, "次")
	}
	//	for k, v := range this.Limit {
	//		fmt.Println("第", k, "球:大已开出", v.Big, "期,小已开出",
	//			v.Small, "期,单已开出", v.Odd, "期,双已开出", v.Even, "期")
	//	}
	//	for k, v := range this.Stars {
	//		fmt.Println("号码", k, "未出次数", v.No, "次 已出次数", v.Open, "次")
	//	}
}
