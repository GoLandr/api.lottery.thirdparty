package model

import "encoding/xml"

//彩票统计
type LotteryStat struct {
	Odd_Even  map[int]*Limit //1单 2双
	Big_Small map[int]*Limit //1大 2小
	Star      *Limit         //五星
	Num       *Limit
}

type Limit struct {
	Odd   int //单
	Even  int //双
	Big   int //大
	Small int //小
}
type StarsLimt struct {
	Open int //开出
	No   int //未开
}

type FlagObj struct {
	Flag int         `json:"flag"`
	Msg  string      `json:"msg"`
	Val  interface{} `json:"val"`
}
type FlagObjOfWft struct {
	Flag int         `json:"error"`
	Msg  string      `json:"msg"`
	Val  interface{} `json:"ok"`
}

//wxid,unionid,iconurl,nickname,money,registertime,agent_id,lastlogintime,todaywinlose
type User struct {
	Id            int    `json:"id" sql:"id"`
	Wxid          string `json:"wxid" sql:"wxid"`
	Unionid       string `json:"unionid" sql:"unionid"`
	Iconurl       string `json:"iconurl" sql:"iconurl"`
	Money         int    `json:"money" sql:"money"`
	Ip            string `json:"ip" sql:"ip"`
	Nickname      string `json:"nickname" sql:"nickname"`
	Registertime  int    `json:"registertime" sql:"registertime"`
	Agent_id      int    `json:"agent_id" sql:"agent_id"`
	Lastlogintime int    `json:"lastlogintime" sql:"lastlogintime"`
	TotalPlayCnt  int    `json:"total_play_cnt" sql:"total_play_cnt"`
	TotalPlaySort int    `json:"total_play_sort" sql:"-"`
	Jdcount       int    `json:"Jdcount" sql:"-"`
	Sjcount       int    `json:"Sjcount" sql:"-"`
	//	Todaywinlose  string `json:"todaywinlose" sql:"todaywinlose"`
}
type CtrlCardRate struct {
	Id    string  `json:"id" sql:"id"`
	Score int     `json:"score" sql:"score"`
	Rate  float64 `json:"rate" sql:"rate"`
}
type GoodCnfRate struct {
	Five_of_akind  float64 `json:"five_of_akind" sql:"five_of_akind"`
	Four_of_akind  float64 `json:"four_of_akind" sql:"four_of_akind"`
	Flush_straight float64 `json:"flush_straight" sql:"flush_straight"`
	All_straight   float64 `json:"all_straight" sql:"all_straight"`
}
type GmConfig struct {
	GmGoodCardInitial  int `json:"gm_good_card_initial" sql:"gm_good_card_initial"`
	GmGoodCardIncreace int `json:"gm_good_card_increace" sql:"gm_good_card_increace"`
	NormalCardRate     int `json:"normal_card_rate" sql:"normal_card_rate"`
}

type StoreConfigs struct {
	Sku       string `json:"sku" sql:"sku"`
	Name      string `json:"name" sql:"name"`
	NeedMoney int    `json:"needMoney" sql:"needMoney"`
	AddMoney  int    `json:"addMoney" sql:"addMoney"`
}

type PageObj struct {
	Size     int         `json:"size"`
	PageNo   int         `json:"pageNo"`
	PageSize int         `json:"pageSize"`
	Vals     interface{} `json:"vals"`
}

type PageSumObj struct {
	PageObj
	SumInfo map[string]interface{} `json:"sumObj"`
}

// 微信支付
type XmlReq struct {
	XMLName          xml.Name
	Mch_id           string `xml:"mch_id"`
	Version          string `xml:"version"`
	Service          string `xml:"service"`
	Out_trade_no     string `xml:"out_trade_no"`
	Body             string `xml:"body"`
	Attach           string `xml:"attach"`
	Total_fee        string `xml:"total_fee"`
	Mch_create_ip    string `xml:"mch_create_ip"`
	Notify_url       string `xml:"notify_url"`
	Limit_credit_pay int    `xml:"limit_credit_pay"`
	Nonce_str        string `xml:"nonce_str"`
	Sign             string `xml:"sign"`
}
type XmlResp struct {
	XMLName    xml.Name
	Status     []string `xml:"status"`
	ResultCode []string `xml:"result_code"`
	Attach     string   `xml:"attach"`
}

type StoreConfig struct {
	Sku       string `json:"sku" sql:"sku"`
	Name      string `json:"name" sql:"name"`
	NeedMoney int    `json:"needMoney" sql:"needMoney"`
	AddMoney  int    `json:"addMoney" sql:"addMoney"`
}

type TaskConfig struct {
	Id    int `json:"sku",sql:"sku"`
	Type  int `json:"type",sql:"type"`
	Cnt   int `json:"cnt",sql:"cnt"`
	Prize int `json:"prize",sql:"prize"`
}

type USERPLAY_RECORD_ACCS struct {
	UsrplayRecId string `sql:"column:id" json:"id"`
	Userid       int    `sql:"column:userid"`
	Ismaster     int    `sql:"column:ismaster"`
	Score        int    `sql:"column:score"`
	AccountsId   string `sql:"column:accounts_id"`
	Cdate        string `sql:"column:c_date"`
	Adate        string `sql:"column:a_date"`

	Createtime int    `sql:"column:createtime" json:"ct"`
	Tableid    int    `sql:"column:tableid"`
	Stage      string `sql:"column:stage"`
	Ownerid    int    `sql:"column:owner_id" json:"owner_id"`
	Type       int    `sql:"column:type"`
	Dodate     string `sql:"column:do_date"`
}
type SSS_ROOM_RECORD struct {
	Id          int `sql:"id"`
	Roomid      int
	Tableid     int
	RoundDetail string // column name will be `round_detail`
	Playtime    int
	SPlaytime   string `sql:"column:play_time"`
	Createtime  int
	SCreatetime string `sql:"column:create_time"`
	Owner_id    int    `sql:"owner_id"`
}
