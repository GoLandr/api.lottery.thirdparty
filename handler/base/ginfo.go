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
var GLimit map[int][]*model.TLimit
var GOddEvenLimit map[int]map[int][]*model.TLimit

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
	for k, v := range GLimit {
		oddEvenLst := make(map[int][]*model.TLimit)
		for _, entry := range v {
			oddEvenLst[entry.Odd_even_limit] = append(oddEvenLst[entry.Odd_even_limit], entry)
		}
		GOddEvenLimit[k] = oddEvenLst
	}
	logs.Debug("GOddEvenLimit_", mathstr.GetJsonPlainStr(GOddEvenLimit))
}
