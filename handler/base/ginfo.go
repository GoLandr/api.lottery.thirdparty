package base

import (
	"mycommon/logs"
	"mycommon/mathstr"
	"mycommon/utils"

	"api.lottery.thirdparty/global"

	"api.lottery.thirdparty/model"
)

var GLotteryMgr *LotteryMgr //全局变量
var GMenber map[int]*model.Menber
var GLimit map[int][]*model.TLimit                 //k->mode
var GOddEvenLimit map[int]map[int][]*model.TLimit  //单双 k->mode ->odd_even
var GBigSmallLimit map[int]map[int][]*model.TLimit //大小 k->mode ->big_small
var GTotalBSLimit map[int]map[int][]*model.TLimit  //总和大小 k->mode ->total
var GTotalOELimit map[int]map[int][]*model.TLimit  //总和单双 k->mode ->total
var GStarsLimit map[int]map[int][]*model.TLimit    //五星 k->mode ->star
var GPredLimit map[int]map[int][]*model.TLimit     //龙虎 k->mode ->pred

func InitConfigs() {
	if GLotteryMgr == nil {
		mgr := new(LotteryMgr)
		GLotteryMgr = mgr
	}
	InitMenber()
	InitLimit()
}

func InitMenber() {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(e)
		}
	}()
	GMenber = make(map[int]*model.Menber)
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	var result []model.Menber
	err := orm.SetTable("t_menber").FindAll(&result)
	utils.ThrowError(err)
	for k, v := range result {
		GMenber[v.Id] = &result[k]
	}
	logs.Debug("GMenber_", mathstr.GetJsonPlainStr(GMenber))
}

func InitLimit() {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(e)
		}
	}()
	GLimit = make(map[int][]*model.TLimit)
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	var result []model.TLimit
	err := orm.SetTable("t_limit").FindAll(&result)
	utils.ThrowError(err)
	logs.Debug("Limit_", mathstr.GetJsonPlainStr(result))
	for k, v := range result {
		GLimit[v.Mode] = append(GLimit[v.Mode], &result[k])
	}
	logs.Debug("GLimit_", mathstr.GetJsonPlainStr(GLimit))
	GOddEvenLimit = make(map[int]map[int][]*model.TLimit)
	GBigSmallLimit = make(map[int]map[int][]*model.TLimit)
	GTotalBSLimit = make(map[int]map[int][]*model.TLimit)
	GTotalOELimit = make(map[int]map[int][]*model.TLimit)
	GStarsLimit = make(map[int]map[int][]*model.TLimit)
	GPredLimit = make(map[int]map[int][]*model.TLimit)
	for k, v := range GLimit {
		oddEvenLst := make(map[int][]*model.TLimit)
		bigSmallLst := make(map[int][]*model.TLimit)
		totalOELst := make(map[int][]*model.TLimit)
		totalBSLst := make(map[int][]*model.TLimit)
		starsLst := make(map[int][]*model.TLimit)
		predLst := make(map[int][]*model.TLimit)
		for _, entry := range v {
			oddEvenLst[entry.Odd_even_limit] = append(oddEvenLst[entry.Odd_even_limit], entry)
			bigSmallLst[entry.Big_small_limit] = append(bigSmallLst[entry.Big_small_limit], entry)
			totalOELst[entry.Total_oe_limit] = append(totalOELst[entry.Total_oe_limit], entry)
			totalBSLst[entry.Total_bs_limit] = append(totalBSLst[entry.Total_bs_limit], entry)
			starsLst[entry.Star_limit] = append(starsLst[entry.Star_limit], entry)
			predLst[entry.Pred_limit] = append(predLst[entry.Pred_limit], entry)

		}
		GOddEvenLimit[k] = oddEvenLst
		GBigSmallLimit[k] = bigSmallLst
		GTotalOELimit[k] = totalOELst
		GTotalBSLimit[k] = totalBSLst
		GStarsLimit[k] = starsLst
		GPredLimit[k] = predLst
	}
	logs.Debug("GOddEvenLimit_", mathstr.GetJsonPlainStr(GOddEvenLimit))
}
