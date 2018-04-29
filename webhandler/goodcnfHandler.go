package webhandler

import (
	"fmt"
	"mycommon/mathstr"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/global/myconst"
	"api.lottery.thirdparty/model"
	"api.lottery.thirdparty/sssutils"
)

type GoodCnf struct{}

func (this *GoodCnf) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "save":
		return this.saveInfo(param)

	case "select":
		return this.getInfo(param)

	}

	return "no func name:" + do, 999
}

func (this *GoodCnf) getInfo(vals map[string]interface{}) (interface{}, int) {
	// 昵称，手机号，微信open_id  钻石余额
	gt := mathstr.Math2intDefault0(vals["gt"])
	if 0 == gt {
		gt = myconst.GT_SSS
	}

	switch gt {

	case myconst.GT_SSS:
		return this.selectSSS()

	default:
		return this.selectSSS()

	}
}
func (this *GoodCnf) selectSSS() (interface{}, int) {

	orm := global.GetOrm()
	var res model.GoodCnfRate
	err := orm.SetTable("sss_configs").Find(&res)
	utils.ThrowError(err)

	return res, -1
}

func (this *GoodCnf) saveInfo(vals map[string]interface{}) (interface{}, int) {
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

	default:
		res, err = this.saveSSS(vals)

	}

	return res, err
}

func (this *GoodCnf) saveSSS(vals map[string]interface{}) (interface{}, int) {
	// 接收参数 id agent_id
	five_of_akind := fmt.Sprint(vals["five_of_akind"])
	four_of_akind := fmt.Sprint(vals["four_of_akind"])
	flush_straight := fmt.Sprint(vals["flush_straight"])
	all_straight := fmt.Sprint(vals["all_straight"])

	utils.S_CRPM(vals)

	nmp := map[string]interface{}{}
	if "" != five_of_akind {
		nmp["five_of_akind"] = five_of_akind
	}
	if "" != four_of_akind {
		nmp["four_of_akind"] = four_of_akind
	}
	if "" != flush_straight {
		nmp["flush_straight"] = flush_straight
	}
	if "" != all_straight {
		nmp["all_straight"] = all_straight
	}

	if len(nmp) > 0 {
		orm := global.GetOrm()
		_, err := orm.SetTable("sss_configs").UpdateOnly(nmp)
		utils.ThrowError(err)
	}

	return nil, -1
}
