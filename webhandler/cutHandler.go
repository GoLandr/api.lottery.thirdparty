package webhandler

import (
	"mycommon/logs"
	"mycommon/mathstr"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/config"
	"api.lottery.thirdparty/global"
)

type Cutfunc struct{}

func (this *Cutfunc) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "select":
		return this.selectWhere(param)

	case "set":
		return this.setValue(param)

	}

	return "no func name:" + do, 999
}

func (this *Cutfunc) selectWhere(vals map[string]interface{}) (interface{}, int) {
	orm := global.GetOrm()
	mp, err := orm.SetTable("sss_configs").FindOneMap()
	utils.ThrowError(err)

	cutCost := mathstr.Math2intDefault0(mp["cut_cost"])

	return cutCost, -1
}

func (this *Cutfunc) setValue(vals map[string]interface{}) (interface{}, int) {
	cost := mathstr.Math2intDefault0(vals["cost"])
	utils.S_CRPM(vals, "cost")

	if cost < 0 {
		cost = 0
	}

	orm := global.GetOrm()
	sqlupdate := mathstr.S_SFT(`update sss_configs set cut_cost={0}`, cost)
	_, err := orm.Exec(sqlupdate)
	utils.ThrowError(err)

	go this.LocalSyncConigs()

	return nil, -1
}
func (this *Cutfunc) LocalSyncConigs() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	server := config.GetLogicServer() + "/ReloadDbCfg"
	logs.Debug("__restartconfigs:", server)
	resp, err := utils.HttpGet(server)
	utils.ThrowError(err)
	logs.Debug("__resp:", resp)
}
