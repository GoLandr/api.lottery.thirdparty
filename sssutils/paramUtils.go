package sssutils

import (
	"mycommon/mathstr"
	fmt "mycommon/myinherit/myfmt"
	"mycommon/utils"
	"strings"
)

// 检测上传文件 k-v 是否含必须
func CheckNotEmpty(detailItems []map[string]interface{}, colItemMap map[string]string, specialkeys []string) {

	if len(detailItems) <= 0 || len(specialkeys) <= 0 || len(colItemMap) <= 0 {
		return
	}

	enCnMap := map[string]string{}
	for enkey, cnkey := range colItemMap {
		if strings.HasSuffix(cnkey, "*") {
			cnkey = cnkey[:len(cnkey)-1]
		}
		enCnMap[cnkey] = enkey
	}

	for i, e := range detailItems {
		if nil == e {
			continue
		}

		for _, cnkey := range specialkeys {
			if strings.HasSuffix(cnkey, "*") {
				cnkey = cnkey[:len(cnkey)-1]
			}
			enkey := enCnMap[cnkey]
			tval := fmt.Sprint(e[enkey])
			if "" == tval {
				utils.ThrowErrorStr(mathstr.S_SFT(`第{1}行[{0}]为必填不能为空`, cnkey, (i + 2)))
			}
		}
	}
}
