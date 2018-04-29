package webhandler

import (
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/model"
	"api.lottery.thirdparty/sssutils"
)

type Users struct{}

func (this *Users) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "update":
		return this.update(param)

	case "select":
		return this.selectWhere(param)

	case "selectone":
		return this.selectOne(param)

	}

	return "no func name:" + do, 999
}

func (this *Users) selectWhere(vals map[string]interface{}) (interface{}, int) {
	// 昵称，手机号，微信open_id  钻石余额
	nickname := fmt.Sprint(vals["nick_name"])
	sex := fmt.Sprint(vals["sex"])
	unionid := fmt.Sprint(vals["unionid"])
	wxid := fmt.Sprint(vals["wxid"])
	money := fmt.Sprint(vals["money"])
	ip := fmt.Sprint(vals["ip"])
	id := fmt.Sprint(vals["id"])
	utils.S_CRPM(vals, "pageNo", "pageSize")
	pageNo, pageSize, offset, _ := mathstr.ParamGetPageInfoSql(vals)

	wherestr := " 1=1 "
	if "" != nickname {
		wherestr += mathstr.S_SFT(` and nicknme like '%{0}%' `, nickname)
	}
	if "" != nickname {
		wherestr += mathstr.S_SFT(` and ip like '%{0}%' `, ip)
	}
	if "" != id {
		wherestr += mathstr.S_SFT(` and id like '%{0}%' `, id)
	}
	if "" != sex {
		wherestr += mathstr.S_SFT(` and sex={0} `, sex)
	}
	if "" != unionid {
		wherestr += mathstr.S_SFT(` and unionid like '%{0}%' `, unionid)
	}
	if "" != wxid {
		wherestr += mathstr.S_SFT(` and wxid like '%{0}%' `, wxid)
	}
	if "" != money {
		wherestr += mathstr.S_SFT(` and money like '%{0}%' `, money)
	}

	orm := global.GetOrm()
	var reslst []model.User
	err := orm.SetTable("users").
		Where(wherestr).
		Offset(offset).
		Limit(pageSize).
		FindAll(&reslst)
	utils.ThrowError(err)

	sum, err := orm.SetTable("users").Select("count(1) size,sum(money) total_money").
		Where(wherestr).FindOneMap()
	utils.ThrowError(err)
	mathstr.RoundMap4(sum, "size", "total_money")

	var pageobj model.PageSumObj
	pageobj.PageNo = pageNo
	pageobj.PageSize = pageSize
	pageobj.Size = mathstr.Math2intDefault0(sum["size"])
	pageobj.SumInfo = sum
	pageobj.Vals = reslst

	return pageobj, -1
}

func (this *Users) selectOne(vals map[string]interface{}) (interface{}, int) {
	id := fmt.Sprint(vals["id"])
	utils.S_CRPM(vals, "id")

	orm := global.GetOrm()
	var resobj model.User
	err := orm.SetTable("users").Where("id=?", id).Find(&resobj)
	utils.ThrowError(err)

	return resobj, -1
}

func (this *Users) update(vals map[string]interface{}) (interface{}, int) {
	// 接收参数 id agent_id
	utils.S_CRPM(vals, "id")

	orm := global.GetOrm()

	nmp := map[string]interface{}{}
	mathstr.S_CME(vals, nmp, "id", "agent_id")
	if len(nmp) <= 1 {
		return nil, -1
	}

	_, err := orm.SetTable("users").SetPK("id").UpdateOnly(nmp)
	utils.ThrowError(err)

	go sssutils.SyncRedis(nmp)

	return nil, -1
}
