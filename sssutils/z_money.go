package sssutils

import (
	"mycommon/encode"
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"

	"api.lottery.thirdparty/global"

	"sync"
)

const (
	MONEY_STATUS_WAITING  = 1
	MONEY_STATUS_PAYING   = 10
	MONEY_STATUS_FINISHED = 100

	MONEY_FINISH_CODE_SUCCESS = 1
	MONEY_FINISH_CODE_FAILUR  = 2
	MONEY_FINISH_CODE_CANCEL  = 3

	MONEY_REASON_TYPE_ACHIEVEMENT  = 1  // 成就奖励加钱
	MONEY_REASON_TYPE_NORMALTASK   = 2  // 普通任务加钱
	MONEY_REASON_TYPE_RECHARGE     = 3  // 充值加钱
	MONEY_REASON_TYPE_BACKRECHARGE = 4  // 后台充值加钱
	MONEY_REASON_TYPE_BACKPRESENT  = 5  // 后台赠送
	MONEY_REASON_TYPE_SHAREWECHAT  = 6  // 分享朋友圈加钻
	MONEY_REASON_TYPE_CREATECOST   = 10 // 房间创建消费

	MONEY_REASON_ACHIEVEMENT  = "ACHIEVEMENT"  // 成就奖励加钱
	MONEY_REASON_NORMALTASK   = "NORMALTASK"   // 普通任务加钱
	MONEY_REASON_RECHARGE     = "RECHARGE"     // 充值加钱
	MONEY_REASON_BACKRECHARGE = "BACKRECHARGE" // 后台充值加钱
	MONEY_REASON_BACKPRESENT  = "BACKPRESENT"  // 后台奖励
	MONEY_REASON_CREATECOST   = "CREATECOST"   // 房间创建消费

	ROOM_SSS = 1
)

type MoneyRecord struct {
	Id         string `sql:"id"`
	Room       int    `sql:"room"`
	SerialNo   string `sql:"serial_no"`
	Userid     int    `sql:"userid"`
	ReasonType int    `sql:"reason_type"`
	Reason     string `sql:"reason"`
	Flag       int    `sql:"flag"`
	Money      int    `sql:"money"`
	Status     int    `sql:"status"`
	CDate      string `sql:"c_date"`
	FinishCode int    `sql:"finish_code"`
	FinishDate string `sql:"finish_date"`

	M sync.Mutex `sql:"-"`
}

func NewMoneyRecord(usrid int, cdate string, flag int, money int) *MoneyRecord {
	// 获取当前最大流水号
	maxSerialNo := GetMaxSerialNo(ROOM_SSS)
	var curNo string
	if "" == maxSerialNo {
		curNo = fmt.Sprint(ROOM_SSS, "000", utils.NowStr(), "000001")
	} else {
		curNo = mathstr.StrIdentity(maxSerialNo)
	}

	mr := MoneyRecord{
		Id:       encode.UUID(),
		Room:     ROOM_SSS,
		SerialNo: curNo,
		Userid:   usrid,
		CDate:    cdate,
		Flag:     flag,
		Money:    money,
		Status:   MONEY_STATUS_WAITING,
	}

	return &mr
}

// DB 创建充值记录
func (this *MoneyRecord) Create() {
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)

	if this.Reason == "" || this.ReasonType != 0 {
		switch this.ReasonType {

		case MONEY_REASON_TYPE_ACHIEVEMENT:
			this.Reason = MONEY_REASON_ACHIEVEMENT

		case MONEY_REASON_TYPE_NORMALTASK:
			this.Reason = MONEY_REASON_NORMALTASK

		case MONEY_REASON_TYPE_RECHARGE:
			this.Reason = MONEY_REASON_RECHARGE

		case MONEY_REASON_TYPE_BACKRECHARGE:
			this.Reason = MONEY_REASON_BACKRECHARGE

		case MONEY_REASON_TYPE_BACKPRESENT:
			this.Reason = MONEY_REASON_BACKPRESENT

		case MONEY_REASON_TYPE_CREATECOST:
			this.Reason = MONEY_REASON_CREATECOST
		}
	}

	err := orm.SetTable("t_money_record").InsertModel(this)
	utils.ThrowError(err)
}

// DB 应用充值记录，使记录生效
func (this *MoneyRecord) Apply() {
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)

	this.FinishCode = MONEY_FINISH_CODE_SUCCESS

	mp := map[string]interface{}{
		"id":          this.Id,
		"status":      MONEY_STATUS_FINISHED,
		"finish_code": this.FinishCode,
		"finish_date": utils.Now(),
	}
	_, err := orm.SetTable("t_money_record").SetPK("id").UpdateOnly(mp)
	utils.ThrowError(err)
}

// 取消单据
func (this *MoneyRecord) Cancel() {
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)

	this.FinishCode = MONEY_FINISH_CODE_CANCEL

	mp := map[string]interface{}{
		"id":          this.Id,
		"status":      MONEY_STATUS_FINISHED,
		"finish_code": this.FinishCode,
		"finish_date": utils.Now(),
	}
	_, err := orm.SetTable("t_money_record").SetPK("id").UpdateOnly(mp)
	utils.ThrowError(err)
}

func GetMaxSerialNo(room int) string {
	orm := global.GetNewOrm()
	defer global.CloseOrm(orm)

	res, err := orm.SetTable("t_money_record").Select("max(serial_no) max_no").Where("room=?", room).FindOneMap()
	utils.ThrowError(err)

	return fmt.Sprint(res["max_no"])
}
