package webhandler

import (
	"mycommon/mathstr"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/global/myconst"
	"api.lottery.thirdparty/model"
)

type Record struct{}

func (this *Record) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "userplay":
		return this.userplay(param)

	case "tabledetail":
		return this.tabledetail(param)

	}

	return "no func name:" + do, 999
}

func (this *Record) userplay(vals map[string]interface{}) (interface{}, int) {
	// 昵称，手机号，微信open_id  钻石余额
	gt := mathstr.Math2intDefault0(vals["gt"])
	uid := mathstr.Math2intDefault0(vals["uid"])
	pageNo := mathstr.Math2intDefault0(vals["pageNo"])
	pageSize := mathstr.Math2intDefault0(vals["pageSize"])
	if 0 == gt {
		gt = myconst.GT_SSS
	}

	switch gt {

	case myconst.GT_SSS:
		return this.userplaySSS(uid, pageNo, pageSize)

	default:
		return this.userplaySSS(uid, pageNo, pageSize)

	}
}
func (this *Record) userplaySSS(uid, pageNo, pageSize int) (interface{}, int) {
	orm := global.GetOrm()
	offset := (pageNo - 1) * pageSize
	var reslst []model.USERPLAY_RECORD_ACCS
	err := orm.SetTable("t_userplay_record ur").
		Join(`left`, "t_accounts_record ar", `ur.accounts_id=ar.id`).
		Select(`ur.*,ar.createtime,ar.tableid,ar.stage,ar.owner_id,ar.type,ar.do_date,ar.round`).
		Where("ur.userid=?", uid).
		Offset(offset).Limit(pageSize).
		OrderBy("c_date desc").FindAll(&reslst)
	utils.ThrowError(err)

	return reslst, -1

}

func (this *Record) tabledetail(vals map[string]interface{}) (interface{}, int) {
	// 昵称，手机号，微信open_id  钻石余额
	ct := mathstr.Math2intDefault0(vals["ct"])
	tableid := mathstr.Math2intDefault0(vals["tableid"])
	pageNo := mathstr.Math2intDefault0(vals["pageNo"])
	pageSize := mathstr.Math2intDefault0(vals["pageSize"])

	orm := global.GetOrm()
	offset := (pageNo - 1) * pageSize
	var reslst []model.SSS_ROOM_RECORD
	err := orm.SetTable("t_userplay_record ur").
		Where("tableid=? and createtime=?", tableid, ct).
		Offset(offset).Limit(pageSize).
		OrderBy("id desc").FindAll(&reslst)
	utils.ThrowError(err)

	return reslst, -1
}
