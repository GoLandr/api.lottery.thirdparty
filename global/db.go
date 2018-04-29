package global

import (
	"database/sql"
	"fmt"
	"mycommon/logs"
	"mycommon/utils"
	"os"
	"sync"

	"api.lottery.thirdparty/config"
	"api.lottery.thirdparty/model"

	"github.com/astaxie/beedb"
	_ "github.com/go-sql-driver/mysql"
)

var orm beedb.Model
var storeConfigs map[string]model.StoreConfig
var taskConfigs map[int][]model.TaskConfig
var taskMemRecords TaskMemRecord
var syscnfMap map[string]string

type TaskMemRecord struct {
	Data map[int]map[string]int
	M    sync.Mutex
}

func (this *TaskMemRecord) UpdateMemRecord(uid int, val map[string]int) {
	this.M.Lock()
	defer this.M.Unlock()
	this.Data[uid] = val
}
func (this *TaskMemRecord) GetMemRecord(uid int) map[string]int {
	this.M.Lock()
	defer this.M.Unlock()
	return this.Data[uid]
}

func init() {
	if len(os.Args) > 1 {
		return
	}
	initMethod()
	//	initStoreConfig()
	//	initTaskConfig()
	//	initTaskMemRecord()
	//	InitSyscnf()
}

func initMethod() {
	// 初始化mysql
	db, err := sql.Open("mysql", config.GetMysqlDriver())

	orm = beedb.New(db)
	beedb.OnDebug = true

	logs.Debug("mysqldriver:", config.GetMysqlDriver())
	if nil != err {
		logs.Error(err.Error())
		return
	}

	logs.Info("connect mysql successed")
}

func initStoreConfig() {
	var lst []model.StoreConfig
	err := orm.SetTable("sss_store_configs").FindAll(&lst)
	utils.ThrowError(err)
	storeConfigs = map[string]model.StoreConfig{}
	for _, e := range lst {
		storeConfigs[e.Sku] = e
	}
}

func initTaskConfig() {
	var lst []model.TaskConfig
	err := orm.SetTable("sss_task_configs").FindAll(&lst)
	utils.ThrowError(err)
	taskConfigs = map[int][]model.TaskConfig{}
	for _, e := range lst {
		taskConfigs[e.Type] = append(taskConfigs[e.Type], e)
	}
}

func initTaskMemRecord() {
	tmp := map[int]map[string]int{}
	taskMemRecords = TaskMemRecord{
		Data: tmp,
	}
}

func InitSyscnf() {
	reslst, err := orm.SetTable("t_sys_config").FindStringInterface()
	utils.ThrowError(err)
	syscnfMap = map[string]string{}
	for _, e := range reslst {
		syscnfMap[fmt.Sprint(e["ckey"])] = fmt.Sprint(e["cval"])
	}
}

func GetOrm() *beedb.Model {
	return &orm
}
func GetNewOrm() *beedb.Model {
	db, err := sql.Open("mysql", config.GetMysqlDriver())
	utils.ThrowError(err)
	tmpOrm := beedb.New(db)
	return &tmpOrm
}
func CloseOrm(orm *beedb.Model) {
	if nil == orm {
		return
	}
	if nil == orm.Db {
		return
	}
	err := orm.Db.Close()
	if nil != err {
		utils.ThrowError(err)
	}
}

func GetStoreConfig() map[string]model.StoreConfig {
	return storeConfigs
}
func GetTaskConfig() map[int][]model.TaskConfig {
	return taskConfigs
}
func GetTaskMemRecords() *TaskMemRecord {
	return &taskMemRecords
}
func GetSyscnf(ckey string) string {
	return syscnfMap[ckey]
}
