package base

import (
	"fmt"
	"mycommon/logs"
	"mycommon/utils"
	"net/url"

	"api.lottery.thirdparty/common"
	"api.lottery.thirdparty/model"
)

func GetPushMenber(limit map[int][]*model.TLimit, maxVal int) []int {
	var retLst []int
	for key, lst := range limit {
		if maxVal >= key {
			for _, entry := range lst {
				//判断是否会员
				if entry.Is_valid == 1 && CheckValid(entry.Valid_date) {
					if menber, tok := GMenber[entry.Menber_id]; tok {
						retLst = append(retLst, menber.Code)
					}
				}
			}
		}
	}

	return retLst
}
func CheckValid(valid_data string) bool {
	flag := false
	if utils.NowDate() <= valid_data {
		flag = true
	}
	//	logs.Debug("flag_", flag)
	return flag
}
func sendMsgToFriend(retLst []int, msg string) {
	logs.Debug("sendMsgToFriend_", msg, "_retLst_", retLst)
	for _, v := range retLst {
		param := make(url.Values)
		param.Add("key", "123456")
		param.Add("code", fmt.Sprint(v))
		param.Add("msg", msg)
		_, err := common.Httppost(SEND_MSG_TO_FRIEND_URL, param)
		if err != nil {
			//panic()
		}
	}
}
