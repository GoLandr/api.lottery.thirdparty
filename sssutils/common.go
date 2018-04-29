package sssutils

import (
	"mycommon/logs"
	"mycommon/mathstr"
	"mycommon/utils"
	"net/http"
	"os/exec"

	"api.lottery.thirdparty/config"
	"api.lottery.thirdparty/global/myconst"

	"github.com/astaxie/beedb"
)

func HttpPost(postIp string, param map[string]interface{}) (map[string]interface{}, error, *http.Response) {
	jstr := mathstr.GetJsonPlainStr(param)
	res, err, resp := utils.HttpPost2(postIp, jstr)
	if nil != err {
		return nil, err, resp
	}
	var resMap map[string]interface{}
	mathstr.JsonUnmarsh(res, &resMap)

	return resMap, nil, resp
}

func CloseDB(orm *beedb.Model) {
	if nil != orm {
		if nil != orm.Db {
			err := orm.Db.Close()
			utils.ThrowError(err)
		}
	}
}

func ReloadCnf(gt int) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	shFileName := "reloadsss.sh"
	switch gt {
	case myconst.GT_SSS:
		shFileName = "reloadsss.sh"

	case myconst.GT_NN:
		shFileName = "reloadnn.sh"

	default:
		shFileName = "reloadsss.sh"
	}

	shDir := config.GeShPath()
	resbt, err := ExeShell(shDir + shFileName)
	utils.ThrowError(err)
	logs.Info("__result:", string(resbt))
}
func ExeShell(cmdPath string) ([]byte, error) {
	cmd := exec.Command(cmdPath)
	return cmd.Output()
}
