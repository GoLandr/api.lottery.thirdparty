package sssutils

import (
	"fmt"
	"mycommon/logs"
	"mycommon/mathstr"
	"mycommon/mymsg"
	"mycommon/utils"
	"runtime"
	"time"

	"api.lottery.thirdparty/global"
	"api.lottery.thirdparty/model"
)

const (
	REDIS_PRE_HALL_USR_ID       = "_HALLUSR_ID_"
	REDIS_PRE_HALL_USR_WS       = "_HALLUSR_WS_"
	REDIS_PRE_HALL_TID_GT_MP    = "_G_HALL_TIDGTMP" // 单独的名字
	REDIS_PRE_HALL_UID_LASTINFO = "_HALLUSR_LASTINFO_"
)

func RedisSetUsr(m *model.User, wsid string) {
	key := fmt.Sprint(REDIS_PRE_HALL_USR_ID, m.Id)
	val := mathstr.GetJsonPlainStr(m)
	logs.Debug("__存入 key", key, ",key:", val)
	err := global.RedisSet(key, val)
	if nil != err {
		logs.Error(err)
	}

	if wsid != "" {
		key := fmt.Sprint(REDIS_PRE_HALL_USR_WS, wsid)
		err := global.RedisSet(key, fmt.Sprint(m.Id))
		if nil != err {
			logs.Error(err)
		}
	}
}

func RedisUpdateUsr(m mymsg.Map) {
	if nil == m["id"] {
		return
	}
	key := fmt.Sprint(REDIS_PRE_HALL_USR_ID, m["id"])
	val, err := global.RedisGet(key)
	utils.ThrowError(err)

	var valMap map[string]interface{}
	if val != "" {
		mathstr.JsonUnmarsh(val, &valMap)
	} else {
		valMap = make(map[string]interface{})
	}
	for k, v := range m {
		valMap[k] = v
	}
	val = mathstr.GetJsonPlainStr(valMap)
	err = global.RedisSet(key, val)
	if nil != err {
		logs.Error(err)
	}
}

func RedisGetUsrById(id int) *model.User {
	key := fmt.Sprint(REDIS_PRE_HALL_USR_ID, id)
	val, err := global.RedisGet(key)
	if nil != err {
		logs.Error(err)
		return nil
	}
	if val == "" {
		return nil
	}
	var usr model.User
	mathstr.JsonUnmarsh(val, &usr)
	return &usr
}

func RedisGetUsrByWs(wsid string) *model.User {
	key := fmt.Sprint(REDIS_PRE_HALL_USR_WS, wsid)
	val, err := global.RedisGet(key)
	if nil != err {
		logs.Error(err)
		return nil
	}
	if val == "" {
		return nil
	}
	id := mathstr.Math2intDefault0(val)
	return RedisGetUsrById(id)
}

// REDIS_PRE_HALL_TID_GT_MP
func RedisPutTidGt(tid int, gt int) {
	key := REDIS_PRE_HALL_TID_GT_MP
	val, err := global.RedisGet(key)
	if nil != err {
		logs.Error(err)
	}

	var tidGtMap map[int]int
	if val != "" {
		mathstr.JsonUnmarsh(val, &tidGtMap)
	} else {
		tidGtMap = make(map[int]int)
	}
	tidGtMap[tid] = gt
	val = mathstr.GetJsonPlainStr(tidGtMap)
	global.RedisSet(key, val)
}

// true 存在 false 不存在
func RedisCheckTid(tid int) bool {
	key := REDIS_PRE_HALL_TID_GT_MP
	val, err := global.RedisGet(key)
	if nil != err {
		logs.Error(err)
	}

	var tidGtMap map[int]int
	if val != "" {
		mathstr.JsonUnmarsh(val, &tidGtMap)
		_, ok := tidGtMap[tid]
		return ok
	} else {
		return false
	}
}

func RedisDelTid(tid int) {
	key := REDIS_PRE_HALL_TID_GT_MP
	val, err := global.RedisGet(key)
	if nil != err {
		logs.Error(err)
	}

	var tidGtMap map[int]int
	if val != "" {
		mathstr.JsonUnmarsh(val, &tidGtMap)
	} else {
		tidGtMap = make(map[int]int)
	}
	delete(tidGtMap, tid)
	val = mathstr.GetJsonPlainStr(tidGtMap)
	global.RedisSet(key, val)
}

func GetGtByTid(tid int) int {
	key := REDIS_PRE_HALL_TID_GT_MP
	val, err := global.RedisGet(key)
	if nil != err {
		logs.Error(err)
	}

	var tidGtMap map[int]int
	if val != "" {
		mathstr.JsonUnmarsh(val, &tidGtMap)
	} else {
		tidGtMap = make(map[int]int)
	}
	return tidGtMap[tid]
}

func RedisAddLastInfo(uid int, info map[string]interface{}) {
	key := fmt.Sprint(REDIS_PRE_HALL_UID_LASTINFO, uid)
	val := mathstr.GetJsonPlainStr(info)
	err := global.RedisTimeSet(key, val, 2*24*time.Hour)
	if nil != err {
		logs.Error(err)
	}
}
func RedisDelLastInfo(uid int) {
	key := fmt.Sprint(REDIS_PRE_HALL_UID_LASTINFO, uid)
	err := global.RedisDel(key)
	if nil != err {
		logs.Error(err)
	}
}
func RedisGetLastInfo(uid int) map[string]interface{} {
	key := fmt.Sprint(REDIS_PRE_HALL_UID_LASTINFO, uid)
	val, err := global.RedisGet(key)
	if nil != err {
		logs.Error(err)
	}

	var info map[string]interface{}
	if val != "" {
		mathstr.JsonUnmarsh(val, &info)
	} else {
		info = make(map[string]interface{})
	}

	return info
}

func SyncRedis(nmp map[string]interface{}) {
	defer func() {
		if e := recover(); e != nil {
			err, ok := e.(error)
			if ok {
				// 日志记录
				for i := 3; i <= 7; i++ {
					_, f, line, ok := runtime.Caller(i)
					if !ok {
						continue
					}
					if i == 3 {
						logs.Error("__err:[", err, i, "]__fname:[", f, "]__line:[", line, "]")
					} else {
						logs.Error("__fname:[", f, "]__line:[", line, "]")
					}
				}
			}
		}
	}()

	ChangeStr2int(nmp, "id", "curroomid", "curtableid", "curseatid", "chips", "sourceid",
		"sex", "money", "cardcount", "registertime", "platformid", "level", "nn_level", "lastroomstatus",
		"todaywinlose", "nn_todaywinlose", "shared", "online")

	RedisUpdateUsr(mymsg.Map(nmp))
}

func ChangeStr2int(nmp map[string]interface{}, attributes ...string) {
	for i, _ := range attributes {
		tname := attributes[i]
		if nil == nmp[tname] {
			continue
		}
		nmp[tname] = mathstr.Math2intDefault0(nmp[tname])
	}
}
