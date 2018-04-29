package webhandler

import (
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/global"
)

type SysCnf struct{}

func (this *SysCnf) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "set":
		return this.setInfo(param)

	case "get":
		return this.getInfo(param)

	}

	return "no func name:" + do, 999
}

func (this *SysCnf) getInfo(vals map[string]interface{}) (interface{}, int) {
	return global.GetSyscnf(fmt.Sprint(vals["key"])), -1
}

func (this *SysCnf) setInfo(vals map[string]interface{}) (interface{}, int) {
	ckey := fmt.Sprint(vals["key"])
	cval := fmt.Sprint(vals["val"])
	utils.S_CRPM(vals, "key", "val")

	orm := global.GetOrm()
	resmap, err := orm.SetTable("t_sys_config").Where("ckey=?", ckey).FindOneMap()
	utils.ThrowError(err)

	nmp := map[string]interface{}{
		"ckey": ckey,
		"cval": cval,
	}
	if len(resmap) == 0 {
		nmp["c_date"] = utils.Now()
		// 不存在
		_, err := orm.SetTable("t_sys_config").Insert(nmp)
		utils.ThrowError(err)
	} else {
		nmp["a_date"] = utils.Now()
		// 存在
		_, err := orm.SetTable("t_sys_config").SetPK("ckey").UpdateOnly(nmp)
		utils.ThrowError(err)
	}

	global.InitSyscnf()

	return nil, -1
}
