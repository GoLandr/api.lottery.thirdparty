package webhandler

import (
	"mycommon/logs"
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/global/myconst"
	"api.lottery.thirdparty/model"
	"api.lottery.thirdparty/sssutils"
)

type Customer struct{}

func (this *Customer) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "save":
		return this.saveInfo(param)

	case "select":
		return this.selectWhere(param)

	}

	return "no func name:" + do, 999
}

func (this *Customer) selectWhere(vals map[string]interface{}) (interface{}, int) {
	// 昵称，手机号，微信open_id  钻石余额
	gt := mathstr.Math2intDefault0(vals["gt"])
	if 0 == gt {
		gt = myconst.GT_SSS
	}

	switch gt {

	case myconst.GT_SSS:
		return this.selectSSS()

	case myconst.GT_NN:
		return this.selectNN()

	default:
		return this.selectSSS()

	}
}
func (this *Customer) selectSSS() (interface{}, int) {

	orm := global.GetOrm()
	var res model.GmConfig
	err := orm.SetTable("sss_gm_configs").Find(&res)
	utils.ThrowError(err)

	return res, -1
}
func (this *Customer) selectNN() (interface{}, int) {

	orm := global.GetOrm()
	var res model.GmConfig
	err := orm.SetTable("nn_gm_configs").Find(&res)
	utils.ThrowError(err)

	return res, -1
}

func (this *Customer) saveInfo(vals map[string]interface{}) (interface{}, int) {
	gt := mathstr.Math2intDefault0(vals["gt"])
	if 0 == gt {
		gt = myconst.GT_SSS
	}
	defer func() {
		go sssutils.ReloadCnf(gt)
	}()

	var res interface{}
	var err int
	switch gt {

	case myconst.GT_SSS:
		res, err = this.saveSSS(vals)

	case myconst.GT_NN:
		res, err = this.saveNN(vals)

	default:
		res, err = this.saveSSS(vals)

	}

	return res, err
}

func (this *Customer) saveSSS(vals map[string]interface{}) (interface{}, int) {
	// 接收参数 id agent_id
	// gm_good_card_increase
	// gm_good_card_increace
	orm := global.GetOrm()
	sqlstr := mathstr.S_SFT(`update sss_gm_configs set table_good_card_initial=0`)
	if "" != fmt.Sprint(vals["gm_good_card_initial"]) {
		sqlstr += mathstr.S_SFT(`,gm_good_card_initial={0}`, vals["gm_good_card_initial"])
	}
	if "" != fmt.Sprint(vals["gm_good_card_increace"]) {
		sqlstr += mathstr.S_SFT(`,gm_good_card_increace={0}`, vals["gm_good_card_increace"])
	}
	if "" != fmt.Sprint(vals["normal_card_rate"]) {
		sqlstr += mathstr.S_SFT(`,normal_card_rate={0}`, vals["normal_card_rate"])
	}
	logs.Debug("__saveSSS:", sqlstr)
	_, err := orm.Exec(sqlstr)
	utils.ThrowError(err)

	return nil, -1
}

func (this *Customer) saveNN(vals map[string]interface{}) (interface{}, int) {
	return nil, -1
}
