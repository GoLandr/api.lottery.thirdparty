package base

import (
	"mycommon/logs"
	"mycommon/utils"
	"nn/global"

	"api.lottery.thirdparty/model"
)

var GLotteryMgr *LotteryMgr //全局变量
func InitConfigs() {
	if GLotteryMgr == nil {
		mgr := new(LotteryMgr)
		GLotteryMgr = mgr
	}
}

func InitMenber() {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(e)
		}
	}()
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	var menber []model.Menber
	err := orm.SetTable("t_menber").Find(&menber)
	utils.ThrowError(err)
}
