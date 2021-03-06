package webhandler

import (
	"mycommon/mathstr"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/global/myconst"
	"api.lottery.thirdparty/model"
	"api.lottery.thirdparty/sssutils"
)

type Probability struct{}

func (this *Probability) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "save":
		return this.saveInfo(param)

	case "select":
		return this.selectWhere(param)

	}

	return "no func name:" + do, 999
}

func (this *Probability) selectWhere(vals map[string]interface{}) (interface{}, int) {
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
func (this *Probability) selectSSS() (interface{}, int) {

	orm := global.GetOrm()
	var reslst []model.CtrlCardRate
	err := orm.SetTable("t_card_ctrls").OrderBy("score asc").FindAll(&reslst)
	utils.ThrowError(err)

	return reslst, -1
}
func (this *Probability) selectNN() (interface{}, int) {

	orm := global.GetOrm()
	var reslst []model.CtrlCardRate
	err := orm.SetTable("nn_card_ctrls").OrderBy("score asc").FindAll(&reslst)
	utils.ThrowError(err)

	return reslst, -1
}

func (this *Probability) saveInfo(vals map[string]interface{}) (interface{}, int) {
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

func (this *Probability) saveSSS(vals map[string]interface{}) (interface{}, int) {
	// 接收参数 id agent_id
	details := mathstr.JsonMashDetails(vals["info"])
	utils.S_CRPM(vals, "info")

	if len(details) <= 0 {
		return nil, -1
	}

	orm := global.GetOrm()
	sqlTruncate := `truncate table t_card_ctrls`
	_, err := orm.Exec(sqlTruncate)
	utils.ThrowError(err)

	sqlInsert := `insert into t_card_ctrls (id,score,rate) values`
	for i, e := range details {
		tsql := ""
		if i == 0 {
			tsql = mathstr.S_SFT(`({0},{1},{2})`, e["id"], e["score"], e["rate"])
		} else {
			tsql = mathstr.S_SFT(`,({0},{1},{2})`, e["id"], e["score"], e["rate"])
		}
		sqlInsert += tsql
	}
	sqlInsert += ";"

	_, err = orm.Exec(sqlInsert)
	utils.ThrowError(err)

	return nil, -1
}

func (this *Probability) saveNN(vals map[string]interface{}) (interface{}, int) {
	// 接收参数 id agent_id
	details := mathstr.JsonMashDetails(vals["info"])
	utils.S_CRPM(vals, "info")

	if len(details) <= 0 {
		return nil, -1
	}

	orm := global.GetOrm()
	sqlTruncate := `truncate table nn_card_ctrls`
	_, err := orm.Exec(sqlTruncate)
	utils.ThrowError(err)

	sqlInsert := `insert into nn_card_ctrls (id,score,rate) values`
	for i, e := range details {
		tsql := ""
		if i == 0 {
			tsql = mathstr.S_SFT(`({0},{1},{2})`, e["id"], e["score"], e["rate"])
		} else {
			tsql = mathstr.S_SFT(`,({0},{1},{2})`, e["id"], e["score"], e["rate"])
		}
		sqlInsert += tsql
	}
	sqlInsert += ";"

	_, err = orm.Exec(sqlInsert)
	utils.ThrowError(err)

	return nil, -1
}
