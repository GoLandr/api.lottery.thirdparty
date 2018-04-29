package wft

import (
	"encoding/json"
	"encoding/xml"
	"mycommon/encode"
	"mycommon/logs"
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"
	"net/http"

	"api.lottery.thirdparty/config"
	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/global/myconst"
	"api.lottery.thirdparty/model"
	"api.lottery.thirdparty/sssutils"
)

type Weifutong struct{}

func (this *Weifutong) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	switch do {

	case "setgm":
		fallthrough
	case "setGM":
		fallthrough
	case "setGm":
		return this.setGM(param)

	case "getgm":
		fallthrough
	case "getGM":
		fallthrough
	case "getGm":
		return this.getGM(param)

	case "getmoney":
		fallthrough
	case "getMoney":
		return this.getMoney(param)

	case "getmoney2":
		fallthrough
	case "getMoney2":
		return this.getMoney2(param)

	case "getmoney3":
		fallthrough
	case "getMoney3":
		return this.getMoney3(param)

	case "updateuserinfo":
		fallthrough
	case "updateUserInfo":
		return this.updateUserInfo(param)

	case "verifybuy":
		return this.verifybuy(param)

	case "submitbuy":
		return this.submitbuy(param)

	case "recordtasktraces":
		fallthrough
	case "recordTaskTraces":
		return this.recordTaskTraces(param)

	case "setgamenotice":
		fallthrough
	case "setGameNotice":
		return this.setGameNotice(param)

	case "getgamenotice":
		fallthrough
	case "getGameNotice":
		return this.getGameNotice(param)

	case "addgamemsg":
		fallthrough
	case "addGameMsg":
		return this.addGameMsg(param)

	case "delgamemsg":
		fallthrough
	case "delGameMsg":
		return this.delGameMsg(param)

	case "updategamemsg":
		fallthrough
	case "updateGameMsg":
		return this.updateGameMsg(param)

	case "getgamemsg":
		fallthrough
	case "getGameMsg":
		return this.getGameMsg(param)

	case "addblack":
		fallthrough
	case "addBlack":
		return this.addBlack(param)

	case "removeblack":
		fallthrough
	case "removeBlack":
		return this.removeBlack(param)

	case "getblack":
		fallthrough
	case "getBlack":
		return this.getBlack(param)

	case "getUserLst":
		return this.getUserLst(param)

	case "getUserOne":
		return this.getUserOne(param)

	}

	return "no func name:" + do, 999
}

func (this *Weifutong) addGameMsg(vals map[string]interface{}) (interface{}, int) {
	content := fmt.Sprint(vals["content"])
	sort := fmt.Sprint(vals["sort"])
	utils.S_CRPM(vals, "content")

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	ump := map[string]interface{}{
		"id":      encode.UUID(),
		"content": content,
		"sort":    sort,
		"c_date":  utils.Now(),
	}
	_, err := orm.SetTable("t_immediate_notice").Insert(ump)
	utils.ThrowError(err)

	return ump["id"], -1
}
func (this *Weifutong) delGameMsg(vals map[string]interface{}) (interface{}, int) {
	id := fmt.Sprint(vals["id"])
	utils.S_CRPM(vals, "id")

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	sqlupdate := mathstr.S_SFT(`delete from t_immediate_notice where id='{0}'`, id)
	_, err := orm.Exec(sqlupdate)
	utils.ThrowError(err)

	return nil, -1
}
func (this *Weifutong) updateGameMsg(vals map[string]interface{}) (interface{}, int) {
	id := fmt.Sprint(vals["id"])
	utils.S_CRPM(vals, "id")

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	ump := map[string]interface{}{
		"id":     id,
		"a_date": utils.Now(),
	}
	mathstr.S_CME(vals, ump, "content", "sort")
	_, err := orm.SetTable("t_immediate_notice").SetPK("id").UpdateOnly(ump)
	utils.ThrowError(err)

	return nil, -1
}
func (this *Weifutong) getGameMsg(vals map[string]interface{}) (interface{}, int) {
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	reslst, err := orm.SetTable("t_immediate_notice").
		Select("id,content,sort,c_date,a_date").OrderBy("sort asc").FindStringInterface()
	utils.ThrowError(err)

	return reslst, -1
}

func (this *Weifutong) setGameNotice(vals map[string]interface{}) (interface{}, int) {
	title := fmt.Sprint(vals["title"])
	content := fmt.Sprint(vals["content"])

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	ump := map[string]interface{}{
		"id":      encode.UUID(),
		"title":   title,
		"content": content,
		"status":  1,
		"c_date":  utils.Now(),
	}
	_, err := orm.SetTable("t_game_notice").SetPK("id").Insert(ump)
	utils.ThrowError(err)

	sqlupdate := mathstr.S_SFT(`delete from t_game_notice where id<>'{0}'`, ump["id"])
	_, err = orm.Exec(sqlupdate)
	utils.ThrowError(err)

	return nil, -1
}
func (this *Weifutong) getGameNotice(vals map[string]interface{}) (interface{}, int) {

	awCnf := config.GetAWInterface()
	if awCnf.Start {
		return getAWNotice(awCnf.GetNoticeUrl), -1
	} else {

		orm := global.GetNewOrm()
		defer global.CloseOrm(orm)
		gameNotice, err := orm.SetTable("t_game_notice").Where(`status=1`).FindOneMap()
		utils.ThrowError(err)

		res := map[string]interface{}{
			"title":   gameNotice["title"],
			"content": gameNotice["content"],
		}

		return res, -1
	}
}

func getAWNotice(getUrl string) map[string]interface{} {
	if "" == getUrl {
		return nil
	}

	param := mathstr.S_SFT(`m={0}`, "notice")
	respStr, err := utils.HttpPost(getUrl, param)
	utils.ThrowError(err)
	logs.Debug("getAWNotice__respstr:", respStr)
	var resp model.AWFlag
	mathstr.JsonUnmarsh(respStr, &resp)

	if !resp.Flag {
		logs.Error("获取返回空，返回值为：", respStr)
		return nil
	}

	if nil == resp.Data {
		logs.Error("对象为空:", respStr)
		return nil
	}

	var noticeMap map[string]interface{}
	mathstr.JsonUnmarshInterface(resp.Data, &noticeMap)

	return noticeMap
}

func (this *Weifutong) addBlack(vals map[string]interface{}) (interface{}, int) {
	uid := fmt.Sprint(vals["uid"])
	wxid := fmt.Sprint(vals["wxid"])
	utils.S_CRPM(vals, "uid")

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	emap, err := orm.SetTable("t_black").Where("user_id=?", uid).FindOneMap()
	utils.ThrowError(err)
	if "" != fmt.Sprint(emap["id"]) {
		return 1, -1
	}

	mp := map[string]interface{}{
		"user_id": uid,
		"wxid":    wxid,
		"id":      encode.UUID(),
	}
	_, err = orm.SetTable("t_black").Insert(mp)
	utils.ThrowError(err)

	return 1, -1
}
func (this *Weifutong) removeBlack(vals map[string]interface{}) (interface{}, int) {
	uid := fmt.Sprint(vals["uid"])
	wxid := fmt.Sprint(vals["wxid"])
	utils.S_CRPM(vals, "uid")

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)

	sqlstr := mathstr.S_SFT(`delete from t_black where user_id='{0}' or wxid='{1}'`, uid, wxid)
	if wxid == "" {
		sqlstr = mathstr.S_SFT(`delete from t_black where user_id='{0}'`, uid)
	}
	_, err := orm.Exec(sqlstr)
	utils.ThrowError(err)

	return 1, -1
}
func (this *Weifutong) getBlack(vals map[string]interface{}) (interface{}, int) {
	utils.S_CRPM(vals)

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	var reslst []model.User
	err := orm.SetTable("t_black b").Join("left", "users u", "b.user_id=u.id").Select("u.*").FindAll(&reslst)
	utils.ThrowError(err)

	return reslst, -1
}

func (this *Weifutong) setGM(vals map[string]interface{}) (interface{}, int) {
	uid := fmt.Sprint(vals["uid"])
	set := mathstr.Math2intDefault0(vals["set"])
	utils.S_CRPM(vals, "uid", "set")

	level := 0
	if 1 == set {
		level = 99
	}

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	ump := map[string]interface{}{
		"id":    uid,
		"level": level,
	}
	_, err := orm.SetTable("users").SetPK("id").UpdateOnly(ump)
	utils.ThrowError(err)
	go sssutils.SyncRedis(ump)

	go this.LocalSyncGM(mathstr.Math2intDefault0(uid), level)

	return 1, -1
}
func (this *Weifutong) LocalSyncGM(uid int, level int) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	server := config.GetLogicServer() + "/localSyncGM"
	param := map[string]interface{}{
		"uid":   uid,
		"level": level,
		"auth":  config.GetLocalauth(),
	}
	strParam := mathstr.GetJsonPlainStr(param)
	headerMap := map[string]string{
		"Connection": "close",
	}
	resp, err := utils.HttpPostWithHeader(server, strParam, headerMap)
	utils.ThrowError(err)
	logs.Debug("__resp:", resp)
}

func (this *Weifutong) getGM(vals map[string]interface{}) (interface{}, int) {
	key := fmt.Sprint(vals["key"])
	utils.S_CRPM(vals, "pageSize", "pageNo")

	pageNo, pageSize, offset := mathstr.ParamGetPageInfo(vals)

	wherestr := ` level=99 `
	if "" != key {
		wherestr += mathstr.S_SFT(` and (nickname like '%{0}%' or id like '%{0}%') `, key)
	}

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	var userlst []model.User
	err := orm.SetTable("users").Where(wherestr).Offset(offset).Limit(pageSize).FindAll(&userlst)
	utils.ThrowError(err)

	sum, err := orm.SetTable("users").Where(wherestr).Select("count(1) size").FindOneMap()
	utils.ThrowError(err)

	var pageObj model.PageObj
	pageObj.PageNo = pageNo
	pageObj.PageSize = pageSize
	pageObj.Size = mathstr.Math2intDefault0(sum["size"])
	pageObj.Vals = userlst

	return pageObj, -1
}

func (this *Weifutong) getMoney(vals map[string]interface{}) (interface{}, int) {
	uid := fmt.Sprint(vals["uid"])
	money := mathstr.Math2float64Default0(vals["money"])
	gameKind := fmt.Sprint(vals["gameKind"])
	utils.S_CRPM(vals, "uid", "money")

	if "" == gameKind {
		gameKind = "1"
	}

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	usr, err := orm.SetTable("users").Where("id=?", uid).FindOneMap()
	utils.ThrowError(err)
	if "" == fmt.Sprint(usr["id"]) {
		return nil, myconst.WRONG_USER
	}

	sqlupdate := mathstr.S_SFT(`
		update users set money=money+{0} where id={1}
		`, money, uid)
	_, err = orm.Exec(sqlupdate)
	utils.ThrowError(err)

	now := utils.NowTimeObj().Unix() / 1000
	paysign := fmt.Sprint("&&addByAdmin_", now)
	payRecord := map[string]interface{}{}
	payRecord["uid"] = usr["id"]
	payRecord["pay"] = money
	payRecord["at"] = now
	payRecord["paysign"] = paysign
	_, err = orm.SetTable("pay_records").Insert(payRecord)
	utils.ThrowError(err)

	// 通知逻辑服
	go this.AddMoney2LogicServer(mathstr.Math2intDefault0(usr["id"]), money)

	return 1, -1
}
func (this *Weifutong) getMoney2(vals map[string]interface{}) (interface{}, int) {
	wxid := fmt.Sprint(vals["uid"]) // 微信id
	money := mathstr.Math2float64Default0(vals["money"])
	oid := fmt.Sprint(vals["oid"])
	gameKind := fmt.Sprint(vals["gameKind"])
	utils.S_CRPM(vals, "uid", "money", "oid")

	if "" == gameKind {
		gameKind = "1"
	}

	paysign := fmt.Sprint("&&wanjun_", oid)
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	ePayRec, err := orm.SetTable("pay_records").Where("paysign=?", paysign).FindOneMap()
	utils.ThrowError(err)
	if "" != fmt.Sprint(ePayRec["id"]) {
		return nil, myconst.ORDER_REPEAT
	}

	usr, err := orm.SetTable("users").Where("wxid=?", wxid).FindOneMap()
	utils.ThrowError(err)
	if "" == fmt.Sprint(usr["id"]) {
		return nil, myconst.WRONG_USER
	}

	sqlupdate := mathstr.S_SFT(`
		update users set money=money+{0} where id={1}
		`, money, usr["id"])
	_, err = orm.Exec(sqlupdate)
	utils.ThrowError(err)

	go MoneyChangeHandler(mathstr.Math2intDefault0(usr["id"]), int(money), mathstr.Math2intDefault0(usr["money"]))

	now := utils.NowTimeObj().Unix() / 1000
	payRecord := map[string]interface{}{}
	payRecord["uid"] = usr["id"]
	payRecord["pay"] = money
	payRecord["at"] = now
	payRecord["paysign"] = paysign
	_, err = orm.SetTable("pay_records").Insert(payRecord)
	utils.ThrowError(err)

	// 通知逻辑服
	go this.AddMoney2LogicServer(mathstr.Math2intDefault0(usr["id"]), money)

	return 1, -1
}

func (this *Weifutong) getMoney3(vals map[string]interface{}) (interface{}, int) {
	uid := fmt.Sprint(vals["uid"]) // 用户id，自增那个
	money := mathstr.Math2float64Default0(vals["money"])
	oid := fmt.Sprint(vals["oid"])
	gameKind := fmt.Sprint(vals["gameKind"])
	utils.S_CRPM(vals, "uid", "money", "oid")

	if "" == gameKind {
		gameKind = "1"
	}

	paysign := fmt.Sprint("&&wanjun_", oid)
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	ePayRec, err := orm.SetTable("pay_records").Where("paysign=?", paysign).FindOneMap()
	utils.ThrowError(err)
	if "" != fmt.Sprint(ePayRec["id"]) {
		return nil, myconst.ORDER_REPEAT
	}

	usr, err := orm.SetTable("users").Where("id=?", uid).FindOneMap()
	utils.ThrowError(err)
	if "" == fmt.Sprint(usr["id"]) {
		return nil, myconst.WRONG_USER
	}

	sqlupdate := mathstr.S_SFT(`
		update users set money=money+{0} where id={1}
		`, money, usr["id"])
	_, err = orm.Exec(sqlupdate)
	utils.ThrowError(err)

	go MoneyChangeHandler(mathstr.Math2intDefault0(uid), int(money), mathstr.Math2intDefault0(usr["money"]))

	now := utils.NowTimeObj().Unix()
	payRecord := map[string]interface{}{}
	payRecord["uid"] = usr["id"]
	payRecord["pay"] = money
	payRecord["at"] = now
	payRecord["paysign"] = paysign
	_, err = orm.SetTable("pay_records").Insert(payRecord)
	utils.ThrowError(err)

	// 通知逻辑服
	go this.AddMoney2LogicServer(mathstr.Math2intDefault0(usr["id"]), money)

	return 1, -1
}
func MoneyChangeHandler(uid int, moneyChange int, moneyBase int) {
	sssutils.SyncRedis(map[string]interface{}{
		"id":    uid,
		"money": moneyBase + moneyChange,
	})
	if config.IsHall() {
		// 大厅模式下，要进行记录新增
		moneyRec := sssutils.NewMoneyRecord(uid, utils.Now(), 1, moneyChange)
		moneyRec.ReasonType = sssutils.MONEY_REASON_TYPE_RECHARGE
		moneyRec.Create()
		moneyRec.Apply()
	}
}

func (this *Weifutong) AddMoney2LogicServer(uid int, moneyflo float64) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	money := int(moneyflo)

	server := config.GetLogicServer() + "/localSyncMoney"
	param := map[string]interface{}{
		"uid":   uid,
		"money": money,
		"auth":  config.GetLocalauth()}
	strParam := mathstr.GetJsonPlainStr(param)
	headerMap := map[string]string{
		"Connection": "close",
	}
	resp, err := utils.HttpPostWithHeader(server, strParam, headerMap)
	utils.ThrowError(err)
	logs.Debug("__resp:", resp)
}

func (this *Weifutong) updateUserInfo(vals map[string]interface{}) (interface{}, int) {
	uid := fmt.Sprint(vals["uid"])
	agentid := fmt.Sprint(vals["agent_id"])
	utils.S_CRPM(vals, "uid", "agent_id")

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	tmp := map[string]interface{}{
		"id":       uid,
		"agent_id": agentid,
	}
	_, err := orm.SetTable("users").SetPK("id").UpdateOnly(tmp)
	utils.ThrowError(err)

	go sssutils.SyncRedis(tmp)

	return nil, -1
}

func (this *Weifutong) verifybuy(vals map[string]interface{}) (interface{}, int) {
	outTradeNo := fmt.Sprint(vals["out_trade_no"])
	utils.S_CRPM(vals, "out_trade_no")

	paySign := "&&weipay_" + outTradeNo
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	payRecord, err := orm.SetTable("pay_records").Where("paysign=?", paySign).FindOneMap()
	utils.ThrowError(err)
	if "" == fmt.Sprint(payRecord["id"]) {
		return nil, myconst.SQL_ERROR
	}

	res := map[string]interface{}{
		"addedMoney": payRecord["pay"],
	}

	return res, -1
}

func (this *Weifutong) submitbuy(vals map[string]interface{}) (interface{}, int) {
	outTradeNo := fmt.Sprint(vals["out_trade_no"])
	total_fee := fmt.Sprint(vals["total_fee"])
	body := fmt.Sprint(vals["body"])
	attach := fmt.Sprint(vals["attach"])
	mch_create_ip := fmt.Sprint(vals["mch_create_ip"])
	utils.S_CRPM(vals, "out_trade_no", "total_fee", "body")

	attachJson := mathstr.GetJsonPlainStr(attach)

	paramsMap := map[string]interface{}{
		"mch_id":           config.GetMchid(),
		"version":          config.GetVersion(),
		"service":          "unified.trade.pay",
		"out_trade_no":     outTradeNo,
		"body":             body,
		"attach":           attachJson,
		"total_fee":        total_fee,
		"mch_create_ip":    mch_create_ip,
		"notify_url":       config.GetWxzfNotifyUrl(),
		"limit_credit_pay": 1,
		"nonce_str":        mathstr.RandChar(13),
	}
	paramKeys := []string{"mch_id", "version", "service", "out_trade_no",
		"body", "attach", "total_fee", "mch_create_ip", "notify_url",
		"limit_credit_pay", "nonce_str"}
	sign := sssutils.GenerageWXZFSign(paramsMap, paramKeys, config.GetWxzfSignkey())
	params := model.XmlReq{
		Mch_id:           config.GetMchid(),
		Version:          config.GetVersion(),
		Service:          "unified.trade.pay",
		Out_trade_no:     outTradeNo,
		Body:             body,
		Attach:           attachJson,
		Total_fee:        total_fee,
		Mch_create_ip:    mch_create_ip,
		Notify_url:       config.GetWxzfNotifyUrl(),
		Limit_credit_pay: 1,
		Nonce_str:        mathstr.RandChar(13),
		Sign:             sign,
	}

	xmlstr := mathstr.GetXmlPlainStr(params)
	resp, err := utils.HttpPost(config.GetWxzfRequestUrl(), xmlstr)
	utils.ThrowError(err)

	// 将resp转换为xml格式
	var respMap model.XmlResp
	err = xml.Unmarshal([]byte(resp), &respMap)
	utils.ThrowError(err)

	if 0 == len(respMap.Status) {
		return nil, myconst.WRONG_ON_PAY
	}

	if "0" == respMap.Status[0] {
		return respMap, -1
	}

	return nil, myconst.WRONG_ON_PAY
}

// text 返回
func (this *Weifutong) BuyCallback(resp *model.XmlResp) string {
	if len(resp.Status) <= 0 || len(resp.ResultCode) <= 0 {
		return "fail"
	}
	if resp.Status[0] != "0" || "0" != resp.ResultCode[0] {
		return "fail"
	}

	attach := resp.Attach
	var attMap map[string]interface{}
	err := json.Unmarshal([]byte(attach), &attMap)
	utils.ThrowError(err)
	uid := fmt.Sprint(attMap["uid"])
	sku := fmt.Sprint(attMap["sku"])
	outTradeNo := fmt.Sprint(attMap["out_trade_no"])
	code := this.onReceipt(uid, sku, outTradeNo)
	if 0 == code {
		return "success"
	}
	return "fail"
}

func (this *Weifutong) onReceipt(uid string, sku string, orderno string) int {
	// 验证是否重复订单
	var paySign = "&&weipay_" + orderno
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	payRecord, err := orm.SetTable("pay_records").Where("paysign=?", paySign).FindOneMap()
	utils.ThrowError(err)

	if "" != fmt.Sprint(payRecord) {
		return 1
	}

	stoConfMap := global.GetStoreConfig()
	eConf := stoConfMap[sku]
	if "" == eConf.Sku {
		return 1
	}
	addMoney := eConf.AddMoney
	usr, err := orm.SetTable("users").Where("id=?", uid).FindOneMap()
	utils.ThrowError(err)
	if "" == fmt.Sprint(usr["id"]) {
		return 1
	}

	// 加钱
	sqlstr := mathstr.S_SFT(`
		update user set money=money+{0} where id={1}
		`, addMoney, uid)
	_, err = orm.Exec(sqlstr)
	utils.ThrowError(err)

	go MoneyChangeHandler(mathstr.Math2intDefault0(usr["id"]), addMoney, mathstr.Math2intDefault0(usr["money"]))

	// 加记录
	newPayRec := map[string]interface{}{
		"uid":     uid,
		"pay":     addMoney,
		"at":      utils.NowTimeObj().Unix() / 1000,
		"paysign": paySign,
	}
	_, err = orm.SetTable("pay_records").Insert(newPayRec)
	utils.ThrowError(err)

	return 0
}

func (this *Weifutong) recordTaskTraces(vals map[string]interface{}) (interface{}, int) {
	uid := mathstr.Math2intDefault0(vals["uid"])
	t := mathstr.Math2intDefault0(vals["t"])
	utils.S_CRPM(vals, "uid", "t")

	taskConflstMap := global.GetTaskConfig()
	taskCfg := taskConflstMap[t]
	if len(taskCfg) <= 0 {
		return nil, myconst.WRONG_PARAM
	}

	taskMemRecord := global.GetTaskMemRecords()
	memTrace := taskMemRecord.GetMemRecord(uid)
	if nil == memTrace {
		orm := global.GetNewOrm()
		defer global.CloseOrm(orm)
		taskRecTmp, err := orm.SetTable("user_task_records").Where("uid=?", uid).FindOneMap()
		utils.ThrowError(err)
		taskRec := map[string]int{}
		for k, v := range taskRecTmp {
			intv := mathstr.Math2intDefault0(v)
			taskRec[k] = intv
		}
		if "" != fmt.Sprint(taskRec["id"]) {
			taskMemRecord.UpdateMemRecord(uid, taskRec)
			memTrace = taskMemRecord.GetMemRecord(uid)
			return this.handleTask(t, uid, memTrace)
		}

		return nil, 709
	}

	return this.handleTask(t, uid, memTrace)
}
func (this *Weifutong) handleTask(tInt int, uid int, memTrace map[string]int) (interface{}, int) {
	t := fmt.Sprint(tInt)
	var k = "t" + t
	var kget = "t" + t + "get"
	var ktar = "t" + t + "target"
	if memTrace[kget] == 0 {
		if memTrace[k] == memTrace[ktar] {
			return nil, -1
		} else {
			memTrace[k] += 1
			if memTrace[k] > memTrace[ktar] {
				memTrace[k] = memTrace[ktar]
			} else {
				// update db now
				updateParam := map[string]interface{}{
					k:     memTrace[k],
					"uid": uid,
				}

				orm := global.GetNewOrm()
				defer global.CloseOrm(orm)
				_, err := orm.SetTable("user_task_records").SetPK("uid").UpdateOnly(updateParam)
				utils.ThrowError(err)

				//	db.user_task_records.update(updateParam,{where:{uid:uid}}).then(function(){})
			}
			return nil, -1
		}
	}

	return nil, -1
}

func (this *Weifutong) getUserLst(vals map[string]interface{}) (interface{}, int) {
	key := fmt.Sprint(vals["key"])
	ip := fmt.Sprint(vals["ip"])
	agent_id := fmt.Sprint(vals["agent_id"])
	id := fmt.Sprint(vals["id"])

	RsTime := fmt.Sprint(vals["RsTime"]) // 2016-01-02
	ReTime := fmt.Sprint(vals["ReTime"]) // 2016-01-02

	if "" == ReTime {
		ReTime = utils.NowDate()
	}
	if RsTime > ReTime {
		return "起始时间必须小于结束时间", 20
	}

	// 数据格式检测

	orderBy := fmt.Sprint(vals["order"])

	utils.S_CRPM(vals, "pageSize", "pageNo")

	pageNo, pageSize, offset := mathstr.ParamGetPageInfo(vals)

	wherestr := ` 1=1 `
	if "" != key {
		wherestr += mathstr.S_SFT(` and (nickname like '%{0}%' or id like '%{0}%') `, key)
	}
	if "" != ip {
		wherestr += mathstr.S_SFT(` and (ip like '%{0}%') `, ip)
	}
	if "" != agent_id {
		wherestr += mathstr.S_SFT(` and ifnull(agent_id,0)={0} `, agent_id)
	}
	if "" != id {
		wherestr += mathstr.S_SFT(` and id={0} `, id)
	}
	if "" != RsTime {
		startTimeObj := utils.GetTimeFromStr(RsTime + " 00:00:00")
		RsTime := startTimeObj.Unix()
		wherestr += mathstr.S_SFT(` and registertime>={0} `, RsTime)
	}
	if "" != ReTime {
		endTimeObj := utils.GetTimeFromStr(ReTime + " 23:59:59")
		ReTime := endTimeObj.Unix()
		wherestr += mathstr.S_SFT(` and registertime<={0} `, ReTime)
	}
	// total_play_cnt id money registertime todaywinlose
	if "" == orderBy {
		//		orderBy = "total_play_cnt desc"
		orderBy = "registertime desc"
	} else if "total_play_cnt" == orderBy {
		orderBy = "total_play_cnt desc"
	} else if "id" == orderBy {
		orderBy = "id asc"
	} else if "money" == orderBy {
		orderBy = "money desc"
	} else if "registertime" == orderBy {
		orderBy = "registertime desc"
	} else if "todaywinlose" == orderBy {
		orderBy = "todaywinlose desc"
	} else {
		orderBy = "total_play_cnt desc"
	}

	logs.Debug("__where:", wherestr)

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	sqlstr := mathstr.S_SFT(`
		SELECT u.id,u.wxid,u.unionid,u.iconurl,u.nickname,u.money,u.registertime,u.agent_id,u.lastlogintime,b.jdCount,b.sjcount
		from users u LEFT JOIN bb_playerroom b on u.id=b.userID     
		where {0}
		order by {1}
		limit {2},{3}
		`, wherestr, orderBy, offset, pageSize)
	logs.Debug("__sqlstr:", sqlstr)
	objLst, err := orm.QueryMap(sqlstr)
	utils.ThrowError(err)

	var userlst []model.User
	for _, e := range objLst {
		tusr := model.User{}
		tusr.Id = mathstr.Math2intDefault0(e["id"])
		tusr.Ip = fmt.Sprint(e["ip"])
		tusr.Wxid = fmt.Sprint(e["wxid"])
		tusr.Unionid = fmt.Sprint(e["unionid"])
		tusr.Iconurl = fmt.Sprint(e["iconurl"])
		tusr.Nickname = fmt.Sprint(e["nickname"])
		tusr.Money = mathstr.Math2intDefault0(e["money"])
		tusr.Registertime = mathstr.Math2intDefault0(e["registertime"])
		tusr.Agent_id = mathstr.Math2intDefault0(e["agent_id"])
		tusr.Lastlogintime = mathstr.Math2intDefault0(e["lastlogintime"])
		tusr.Jdcount = mathstr.Math2intDefault0(e["jdcount"])
		tusr.Sjcount = mathstr.Math2intDefault0(e["sjcount"])
		userlst = append(userlst, tusr)
	}

	//	err := orm.SetTable("users").Where(wherestr).OrderBy(orderBy).Offset(offset).Limit(pageSize).FindAll(&userlst)
	//	utils.ThrowError(err)
	//	sum, err := orm.SetTable("users").Where(wherestr).Select("count(1) size").FindOneMap()
	//	utils.ThrowError(err)

	sqlsum := mathstr.S_SFT(`
		SELECT count(*)
		from users u LEFT JOIN bb_playerroom b on u.id=b.userID     
		where {0}
		`, wherestr, orderBy, pageSize, offset)
	logs.Debug("__sqlsum:", sqlsum)
	sum, err := orm.QueryOne(sqlsum)
	utils.ThrowError(err)

	var pageObj model.PageObj
	pageObj.PageNo = pageNo
	pageObj.PageSize = pageSize
	pageObj.Size = mathstr.Math2intDefault0(sum["size"])
	pageObj.Vals = userlst

	return pageObj, -1
}

func (this *Weifutong) getUserOne(vals map[string]interface{}) (interface{}, int) {
	wxid := fmt.Sprint(vals["wxid"])
	id := fmt.Sprint(vals["id"])

	if "" == wxid && "" == id {
		return "params error", myconst.WRONG_PARAM
	}

	wherestr := ` 1=1 `
	if "" != wxid {
		wherestr += mathstr.S_SFT(` and wxid='{0}' `, wxid)
	}
	if "" != id {
		wherestr += mathstr.S_SFT(` and id={0} `, id)
	}

	logs.Debug("__where:", wherestr)

	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)
	var user model.User
	err := orm.SetTable("users").Where(wherestr).Find(&user)
	utils.ThrowError(err)

	return user, -1
}
