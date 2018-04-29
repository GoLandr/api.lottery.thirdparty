package config

import (
	"encoding/json"
	"log"
	"mycommon/logs"
	"mycommon/myfile"
	fmt "mycommon/myinherit/myfmt"
	"os"

	"gopkg.in/redis.v3"
)

var cnfListen int
var cnfLogs ConfLogs
var cnfMysql ConfMysql
var cnfRedis ConfRedis
var cnfAll ConfWebapi

var shpath string
var logicServer string
var localauth string
var wxzfsignkey string
var wxzfNotifyUrl string
var requestUrl string
var mchid string
var version string

type ConfWebapi struct {
	Listen        int           `json:"listen"`
	Logs          ConfLogs      `json:"logs"`
	Mysql         ConfMysql     `json:"mysql"`
	Redis         ConfRedis     `json:"redis"`
	LogicServer   string        `json:"logicServer"`
	Localauth     string        `json:"localauth"`
	Wxzfsignkey   string        `json:"wxzfsignkey"`
	WxzfNotifyUrl string        `json:"wxzfNotifyUrl"`
	requestUrl    string        `json:"requestUrl"`
	Mch_id        string        `json:"mch_id"`
	Version       string        `json:"version"`
	ShPath        string        `json:"shpath"`  // 调用脚本文件夹路径
	IsHall        bool          `json:"is_hall"` // 是否大厅模式
	AWSyncCnf     ConfInterface `json:"AWSyncCnf"`
}

/*
	对接接口相关,目前适用于
	1. 爱玩十三水
*/
type ConfInterface struct {
	Start           bool   `json:"start"`
	AddUsrUrl       string `json:"addUsrUrl"`
	GetUsrUrl       string `json:"getUsrUrl"`
	GetNoticeUrl    string `json:"getNoticeUrl"`
	GetTumpetMsgUrl string `json:"getTumpetMsgUrl"`
	MinusMoneyUrl   string `json:"minusMoneyUrl"`
}

type ConfLogs struct {
	Dir      string `json:"dir"`
	File     string `json:"file"`
	Level    int    `json:"level"`
	Savefile bool   `json:"savefile"`
}

type ConfMysql struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
	Db     string `json:"db"`
}

type ConfRedis struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Passwd string `json:"passwd"`
}

/*
 * 初始化配置
 */
func init() {
	if len(os.Args) > 1 {
		return
	}
	var conf ConfWebapi
	str := myfile.ReadConfFile("api-lottery-thirdparty.json")
	log.Println(str)
	err := json.Unmarshal([]byte(str), &conf)
	fmt.Println("--api.lottery.3party config init-err:", err)

	if nil == err {
		cnfListen = conf.Listen
		cnfLogs = conf.Logs
		cnfMysql = conf.Mysql
		cnfRedis = conf.Redis

		shpath = conf.ShPath
		logicServer = conf.LogicServer
		localauth = conf.Localauth
		wxzfsignkey = conf.Wxzfsignkey
		wxzfNotifyUrl = conf.WxzfNotifyUrl
		mchid = conf.Mch_id
		version = conf.Version

		cnfAll = conf
	} else {
		fmt.Println(err.Error())
		cnfListen = 443
	}

	if nil == err {
		logs.Init(myfile.GetConfPath(conf.Logs.Dir),
			conf.Logs.File, conf.Logs.Level, conf.Logs.Savefile)
	} else {
		logs.Init(myfile.GetExePath("../log/"), "api-lottery-3pty-",
			logs.LOG_DEBUG, false)
	}

}

func GetListen() int {
	return cnfListen
}
func GeShPath() string {
	return shpath
}
func GetLogicServer() string {
	return logicServer
}
func GetLocalauth() string {
	return localauth
}
func GetWxzfSignkey() string {
	return wxzfsignkey
}
func GetWxzfRequestUrl() string {
	return requestUrl
}
func GetMchid() string {
	return mchid
}
func GetVersion() string {
	return version
}
func GetWxzfNotifyUrl() string {
	return wxzfNotifyUrl
}
func GetAWInterface() ConfInterface {
	return cnfAll.AWSyncCnf
}
func IsHall() bool {
	return cnfAll.IsHall
}

func GetMysqlDriver() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		cnfMysql.User, cnfMysql.Passwd, cnfMysql.Ip, cnfMysql.Port, cnfMysql.Db)
}

func GetSqlAdminDriver() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		cnfMysql.User, cnfMysql.Passwd, cnfMysql.Ip, cnfMysql.Port, cnfMysql.Db)
}

func GetRedisOption() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cnfRedis.Ip, cnfRedis.Port),
		Password: cnfRedis.Passwd, // no password set
		DB:       0,               // use default DB
	}
}
