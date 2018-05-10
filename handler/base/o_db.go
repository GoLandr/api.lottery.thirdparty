package base

import (
	"mycommon/logs"
	"mycommon/mathstr"
	"mycommon/utils"

	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/model"
)

func SaveSSC(tablename string, lottery interface{}, mode int) {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(e)
		}
	}()
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	if mode == CQSSC_TYPE || mode == XJSSC_TYPE ||
		mode == TJSSC_TYPE || mode == YNSSC_TYPE {
		ssc := lottery.(model.SSC)
		err := orm.SetTable(tablename).SetPK("flowid").InsertModel(ssc)
		utils.ThrowError(err)
	}

}

//檢查是否有数据
func CheckLottery(tablename string, key string) int {
	defer func() {
		if e := recover(); e != nil {
			logs.Error(e)
		}
	}()
	haveCount := 0
	//	str := fmt.Sprint("Select count(1) num from ", tablename,
	//		" where flowid='", key, "'")
	//	fmt.Println("str_____-", str)
	//	sqlstr := mathstr.S_SFT(str)
	sqlstr := mathstr.S_SFT(`
		Select count(1) num from {0}
		where flowid='{1}'
	`, tablename, key)
	//	fmt.Println("sqlstr_____-", sqlstr)
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	play, err := orm.QueryOne(sqlstr)
	utils.ThrowError(err)
	//	log.Println("play_", play)
	if err == nil {
		haveCount = mathstr.Math2intDefault0(play["num"])
	}
	return haveCount
}
